package service

import (
	"time"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// CheckRaceStatus determines the status of a race based on its advertised start time.
// Returns CLOSED if the advertised start time is in the past,
// OPEN otherwise.
func CheckRaceStatus(race *racing.Race) racing.RaceStatus {
	if race.AdvertisedStartTime == nil {
		return racing.RaceStatus_UNSPECIFIED
	}

	advertisedTime := race.AdvertisedStartTime.AsTime()
	if advertisedTime.Before(time.Now()) {
		return racing.RaceStatus_CLOSED
	}
	return racing.RaceStatus_OPEN
}

// UpdateRacesStatus updates the status field for all races in the provided slice.
// This is a pure function that returns a new slice with updated statuses.
func UpdateRacesStatus(races []*racing.Race) []*racing.Race {
	updatedRaces := make([]*racing.Race, len(races))
	for i, race := range races {
		// Create a new Race instance
		updatedRace := &racing.Race{
			Id:                  race.Id,
			MeetingId:           race.MeetingId,
			Name:                race.Name,
			Number:              race.Number,
			Visible:             race.Visible,
			AdvertisedStartTime: race.AdvertisedStartTime,
			Status:              CheckRaceStatus(race),
		}
		updatedRaces[i] = updatedRace
	}
	return updatedRaces
}
