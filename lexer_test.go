package logen_test

import (
	"fmt"
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

	t.Run("single type: field", func(t *testing.T) {
		t.Run("should return tokens with correct field", func(t *testing.T) {
			lexer := logen.NewLexer("fieldA")
			expected := []logen.Token{
				{Type: logen.TokenTypeField, Value: "fieldA"},
				{Type: logen.TokenTypeEOF, Value: ""},
			}

			actual := collectTokens(t, lexer)

			assert.Equal(t, expected, actual)
		})

		t.Run("should return tokens with multiple correct fields", func(t *testing.T) {
			lexer := logen.NewLexer("fieldA fieldB")
			expected := []logen.Token{
				{Type: logen.TokenTypeField, Value: "fieldA"},
				{Type: logen.TokenTypeField, Value: "fieldB"},
				{Type: logen.TokenTypeEOF, Value: ""},
			}

			actual := collectTokens(t, lexer)

			assert.Equal(t, expected, actual)
		})
	})

	t.Run("combination", func(t *testing.T) {
		t.Run("should return tokens for basic expression", func(t *testing.T) {
			lexer := logen.NewLexer("fieldA=\"valueA\"")
			expected := []logen.Token{
				{Type: logen.TokenTypeField, Value: "fieldA"},
				{Type: logen.TokenTypeOperator, Value: "="},
				{Type: logen.TokenTypeValue, Value: "valueA"},
				{Type: logen.TokenTypeEOF, Value: ""},
			}

			actual := collectTokens(t, lexer)

			assert.Equal(t, expected, actual)
		})

		t.Run("should return tokens for basic expression where each token seperated by white space", func(t *testing.T) {
			lexer := logen.NewLexer("fieldA = \"valueA\"")
			expected := []logen.Token{
				{Type: logen.TokenTypeField, Value: "fieldA"},
				{Type: logen.TokenTypeOperator, Value: "="},
				{Type: logen.TokenTypeValue, Value: "valueA"},
				{Type: logen.TokenTypeEOF, Value: ""},
			}

			actual := collectTokens(t, lexer)

			assert.Equal(t, expected, actual)
		})
	})

	t.Run("invalid char position error", func(t *testing.T) {
		type testTable struct {
			text string
			pos  int
		}

		table := []testTable{
			{text: "123Field", pos: 0},
			{text: "123\"Field\"", pos: 0},
			{text: "Field123", pos: 5},
		}

		for _, row := range table {
			t.Run(fmt.Sprintf("should return error at pos %d for %s", row.pos, row.text), func(t *testing.T) {
				lexer := logen.NewLexer(row.text)
				maxIteration := 10
				var err error

				for range maxIteration {
					_, err = lexer.GetNextToken()

					if err != nil {
						break
					}
				}

				assert.EqualError(t, err, fmt.Sprintf("failed to get next token: invalid char at pos %d", row.pos))
			})
		}
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
