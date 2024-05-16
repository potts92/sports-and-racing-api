package service

import (
	"github.com/potts92/sports-and-racing-api/sports/db"
	"github.com/potts92/sports-and-racing-api/sports/proto/sports"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Sports interface {
	// ListEvents will return a collection of events.
	ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error)
	// GetEvent will return a single event.
	GetEvent(ctx context.Context, in *sports.GetEventRequest) (*sports.Event, error)
	// UpdateScore will update an event's score.
	UpdateScore(ctx context.Context, in *sports.UpdateScoreRequest) (*sports.Event, error)
}

// sportsService implements the Sports interface.
type sportsService struct {
	eventsRepo db.EventsRepo
}

// NewSportsService instantiates and returns a new sportsService.
func NewSportsService(eventsRepo db.EventsRepo) Sports {
	return &sportsService{eventsRepo}
}

func (s *sportsService) ListEvents(ctx context.Context, in *sports.ListEventsRequest) (*sports.ListEventsResponse, error) {
	events, err := s.eventsRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &sports.ListEventsResponse{Events: events}, nil
}

func (s *sportsService) GetEvent(ctx context.Context, in *sports.GetEventRequest) (*sports.Event, error) {
	event, err := s.eventsRepo.Get(in.Id)

	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, status.Error(codes.NotFound, "event not found")
	}

	return event, nil
}

func (s *sportsService) UpdateScore(ctx context.Context, in *sports.UpdateScoreRequest) (*sports.Event, error) {
	event, err := s.eventsRepo.UpdateScore(in.Id, in.HomeScore, in.AwayScore, in.Finalised)

	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, status.Error(codes.NotFound, "event not found")
	}

	return event, nil
}
