package db

const (
	eventsList  = "list"
	eventUpdate = "update"
	eventGet    = "get"
	defaultGet  = `
			SELECT
				e.id,
			    s.name,
			    c.name 'competition',
			    ht.name 'home_team',
			    at.name 'away_team',
			    e.home_score,
			    e.away_score,
				e.advertised_start_time ,
				e.score_finalised
			FROM events e
			LEFT JOIN sports s
			ON e.sport_id = s.id
			LEFT JOIN competitions c
			ON e.competition_id = c.id
			LEFT JOIN teams ht
			ON e.home_team = ht.id
			LEFT JOIN teams at
			ON e.away_team = at.id
		`
)

func getEventQueries() map[string]string {
	return map[string]string{
		eventsList: defaultGet + ";",
		eventUpdate: `
			UPDATE events
			SET home_score = ?, away_score = ?, score_finalised = ?
			WHERE id = ? AND score_finalised = FALSE;
		`,
		eventGet: defaultGet + "WHERE e.id = ?;",
	}
}
