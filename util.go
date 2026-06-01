package logen

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

const (
	TimestampLayout = "2006-01-02T15:04:05"
)

func GetRawLogs(filename string) ([]string, error) {
	result := make([]string, 0)
	file, err := os.Open(filename)

	if err != nil {
		return result, fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}

	return result, nil
}

func ParseRawLogs(rawLogs []string) ([]Log, error) {
	result := make([]Log, 0)

	for _, rawLog := range rawLogs {
		parts := strings.SplitN(rawLog, " ", 3)

		if len(parts) < 3 {
			return result, fmt.Errorf("failed to parse text: %s", rawLog)
		}

		timestamp, err := time.Parse(TimestampLayout, parts[0])

		if err != nil {
			return result, fmt.Errorf("failed to parse time: %w", err)
		}

		log := Log{
			Timestamp: timestamp,
			Level:     parts[1],
			Message:   parts[2],
		}

		result = append(result, log)
	}

	return result, nil
}

func ParseRawQuery(raw string) (Query, error) {
	var result Query

	parts := strings.SplitN(raw, "=", 2)

	if len(parts) < 2 {
		return result, fmt.Errorf("failed to parse query")
	}

	filters := make([]QueryFilter, 0)
	filter := QueryFilter{
		Field:    parts[0],
		Operator: "=",
		Value:    parts[1],
	}

	acceptedFields := []string{"timestamp", "level", "message"}

	if !slices.Contains(acceptedFields, filter.Field) {
		return result, fmt.Errorf("invalid field \"%s\" is given", filter.Field)
	}

	// TODO not working for a while because we use split by "="
	// acceptedOperator := []string{"="}

	// if !slices.Contains(acceptedOperator, filter.Operator) {
	// 	return result, fmt.Errorf("invalid operator \"%s\" is given", filter.Field)
	// }

	filters = append(filters, filter)
	result.Filters = filters

	return result, nil
}

func QueryLogs(logs []Log, query Query) ([]Log, error) {
	result := make([]Log, 0)

	for _, log := range logs {
		for _, filter := range query.Filters {
			// TODO filter timestamp
			selectedField := ""

			switch filter.Field {
			case "level":
				selectedField = log.Level
			case "message":
				selectedField = log.Message
			}

			if filter.Operator == "=" && selectedField == filter.Value {
				result = append(result, log)
			}
		}
	}

	return result, nil
}

func StringifyLogs(logs []Log) string {
	if len(logs) == 0 {
		return "No logs found"
	}

	var result strings.Builder

	for i, log := range logs {
		newLine := "\n"

		if i == len(logs)-1 {
			newLine = ""
		}

		timestamp := log.Timestamp.Format(TimestampLayout)
		fmt.Fprintf(&result, "%s %s %s%s", timestamp, log.Level, log.Message, newLine)
	}

	return result.String()
}

type Log struct {
	Timestamp time.Time
	Level     string
	Message   string
}

type Query struct {
	Filters []QueryFilter
}

type QueryFilter struct {
	Field    string
	Operator string
	Value    string
}
