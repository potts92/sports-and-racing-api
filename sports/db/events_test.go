package db

import (
	"database/sql"
	"github.com/potts92/sports-and-racing-api/sports/proto/sports"
	"testing"
)

var eventsRepoInstance EventsRepo

// TestEventsRepo_List tests the List method of the EventsRepo with the score_finalised filter
func TestEventsRepo_List(t *testing.T) {
	eventsRepoInstance = initEventsRepo()

	scoreFinalised := true
	scoreNotFinalised := false

	testCases := []struct {
		name     string
		filter   *sports.ListEventsRequestFilter
		validate func([]*sports.Event) bool
	}{
		{
			name: "Score finalised",
			filter: &sports.ListEventsRequestFilter{
				ScoreFinalised: &scoreFinalised,
			},
			validate: func(events []*sports.Event) bool {
				return len(events) != 0 && checkCondition(events, func(event *sports.Event) bool {
					return event.ScoreFinalised
				})
			},
		},
		{
			name: "Score not finalised",
			filter: &sports.ListEventsRequestFilter{
				ScoreFinalised: &scoreNotFinalised,
			},
			validate: func(events []*sports.Event) bool {
				return len(events) != 0 && checkCondition(events, func(event *sports.Event) bool {
					return !event.ScoreFinalised
				})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			events, err := eventsRepoInstance.List(tc.filter)
			if err != nil {
				t.Errorf("Error listing events: %v", err)
			}

			if !tc.validate(events) {
				t.Errorf("Events did not pass validation")
			}
		})

	}
}

func initEventsRepo() EventsRepo {
	if eventsRepoInstance == nil {
		sportsDB, _ := sql.Open("sqlite", "./sports.db")
		eventsRepoInstance := NewEventsRepo(sportsDB)
		eventsRepoInstance.Init()
	}

	return eventsRepoInstance
}

func checkCondition(events []*sports.Event, condition func(*sports.Event) bool) bool {
	for _, event := range events {
		if !condition(event) {
			return false
		}
	}

	return true
}
