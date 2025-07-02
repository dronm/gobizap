package ws

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dronm/gobizap/v2/logger"

	"github.com/dronm/session"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins
		},
	}
	clients   = make(map[string]Client)
	clientsMu sync.RWMutex
)

func Upgrade(w http.ResponseWriter, r *http.Request, sess session.Session) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return fmt.Errorf("upgrader.Upgrade() failed: %v", err)
	}

	clientID := ""

	clientsMu.Lock()
	clients[clientID] = Client{ID: clientID, Conn: conn}
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, clientID)
		clientsMu.Unlock()
		conn.Close()
	}()

	done := make(chan struct{})
	go func() {
		<-r.Context().Done()
		logger.Logger.Info("Request context closed, terminating WebSocket connection")
		conn.Close()
		close(done)
	}()

	for {
		select {
		case <-done:
			return nil

		default:
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("conn.ReadMessage() failed: %v", err)
			}

			logger.Logger.Debugf("Received: type:%d, msg:%s\n", msgType, msg)

			// Echo the message back
			if err := conn.WriteMessage(msgType, msg); err != nil {
				return fmt.Errorf("conn.WriteMessage() failed: %v", err)
			}
		}
	}
}

func CleanupConnections(interval int) {
	for {
		time.Sleep(time.Duration(interval) * time.Second) // Periodic cleanup every 30 seconds
		clientsMu.Lock()
		for _, client := range clients {
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				client.Conn.Close()
				delete(clients, client.ID)
			}
		}
		clientsMu.Unlock()
	}
}

func ShutdownWebSockets(ctx context.Context, appShutdownTimeout time.Duration) {
	clientsMu.Lock()
	conns := make([]*websocket.Conn, 0, len(clients))
	for _, c := range clients {
		conns = append(conns, c.Conn)
	}
	clientsMu.Unlock()

	for _, conn := range conns {
		if ctx.Err() != nil {
			break
		}

		// Per-connection timeout context
		connCtx, cancel := context.WithTimeout(ctx, appShutdownTimeout)
		errChan := make(chan error, 1)

		go func() {
			err := conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server shutting down"),
			)
			conn.Close()
			errChan <- err
		}()

		select {
		case <-connCtx.Done():
			// Timed out or canceled
		case <-ctx.Done():
			// Global shutdown canceled
			cancel()
			return
		case <-errChan:
			// Close completed
		}

		cancel()
	}

}

func PublishEvent(publisherID string, payload []byte) {
	clientsMu.RLock()
	for _, client := range clients {
		if client.ID == publisherID {
			continue
		}
		if err := client.Conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			logger.Logger.Errorf("PublishEvent conn.WriteMessage: %v", err)
		}
	}
	clientsMu.RUnlock()
}
