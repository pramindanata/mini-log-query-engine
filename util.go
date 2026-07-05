package logen

import (
	"bufio"
	"cmp"
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

func QueryLogs(logs []Log, rawQuery string) ([]Log, error) {
	result := make([]Log, 0)
	lexer := NewLexer(rawQuery)
	tokens := make([]Token, 0)

	for {
		token, err := lexer.GetNextToken()

		if err != nil {
			return result, fmt.Errorf("failed to get next token: %w", err)
		}

		tokens = append(tokens, token)

		if token.Type == TokenTypeEOF {
			break
		}
	}

	parser := NewParser(tokens)
	query, err := parser.Parse()

	if err != nil {
		return result, fmt.Errorf("failed to parse tokens: %w", err)
	}

	for _, log := range logs {
		kept := true

		for _, filter := range query.Filter.Items {
			selectedField := ""

			switch filter.Field {
			case "level":
				selectedField = log.Level
			case "message":
				selectedField = log.Message
			}

			if filter.Operator == "=" && selectedField != filter.Value {
				kept = false
			}
		}

		if kept {
			result = append(result, log)
		}
	}

	if query.Sort.Field != "" {
		slices.SortFunc(result, func(a, b Log) int {
			switch query.Sort.Field {
			case "level":
				if query.Sort.Direction == "DESC" {
					return cmp.Compare(b.Level, a.Level)
				}

				return cmp.Compare(a.Level, b.Level)
			case "message":
				if query.Sort.Direction == "DESC" {
					return cmp.Compare(b.Message, a.Message)
				}

				return cmp.Compare(a.Message, b.Message)
			case "timestamp":
				if query.Sort.Direction == "DESC" {
					return b.Timestamp.Compare(a.Timestamp)
				}

				return a.Timestamp.Compare(b.Timestamp)
			}

			return 0
		})
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
