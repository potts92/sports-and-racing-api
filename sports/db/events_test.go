package db

import (
	"database/sql"
	"github.com/potts92/sports-and-racing-api/sports/proto/sports"
	"testing"
)

var eventsRepoInstance EventsRepo

var finalisedEvent *sports.Event
var notFinalisedEvent *sports.Event

// TestEventsRepo_List tests the List method of the EventsRepo and saves the finalised and not finalised events for use in the UpdateScore test
// If this test fails (and these events are unable to be found and stored), the UpdateScore test will not run anyway
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
				if len(events) != 0 {
					// Saved for use by the UpdateScore test
					finalisedEvent = events[0]
				}
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
				if len(events) != 0 {
					// Saved for use by the UpdateScore test
					notFinalisedEvent = events[0]
				}
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

// todo: test update score - valid and invalid updates
func TestEventsRepo_UpdateScore(t *testing.T) {
	// Use real database connection to test that scores are actually updated
	eventsRepoInstance = initEventsRepo()

	testCases := []struct {
		name      string
		id        int64
		homeScore int32
		awayScore int32
		finalised bool
		validate  func(*sports.Event) bool
	}{
		{
			name:      "Update finalised event",
			id:        finalisedEvent.Id,
			homeScore: finalisedEvent.HomeScore + 1,
			awayScore: finalisedEvent.AwayScore + 1,
			finalised: true,
			validate: func(event *sports.Event) bool {
				return event == nil
			},
		},
		{
			name:      "Update not finalised event",
			id:        notFinalisedEvent.Id,
			homeScore: notFinalisedEvent.HomeScore + 1,
			awayScore: notFinalisedEvent.AwayScore + 1,
			finalised: true,
			validate: func(event *sports.Event) bool {
				return event.ScoreFinalised && event.HomeScore == notFinalisedEvent.HomeScore+1 && event.AwayScore == notFinalisedEvent.AwayScore+1
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event, _ := eventsRepoInstance.UpdateScore(tc.id, tc.homeScore, tc.awayScore, tc.finalised)
			if !tc.validate(event) {
				t.Errorf("Test case failed, expected event to be updated: %s", tc.name)
			}
		})

	}
}

func initEventsRepo() EventsRepo {
	if eventsRepoInstance == nil {
		sportsDB, _ := sql.Open("sqlite3", "./sports.db")
		eventsRepoInstance = NewEventsRepo(sportsDB)
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
