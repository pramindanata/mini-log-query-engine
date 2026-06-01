# Mini Log Query Engine with Golang

This project is a mini log query engine implemented in Golang. It allows you to query log data using a simple SQL-like syntax. The engine supports basic operations such as filtering & sorting on log data.

> This is a challenge project, and the code is not fully optimized or production-ready. It is intended for learning and demonstration purposes.

## Todo

- [x] Implement basic equal match operation
- [ ] Implement sorting
- [ ] Implement contains match operation
- [ ] Implement equal match by timestamp field
- [ ] Implement gte & lte match by timestamp field
- [ ] Implement multiple conditions with AND/OR
- [ ] Implement `stats` command

## Usage

```bash
# Run the log query engine with a log file
go run cmd/main.go logs.txt

# Filter logs by level
> level=ERROR

#  Quit the engine
> quit
```
