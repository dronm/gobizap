package eventServer

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dronm/session"
	"github.com/dronm/gobizap/v2/api"

	"github.com/dronm/gobizap/v2/database"
	"github.com/dronm/gobizap/v2/logger"
	"github.com/dronm/gobizap/v2/ws"
)

const (
	DB_ACQUIRE_CONN_WAIT     = 100
	DB_MAX_ACQUIRE_CONN_WAIT = 60000
)

//---- UniqEvents -----

// UniqEvents holds counters for events.
// This stucture is never cleanedup, as
// it holds event id`s and there number is finite.
// If an event is not needed any more (not used by any connection),
// it is removed from server.
type UniqEvents struct {
	mx sync.Mutex
	m  map[string]int // unique event counter
}

// AddEvent puts new event with the given ID to database channel.
// In pg it executes listen command.
func (e *UniqEvents) AddEvent(dbEventID string, qChan chan string) {
	e.mx.Lock()
	cnt, ok := e.m[dbEventID]
	if !ok {
		qChan <- `LISTEN "` + dbEventID + `"`
		e.m[dbEventID] = 1
	} else {
		e.m[dbEventID] = cnt + 1
	}
	e.mx.Unlock()
}

// RemoveEvent puts unlisten command to database channel.
func (e *UniqEvents) RemoveEvent(dbEventID string, qChan chan string) {
	e.mx.Lock()
	if cnt, ok := e.m[dbEventID]; ok {
		if cnt == 1 {
			qChan <- `UNLISTEN "` + dbEventID + `"`
			delete(e.m, dbEventID)
		} else {
			e.m[dbEventID] = cnt - 1
		}
	}
	e.mx.Unlock()
}

// EventCount returns count for a specific event ID.
// The second parameter is true, if the event exists.
// Count can be 0, in nobody is listening to this event.
func (e *UniqEvents) EventCount(eventID string) (int, bool) {
	e.mx.Lock()
	defer e.mx.Unlock()

	value, ok := e.m[eventID]

	return value, ok
}

// TotalEventCount return total event count in a map.
func (e *UniqEvents) TotalEventCount() int {
	e.mx.Lock()
	defer e.mx.Unlock()

	return len(e.m)
}

// ---- EventServer ----
type EventServer struct {
	DbPool      *pgxpool.Pool //
	DbQuery     chan string   // for notification queries
	Events      *UniqEvents   // count of unique events for db
	LocalEvents map[string]struct{}

	ctx        context.Context
	cancel     context.CancelFunc
	cancelDone chan struct{}

	sess session.Session

	loopPause time.Duration
	// ReconnectParams waitStrat.WaitStrategy
}

func NewEventServer(localEvents map[string]struct{}) *EventServer {
	return &EventServer{LocalEvents: localEvents}
}

func (s *EventServer) OnNotification(_ *pgconn.PgConn, n *pgconn.Notification) {
	logger.Logger.Debugf("OnNotification Channel:%s, Payload:%s", n.Channel, n.Payload)
	srvMeth := strings.Split(n.Channel, ".")
	if len(srvMeth) < 2 {
		logger.Logger.Errorf("OnNotification invalid service.method signature for:", n.Channel)
		return
	}

	// local event
	if s.LocalEvents != nil {
		if _, ok := s.LocalEvents[n.Channel]; ok {
			// local cosumer, execute service function
			params, err := api.UnmarshalParams([]byte(n.Payload))
			if err != nil {
				logger.Logger.Errorf("OnNotification api.UnmarshalParams: %v", err)
				return
			}

			logger.Logger.Debugf("executing local service call to %s.%s with params: %v", srvMeth[0], srvMeth[1], params)

			api.CallMethod(s.ctx, srvMeth[0], srvMeth[1], params,
				&api.ServiceContext{DB: database.DB, Session: s.sess},
			)
			return
		}
	}

	// publish event for all client consumers
	evPayload := fmt.Sprintf(`{"id": "%s", "params": %s}`, n.Channel, n.Payload)
	ws.PublishEvent("", []byte(evPayload)) // send to all
}

func (s *EventServer) Start() {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.cancelDone = make(chan struct{})
	defer close(s.cancelDone)

	s.DbQuery = make(chan string, 10)
	s.Events = &UniqEvents{m: make(map[string]int, 0)}

	if s.LocalEvents != nil {
		s.Events.mx.Lock()
		for evntId := range s.LocalEvents {
			s.Events.m[evntId] = 1 // one instance only
		}
		s.Events.mx.Unlock()
	}

	logger.Logger.Info("EventServer: started, loop pause: %v", s.loopPause)

	dbAcquireWait := DB_ACQUIRE_CONN_WAIT

	for {
		var conn *pgxpool.Conn

		select {
		case <-s.ctx.Done():
			logger.Logger.Debug("EventServer breaking loop on stop request")
			return
		default:
			var err error
			conn, err = s.DbPool.Acquire(s.ctx)
			if err != nil {
				if dbAcquireWait > DB_MAX_ACQUIRE_CONN_WAIT {
					dbAcquireWait = DB_MAX_ACQUIRE_CONN_WAIT
				}
				logger.Logger.Errorf("EventSrv DbPool.Acquire(): %v", err)

				time.Sleep(time.Duration(dbAcquireWait) * time.Millisecond)
				dbAcquireWait = dbAcquireWait * 2
				continue
			}
		}

		for evnt := range s.Events.m {
			logger.Logger.Debugf("EventSrv LocalEvent: %s", evnt)
			conn.Exec(s.ctx, `LISTEN "`+evnt+`"`)
		}

		logger.Logger.Debug("EventSrv acquired connection")

		dbAcquireWait = DB_ACQUIRE_CONN_WAIT

		var q string
		for {
			select {
			case <-s.ctx.Done():
				return
			case q = <-s.DbQuery:
			default:
				q = ";"
			}

			if _, err := conn.Exec(s.ctx, q); err != nil {
				if s.ctx.Err() == context.Canceled {
					conn.Release()
					return
				}
				logger.Logger.Errorf("EventSrv conn.Exec(): %v on query: %s", err, q)

				conn.Release()
				break
			}

			// paause
			select {
			case <-s.ctx.Done():
			case <-time.After(s.loopPause):
			}
		}
	}
}

func (s *EventServer) Stop(ctx context.Context) {
	if s.cancel == nil {
		return
	}
	logger.Logger.Debug("EventServer stopping on request...")
	s.cancel()

	select {
	case <-ctx.Done():
	case <-s.cancelDone:
	}
	logger.Logger.Info("EventServer stopped")
}
