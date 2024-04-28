# Task Five
Create a `sports` service that for sake of simplicity, implements a similar API to racing. This sports API can be called `ListEvents`. We'll leave it up to you to determine what you might think a sports event is made up off, but it should at minimum have an `id`, a `name` and an `advertised_start_time`.

> Note: this should be a separate service, not bolted onto the existing racing service. At an extremely high-level, the diagram below attempts to provide a visual representation showing the separation of services needed and flow of requests.
>
> ![](example.png)

## Considerations
### Structure of an `event` message
- An `event` message should have an `id`, a `name` (of the sport) and an `advertised_start_time` at minimum
- The `id` should be a unique identifier for the `event`
- We need a `home_team` and `away_team` for the `event` (and these should not be the same team when generating mock data)
- We cannot allow betting during a game, so we will derive a `status` field based on the `advertised_start_time` of the `event`. If the `advertised_start_time` is in the past, the `status` should be `CLOSED`, otherwise it should be `OPEN`
- We need to be able to update the `home_score` and `away_score` for the `event`, until the score is finalised. We will add a `score_finalised` flag to the `event` to indicate that the score is finalised. This can be passed as an optional parameter when updating the score to denote that the game has finished
  - A score may be updated after it has finished, but only once. We will set a `score_finalised` flag to `true` after the score has been updated when the `event` is `CLOSED` 
- Relationships between the tables are as follows:
  - A `sport` can have many `competition`s and multiple `team`s
  - A `competition` can have many teams
  - A `team` can have many `competition`s
  - An `event` can have one `sport`, one `competition`, one `home_team`, one `away_team`, one `home_score` and one `away_score`

### Updating an event's score
- Google best practices says that we need to return the updated `event` after updating the score
- Need to ensure we use prepared statements to avoid SQL injection
- Need to run the `UPDATE` query followed by a `SELECT` query to return the updated `event` if the `UPDATE` was successful
  - Could achieve this with a transaction, but using the Golang `sql` module we'll just run them sequentially
- We need to ensure that the `score_finalised` flag is set to `true` after the score has been updated, and that the score cannot be updated again after this flag is set (and that if it already was set to `false`, the score is not updated)

## Future Considerations
- If performance deteriorates as the number of `event`s increases, we can consider:
  - Foreign keys on `sport`, `competition` and `team` names rather than IDs so we don't need to `JOIN` on the `sport`, `competition` and `team` tables to get the names
    - Drawback would be having to update the `sport`, `competition` and `team` names in multiple places if they change
  - Setting up material views to store commonly run complex query (like our multi-`JOIN` query) results in materialised views to allow the database to retrieve the data without having to recompute the entire query
  - Using a caching layer to store the results of the query to avoid having to run the query every time (this is more useful if the data is read more often than it is written which may not be the case for live sports) 
  - Using a NoSQL database like MongoDB to store the data as it is more scalable than a relational database like MySQL -> given that we control the access patterns for this data, this would be a good fit
  - Partitioning the `event` data across multiple tables to improve performance (potentially by `sport`, `competition` or even on the derived `status` field if we move this into the database)
- Return different errors for `UpdateScores` RPC if the `event` is not found or if the `score_finalised` flag is set to `true`

## Final Solution
- Track `competition`s, `team`s and `sport`s in separate tables, joined to the `event`s table with indices added to the foreign keys to sure performance is optimised
- `faker` module used to generate all of these fields except for the `sport_name` - which we get around by defining a hardcoded map of sports
- Given that we are randomising the amount of data in the database for this service (as opposed to the racing service which only ever attempts to `INSERT IGNORE INTO` 100 rows in one table), we first check if the `sports.db` SQLite database exists and if it does, we do not run the `seed` function to repopulate it
- Rather than scanning returned rows to map teams/ sport/ competition IDs to their relevant names, we handle this with a `JOIN` in the SQL query
- A `Get` RPC on the `sports` service allows for fetching a single `event` by its `id`
- The `ListEvents` RPC on the `sports` service allows for fetching all `event`s
>I've not added filtering and sorting to this RPC as we've shown how we can do it previously, but in the future we could:
>    - Allow filtering by `sport` and `competition` with simple `WHERE` clauses in the SQL query (utilising relevant `JOIN`s)
>    - Allow filtering by team (regardless of home/ away) by adding an `OR` clause to the SQL query's `WHERE` clause (utilising relevant `JOIN`s) e.g. `WHERE home_team = 'Wildcats' OR away_team = 'Wildcats'`
- The `UpdateScore` RPC on the `sports` service allows for updating the score of an `event`, and setting the `score_finalised` flag to lock further updates to the score
  - We sequentially run the `UPDATE` and `SELECT` (if the `UPDATE` is successful) queries if the `score_finalised` flag is set to `false`
  - We return the updated `event` after the score is updated

## Sending a Request
Make a request for all events...
```bash
curl -X "POST" "http://localhost:8000/v1/list-events" \
     -H 'Content-Type: application/json' \
     -d $'{}'
```

Make a request for a specific race (ID of 11)...
```bash
curl "http://localhost:8000/v1/event/11" \
     -H 'Content-Type: application/json'
```

Make a request to update (and lock) a specific event's scores...
```bash
curl -X "PATCH" "http://localhost:8000/v1/event/11/update-score" \
     -H 'Content-Type: application/json' \
     -d $'{
  "away_score": "4",
  "finalised": true,
  "home_score": 5
}'
```

## Testing
- Tests to check that the new sports service's `List` and `UpdateScore` RPCs work as expected have been added to the [races_test.go](racing/db/races_test.go) file. We do so by testing that the new `List` function from the `eventsRepo` can filter by `score_finalised` (`true` and `false`) and that the `UpdateScore` function from the `eventsRepo` can correctly update the score of an `event` and set the `score_finalised` flag if it was not already finalised.
- `UpdateScore` test uses rows returned from the `List` function to ensure that we can test with finalised and non-finalised scores
- Run all tests from the root of the project with:
```bash 
cd ./sports
go test ./...
```