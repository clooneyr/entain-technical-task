// sports/service/status_test.go

package service

import (
	"testing"
	"time"

	"git.neds.sh/matty/entain/sports/proto/sports"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCheckEventStatus(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		startTime time.Time
		want      sports.EventStatus
	}{
		{
			name:      "past event should be CLOSED",
			startTime: now.Add(-1 * time.Hour),
			want:      sports.EventStatus_CLOSED,
		},
		{
			name:      "future event should be OPEN",
			startTime: now.Add(1 * time.Hour),
			want:      sports.EventStatus_OPEN,
		},
		// Add more test cases if you have additional status logic
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &sports.Event{
				AdvertisedStartTime: timestamppb.New(tt.startTime),
			}

			status := CheckEventStatus(event)
			assert.Equal(t, tt.want, status)
		})
	}
}

func TestUpdateEventsStatus(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name   string
		events []*sports.Event
		want   []sports.EventStatus
	}{
		{
			name: "updates statuses correctly",
			events: []*sports.Event{
				{
					Id:                  1,
					Name:                "Past Event",
					AdvertisedStartTime: timestamppb.New(now.Add(-1 * time.Hour)),
				},
				{
					Id:                  2,
					Name:                "Future Event",
					AdvertisedStartTime: timestamppb.New(now.Add(1 * time.Hour)),
				},
				{
					Id:   3,
					Name: "No Time Event",
					// No advertised start time
				},
			},
			want: []sports.EventStatus{
				sports.EventStatus_CLOSED,
				sports.EventStatus_OPEN,
				sports.EventStatus_UNSPECIFIED,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedEvents := UpdateEventsStatus(tt.events)

			assert.Equal(t, len(tt.events), len(updatedEvents))

			for i, status := range tt.want {
				assert.Equal(t, status, updatedEvents[i].Status)
			}

			// Also check that other fields were copied correctly
			for i, original := range tt.events {
				assert.Equal(t, original.Id, updatedEvents[i].Id)
				assert.Equal(t, original.Name, updatedEvents[i].Name)
			}
		})
	}
}
