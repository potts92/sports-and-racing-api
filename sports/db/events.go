package db

import (
	"database/sql"
	"github.com/potts92/sports-and-racing-api/sports/proto/sports"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type EventsRepo interface {
	// Init will initialise our events repository.
	Init() error

	// List will return a list of events.
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
	return nil, nil
}
