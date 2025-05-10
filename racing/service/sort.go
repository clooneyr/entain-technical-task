package service

import (
	"sort"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// SortRaces sorts the provided races based on the filter's sort criteria.
// If no sort criteria is specified, races are sorted by advertised start time in ascending order.
// Returns a new slice with the races sorted according to the specified criteria.
func SortRaces(races []*racing.Race, filter *racing.ListRacesRequestFilter) []*racing.Race {
	if filter == nil {
		return sortByAdvertisedStartTime(races, racing.SortOrder_SORT_ORDER_ASC)
	}

	// Default to advertised start time if no sort field specified
	sortBy := racing.SortBy_SORT_BY_ADVERTISED_START_TIME
	if filter.SortBy != nil {
		sortBy = *filter.SortBy
	}

	// Default to ascending if no sort order specified
	sortOrder := racing.SortOrder_SORT_ORDER_ASC
	if filter.SortOrder != nil {
		sortOrder = *filter.SortOrder
	}

	switch sortBy {
	case racing.SortBy_SORT_BY_ADVERTISED_START_TIME:
		return sortByAdvertisedStartTime(races, sortOrder)
	case racing.SortBy_SORT_BY_NAME:
		return sortByName(races, sortOrder)
	case racing.SortBy_SORT_BY_NUMBER:
		return sortByNumber(races, sortOrder)
	default:
		return sortByAdvertisedStartTime(races, sortOrder)
	}
}

// sortByAdvertisedStartTime sorts the provided races by their advertised start time.
// The order parameter determines whether races are sorted in ascending (SORT_ORDER_ASC) or
// descending (SORT_ORDER_DESC) order.
// Returns a new slice with the races sorted by advertised start time.
func sortByAdvertisedStartTime(races []*racing.Race, order racing.SortOrder) []*racing.Race {
	sortedRaces := make([]*racing.Race, len(races))
	copy(sortedRaces, races)

	sort.Slice(sortedRaces, func(i, j int) bool {
		// Handle nil cases
		if sortedRaces[i].AdvertisedStartTime == nil {
			return false
		}
		if sortedRaces[j].AdvertisedStartTime == nil {
			return true
		}

		iTime := sortedRaces[i].AdvertisedStartTime.AsTime()
		jTime := sortedRaces[j].AdvertisedStartTime.AsTime()

		if order == racing.SortOrder_SORT_ORDER_DESC {
			return iTime.After(jTime)
		}
		return iTime.Before(jTime)
	})
	return sortedRaces
}

// sortByName sorts the provided races by their name.
// The order parameter determines whether races are sorted in ascending (SORT_ORDER_ASC) or
// descending (SORT_ORDER_DESC) order.
// Returns a new slice with the races sorted by name.
func sortByName(races []*racing.Race, order racing.SortOrder) []*racing.Race {
	sortedRaces := make([]*racing.Race, len(races))
	copy(sortedRaces, races)

	sort.Slice(sortedRaces, func(i, j int) bool {
		if order == racing.SortOrder_SORT_ORDER_DESC {
			return sortedRaces[i].Name > sortedRaces[j].Name
		}
		return sortedRaces[i].Name < sortedRaces[j].Name
	})
	return sortedRaces
}

// sortByNumber sorts the provided races by their number.
// The order parameter determines whether races are sorted in ascending (SORT_ORDER_ASC) or
// descending (SORT_ORDER_DESC) order.
// Returns a new slice with the races sorted by number.
func sortByNumber(races []*racing.Race, order racing.SortOrder) []*racing.Race {
	sortedRaces := make([]*racing.Race, len(races))
	copy(sortedRaces, races)

	sort.Slice(sortedRaces, func(i, j int) bool {
		if order == racing.SortOrder_SORT_ORDER_DESC {
			return sortedRaces[i].Number > sortedRaces[j].Number
		}
		return sortedRaces[i].Number < sortedRaces[j].Number
	})
	return sortedRaces
}
