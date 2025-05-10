// racing/service/sort_test.go

package service

import (
	"testing"
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSortRaces(t *testing.T) {
	tests := []struct {
		name     string
		races    []*racing.Race
		filter   *racing.ListRacesRequestFilter
		wantDesc bool
		sortBy   racing.SortBy
	}{
		{
			name: "sorts by advertised start time ascending by default",
			races: []*racing.Race{
				{AdvertisedStartTime: timestamppb.New(time.Now().Add(time.Hour))},
				{AdvertisedStartTime: timestamppb.New(time.Now())},
			},
			filter:   nil,
			wantDesc: false,
			sortBy:   racing.SortBy_SORT_BY_ADVERTISED_START_TIME,
		},
		{
			name: "sorts by name ascending",
			races: []*racing.Race{
				{Name: "Race B"},
				{Name: "Race A"},
			},
			filter: &racing.ListRacesRequestFilter{
				SortBy:    racing.SortBy_SORT_BY_NAME.Enum(),
				SortOrder: racing.SortOrder_SORT_ORDER_ASC.Enum(),
			},
			wantDesc: false,
			sortBy:   racing.SortBy_SORT_BY_NAME,
		},
		{
			name: "sorts by number descending",
			races: []*racing.Race{
				{Number: 1},
				{Number: 2},
			},
			filter: &racing.ListRacesRequestFilter{
				SortBy:    racing.SortBy_SORT_BY_NUMBER.Enum(),
				SortOrder: racing.SortOrder_SORT_ORDER_DESC.Enum(),
			},
			wantDesc: true,
			sortBy:   racing.SortBy_SORT_BY_NUMBER,
		},
		{
			name: "handles nil start times",
			races: []*racing.Race{
				{AdvertisedStartTime: nil},
				{AdvertisedStartTime: timestamppb.New(time.Now())},
			},
			filter:   nil,
			wantDesc: false,
			sortBy:   racing.SortBy_SORT_BY_ADVERTISED_START_TIME,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortRaces(tt.races, tt.filter)

			// Verify sorting based on sort type
			for i := 1; i < len(got); i++ {
				switch tt.sortBy {
				case racing.SortBy_SORT_BY_ADVERTISED_START_TIME:
					if got[i].AdvertisedStartTime == nil {
						continue
					}
					if got[i-1].AdvertisedStartTime == nil {
						t.Errorf("nil start time found before non-nil start time")
						continue
					}

					iTime := got[i].AdvertisedStartTime.AsTime()
					prevTime := got[i-1].AdvertisedStartTime.AsTime()

					if tt.wantDesc {
						if iTime.After(prevTime) {
							t.Errorf("expected descending order, got ascending at index %d", i)
						}
					} else {
						if iTime.Before(prevTime) {
							t.Errorf("expected ascending order, got descending at index %d", i)
						}
					}

				case racing.SortBy_SORT_BY_NAME:
					if tt.wantDesc {
						if got[i].Name > got[i-1].Name {
							t.Errorf("expected descending order, got ascending at index %d", i)
						}
					} else {
						if got[i].Name < got[i-1].Name {
							t.Errorf("expected ascending order, got descending at index %d", i)
						}
					}

				case racing.SortBy_SORT_BY_NUMBER:
					if tt.wantDesc {
						if got[i].Number > got[i-1].Number {
							t.Errorf("expected descending order, got ascending at index %d", i)
						}
					} else {
						if got[i].Number < got[i-1].Number {
							t.Errorf("expected ascending order, got descending at index %d", i)
						}
					}
				}
			}
		})
	}
}
