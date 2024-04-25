package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

// RacesRepo provides repository access to races.
type RacesRepo interface {
	// Init will initialise our races repository.
	Init() error

	// List will return a list of races.
	List(filter *racing.ListRacesRequestFilter, sort *racing.ListRacesRequestSortOrder) ([]*racing.Race, error)

	// Get will return a single race.
	Get(id int64) (*racing.Race, error)
}

type racesRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewRacesRepo creates a new races repository.
func NewRacesRepo(db *sql.DB) RacesRepo {
	return &racesRepo{db: db}
}

// Init prepares the race repository dummy data.
func (r *racesRepo) Init() error {
	var err error

	r.init.Do(func() {
		// For test/example purposes, we seed the DB with some dummy races.
		err = r.seed()
	})

	return err
}

func (r *racesRepo) List(filter *racing.ListRacesRequestFilter, sort *racing.ListRacesRequestSortOrder) ([]*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[racesList]

	query, args = r.applyFilter(query, filter)

	query = r.setSortOrder(query, sort)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanRaces(rows)
}

func (r *racesRepo) Get(id int64) (*racing.Race, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getRaceQueries()[raceGet]

	args = append(args, id)

	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	race, err := r.scanRaces(rows)

	if err != nil {
		return nil, err
	}

	if len(race) == 0 {
		return nil, nil
	}

	return race[0], nil
}

func (r *racesRepo) applyFilter(query string, filter *racing.ListRacesRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.MeetingIds) > 0 {
		clauses = append(clauses, "meeting_id IN ("+strings.Repeat("?,", len(filter.MeetingIds)-1)+"?)")

		for _, meetingID := range filter.MeetingIds {
			args = append(args, meetingID)
		}
	}

	// Need to check if visible is in filter so that it won't default to false
	// If not set, visibility won't be set and SQL won't be affected
	// Go struct field is a pointer to a boolean (set as optional in protobuf) to differentiate between unset and false
	if filter.Visible != nil {
		var visibility string
		if *filter.Visible == true {
			visibility = "true"
		} else {
			visibility = "false"
		}
		clauses = append(clauses, "visible = "+visibility)
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *racesRepo) scanRaces(
	rows *sql.Rows,
) ([]*racing.Race, error) {
	var races []*racing.Race

	for rows.Next() {
		var race racing.Race
		var advertisedStart time.Time

		if err := rows.Scan(&race.Id, &race.MeetingId, &race.Name, &race.Number, &race.Visible, &advertisedStart); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}

			return nil, err
		}

		ts := timestamppb.New(advertisedStart)

		race.AdvertisedStartTime = ts

		//faker generates dates in UTC for the database, so we need to compare to UTC time
		open := ts.AsTime().After(time.Now().UTC())
		if open {
			race.Status = racing.Status_OPEN
		} else {
			race.Status = racing.Status_CLOSED
		}

		races = append(races, &race)
	}

	return races, nil
}

// Sorts races by the given sort attribute (and sort direction if provided - otherwise defaults to ASC)
func (r *racesRepo) setSortOrder(query string, sort *racing.ListRacesRequestSortOrder) string {
	if sort == nil {
		return query
	}

	sortAttribute := sort.SortAttribute
	sortDirection := sort.SortDirection

	query += " ORDER BY " + sortAttribute + " " + sortDirection.String()

	return query
}
