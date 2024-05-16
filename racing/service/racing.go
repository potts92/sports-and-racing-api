package service

import (
	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
	// GetRace will return a single race.
	GetRace(ctx context.Context, in *racing.GetRaceRequest) (*racing.Race, error)
}

// racingService implements the Racing interface.
type racingService struct {
	racesRepo db.RacesRepo
}

// NewRacingService instantiates and returns a new racingService.
func NewRacingService(racesRepo db.RacesRepo) Racing {
	return &racingService{racesRepo}
}

func (s *racingService) ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error) {
	races, err := s.racesRepo.List(in.Filter, in.Sort)
	if err != nil {
		return nil, err
	}

	return &racing.ListRacesResponse{Races: races}, nil
}

func (s *racingService) GetRace(ctx context.Context, in *racing.GetRaceRequest) (*racing.Race, error) {
	race, err := s.racesRepo.Get(in.Id)

	if err != nil {
		return nil, err
	}

	if race == nil {
		return nil, status.Error(codes.NotFound, "race not found")
	}

	return race, nil
}
