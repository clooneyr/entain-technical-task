package service

import (
	"sort"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

// SortEvents sorts the provided events based on the filter's sort criteria.
// If no sort criteria is specified, events are sorted by advertised start time in ascending order.
// Returns a new slice with the events sorted according to the specified criteria.
func SortEvents(events []*sports.Event, filter *sports.ListEventsRequestFilter) []*sports.Event {
	if filter == nil {
		return sortByAdvertisedStartTime(events, sports.SortOrder_SORT_ORDER_ASC)
	}

	// Default to advertised start time if sort field is unspecified
	sortBy := sports.SortBy_SORT_BY_ADVERTISED_START_TIME
	if filter.SortBy != sports.SortBy_SORT_BY_UNSPECIFIED {
		sortBy = filter.SortBy
	}

	// Default to ascending if sort order is unspecified
	sortOrder := sports.SortOrder_SORT_ORDER_ASC
	if filter.SortOrder != sports.SortOrder_SORT_ORDER_UNSPECIFIED {
		sortOrder = filter.SortOrder
	}

	switch sortBy {
	case sports.SortBy_SORT_BY_ADVERTISED_START_TIME:
		return sortByAdvertisedStartTime(events, sortOrder)
	case sports.SortBy_SORT_BY_NAME:
		return sortByName(events, sortOrder)
	case sports.SortBy_SORT_BY_VENUE:
		return sortByVenue(events, sortOrder)
	default:
		return sortByAdvertisedStartTime(events, sortOrder)
	}
}

// sortByAdvertisedStartTime sorts the provided events by their advertised start time.
// The order parameter determines whether events are sorted in ascending (SORT_ORDER_ASC) or
// descending (SORT_ORDER_DESC) order.
// Returns a new slice with the events sorted by advertised start time.
func sortByAdvertisedStartTime(events []*sports.Event, order sports.SortOrder) []*sports.Event {
	sortedEvents := make([]*sports.Event, len(events))
	copy(sortedEvents, events)

	sort.Slice(sortedEvents, func(i, j int) bool {
		// Handle nil cases
		if sortedEvents[i].AdvertisedStartTime == nil {
			return false
		}
		if sortedEvents[j].AdvertisedStartTime == nil {
			return true
		}

		iTime := sortedEvents[i].AdvertisedStartTime.AsTime()
		jTime := sortedEvents[j].AdvertisedStartTime.AsTime()

		if order == sports.SortOrder_SORT_ORDER_DESC {
			return iTime.After(jTime)
		}
		return iTime.Before(jTime)
	})
	return sortedEvents
}

// sortByName sorts the provided events by their name.
// The order parameter determines whether events are sorted in ascending (SORT_ORDER_ASC) or
// descending (SORT_ORDER_DESC) order.
// Returns a new slice with the events sorted by name.
func sortByName(events []*sports.Event, order sports.SortOrder) []*sports.Event {
	sortedEvents := make([]*sports.Event, len(events))
	copy(sortedEvents, events)

	sort.Slice(sortedEvents, func(i, j int) bool {
		if order == sports.SortOrder_SORT_ORDER_DESC {
			return sortedEvents[i].Name > sortedEvents[j].Name
		}
		return sortedEvents[i].Name < sortedEvents[j].Name
	})
	return sortedEvents
}

// sortByVenue sorts the provided events by their venue.
// The order parameter determines whether events are sorted in ascending (SORT_ORDER_ASC) or
// descending (SORT_ORDER_DESC) order.
// Returns a new slice with the events sorted by venue.
func sortByVenue(events []*sports.Event, order sports.SortOrder) []*sports.Event {
	sortedEvents := make([]*sports.Event, len(events))
	copy(sortedEvents, events)

	sort.Slice(sortedEvents, func(i, j int) bool {
		if order == sports.SortOrder_SORT_ORDER_DESC {
			return sortedEvents[i].Venue > sortedEvents[j].Venue
		}
		return sortedEvents[i].Venue < sortedEvents[j].Venue
	})
	return sortedEvents
}
