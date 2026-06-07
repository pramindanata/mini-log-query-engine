package logen_test

import (
	"logen"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("should return correct AST for a single condition", func(t *testing.T) {
		tokens := []logen.Token{
			{Type: logen.TokenTypeField, Value: "fieldA"},
			{Type: logen.TokenTypeOperator, Value: "="},
			{Type: logen.TokenTypeValue, Value: "valueA"},
		}

		parser := logen.NewParser(tokens)
		actual, err := parser.Parse()

		expected := logen.ASTNode[logen.ASTNodeValueCondition]{
			Type: logen.ASTNodeTypeCondition,
			Value: logen.ASTNodeValueCondition{
				Field:    "fieldA",
				Operator: "=",
				Value:    "valueA",
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
