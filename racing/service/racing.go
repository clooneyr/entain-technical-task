package service

import (
	"database/sql"
	"errors"

	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Racing interface {
	// ListRaces will return a collection of races.
	ListRaces(ctx context.Context, in *racing.ListRacesRequest) (*racing.ListRacesResponse, error)
	// GetRace will return a single race by its ID.
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
	races, err := s.racesRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	// Update statuses
	races = UpdateRacesStatus(races)

	// Apply sorting
	sortedRaces := SortRaces(races, in.Filter)

	return &racing.ListRacesResponse{Races: sortedRaces}, nil
}

// GetRace retrieves a single race by its ID
func (s *racingService) GetRace(ctx context.Context, req *racing.GetRaceRequest) (*racing.Race, error) {
	// Validate the request
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid race ID")
	}

	// Get the race from the database
	race, err := s.racesRepo.GetRace(req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "race not found")
		}
		return nil, status.Error(codes.Internal, "failed to get race")
	}

	// Update the race status
	race.Status = CheckRaceStatus(race)

	return race, nil
}
