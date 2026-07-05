# Mini Log Query Engine with Golang

This project is a mini log query engine implemented in Golang. It allows you to query log data using a simple SQL-like syntax. The engine supports basic operations such as filtering & sorting on log data.

> This is a challenge project, and the code is not fully optimized or production-ready. It is intended for learning and demonstration purposes.

## Todo

- [x] Implement basic equal match operation
- [x] Implement sorting
- [ ] Make CLI less painful to use (allow user to move cursor to the left)
- [ ] Implement OR conditions
- [ ] Implement multiple conditions with AND/OR
- [ ] Implement equal match by timestamp field
- [ ] Implement gte & lte match by timestamp field
- [ ] Implement contains match operation
- [ ] Implement `stats` command

## Usage

```bash
# Run the log query engine with a log file
go run cmd/main.go logs.txt

# Apply query
> level="ERROR" AND message="Database connection failed" sort timestamp desc

#  Quit the engine
> quit
```
