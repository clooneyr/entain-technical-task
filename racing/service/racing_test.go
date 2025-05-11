package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.neds.sh/matty/entain/racing/db"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultTimeout = 5 * time.Second
)

// RacesRepoMock is a mock implementation of db.RacesRepo
type RacesRepoMock struct {
	mock.Mock
}

// RacesRepoMock implements db.RacesRepo
var _ db.RacesRepo = (*RacesRepoMock)(nil)

func (m *RacesRepoMock) Init() error {
	args := m.Called()
	return args.Error(0)
}

func (m *RacesRepoMock) List(filter *racing.ListRacesRequestFilter) ([]*racing.Race, error) {
	args := m.Called(filter)
	return args.Get(0).([]*racing.Race), args.Error(1)
}

func (m *RacesRepoMock) GetRace(id int64) (*racing.Race, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*racing.Race), args.Error(1)
}

// assertRacesEqual compares two slices of races, ignoring timestamps
func assertRacesEqual(t *testing.T, expected, actual []*racing.Race) {
	t.Helper()
	require.Equal(t, len(expected), len(actual), "number of races should match")

	for i := range expected {
		assert.Equal(t, expected[i].Id, actual[i].Id, "race ID should match")
		assert.Equal(t, expected[i].Visible, actual[i].Visible, "race visibility should match")
		assert.Equal(t, expected[i].Name, actual[i].Name, "race name should match")
		assert.Equal(t, expected[i].MeetingId, actual[i].MeetingId, "meeting ID should match")
		assert.Equal(t, expected[i].Number, actual[i].Number, "race number should match")
	}
}

func TestListRaces(t *testing.T) {
	// Create a fixed timestamp for testing
	now := timestamppb.Now()

	tests := []struct {
		desc      string
		filter    *racing.ListRacesRequestFilter
		mockRaces []*racing.Race
		wantRaces []*racing.Race
		wantErr   error
	}{
		{
			desc: "when visible_only is true, returns only visible races",
			filter: &racing.ListRacesRequestFilter{
				VisibleOnly: boolPtr(t, true),
			},
			mockRaces: []*racing.Race{
				{
					Id:                  1,
					Visible:             true,
					Name:                "Race 1",
					MeetingId:           100,
					Number:              1,
					AdvertisedStartTime: now,
				},
			},
			wantRaces: []*racing.Race{
				{
					Id:                  1,
					Visible:             true,
					Name:                "Race 1",
					MeetingId:           100,
					Number:              1,
					AdvertisedStartTime: now,
				},
			},
			wantErr: nil,
		},
		{
			desc: "when visible_only is false, returns only non-visible races",
			filter: &racing.ListRacesRequestFilter{
				VisibleOnly: boolPtr(t, false),
			},
			mockRaces: []*racing.Race{
				{
					Id:                  2,
					Visible:             false,
					Name:                "Race 2",
					MeetingId:           100,
					Number:              2,
					AdvertisedStartTime: now,
				},
			},
			wantRaces: []*racing.Race{
				{
					Id:                  2,
					Visible:             false,
					Name:                "Race 2",
					MeetingId:           100,
					Number:              2,
					AdvertisedStartTime: now,
				},
			},
			wantErr: nil,
		},
		{
			desc:   "when visible_only is not provided, returns all races",
			filter: &racing.ListRacesRequestFilter{},
			mockRaces: []*racing.Race{
				{
					Id:                  1,
					Visible:             true,
					Name:                "Race 1",
					MeetingId:           100,
					Number:              1,
					AdvertisedStartTime: now,
				},
				{
					Id:                  2,
					Visible:             false,
					Name:                "Race 2",
					MeetingId:           100,
					Number:              2,
					AdvertisedStartTime: now,
				},
			},
			wantRaces: []*racing.Race{
				{
					Id:                  1,
					Visible:             true,
					Name:                "Race 1",
					MeetingId:           100,
					Number:              1,
					AdvertisedStartTime: now,
				},
				{
					Id:                  2,
					Visible:             false,
					Name:                "Race 2",
					MeetingId:           100,
					Number:              2,
					AdvertisedStartTime: now,
				},
			},
			wantErr: nil,
		},
		{
			desc: "when repository returns error, propagates error",
			filter: &racing.ListRacesRequestFilter{
				VisibleOnly: boolPtr(t, true),
			},
			mockRaces: nil,
			wantRaces: nil,
			wantErr:   assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// Setup
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()

			mockRepo := new(RacesRepoMock)
			mockRepo.On("List", mock.MatchedBy(func(filter *racing.ListRacesRequestFilter) bool {
				if tt.filter.VisibleOnly == nil {
					return filter.VisibleOnly == nil
				}
				return filter.VisibleOnly != nil && *filter.VisibleOnly == *tt.filter.VisibleOnly
			})).Return(tt.mockRaces, tt.wantErr)

			svc := NewRacingService(mockRepo)

			// Execute
			got, err := svc.ListRaces(ctx, &racing.ListRacesRequest{
				Filter: tt.filter,
			})

			// Verify
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			assertRacesEqual(t, tt.wantRaces, got.Races)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetRace(t *testing.T) {
	// Create a fixed timestamp for testing
	now := timestamppb.Now()

	tests := []struct {
		desc     string
		id       int64
		mockRace *racing.Race
		wantRace *racing.Race
		wantErr  error
	}{
		{
			desc: "when race exists, returns race",
			id:   1,
			mockRace: &racing.Race{
				Id:                  1,
				Visible:             true,
				Name:                "Race 1",
				MeetingId:           100,
				Number:              1,
				AdvertisedStartTime: now,
			},
			wantRace: &racing.Race{
				Id:                  1,
				Visible:             true,
				Name:                "Race 1",
				MeetingId:           100,
				Number:              1,
				AdvertisedStartTime: now,
			},
			wantErr: nil,
		},
		{
			desc:     "when race doesn't exist, returns error",
			id:       999,
			mockRace: nil,
			wantRace: nil,
			wantErr:  sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// Setup
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()

			mockRepo := new(RacesRepoMock)
			mockRepo.On("GetRace", tt.id).Return(tt.mockRace, tt.wantErr)

			svc := NewRacingService(mockRepo)

			// Execute
			got, err := svc.GetRace(ctx, &racing.GetRaceRequest{
				Id: tt.id,
			})

			// Verify
			if tt.wantErr != nil {
				require.Error(t, err)
				if tt.wantErr == sql.ErrNoRows {
					assert.Equal(t, codes.NotFound, status.Code(err))
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantRace.Id, got.Id)
			assert.Equal(t, tt.wantRace.Visible, got.Visible)
			assert.Equal(t, tt.wantRace.Name, got.Name)
			assert.Equal(t, tt.wantRace.MeetingId, got.MeetingId)
			assert.Equal(t, tt.wantRace.Number, got.Number)
			mockRepo.AssertExpectations(t)
		})
	}
}

// boolPtr returns a pointer to the given bool.
func boolPtr(t *testing.T, b bool) *bool {
	t.Helper()
	return &b
}
