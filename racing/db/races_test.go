package db

import (
	"database/sql"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"sort"
	"testing"
)

// TestRacesRepo_List tests the List method of the RacesRepo
func TestRacesRepo_List(t *testing.T) {
	//No need to mock the database as we are using an in-memory "mocked" (faker) database
	racingDB, _ := sql.Open("sqlite3", "./racing.db")
	racesRepo := NewRacesRepo(racingDB)
	racesRepo.Init()

	//Test the List method without a filter
	if allRaces, _ := racesRepo.List(nil, nil); len(allRaces) == 0 {
		t.Errorf("Expected to get a list of races")
	}

	//Test the List method with a meeting_id filter
	meetingIdFilter := &racing.ListRacesRequestFilter{
		MeetingIds: []int64{1, 5},
	}
	if races, _ := racesRepo.List(meetingIdFilter, nil); len(races) == 0 || !checkCondition(races, func(r *racing.Race) bool {
		return r.MeetingId == 1 || r.MeetingId == 5
	}) {
		t.Errorf("Expected to get a list of races with a meeting id of 1 or 5")
	}

	//Test the List method with a visibility filter
	visibility := true
	visibilityFilter := &racing.ListRacesRequestFilter{
		Visible: &visibility,
	}
	if races, _ := racesRepo.List(visibilityFilter, nil); len(races) == 0 || !checkCondition(races, func(r *racing.Race) bool {
		return r.Visible
	}) {
		t.Errorf("Expected to get a list of visible races")
	}

	//Test the List method with a visibility filter and a meeting_id filter
	filter := &racing.ListRacesRequestFilter{
		MeetingIds: []int64{1, 5},
		Visible:    &visibility,
	}
	if races, _ := racesRepo.List(filter, nil); len(races) == 0 || !checkCondition(races, func(r *racing.Race) bool {
		return r.MeetingId == 1 || r.MeetingId == 5 && r.Visible
	}) {
		t.Errorf("Expected to get a list of races with a meeting id of 1 or 5 and visible")
	}

	//Test the List method with a time based sort order (ascending)
	sortOrder := &racing.ListRacesRequestSortOrder{
		SortAttribute: "advertised_start_time",
		SortDirection: racing.SortDirection_ASC,
	}
	if races, _ := racesRepo.List(nil, sortOrder); len(races) == 0 || sort.SliceIsSorted(races, func(i, j int) bool {
		return races[i].AdvertisedStartTime.AsTime().After(races[j].AdvertisedStartTime.AsTime())
	}) {
		t.Errorf("Expected to get a list of races sorted by advertised start time in ascending order")
	}

	//Test the List method with a time based sort order (descending)
	sortOrder = &racing.ListRacesRequestSortOrder{
		SortAttribute: "advertised_start_time",
		SortDirection: racing.SortDirection_DESC,
	}
	if races, _ := racesRepo.List(nil, sortOrder); len(races) == 0 || sort.SliceIsSorted(races, func(i, j int) bool {
		return races[i].AdvertisedStartTime.AsTime().Before(races[j].AdvertisedStartTime.AsTime())
	}) {
		t.Errorf("Expected to get a list of races sorted by advertised start time in descending order")
	}

	//Test the List method with a boolean based sort order (ascending)
	sortOrder = &racing.ListRacesRequestSortOrder{
		SortAttribute: "visible",
		SortDirection: racing.SortDirection_ASC,
	}
	if races, _ := racesRepo.List(nil, sortOrder); len(races) == 0 || sort.SliceIsSorted(races, func(i, j int) bool {
		return races[i].Visible && !races[j].Visible
	}) {
		t.Errorf("Expected to get a list of of races sorted by visibility in ascending order")
	}

	//Test the List method with a boolean based sort order (descending)
	sortOrder = &racing.ListRacesRequestSortOrder{
		SortAttribute: "visible",
		SortDirection: racing.SortDirection_DESC,
	}
	if races, _ := racesRepo.List(nil, sortOrder); len(races) == 0 || sort.SliceIsSorted(races, func(i, j int) bool {
		return races[j].Visible && !races[i].Visible
	}) {
		t.Errorf("Expected to get a list of of races sorted by visibility in ascending order")
	}

	//Test the List method with an integer based sort order (ascending)
	sortOrder = &racing.ListRacesRequestSortOrder{
		SortAttribute: "meeting_id",
	}
	if races, _ := racesRepo.List(nil, sortOrder); len(races) == 0 || sort.SliceIsSorted(races, func(i, j int) bool {
		return races[i].MeetingId > races[j].MeetingId
	}) {
		t.Errorf("Expected to get a list of races sorted by meeting id in ascending order")
	}

	//Test the List method with an integer based sort order (descending)
	sortOrder = &racing.ListRacesRequestSortOrder{
		SortAttribute: "meeting_id",
		SortDirection: racing.SortDirection_DESC,
	}
	if races, _ := racesRepo.List(nil, sortOrder); len(races) == 0 || sort.SliceIsSorted(races, func(i, j int) bool {
		return races[i].MeetingId < races[j].MeetingId
	}) {
		t.Errorf("Expected to get a list of races sorted by meeting id in descending order")
	}
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
