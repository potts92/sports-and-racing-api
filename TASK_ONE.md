# Task One
## Task
Add another filter to the existing RPC, so we can call `ListRaces` asking for races that are visible only.
> We'd like to continue to be able to fetch all races regardless of their visibility, so try naming your filter as logically as possible. https://cloud.google.com/apis/design/standard_methods#list

## Options
### Filter in `scanRaces`
- Could add a filter to the `scanRaces` function and filter the slice of races passed to the function based on the visibility. Need to ensure if no visibility filter passed, we don't filter at all

### SQL Solution
- Could append to the `WHERE` clause in the SQL query to filter based on the visibility. Need to ensure if no visibility filter passed, by default we still return results regardless of visibility

## Final Solution
- Given the current solution for filtering by the `meeting_id` is SQL based, we will continue to use this method for filtering by visibility.
  - This will allow us to keep the filtering logic in one place and ensure that the filtering is done at the database level rather than in the application code.
  - Filtering in the application code could also result in significant latency (depending on the number of races returned) from having to loop through the full slice of results.
  - This will also allow us to keep the `scanRaces` function clean and simple.
- Protobuf message field is set to an optional boolean so that the generated Go struct field can differentiate between being unset, and set to `false`

## Sending a Request
Make a request for races filtered by meeting_id...

```bash
curl -X "POST" "http://localhost:8000/v1/list-races" \
     -H 'Content-Type: application/json' \
     -d $'{
  "filter": {
    "meeting_ids": [
      1,
      5
    ]
  }
}'
```

## Testing
- A test to ensure that visibility is correctly filtered by has been added to the [races_test.go](racing/db/races_test.go) file.
  - Utilises the in-memory SQL Lite database as this is already using mocked data, so doesn't make sense to add another mocked database (or queries) to memory while also providing valuable test results
- Run the tests using from the root of the module with:
```bash 
cd ./racing
go test ./...
```