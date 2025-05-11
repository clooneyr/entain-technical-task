package service

import (
	"time"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

// CheckEventStatus determines if an event is OPEN or CLOSED based on its advertised start time.
// Events with start times in the past are considered CLOSED, otherwise they are OPEN.
// If no advertised start time is available, the status is UNSPECIFIED.
func CheckEventStatus(event *sports.Event) sports.EventStatus {
	if event.AdvertisedStartTime == nil {
		return sports.EventStatus_UNSPECIFIED
	}

	advertisedTime := event.AdvertisedStartTime.AsTime()
	if advertisedTime.Before(time.Now()) {
		return sports.EventStatus_CLOSED
	}
	return sports.EventStatus_OPEN
}

// UpdateEventsStatus updates the status field for all events in the provided slice.
// It returns a new slice with updated event objects.
func UpdateEventsStatus(events []*sports.Event) []*sports.Event {
	updatedEvents := make([]*sports.Event, len(events))
	for i, event := range events {
		// Create a new event with status updated
		updatedEvents[i] = &sports.Event{
			Id:                  event.Id,
			Name:                event.Name,
			AdvertisedStartTime: event.AdvertisedStartTime,
			Visible:             event.Visible,
			Venue:               event.Venue,
			SportType:           event.SportType,
			Competitors:         event.Competitors,
			Status:              CheckEventStatus(event),
		}
	}
	return updatedEvents
}
