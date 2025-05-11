package db

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"git.neds.sh/matty/entain/sports/proto/sports"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates a new in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	return db
}

func TestEventsRepo_Init(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create repository
	repo := NewEventsRepo(db)

	// Test initialization
	err := repo.Init()
	require.NoError(t, err)

	// Verify table was created
	var tableName string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='events'").Scan(&tableName)
	require.NoError(t, err)
	assert.Equal(t, "events", tableName)

	// Verify seed data was inserted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
	require.NoError(t, err)
	assert.Greater(t, count, 0, "seed data should be inserted")

	// Test idempotence - calling Init() again should not error or duplicate data
	err = repo.Init()
	require.NoError(t, err)

	var newCount int
	err = db.QueryRow("SELECT COUNT(*) FROM events").Scan(&newCount)
	require.NoError(t, err)
	assert.Equal(t, count, newCount, "calling Init() twice should not insert duplicate data")
}

func TestEventsRepo_List(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create and initialize repository
	repo := NewEventsRepo(db)
	err := repo.Init()
	require.NoError(t, err)

	// First, get the actual counts from the database
	var visibleCount, nonVisibleCount int
	err = db.QueryRow("SELECT COUNT(*) FROM events WHERE visible = 1").Scan(&visibleCount)
	require.NoError(t, err)
	err = db.QueryRow("SELECT COUNT(*) FROM events WHERE visible = 0").Scan(&nonVisibleCount)
	require.NoError(t, err)

	totalCount := visibleCount + nonVisibleCount

	tests := []struct {
		desc      string
		filter    *sports.ListEventsRequestFilter
		wantCount int
		wantErr   error
	}{
		{
			desc: "when visible_only is true, returns only visible events",
			filter: &sports.ListEventsRequestFilter{
				VisibleOnly: true,
			},
			wantCount: visibleCount,
			wantErr:   nil,
		},
		{
			desc:      "when filter is nil, returns all events",
			filter:    nil,
			wantCount: totalCount,
			wantErr:   nil,
		},
		{
			desc: "when visible_only is false, returns all events",
			filter: &sports.ListEventsRequestFilter{
				VisibleOnly: false,
			},
			wantCount: totalCount,
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
			assert.Equal(t, tt.wantCount, len(got), "number of events should match")

			// Verify visibility filter
			if tt.filter != nil && tt.filter.VisibleOnly {
				for _, event := range got {
					assert.True(t, event.Visible, "event should be visible when filter.VisibleOnly is true")
				}
			}
		})
	}
}

func TestEventsRepo_ScanEvents(t *testing.T) {
	// Setup test database
	db := setupTestDB(t)
	defer db.Close()

	// Create and initialize repository
	repo := NewEventsRepo(db).(*eventsRepo) // Type assertion to access private methods
	err := repo.Init()
	require.NoError(t, err)

	// Test case for a single event
	t.Run("correctly scans a single event", func(t *testing.T) {
		// Insert a test event with known values
		now := time.Now()
		_, err := db.Exec(
			"INSERT INTO events (id, name, advertised_start, visible, venue, sport_type, competitors) VALUES (?, ?, ?, ?, ?, ?, ?)",
			999, // Use a high ID to avoid conflicts with seed data
			"Test Event",
			now,
			true,
			"Test Venue",
			"Test Sport",
			`["Competitor A", "Competitor B"]`,
		)
		require.NoError(t, err)

		// Query the event
		rows, err := db.Query("SELECT id, name, advertised_start, visible, venue, sport_type, competitors FROM events WHERE id = 999")
		require.NoError(t, err)
		defer rows.Close()

		// Use scanEvents to convert to protobuf objects
		events, err := repo.scanEvents(rows)
		require.NoError(t, err)
		require.Len(t, events, 1, "should return exactly one event")

		// Verify all fields
		event := events[0]
		assert.Equal(t, int64(999), event.Id)
		assert.Equal(t, "Test Event", event.Name)
		assert.True(t, event.Visible)
		assert.Equal(t, "Test Venue", event.Venue)
		assert.Equal(t, "Test Sport", event.SportType)

		// Check competitors
		require.Len(t, event.Competitors, 2)
		assert.Equal(t, "Competitor A", event.Competitors[0])
		assert.Equal(t, "Competitor B", event.Competitors[1])

		// Check timestamp (rough comparison)
		require.NotNil(t, event.AdvertisedStartTime)
		eventTime := event.AdvertisedStartTime.AsTime()
		assert.WithinDuration(t, now, eventTime, 2*time.Second)
	})

	// Test error handling
	t.Run("handles scan errors", func(t *testing.T) {
		// Prepare SQL that will cause a scan error (incompatible type)
		_, err := db.Exec("CREATE TABLE IF NOT EXISTS bad_events (id TEXT)")
		require.NoError(t, err)

		_, err = db.Exec("INSERT INTO bad_events (id) VALUES ('not-a-number')")
		require.NoError(t, err)

		// Try to scan with insufficient columns
		rows, err := db.Query("SELECT id FROM bad_events")
		require.NoError(t, err)
		defer rows.Close()

		_, err = repo.scanEvents(rows)
		assert.Error(t, err, "should return error when scanning incompatible data")
	})
}

func TestEventsRepo_CompetitorsHandling(t *testing.T) {
	// Test cases for different competitor JSON formats
	tests := []struct {
		name        string
		competitors string
		expected    []string
	}{
		{
			name:        "standard JSON array",
			competitors: `["Team A", "Team B"]`,
			expected:    []string{"Team A", "Team B"},
		},
		{
			name:        "empty JSON array",
			competitors: `[]`,
			expected:    []string{},
		},
		{
			name:        "single competitor",
			competitors: `["Solo Competitor"]`,
			expected:    []string{"Solo Competitor"},
		},
		{
			name:        "competitors with commas",
			competitors: `["Team, Inc.", "Other, Ltd."]`,
			expected:    []string{"Team, Inc.", "Other, Ltd."},
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh database for each test case
			db := setupTestDB(t)
			defer db.Close()

			// Create events table
			_, err := db.Exec(`
				CREATE TABLE IF NOT EXISTS events (
					id INTEGER PRIMARY KEY,
					name TEXT,
					advertised_start TIMESTAMP,
					visible BOOLEAN,
					venue TEXT,
					sport_type TEXT,
					competitors TEXT
				)
			`)
			require.NoError(t, err)

			// Create repository
			repo := NewEventsRepo(db).(*eventsRepo)

			// Directly test the scanEvents function instead
			// by preparing a row with the competitors JSON
			_, err = db.Exec(
				"INSERT INTO events (id, name, advertised_start, visible, venue, sport_type, competitors) VALUES (?, ?, ?, ?, ?, ?, ?)",
				i+1,
				"Test Event "+tt.name,
				time.Now(),
				true,
				"Test Venue",
				"Test Sport",
				tt.competitors,
			)
			require.NoError(t, err)

			// Get the row data from the DB
			rows, err := db.Query("SELECT id, name, advertised_start, visible, venue, sport_type, competitors FROM events WHERE id = ?", i+1)
			require.NoError(t, err)
			defer rows.Close()

			// Scan the row into an event
			events, err := repo.scanEvents(rows)
			require.NoError(t, err)
			require.Len(t, events, 1, "should return exactly one event")

			// For verification, let's also parse the JSON directly using Go's JSON parser
			// This will help us understand if the issue is in the database or in the scanner
			var expectedCompetitors []string
			err = json.Unmarshal([]byte(tt.competitors), &expectedCompetitors)
			require.NoError(t, err)

			// Now compare the parsed competitors with what our scanner produced
			assert.Equal(t, len(expectedCompetitors), len(events[0].Competitors),
				"number of competitors should match between direct JSON parsing and scanEvents")

			for j, expected := range expectedCompetitors {
				assert.Equal(t, expected, events[0].Competitors[j],
					"competitor at index %d should match between direct JSON parsing and scanEvents", j)
			}

			// Also verify against our original expected values
			assert.Equal(t, len(tt.expected), len(events[0].Competitors),
				"number of competitors should match the test's expected values")

			for j, expected := range tt.expected {
				assert.Equal(t, expected, events[0].Competitors[j],
					"competitor at index %d should match the test's expected value", j)
			}
		})
	}
}
