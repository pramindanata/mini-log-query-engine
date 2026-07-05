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

func TestQueryLogs(t *testing.T) {
	logs := []logen.Log{
		{
			Timestamp: time.Date(2026, 6, 1, 1, 0, 0, 0, time.UTC),
			Level:     "INFO",
			Message:   "User login success",
		},
		{
			Timestamp: time.Date(2026, 6, 1, 2, 1, 0, 0, time.UTC),
			Level:     "ERROR",
			Message:   "Database connection failed",
		},
		{
			Timestamp: time.Date(2026, 6, 1, 3, 0, 0, 0, time.UTC),
			Level:     "WARN",
			Message:   "Memory usage high",
		},
		{
			Timestamp: time.Date(2026, 6, 1, 4, 0, 0, 0, time.UTC),
			Level:     "INFO",
			Message:   "User logout",
		},
		{
			Timestamp: time.Date(2026, 6, 1, 5, 0, 0, 0, time.UTC),
			Level:     "ERROR",
			Message:   "Payment processing failed",
		},
		{
			Timestamp: time.Date(2026, 6, 1, 6, 0, 0, 0, time.UTC),
			Level:     "ERROR",
			Message:   "Database connection failed",
		},
	}

	t.Run("should return logs that equal match a level", func(t *testing.T) {
		query := "level=\"ERROR\""
		actual, err := logen.QueryLogs(logs, query)
		expected := []logen.Log{logs[1], logs[4], logs[5]}

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should return logs that equal match a message", func(t *testing.T) {
		query := "message=\"User login success\""
		actual, err := logen.QueryLogs(logs, query)
		expected := []logen.Log{logs[0]}

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should return logs that equal match multiple AND filters", func(t *testing.T) {
		query := "level=\"ERROR\" AND message=\"Database connection failed\""
		actual, err := logen.QueryLogs(logs, query)
		expected := []logen.Log{logs[1], logs[5]}

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should sort logs based on the given sort", func(t *testing.T) {
		query := "sort timestamp desc"
		actual, err := logen.QueryLogs(logs, query)
		expected := []logen.Log{logs[5], logs[4], logs[3], logs[2], logs[1], logs[0]}

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should filter & sort logs", func(t *testing.T) {
		query := "level=\"ERROR\" AND message=\"Database connection failed\" sort timestamp desc"
		actual, err := logen.QueryLogs(logs, query)
		expected := []logen.Log{logs[5], logs[1]}

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
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
