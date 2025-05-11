package service

import (
	"context"

	"git.neds.sh/matty/entain/sports/db"
	"git.neds.sh/matty/entain/sports/proto/sports"
)

// Sports defines the interface for sports service operations.
type Sports interface {
	// ListEvents returns a collection of sporting events.
	ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error)
}

// sportsService implements the Sports interface.
type sportsService struct {
	sports.UnimplementedSportsServer
	eventsRepo db.EventsRepo
}

// NewSportsService creates a new sports service.
func NewSportsService(eventsRepo db.EventsRepo) sports.SportsServer {
	return &sportsService{
		eventsRepo: eventsRepo,
	}
}

// ListEvents returns a list of sports events based on the provided filter.
func (s *sportsService) ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error) {
	// Get events from repository
	events, err := s.eventsRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	// Update event statuses (OPEN/CLOSED based on advertised start time)
	events = UpdateEventsStatus(events)

	// Apply sorting if filter provided
	sortedEvents := SortEvents(events, in.Filter)

	return &sports.ListEventsResponse{Events: sortedEvents}, nil
}
