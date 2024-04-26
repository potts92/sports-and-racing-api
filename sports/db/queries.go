package db

const (
	eventsList = "list"
)

func getEventQueries() map[string]string {
	return map[string]string{
		eventsList: `
			SELECT
				e.id,
			    s.name,
			    c.name 'competition',
			    ht.name 'home_team',
			    at.name 'away_team',
			    home_score,
			    away_score,
				advertised_start_time 
			FROM events e
			LEFT JOIN sports s
			ON e.sport_id = s.id
			LEFT JOIN competitions c
			ON e.competition_id = c.id
			LEFT JOIN teams ht
			ON e.home_team = ht.id
			LEFT JOIN teams at
			ON e.away_team = at.id;
		`,
	}
}
