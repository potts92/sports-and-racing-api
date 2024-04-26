## Entain BE Technical Test - Alex Potts
## Getting Started
### Pre-Requisites
- Assumes that the most recent version of Go is installed on your system (I have upgraded the packages to Go 1.22)

### Installation
1. Start the `racing` service...
```bash
cd ./racing
go build && ./racing
```
2. Start the `sports` service...
```bash
cd ./sports
go build && ./sports
```
3. Start the `api` service...
```bash
cd ./api
go build && ./api
```
4. Run `racing` tests...
```bash
cd ./racing
go test ./...
```

## Task Specific Notes/ Instructions
For information about considerations and decisions made for each task, please see the following files:
- [Task One](TASK_ONE.md)
- [Task Two](TASK_TWO.md)
- [Task Three](TASK_THREE.md)
- [Task Four](TASK_FOUR.md)
- [Task Five](TASK_FIVE.md)