package db

import (
	"database/sql"
	"github.com/potts92/sports-and-racing-api/sports/proto/sports"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type EventsRepo interface {
	// Init will initialise our events repository.
	Init() error

	// List will return a list of events.
	List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error)

	// Get will return a single event.
	Get(id int64) (*sports.Event, error)

	// UpdateScore will update a event's score.
	UpdateScore(id int64, homeScore int32, awayScore int32, finalised bool) (*sports.Event, error)
}

type eventsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewEventsRepo creates a new events repository.
func NewEventsRepo(db *sql.DB) EventsRepo {
	return &eventsRepo{db: db}
}

// Init prepares the event repository dummy data.
func (e *eventsRepo) Init() error {
	var err error

	// Check if the database file already exists before seeding it.
	_, err = os.Stat("./db/sports.db")
	if os.IsNotExist(err) {
		e.init.Do(func() {
			// For test/example purposes, we seed the DB with some dummy events.
			err = e.seed()
		})
	}

	return err
}

func (e *eventsRepo) List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
		rows  *sql.Rows
	)

	query = getEventQueries()[eventsList]

	query = e.applyFilter(query, filter)

	rows, err = e.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return e.scanEvents(rows)
}

func (e *eventsRepo) Get(id int64) (*sports.Event, error) {
	var (
		err    error
		query  string
		rows   *sql.Rows
		events []*sports.Event
		args   []interface{}
	)

	query = getEventQueries()[eventGet]

	args = append(args, id)

	rows, err = e.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	events, err = e.scanEvents(rows)

	if len(events) > 0 {
		return events[0], nil
	}

	return nil, nil
}

func (e *eventsRepo) UpdateScore(id int64, homeScore int32, awayScore int32, finalised bool) (*sports.Event, error) {
	var (
		rows sql.Result
		err  error
		args []interface{}
	)

	query := getEventQueries()[eventUpdate]

	args = append(args, homeScore, awayScore, finalised, id)
	rows, err = e.db.Exec(query, args...)

	if err != nil {
		return nil, err
	}

	rowsAffected, _ := rows.RowsAffected()

	if rowsAffected > 0 {
		return e.Get(id)
	}

	return nil, nil
}

func (e *eventsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) string {
	var (
		clauses []string
	)

	if filter == nil {
		return query
	}

	if filter.ScoreFinalised != nil {
		var visibility string
		if *filter.ScoreFinalised {
			visibility = "1"
		} else {
			visibility = "0"
		}
		clauses = append(clauses, "score_finalised = "+visibility)
	}

	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query
}

func (e *eventsRepo) scanEvents(rows *sql.Rows) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStartTime time.Time

		if err := rows.Scan(
			&event.Id,
			&event.Name,
			&event.Competition,
			&event.HomeTeam,
			&event.AwayTeam,
			&event.HomeScore,
			&event.AwayScore,
			&advertisedStartTime,
			&event.ScoreFinalised,
		); err != nil {
			return nil, err
		}

		ts := timestamppb.New(advertisedStartTime)
		event.AdvertisedStartTime = ts

		//faker generates dates in UTC for the database, so we need to compare to UTC time
		open := ts.AsTime().After(time.Now().UTC())
		if open {
			event.Status = sports.Status_OPEN
		} else {
			event.Status = sports.Status_CLOSED
		}

		events = append(events, &event)
	}

	return events, nil
}
