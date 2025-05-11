// sports/service/sort_test.go

package service

import (
	"testing"
	"time"

	"git.neds.sh/matty/entain/sports/proto/sports"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSortEvents(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		events   []*sports.Event
		filter   *sports.ListEventsRequestFilter
		wantDesc bool
		sortBy   sports.SortBy
		checkFn  func(t *testing.T, events []*sports.Event, wantDesc bool)
	}{
		{
			name: "sorts by advertised start time ascending by default",
			events: []*sports.Event{
				{AdvertisedStartTime: timestamppb.New(now.Add(time.Hour))},
				{AdvertisedStartTime: timestamppb.New(now)},
			},
			filter:   nil,
			wantDesc: false,
			sortBy:   sports.SortBy_SORT_BY_ADVERTISED_START_TIME,
			checkFn: func(t *testing.T, events []*sports.Event, wantDesc bool) {
				for i := 1; i < len(events); i++ {
					if events[i].AdvertisedStartTime == nil || events[i-1].AdvertisedStartTime == nil {
						continue
					}

					iTime := events[i].AdvertisedStartTime.AsTime()
					prevTime := events[i-1].AdvertisedStartTime.AsTime()

					if wantDesc {
						assert.True(t, iTime.Before(prevTime) || iTime.Equal(prevTime),
							"expected descending order, got ascending at index %d", i)
					} else {
						assert.True(t, iTime.After(prevTime) || iTime.Equal(prevTime),
							"expected ascending order, got descending at index %d", i)
					}
				}
			},
		},
		{
			name: "sorts by name ascending",
			events: []*sports.Event{
				{Name: "Event B"},
				{Name: "Event A"},
			},
			filter: &sports.ListEventsRequestFilter{
				SortBy:    sports.SortBy_SORT_BY_NAME,
				SortOrder: sports.SortOrder_SORT_ORDER_ASC,
			},
			wantDesc: false,
			sortBy:   sports.SortBy_SORT_BY_NAME,
			checkFn: func(t *testing.T, events []*sports.Event, wantDesc bool) {
				for i := 1; i < len(events); i++ {
					if wantDesc {
						assert.True(t, events[i].Name <= events[i-1].Name,
							"expected descending order, got ascending at index %d", i)
					} else {
						assert.True(t, events[i].Name >= events[i-1].Name,
							"expected ascending order, got descending at index %d", i)
					}
				}
			},
		},
		{
			name: "sorts by venue descending",
			events: []*sports.Event{
				{Venue: "Venue A"},
				{Venue: "Venue B"},
			},
			filter: &sports.ListEventsRequestFilter{
				SortBy:    sports.SortBy_SORT_BY_VENUE,
				SortOrder: sports.SortOrder_SORT_ORDER_DESC,
			},
			wantDesc: true,
			sortBy:   sports.SortBy_SORT_BY_VENUE,
			checkFn: func(t *testing.T, events []*sports.Event, wantDesc bool) {
				for i := 1; i < len(events); i++ {
					if wantDesc {
						assert.True(t, events[i].Venue <= events[i-1].Venue,
							"expected descending order, got ascending at index %d", i)
					} else {
						assert.True(t, events[i].Venue >= events[i-1].Venue,
							"expected ascending order, got descending at index %d", i)
					}
				}
			},
		},
		{
			name: "handles nil start times",
			events: []*sports.Event{
				{AdvertisedStartTime: nil},
				{AdvertisedStartTime: timestamppb.New(now)},
			},
			filter:   nil,
			wantDesc: false,
			sortBy:   sports.SortBy_SORT_BY_ADVERTISED_START_TIME,
			checkFn: func(t *testing.T, events []*sports.Event, wantDesc bool) {
				// Nil timestamps should be at the end
				if len(events) >= 2 && events[0].AdvertisedStartTime == nil {
					assert.Fail(t, "nil start time should be sorted to the end")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SortEvents(tt.events, tt.filter)

			// Skip tests with less than 2 events
			if len(got) < 2 {
				return
			}

			// Use the appropriate check function
			tt.checkFn(t, got, tt.wantDesc)
		})
	}
}

func TestSortByAdvertisedStartTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		events   []*sports.Event
		order    sports.SortOrder
		expected []*sports.Event
	}{
		{
			name: "ascending order",
			events: []*sports.Event{
				{Id: 1, AdvertisedStartTime: timestamppb.New(now.Add(time.Hour))},
				{Id: 2, AdvertisedStartTime: timestamppb.New(now)},
			},
			order: sports.SortOrder_SORT_ORDER_ASC,
			expected: []*sports.Event{
				{Id: 2, AdvertisedStartTime: timestamppb.New(now)},
				{Id: 1, AdvertisedStartTime: timestamppb.New(now.Add(time.Hour))},
			},
		},
		{
			name: "descending order",
			events: []*sports.Event{
				{Id: 1, AdvertisedStartTime: timestamppb.New(now)},
				{Id: 2, AdvertisedStartTime: timestamppb.New(now.Add(time.Hour))},
			},
			order: sports.SortOrder_SORT_ORDER_DESC,
			expected: []*sports.Event{
				{Id: 2, AdvertisedStartTime: timestamppb.New(now.Add(time.Hour))},
				{Id: 1, AdvertisedStartTime: timestamppb.New(now)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortByAdvertisedStartTime(tt.events, tt.order)

			// Compare IDs as times might have slight differences due to serialization
			assert.Equal(t, len(tt.expected), len(result))
			for i := 0; i < len(result); i++ {
				assert.Equal(t, tt.expected[i].Id, result[i].Id)
			}
		})
	}
}

func TestSortByName(t *testing.T) {
	tests := []struct {
		name     string
		events   []*sports.Event
		order    sports.SortOrder
		expected []*sports.Event
	}{
		{
			name: "ascending order",
			events: []*sports.Event{
				{Id: 1, Name: "Event B"},
				{Id: 2, Name: "Event A"},
			},
			order: sports.SortOrder_SORT_ORDER_ASC,
			expected: []*sports.Event{
				{Id: 2, Name: "Event A"},
				{Id: 1, Name: "Event B"},
			},
		},
		{
			name: "descending order",
			events: []*sports.Event{
				{Id: 1, Name: "Event A"},
				{Id: 2, Name: "Event B"},
			},
			order: sports.SortOrder_SORT_ORDER_DESC,
			expected: []*sports.Event{
				{Id: 2, Name: "Event B"},
				{Id: 1, Name: "Event A"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortByName(tt.events, tt.order)

			assert.Equal(t, len(tt.expected), len(result))
			for i := 0; i < len(result); i++ {
				assert.Equal(t, tt.expected[i].Id, result[i].Id)
			}
		})
	}
}

func TestSortByVenue(t *testing.T) {
	tests := []struct {
		name     string
		events   []*sports.Event
		order    sports.SortOrder
		expected []*sports.Event
	}{
		{
			name: "ascending order",
			events: []*sports.Event{
				{Id: 1, Venue: "Venue B"},
				{Id: 2, Venue: "Venue A"},
			},
			order: sports.SortOrder_SORT_ORDER_ASC,
			expected: []*sports.Event{
				{Id: 2, Venue: "Venue A"},
				{Id: 1, Venue: "Venue B"},
			},
		},
		{
			name: "descending order",
			events: []*sports.Event{
				{Id: 1, Venue: "Venue A"},
				{Id: 2, Venue: "Venue B"},
			},
			order: sports.SortOrder_SORT_ORDER_DESC,
			expected: []*sports.Event{
				{Id: 2, Venue: "Venue B"},
				{Id: 1, Venue: "Venue A"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sortByVenue(tt.events, tt.order)

			assert.Equal(t, len(tt.expected), len(result))
			for i := 0; i < len(result); i++ {
				assert.Equal(t, tt.expected[i].Id, result[i].Id)
			}
		})
	}
}
