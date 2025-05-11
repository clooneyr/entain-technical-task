package service

import (
	"testing"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCheckRaceStatus(t *testing.T) {
	tests := []struct {
		name           string
		advertisedTime *timestamppb.Timestamp
		want           racing.RaceStatus
	}{
		{
			name:           "future time returns OPEN",
			advertisedTime: timestamppb.New(time.Now().Add(time.Hour)),
			want:           racing.RaceStatus_OPEN,
		},
		{
			name:           "past time returns CLOSED",
			advertisedTime: timestamppb.New(time.Now().Add(-time.Hour)),
			want:           racing.RaceStatus_CLOSED,
		},
		{
			name:           "nil time returns UNSPECIFIED",
			advertisedTime: nil,
			want:           racing.RaceStatus_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			race := &racing.Race{
				AdvertisedStartTime: tt.advertisedTime,
			}
			got := CheckRaceStatus(race)
			if got != tt.want {
				t.Errorf("CheckRaceStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
