// sports/service/sports_test.go

package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"git.neds.sh/matty/entain/sports/proto/sports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mock repository for testing
type mockEventsRepo struct {
	mock.Mock
}

func (m *mockEventsRepo) Init() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockEventsRepo) List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error) {
	args := m.Called(filter)
	return args.Get(0).([]*sports.Event), args.Error(1)
}

func TestSportsService_ListEvents(t *testing.T) {
	now := time.Now()
	mockEvents := []*sports.Event{
		{
			Id:                  1,
			Name:                "Football Match",
			AdvertisedStartTime: timestamppb.New(now.Add(time.Hour)),
			Visible:             true,
			Venue:               "Anfield",
			SportType:           "Soccer",
			Competitors:         []string{"Liverpool FC", "Manchester United"},
		},
		{
			Id:                  2,
			Name:                "Basketball Game",
			AdvertisedStartTime: timestamppb.New(now.Add(-time.Hour)),
			Visible:             true,
			Venue:               "Staples Center",
			SportType:           "Basketball",
			Competitors:         []string{"Los Angeles Lakers", "Boston Celtics"},
		},
	}

	tests := []struct {
		name    string
		filter  *sports.ListEventsRequestFilter
		setup   func(m *mockEventsRepo)
		wantLen int
		wantErr bool
	}{
		{
			name:   "successful fetch all events",
			filter: &sports.ListEventsRequestFilter{},
			setup: func(m *mockEventsRepo) {
				m.On("List", mock.Anything).Return(mockEvents, nil)
			},
			wantLen: 2,
			wantErr: false,
		},
		{
			name: "successful fetch with visibility filter",
			filter: &sports.ListEventsRequestFilter{
				VisibleOnly: true,
			},
			setup: func(m *mockEventsRepo) {
				m.On("List", mock.Anything).Return([]*sports.Event{mockEvents[0]}, nil)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "error from repository",
			filter: &sports.ListEventsRequestFilter{},
			setup: func(m *mockEventsRepo) {
				m.On("List", mock.Anything).Return([]*sports.Event{}, errors.New("repository error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockEventsRepo)
			tt.setup(mockRepo)

			s := &sportsService{
				eventsRepo: mockRepo,
			}

			resp, err := s.ListEvents(context.Background(), &sports.ListEventsRequest{
				Filter: tt.filter,
			})

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantLen, len(resp.Events))

			// Verify status is correctly set
			for _, event := range resp.Events {
				if event.AdvertisedStartTime != nil {
					startTime := event.AdvertisedStartTime.AsTime()
					expectedStatus := sports.EventStatus_OPEN
					if startTime.Before(time.Now()) {
						expectedStatus = sports.EventStatus_CLOSED
					}
					assert.Equal(t, expectedStatus, event.Status)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
