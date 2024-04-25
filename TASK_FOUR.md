# Task Four
Introduce a new RPC, that allows us to fetch a single race by its ID.
> This link here might help you on your way: https://cloud.google.com/apis/design/standard_methods#get

## Considerations
- Need to follow Google best practices
- Need to return a 404 if no results are found
- We'll return the first result if multiple results are found (as we're fetching by ID)

## Final Solution
- Following the same pattern as the `ListRaces` RPC, a new `GetRace` RPC has been added to the `races.proto` file
  - RPC has been added to the `RacesService` service in the format suggested in https://cloud.google.com/apis/design/standard_methods#get. `race` is included in the RPC resource's path to make the purpose of the endpoint clear and to follow standard method naming conventions. The resource id is mapped to the URL path
  - This RPC accepts an `id` string and returns a `Race` message (note: returned resource maps to the entire response body, not nested within a field ala the `ListRaces` RPC)
- A new prepared statement has been added to the mapped queries returned by the `getRaceQueries` function that expects an `id` string
- `GetRace` RPC calls the `Get` function on the `racesRepo` with the `id` int and returns the result, which uses the above prepared statement to build a query that will retrieve any matching race. We use the same `scanRaces` function the `List` function uses to scan the results and format them into a `Race` message
  - If no results are found, we return a 404 status code with a message indicating no results were found

## Sending a Request
Make a request for a specific race (ID of 5)...

```bash
curl -X "POST" "http://localhost:8000/v1/race/5" \
     -H 'Content-Type: application/json'
```

## Testing
- Tests to check that the new RPC can correctly either return a race when a valid ID is passed, or return no race when an invalid ID is passed has been added to the [races_test.go](racing/db/races_test.go) file. We do so by testing the new `Get` function from the `racesRepo` struct
- Run all tests from the root of the project with:
```bash 
cd ./racing
go test ./...
```