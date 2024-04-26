package db

import (
	"database/sql"
	"log"
	"strconv"
	"syreclabs.com/go/faker"
	"time"
)

type Team struct {
	Name  string
	Sport int
}

type Teams map[int]Team

type Competition struct {
	Name  string
	Sport int
}

type Competitions map[int]Competition

type Sport struct {
	Name string
}

type Sports map[int]Sport

// MockData Set up mock data that maps teams and competitions to their sport
type MockData map[int]struct {
	Teams        []int
	Competitions []int
}

func (e *eventsRepo) seed() error {
	var (
		err          error
		sports       Sports
		teams        Teams
		competitions Competitions
		mockData     MockData
	)

	sports, teams, competitions, mockData, err = e.buildMockData()

	if err != nil {
		return err
	}

	err = e.seedSports(sports)

	if err != nil {
		return err
	}

	err = e.seedTeams(teams)

	if err != nil {
		return err
	}

	err = e.seedCompetitions(competitions)

	if err != nil {
		return err
	}

	err = e.seedEvents(mockData)

	if err != nil {
		return err
	}

	return err
}

// Build mock data for sports, teams, competitions and events
func (e *eventsRepo) buildMockData() (Sports, Teams, Competitions, MockData, error) {
	var err error

	sports := map[int]Sport{
		1: {Name: "NRL"},
		2: {Name: "Basketball"},
		3: {Name: "AFL"},
		4: {Name: "Football"},
		5: {Name: "Cricket"},
	}

	mockData := make(MockData, len(sports))

	teams := map[int]Team{}
	competitions := map[int]Competition{}
	currTeamId := 0
	currCompId := 0

	//Add sports IDs to mockData
	for id := range sports {
		mockData[id] = struct {
			Teams        []int
			Competitions []int
		}{}

		data := mockData[id]

		//Generate between two and ten teams for each sport (and increment the team ID)
		numTeams, _ := strconv.Atoi(faker.Number().Between(2, 10))
		for i := 0; i < numTeams; i++ {
			currTeamId++
			team := faker.Team().Name()
			teams[currTeamId] = Team{Name: team, Sport: id}

			data.Teams = append(data.Teams, currTeamId)
		}

		//Generate between one and three competitions for each sport (and increment the competition ID)
		numCompetitions, _ := strconv.Atoi(faker.Number().Between(1, 3))
		for i := 0; i < numCompetitions; i++ {
			currCompId++
			competition := faker.App().Name() + " Cup"
			competitions[currCompId] = Competition{Name: competition, Sport: id}

			data.Competitions = append(data.Competitions, currCompId)
		}

		mockData[id] = data
	}

	return sports, teams, competitions, mockData, err
}

// Add sports table and populate it with sports
func (e *eventsRepo) seedSports(sports Sports) error {
	statement, err := e.db.Prepare(`CREATE TABLE IF NOT EXISTS sports (id INT PRIMARY KEY, name TEXT)`)

	if err != nil {
		return err
	}
	_, err = statement.Exec()

	if err != nil {
		return err
	}

	//Add all sports to the sports table
	log.Printf("Adding %d sports to the sports table...", len(sports))
	for id, sport := range sports {
		statement, err = e.db.Prepare(`INSERT OR IGNORE INTO sports(id, name) VALUES (?,?)`)

		if err != nil {
			return err
		}

		_, err = statement.Exec(id, sport.Name)

		if err != nil {
			return err
		}
	}

	return err
}

// Add teams table and its index to sports, and populate it with teams
func (e *eventsRepo) seedTeams(teams Teams) error {
	var statement *sql.Stmt
	var err error

	statement, err = e.db.Prepare(`CREATE TABLE IF NOT EXISTS teams (id INT PRIMARY KEY, name TEXT, sport_id INT, FOREIGN KEY(sport_id) REFERENCES sports(id))`)

	if err != nil {
		return err
	}
	_, err = statement.Exec()

	if err != nil {
		return err
	}

	statement, err = e.db.Prepare(`CREATE INDEX IF NOT EXISTS idx_teams_sport_id ON teams (sport_id)`)

	if err != nil {
		return err
	}
	_, err = statement.Exec()

	if err != nil {
		return err
	}

	//Add all teams to the teams table
	log.Printf("Adding %d teams to the teams table...", len(teams))
	for id, team := range teams {
		statement, err = e.db.Prepare(`INSERT OR IGNORE INTO teams(id, name, sport_id) VALUES (?,?,?)`)

		if err != nil {
			return err
		}

		_, err = statement.Exec(id, team.Name, team.Sport)

		if err != nil {
			return err
		}
	}

	return err
}

// Add competitions table and its index to sports, and populate it with competitions
func (e *eventsRepo) seedCompetitions(competitions Competitions) error {
	var statement *sql.Stmt
	var err error

	statement, err = e.db.Prepare(`CREATE TABLE IF NOT EXISTS competitions (id INT PRIMARY KEY, name TEXT, sport_id INT, FOREIGN KEY(sport_id) REFERENCES sports(id))`)

	if err != nil {
		return err
	}
	_, err = statement.Exec()

	if err != nil {
		return err
	}

	statement, err = e.db.Prepare(`CREATE INDEX IF NOT EXISTS idx_competitions_sport_id ON competitions (sport_id)`)

	if err != nil {
		return err
	}
	_, err = statement.Exec()

	if err != nil {
		return err
	}

	//Add all competitions to the competitions table
	log.Printf("Adding %d competitions to the competitions table...", len(competitions))
	for id, competition := range competitions {
		statement, err = e.db.Prepare(`INSERT OR IGNORE INTO competitions(id, name, sport_id) VALUES (?,?,?)`)

		if err != nil {
			return err
		}

		_, err = statement.Exec(id, competition.Name, competition.Sport)

		if err != nil {
			return err
		}
	}

	return err
}

// Add events table and its indices to sports and competitions, and populate it with events randomly generated from the teams, sports and competitions
func (e *eventsRepo) seedEvents(data MockData) error {
	var statement *sql.Stmt
	var err error

	statement, err = e.db.Prepare(`CREATE TABLE IF NOT EXISTS events (id INT PRIMARY KEY, sport_id INT, competition_id INT, home_team INT, away_team INT, home_score INT, away_score INT, advertised_start_time DATETIME , FOREIGN KEY(sport_id) REFERENCES sports(id) FOREIGN KEY(competition_id) REFERENCES competitions(id))`)

	if err != nil {
		return err
	}
	_, err = statement.Exec()

	if err != nil {
		return err
	}

	statement, err = e.db.Prepare(`CREATE INDEX IF NOT EXISTS idx_events_sport_id ON events (sport_id)`)

	if err != nil {
		return err
	}
	_, err = statement.Exec()

	if err != nil {
		return err
	}
	statement, err = e.db.Prepare(`CREATE INDEX IF NOT EXISTS idx_events_competition_id ON events (competition_id)`)

	if err != nil {
		return err
	}
	_, err = statement.Exec()

	if err != nil {
		return err
	}
	statement, err = e.db.Prepare(`CREATE INDEX IF NOT EXISTS idx_events_home_team ON events (home_team)`)

	if err != nil {
		return err
	}
	_, err = statement.Exec()

	if err != nil {
		return err
	}
	statement, err = e.db.Prepare(`CREATE INDEX IF NOT EXISTS idx_events_home_team ON events (away_team)`)

	if err != nil {
		return err
	}
	_, err = statement.Exec()

	if err != nil {
		return err
	}

	//Add 100 events to the events table using randomised data from the mockData map
	log.Print("Adding 100 events to the events table...")
	for i := 1; i <= 100; i++ {
		statement, err = e.db.Prepare(`INSERT OR IGNORE INTO events(id, sport_id, competition_id, home_team, away_team, home_score, away_score, advertised_start_time) VALUES (?,?,?,?,?,?,?,?)`)

		if err != nil {
			return err
		}

		sportId, err := strconv.Atoi(faker.Number().Between(1, len(data)))

		if err != nil {
			return err
		}

		//Ensure home and away teams are different
		homeTeamId := faker.Number().Between(1, len(data[sportId].Teams))
		awayTeamId := faker.Number().Between(1, len(data[sportId].Teams))
		for homeTeamId == awayTeamId {
			awayTeamId = faker.Number().Between(1, len(data[sportId].Teams))
		}

		_, err = statement.Exec(
			i,
			sportId,
			faker.Number().Between(1, len(data[sportId].Competitions)),
			homeTeamId,
			awayTeamId,
			0,
			0,
			faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
		)

		if err != nil {
			return err
		}
	}

	return err
}
