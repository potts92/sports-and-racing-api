package db

import (
	"database/sql"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"github.com/DATA-DOG/go-sqlmock"
	"slices"
	"testing"
	"time"
)

// TestRacesRepo_List tests the List method of the RacesRepo
func TestRacesRepo_List(t *testing.T) {
	//No need to mock the database as we are using an in-memory "mocked" (faker) database
	racingDB, _ := sql.Open("sqlite3", "./racing.db")
	racesRepo := NewRacesRepo(racingDB)
	racesRepo.Init()

	visible := true
	notVisible := false

	testCases := []struct {
		name      string
		filter    *racing.ListRacesRequestFilter
		sortOrder *racing.ListRacesRequestSortOrder
		validate  func([]*racing.Race) bool
		error     bool
	}{
		{
			name:      "Unfiltered, unsorted",
			filter:    nil,
			sortOrder: nil,
			validate: func(races []*racing.Race) bool {
				return len(races) != 0
			},
		},
		{
			name: "Filtered by meeting_id 1 and 5",
			filter: &racing.ListRacesRequestFilter{
				MeetingIds: []int64{1, 5},
			},
			sortOrder: nil,
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 &&
					checkCondition(races, func(r *racing.Race) bool {
						return r.MeetingId == 1 || r.MeetingId == 5
					})
			},
		},
		{
			name: "Filtered by visibility",
			filter: &racing.ListRacesRequestFilter{
				Visible: &visible,
			},
			sortOrder: nil,
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && checkCondition(races, func(r *racing.Race) bool {
					return r.Visible
				})
			},
		},
		{
			name: "Filtered by not visible",
			filter: &racing.ListRacesRequestFilter{
				Visible: &notVisible,
			},
			sortOrder: nil,
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && checkCondition(races, func(r *racing.Race) bool {
					return !r.Visible
				})
			},
		},
		{
			name: "Filtered by meeting_id 1 and 5 and visibility",
			filter: &racing.ListRacesRequestFilter{
				MeetingIds: []int64{1, 5},
				Visible:    &visible,
			},
			sortOrder: nil,
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && checkCondition(races, func(r *racing.Race) bool {
					return (r.MeetingId == 1 || r.MeetingId == 5) && r.Visible
				})
			},
		},
		{
			name: "Filtered by meeting_id 1 and 5 and not visible",
			filter: &racing.ListRacesRequestFilter{
				MeetingIds: []int64{1, 5},
				Visible:    &notVisible,
			},
			sortOrder: nil,
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && checkCondition(races, func(r *racing.Race) bool {
					return (r.MeetingId == 1 || r.MeetingId == 5) && !r.Visible
				})
			},
		},
		{
			name:   "Sorted by advertised start time (ascending)",
			filter: nil,
			sortOrder: &racing.ListRacesRequestSortOrder{
				SortAttribute: "advertised_start_time",
				SortDirection: racing.SortDirection_ASC,
			},
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && isSortedByAdvertisedStartTime(races, true)
			},
		},
		{
			name:   "Sorted by advertised start time (descending)",
			filter: nil,
			sortOrder: &racing.ListRacesRequestSortOrder{
				SortAttribute: "advertised_start_time",
				SortDirection: racing.SortDirection_DESC,
			},
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && isSortedByAdvertisedStartTime(races, false)
			},
		},
		{
			name:   "Sorted by visibility (ascending not explicitly defined)",
			filter: nil,
			sortOrder: &racing.ListRacesRequestSortOrder{
				SortAttribute: "visible",
			},
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && isSortedByVisibility(races, true)
			},
		},
		{
			name:   "Sorted by visibility (descending)",
			filter: nil,
			sortOrder: &racing.ListRacesRequestSortOrder{
				SortAttribute: "visible",
				SortDirection: racing.SortDirection_DESC,
			},
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && isSortedByVisibility(races, false)
			},
		},
		{
			name:   "Sorted by meeting_id (ascending)",
			filter: nil,
			sortOrder: &racing.ListRacesRequestSortOrder{
				SortAttribute: "meeting_id",
			},
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && isSortedBbyMeetingId(races, true)
			},
		},
		{
			name:   "Sorted by meeting_id (descending)",
			filter: nil,
			sortOrder: &racing.ListRacesRequestSortOrder{
				SortAttribute: "meeting_id",
				SortDirection: racing.SortDirection_DESC,
			},
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && isSortedBbyMeetingId(races, false)
			},
		},
		{
			name: "Filtered by meeting_id 1 and 5 and sorted by meeting_id (descending)",
			filter: &racing.ListRacesRequestFilter{
				MeetingIds: []int64{1, 5},
			},
			sortOrder: &racing.ListRacesRequestSortOrder{
				SortAttribute: "meeting_id",
				SortDirection: racing.SortDirection_DESC,
			},
			validate: func(races []*racing.Race) bool {
				return len(races) != 0 && isSortedBbyMeetingId(races, false) && checkCondition(races, func(r *racing.Race) bool {
					return r.MeetingId == 1 || r.MeetingId == 5
				})
			},
		},
		{
			name:   "Expect error with invalid sort attribute",
			filter: nil,
			sortOrder: &racing.ListRacesRequestSortOrder{
				SortAttribute: "fake_column",
				SortDirection: racing.SortDirection_DESC,
			},
			validate: func(races []*racing.Race) bool {
				return len(races) == 0
			},
			error: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			races, err := racesRepo.List(tc.filter, tc.sortOrder)
			if !tc.validate(races) || (err != nil) != tc.error {
				t.Errorf("Test case failed, expected list of races: %s", tc.name)
			}
		})
	}

	t.Run("Derived status", func(t *testing.T) {
		// Create a mock database connection
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		// Create an instance of racesRepo with the mock database connection
		repo := NewRacesRepo(db)

		// Define the expected columns and rows with the three variations of advertised_start_time we need to test
		now := time.Now()
		columns := []string{"id", "meeting_id", "name", "number", "visible", "advertised_start"}
		rows := sqlmock.NewRows(columns).
			AddRow(1, 2, "Race 1", 1, true, now.Add(time.Hour)).
			AddRow(2, 2, "Race 2", 2, false, now.Add(-time.Hour)).
			AddRow(3, 2, "Race 2", 2, false, now)

		// Set the expectation for the mock database connection
		mock.ExpectQuery("SELECT (.+) FROM races").WillReturnRows(rows)

		// Call the scanRaces function
		races, err := repo.List(nil, nil)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// If advertised_start_time is after the current time, status should be OPEN, if not, it should be CLOSED
		//(includes where advertised_start_time is equal to the current time)
		for _, race := range races {
			if race.AdvertisedStartTime.AsTime().After(now) {
				if race.Status != racing.Status_OPEN {
					t.Errorf("Race status was expected to be open")
				}
			} else if race.Status != racing.Status_CLOSED {
				t.Errorf("Race status was expected to be closed")
			}
		}
	})
}

// Checks if a slice of races is sorted by the meeting_id attribute
func isSortedBbyMeetingId(races []*racing.Race, asc bool) bool {
	dirMultiplier := getDirMultiplier(asc)

	return slices.IsSortedFunc(races, func(i, j *racing.Race) int {
		if i.MeetingId < j.MeetingId {
			return -1 * dirMultiplier
		} else if i.MeetingId == j.MeetingId {
			return 0
		} else {
			return 1 * dirMultiplier
		}
	})
}

// Checks if a slice of races is sorted by the visible attribute
func isSortedByVisibility(races []*racing.Race, asc bool) bool {
	dirMultiplier := getDirMultiplier(asc)

	return slices.IsSortedFunc(races, func(i, j *racing.Race) int {
		if i.Visible && !j.Visible {
			return 1 * dirMultiplier
		} else if i.Visible == j.Visible {
			return 0
		} else {
			return -1 * dirMultiplier
		}
	})
}

// Checks if a slice of races is sorted by the advertised start time attribute
func isSortedByAdvertisedStartTime(races []*racing.Race, asc bool) bool {
	return slices.IsSortedFunc(races, func(i, j *racing.Race) int {
		if asc && i.AdvertisedStartTime.AsTime().Before(j.AdvertisedStartTime.AsTime()) || !asc && i.AdvertisedStartTime.AsTime().After(j.AdvertisedStartTime.AsTime()) {
			return -1
		} else if i.AdvertisedStartTime.AsTime().Equal(j.AdvertisedStartTime.AsTime()) {
			return 0
		} else {
			return 1
		}
	})
}

// Returns the direction multiplier based on the sort direction (used to reverse sort check functionality)
func getDirMultiplier(asc bool) int {
	if asc {
		return 1
	}
	return -1
}

// Used to check if all races pass the expected condition
func checkCondition(arr []*racing.Race, condition func(*racing.Race) bool) bool {
	for _, r := range arr {
		if !condition(r) {
			return false
		}
	}
	return true
}
