package services

import (
	"context"

	"github.com/dronm/ds/pgds"
	"github.com/dronm/session"
)

// ProductCatService is a service for managing product categories
type EventService struct {
	DB      *pgds.PgProvider
	Session session.Session
}

func NewEventService(db *pgds.PgProvider, sess session.Session) *EventService {
	return &EventService{DB: db, Session: sess}
}

func (s *NotifAppService) Subscribe(ctx context.Context, events []string)  error {
	return nil
}

func (s *NotifAppService) Unsubscribe(ctx context.Context, events []string)  error {
	return nil
}
