package logen_test

import (
	"logen"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	t.Run("single type: quoted value", func(t *testing.T) {
		t.Run("should return tokens with correct quoted value", func(t *testing.T) {
			lexer := logen.NewLexer("\"value1\"")
			expected := []logen.Token{
				{Type: logen.TokenTypeValue, Value: "value1"},
				{Type: logen.TokenTypeEOF, Value: ""},
			}

			actual := collectTokens(t, lexer)

			assert.Equal(t, expected, actual)
		})

		t.Run("should return tokens with correct multiple quoted values", func(t *testing.T) {
			lexer := logen.NewLexer("\"value1\" \"value2\"")
			expected := []logen.Token{
				{Type: logen.TokenTypeValue, Value: "value1"},
				{Type: logen.TokenTypeValue, Value: "value2"},
				{Type: logen.TokenTypeEOF, Value: ""},
			}

			actual := collectTokens(t, lexer)

			assert.Equal(t, expected, actual)
		})

		t.Run("should return error when no end quote found", func(t *testing.T) {
			lexer := logen.NewLexer("\"value1")
			_, err := lexer.GetNextToken()

			assert.EqualError(t, err, "failed to read quoted value: no end quote found")
		})
	})

	t.Run("single type: operator", func(t *testing.T) {
		t.Run("should return tokens with equal operator", func(t *testing.T) {
			lexer := logen.NewLexer("=")
			expected := []logen.Token{
				{Type: logen.TokenTypeOperator, Value: "="},
				{Type: logen.TokenTypeEOF, Value: ""},
			}

			actual := collectTokens(t, lexer)

			assert.Equal(t, expected, actual)
		})
	})
}

func collectTokens(t *testing.T, lexer *logen.Lexer) []logen.Token {
	tokens := make([]logen.Token, 0)

	for {
		token, err := lexer.GetNextToken()
		assert.NoError(t, err)
		tokens = append(tokens, token)

		if token.Type == logen.TokenTypeEOF {
			break
		}
	}

	return tokens
}
