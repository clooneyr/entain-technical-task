package db

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// setupTestDB creates a new in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	return db
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
		// Timestamps are ignored as they are generated at runtime
	}
}

func TestRacesRepo_List_VisibilityFilter(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create and initialize repository
	repo := NewRacesRepo(db)
	err := repo.Init()
	require.NoError(t, err)

	// First, get the actual counts from the database
	var visibleCount, nonVisibleCount int
	err = db.QueryRow("SELECT COUNT(*) FROM races WHERE visible = 1").Scan(&visibleCount)
	require.NoError(t, err)
	err = db.QueryRow("SELECT COUNT(*) FROM races WHERE visible = 0").Scan(&nonVisibleCount)
	require.NoError(t, err)

	tests := []struct {
		desc      string
		filter    *racing.ListRacesRequestFilter
		wantCount int
		wantErr   error
	}{
		{
			desc: "when visible_only is true, returns only visible races",
			filter: &racing.ListRacesRequestFilter{
				VisibleOnly: boolPtr(t, true),
			},
			wantCount: visibleCount,
			wantErr:   nil,
		},
		{
			desc: "when visible_only is false, returns only non-visible races",
			filter: &racing.ListRacesRequestFilter{
				VisibleOnly: boolPtr(t, false),
			},
			wantCount: nonVisibleCount,
			wantErr:   nil,
		},
		{
			desc:      "when visible_only is not provided, returns all races",
			filter:    &racing.ListRacesRequestFilter{},
			wantCount: visibleCount + nonVisibleCount,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// Execute
			got, err := repo.List(tt.filter)

			// Verify
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantCount, len(got), "number of races should match")

			// Verify visibility filter
			if tt.filter.VisibleOnly != nil {
				for _, race := range got {
					assert.Equal(t, *tt.filter.VisibleOnly, race.Visible, "race visibility should match filter")
				}
			}
		})
	}
}

func TestRacesRepo_List_MeetingIDFilter(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create and initialize repository
	repo := NewRacesRepo(db)
	err := repo.Init()
	require.NoError(t, err)

	// Get a meeting ID that exists in the database
	var meetingID int64
	err = db.QueryRow("SELECT meeting_id FROM races LIMIT 1").Scan(&meetingID)
	require.NoError(t, err)

	// Get the count of races for this meeting
	var meetingCount int
	err = db.QueryRow("SELECT COUNT(*) FROM races WHERE meeting_id = ?", meetingID).Scan(&meetingCount)
	require.NoError(t, err)

	tests := []struct {
		desc      string
		filter    *racing.ListRacesRequestFilter
		wantCount int
		wantErr   error
	}{
		{
			desc: "when meeting_ids is provided, returns only races for those meetings",
			filter: &racing.ListRacesRequestFilter{
				MeetingIds: []int64{meetingID},
			},
			wantCount: meetingCount,
			wantErr:   nil,
		},
		{
			desc: "when meeting_ids is empty, returns all races",
			filter: &racing.ListRacesRequestFilter{
				MeetingIds: []int64{},
			},
			wantCount: 100, // Total number of races in seed data
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// Execute
			got, err := repo.List(tt.filter)

			// Verify
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantCount, len(got), "number of races should match")

			// Verify meeting ID filter
			if len(tt.filter.MeetingIds) > 0 {
				for _, race := range got {
					assert.Contains(t, tt.filter.MeetingIds, race.MeetingId, "race meeting ID should be in filter")
				}
			}
		})
	}
}

func TestRacesRepo_List_CombinedFilters(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create and initialize repository
	repo := NewRacesRepo(db)
	err := repo.Init()
	require.NoError(t, err)

	// Get a meeting ID that exists in the database
	var meetingID int64
	err = db.QueryRow("SELECT meeting_id FROM races LIMIT 1").Scan(&meetingID)
	require.NoError(t, err)

	// Get the count of visible races for this meeting
	var visibleMeetingCount int
	err = db.QueryRow("SELECT COUNT(*) FROM races WHERE meeting_id = ? AND visible = 1", meetingID).Scan(&visibleMeetingCount)
	require.NoError(t, err)

	tests := []struct {
		desc      string
		filter    *racing.ListRacesRequestFilter
		wantCount int
		wantErr   error
	}{
		{
			desc: "when both meeting_ids and visible_only are provided",
			filter: &racing.ListRacesRequestFilter{
				MeetingIds:  []int64{meetingID},
				VisibleOnly: boolPtr(t, true),
			},
			wantCount: visibleMeetingCount,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// Execute
			got, err := repo.List(tt.filter)

			// Verify
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantCount, len(got), "number of races should match")

			// Verify combined filters
			for _, race := range got {
				assert.Contains(t, tt.filter.MeetingIds, race.MeetingId, "race meeting ID should be in filter")
				assert.Equal(t, *tt.filter.VisibleOnly, race.Visible, "race visibility should match filter")
			}
		})
	}
}

func TestRacesRepo_GetRace(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create and initialize repository
	repo := NewRacesRepo(db)
	err := repo.Init()
	require.NoError(t, err)

	tests := []struct {
		desc    string
		id      int64
		wantErr error
	}{
		{
			desc:    "when race exists, returns race",
			id:      1,
			wantErr: nil,
		},
		{
			desc:    "when race doesn't exist, returns error",
			id:      999,
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// Execute
			got, err := repo.GetRace(tt.id)

			// Verify
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.id, got.Id, "race ID should match")
		})
	}
}

// boolPtr returns a pointer to the given bool.
func boolPtr(t *testing.T, b bool) *bool {
	t.Helper()
	return &b
}
