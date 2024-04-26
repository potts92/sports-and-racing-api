# Task Three
## Task
Our races require a new `status` field that is derived based on their `advertised_start_time`'s. The status is simply, `OPEN` or `CLOSED`. All races that have an `advertised_start_time` in the past should reflect `CLOSED`.
> There's a number of ways this could be implemented. Just have a go!


## Options
### SQL Solution
- Derive the status in the SQL query and return it as a column in the result set (adds some latency), example below:
```sql
    SELECT
        created_at,
        CASE
            WHEN DATE(advertised_start_time) < CURRENT_TIMESTAMP() THEN 'OPEN'
            ELSE 'CLOSED'
            END AS status
       FROM races
  ```
### Mapping in `ScanRaces`
- Map the `advertised_start_time` to a `status` in the `ScanRaces` function (may add significant latency depending on size of the returned races array).

## Future Considerations
### Filtering on derived status
- If a future consideration is to filter by the status, this affects the implementation choice:
  - If we derive the status in the SQL query, we can filter by the status in the SQL query.
  - If we derive the status in the `ScanRaces` function, we can filter by the status in the `ScanRaces` function.

#### Filtering on derived status in SQL
- If we derive the status in the SQL query, we can filter by the status in the SQL query but will need to derive a mock table to filter by the derived status (adds significant latency), example below:
```sql
  SELECT *
  FROM (
           SELECT
               created_at,
               CASE
                   WHEN DATE(created_at) < CURRENT_TIMESTAMP() THEN 'true'
                   ELSE 'false'
                   END AS status
           FROM races
       ) AS derived_table
  WHERE
      status = 'OPEN';
  ```
- Alternatively, we could retrieve the current date time in Golang when building the SQL query and use that to filter before/ after to pseudo filter by status, while also deriving the status in the SQL query
```sql
    SELECT
        created_at,
        CASE
            WHEN DATE(advertised_start_time) < ? THEN 'OPEN'
            ELSE 'CLOSED'
            END AS status
    FROM races
    WHERE advertised_start_time < ?;

```

#### Filtering on derived status in `ScanRaces`
- If we derive the status in the `ScanRaces` function, we can filter by the status in the `ScanRaces` function while deriving the status (this would obviously potentially result in a huge amount of unnecessary data being returned, adding significant latency depending on size of the returned races array)

#### A combined approach
- We could combine both approaches by filtering in the SQL using the current date time retrieved in Golang and then deriving the status in the `ScanRaces` function. Would need to pass the date time from the `List` function to the `ScanRaces` function to ensure consistency in filtering and derivation. Depending on the size of expected results, this could be a good compromise between latency and complexity.

## Final Solution
- For simplicity's sake and to avoid adding complexity/ time to the SQL, we will map it in the ScanRaces function to take it "client-side".
- By comparing the `advertised_start_time` with the curren time using the [time](/usr/local/go/src/time/time.go) package's `After` method we ensure that we also set the race to `CLOSED` if the start time exactly matches the current time.
- We will also add a `status` enum field to the `Race` struct to enforce the status field to be either `OPEN` or `CLOSED`.
- If in future we need to filter by status, we can consider the combined approach mentioned above or the SQL approach filtering by a current date time.

## Sending a Request
Make a request for races (derived `status` will be returned in races)...
> If all races in your database are in the past and therefore not `OPEN`, you can remove the [racing.db](racing/db/racing.db) file and run the `racing` service again to generate new races (which will be in the past and future) to test the `OPEN` status.

```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d $'{}'
```

## Testing
- A test to ensure that the status is correctly derived based on the `advertised_start_time` field has been added to the [races_test.go](racing/db/races_test.go) file. We are unable to test the `scanRaces` function directly as it is a private function, so we test the `List` function which calls the `scanRaces` function. This time we do mock the database connection by providing expected rows to be returned when queries structured in a particular way are processed. This ensures we know **_exactly_** what to test for in terms of whether to expect an `OPEN` or `CLOSED` `status`.
- Run the tests using from the root of the module with:
```bash 
cd ./racing
go test ./...
```