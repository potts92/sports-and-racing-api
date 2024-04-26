# Task Two
## Task
We'd like to see the races returned, ordered by their `advertised_start_time`
> Bonus points if you allow the consumer to specify an ORDER/SORT-BY they might be after.

## Options
### Order in `scanRaces`
- Could add sorting to the `scanRaces` function and sort the slice of races according to the `advertised_start_time`. Need to ensure if no order passed, we don't sort at all

### SQL Solution
- Could add an `ORDER BY` clause to the SQL query to order the results by the `advertised_start_time`. Need to ensure if no order passed, we don't order at all

## Final Solution
- Given the current solution for filtering is SQL based, it makes most sense to use this method for ordering as well.
  - Ordering in the application code could also result in significant latency (depending on the number of races returned)
- A `setSortOrder` function has been added to the `racesRepo` receiver that will append to the SQL query based on the sort attribute/ direction passed
  - Errors if a sort attribute is passed that does not match a column in the database
- `List` function calls the `setSortOrder` function after the `applyFilters` function
- `ListRacesRequestSortOrder` message added to the protobuf file is defined as accepting a `sort_attribute` string and `sort_direction` enum (ASC/DESC). With `ASC` being mapped to 0 in the enum, we ensure that this is the default value if no sort direction is passed

## Sending a Request
Make a request for races sorted by advertised_start_time in ascending order...

```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d $'{
  "sort": {
    "sort_direction": "ASC",
    "sort_attribute": "advertised_start_time"
  }
}'
```

## Testing
- Tests have been added to the [races_test.go](racing/db/races_test.go) file to test the correct functionality is applied with different variations of sort inputs
- A test has also been added to the [races_test.go](racing/db/races_test.go) file to ensure that an error is returned if an invalid sort attribute is passed
- Tests have also been refactored to follow table driven testing to follow Go best practices
- - Run the tests using from the root of the module with:
```bash 
cd ./racing
go test ./...
```