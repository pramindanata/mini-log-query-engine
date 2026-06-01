package logen_test

import (
	"logen"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseRawLogs(t *testing.T) {
	t.Run("should return logs when valid input is given", func(t *testing.T) {
		raw := []string{
			"2026-06-01T08:00:00 INFO User login success",
			"2026-06-01T08:01:00 ERROR Database connection failed",
		}

		expected := []logen.Log{
			{
				Timestamp: time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC),
				Level:     "INFO",
				Message:   "User login success",
			},
			{
				Timestamp: time.Date(2026, 6, 1, 8, 1, 0, 0, time.UTC),
				Level:     "ERROR",
				Message:   "Database connection failed",
			},
		}

		result, err := logen.ParseRawLogs(raw)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return error when raw log line is malformed", func(t *testing.T) {
		raw := []string{
			"2026-06-01T08:00:00 INFO",
		}

		_, err := logen.ParseRawLogs(raw)

		assert.EqualError(t, err, "failed to parse text: 2026-06-01T08:00:00 INFO")
	})

	t.Run("should return error when timestamp is invalid", func(t *testing.T) {
		raw := []string{
			"2026-06-01 08:00:00 INFO User login success",
		}

		_, err := logen.ParseRawLogs(raw)

		assert.EqualError(t, err, "failed to parse time: parsing time \"2026-06-01\" as \"2006-01-02T15:04:05\": cannot parse \"\" as \"T\"")
	})
}

func TestParseRawQuery(t *testing.T) {
	t.Run("should handle basic equal query", func(t *testing.T) {
		expected := logen.Query{
			Filters: []logen.QueryFilter{
				{
					Field:    "level",
					Operator: "=",
					Value:    "ERROR",
				},
			},
		}

		result, err := logen.ParseRawQuery("level=ERROR")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return error when invalid field is given", func(t *testing.T) {
		_, err := logen.ParseRawQuery("invalid=something")

		assert.EqualError(t, err, "invalid field \"invalid\" is given")
	})
}

func TestQueryLog(t *testing.T) {
	logs := []logen.Log{
		{
			Timestamp: time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC),
			Level:     "INFO",
			Message:   "User login success",
		},
		{
			Timestamp: time.Date(2026, 6, 1, 8, 1, 0, 0, time.UTC),
			Level:     "ERROR",
			Message:   "Database connection failed",
		},
	}

	t.Run("should return logs that equal match a level", func(t *testing.T) {
		query := logen.Query{
			Filters: []logen.QueryFilter{
				{
					Field:    "level",
					Operator: "=",
					Value:    "ERROR",
				},
			},
		}

		expected := []logen.Log{logs[1]}
		result, err := logen.QueryLogs(logs, query)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("should return logs that equal match a message", func(t *testing.T) {
		query := logen.Query{
			Filters: []logen.QueryFilter{
				{
					Field:    "message",
					Operator: "=",
					Value:    "User login success",
				},
			},
		}

		expected := []logen.Log{logs[0]}
		result, err := logen.QueryLogs(logs, query)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestStringifyLogs(t *testing.T) {
	logs := []logen.Log{
		{
			Timestamp: time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC),
			Level:     "INFO",
			Message:   "User login success",
		},
		{
			Timestamp: time.Date(2026, 6, 1, 8, 1, 0, 0, time.UTC),
			Level:     "ERROR",
			Message:   "Database connection failed",
		},
	}

	t.Run("should return stringified logs", func(t *testing.T) {
		expected := "2026-06-01T08:00:00 INFO User login success\n" +
			"2026-06-01T08:01:00 ERROR Database connection failed"

		result := logen.StringifyLogs(logs)

		assert.Equal(t, expected, result)
	})

	t.Run("should return message when no logs are given", func(t *testing.T) {
		expected := "No logs found"

		result := logen.StringifyLogs([]logen.Log{})

		assert.Equal(t, expected, result)
	})
}
