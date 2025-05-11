package db

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

// EventsRepo provides repository access to sports events.
type EventsRepo interface {
	// Init will initialize our events repository.
	Init() error

	// List will return a list of sporting events.
	List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error)
}

type eventsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewEventsRepo creates a new events repository.
func NewEventsRepo(db *sql.DB) EventsRepo {
	return &eventsRepo{db: db}
}

// Init prepares the events repository with initial data.
func (r *eventsRepo) Init() error {
	var err error

	r.init.Do(func() {
		// Create events table if it doesn't exist
		_, err = r.db.Exec(`
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
		if err != nil {
			return
		}

		// For test/example purposes it seeds the DB with some dummy events.
		err = r.seed()
	})

	return err
}

// seed adds initial test data to the database.
func (r *eventsRepo) seed() error {
	// First check if we already have data
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM events").Scan(&count)
	if err != nil {
		return err
	}

	// If we already have data, don't seed
	if count > 0 {
		return nil
	}

	// Time references for event start times
	now := time.Now()

	// Sample events data
	events := []struct {
		name            string
		advertisedStart time.Time
		visible         bool
		venue           string
		sportType       string
		competitors     string // JSON array of competitors
	}{
		{
			name:            "Premier League: Liverpool vs Manchester United",
			advertisedStart: now.Add(24 * time.Hour),
			visible:         true,
			venue:           "Anfield",
			sportType:       "Soccer",
			competitors:     `["Liverpool FC", "Manchester United"]`,
		},
		{
			name:            "NBA: Lakers vs Celtics",
			advertisedStart: now.Add(48 * time.Hour),
			visible:         true,
			venue:           "Staples Center",
			sportType:       "Basketball",
			competitors:     `["Los Angeles Lakers", "Boston Celtics"]`,
		},
		{
			name:            "Wimbledon: Final",
			advertisedStart: now.Add(-2 * time.Hour),
			visible:         true,
			venue:           "All England Club",
			sportType:       "Tennis",
			competitors:     `["Novak Djokovic", "Rafael Nadal"]`,
		},
		{
			name:            "F1: Monaco Grand Prix",
			advertisedStart: now.Add(72 * time.Hour),
			visible:         false, // Not visible
			venue:           "Circuit de Monaco",
			sportType:       "Formula 1",
			competitors:     `["Red Bull", "Ferrari", "Mercedes"]`,
		},
		{
			name:            "UFC 300",
			advertisedStart: now.Add(12 * time.Hour),
			visible:         true,
			venue:           "T-Mobile Arena",
			sportType:       "MMA",
			competitors:     `["Jon Jones", "Ciryl Gane"]`,
		},
	}

	// Insert events into database
	for i, event := range events {
		_, err = r.db.Exec(
			"INSERT INTO events (id, name, advertised_start, visible, venue, sport_type, competitors) VALUES (?, ?, ?, ?, ?, ?, ?)",
			i+1,
			event.name,
			event.advertisedStart,
			event.visible,
			event.venue,
			event.sportType,
			event.competitors,
		)
		if err != nil {
			return fmt.Errorf("failed to seed events: %w", err)
		}
	}

	return nil
}

// List returns a list of events with optional filtering.
func (r *eventsRepo) List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error) {
	query := "SELECT id, name, advertised_start, visible, venue, sport_type, competitors FROM events"

	var (
		clauses []string
		args    []interface{}
	)

	if filter != nil && filter.VisibleOnly {
		clauses = append(clauses, "visible = ?")
		args = append(args, true)
	}

	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

// scanEvents converts SQL rows to event objects.
func (r *eventsRepo) scanEvents(rows *sql.Rows) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var (
			event           sports.Event
			advertisedStart time.Time
			competitorsJSON string
			competitors     []string
		)

		if err := rows.Scan(
			&event.Id,
			&event.Name,
			&advertisedStart,
			&event.Visible,
			&event.Venue,
			&event.SportType,
			&competitorsJSON,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}

		// Parse JSON competitors list
		competitorsJSON = strings.Trim(competitorsJSON, "[]\"")
		if competitorsJSON != "" {
			competitors = strings.Split(strings.ReplaceAll(competitorsJSON, "\"", ""), ",")
			for i, c := range competitors {
				competitors[i] = strings.TrimSpace(c)
			}
			event.Competitors = competitors
		}

		// Convert time to protobuf timestamp
		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, fmt.Errorf("failed to convert time: %w", err)
		}
		event.AdvertisedStartTime = ts

		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return events, nil
}
