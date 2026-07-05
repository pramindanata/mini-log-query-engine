package logen_test

import (
	"logen"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("should return correct AST for a query with single filter", func(t *testing.T) {
		tokens := []logen.Token{
			{Type: logen.TokenTypeField, Value: "fieldA"},
			{Type: logen.TokenTypeOperator, Value: "="},
			{Type: logen.TokenTypeValue, Value: "valueA"},
			{Type: logen.TokenTypeEOF, Value: ""},
		}

		parser := logen.NewParser(tokens)
		actual, err := parser.Parse()

		expected := logen.ASTNodeQuery{
			Filter: logen.ASTNodeMultiple[logen.ASTNodeFilter]{
				Type: logen.ASTTypeAnd,
				Items: []logen.ASTNodeFilter{
					logen.ASTNodeFilter{
						Type:     logen.ASTTypeFilter,
						Field:    "fieldA",
						Operator: "=",
						Value:    "valueA",
					},
				},
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should return correct AST for a query with multiple filters", func(t *testing.T) {
		tokens := []logen.Token{
			{Type: logen.TokenTypeField, Value: "fieldA"},
			{Type: logen.TokenTypeOperator, Value: "="},
			{Type: logen.TokenTypeValue, Value: "valueA"},
			{Type: logen.TokenTypeLogicalOperator, Value: "AND"},
			{Type: logen.TokenTypeField, Value: "fieldB"},
			{Type: logen.TokenTypeOperator, Value: "="},
			{Type: logen.TokenTypeValue, Value: "valueB"},
			{Type: logen.TokenTypeLogicalOperator, Value: "AND"},
			{Type: logen.TokenTypeField, Value: "fieldC"},
			{Type: logen.TokenTypeOperator, Value: "="},
			{Type: logen.TokenTypeValue, Value: "valueC"},
			{Type: logen.TokenTypeEOF, Value: ""},
		}

		parser := logen.NewParser(tokens)
		actual, err := parser.Parse()

		expected := logen.ASTNodeQuery{
			Filter: logen.ASTNodeMultiple[logen.ASTNodeFilter]{
				Type: logen.ASTTypeAnd,
				Items: []logen.ASTNodeFilter{
					logen.ASTNodeFilter{
						Type:     logen.ASTTypeFilter,
						Field:    "fieldA",
						Operator: "=",
						Value:    "valueA",
					},
					logen.ASTNodeFilter{
						Type:     logen.ASTTypeFilter,
						Field:    "fieldB",
						Operator: "=",
						Value:    "valueB",
					},
					logen.ASTNodeFilter{
						Type:     logen.ASTTypeFilter,
						Field:    "fieldC",
						Operator: "=",
						Value:    "valueC",
					},
				},
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should return error when a filter is followed by non logical operator token", func(t *testing.T) {
		tokens := []logen.Token{
			{Type: logen.TokenTypeField, Value: "fieldA"},
			{Type: logen.TokenTypeOperator, Value: "="},
			{Type: logen.TokenTypeValue, Value: "valueA"},
			{Type: logen.TokenTypeField, Value: "invalid"},
			{Type: logen.TokenTypeEOF, Value: ""},
		}

		parser := logen.NewParser(tokens)
		_, err := parser.Parse()

		assert.EqualError(t, err, "failed to get logical operator: failed to consume logical operator: expected token type LOGICAL_OPERATOR got `invalid` (FIELD) instead at pos 3")
	})

	t.Run("should return error when a logical operator is not followed by anything", func(t *testing.T) {
		tokens := []logen.Token{
			{Type: logen.TokenTypeField, Value: "fieldA"},
			{Type: logen.TokenTypeOperator, Value: "="},
			{Type: logen.TokenTypeValue, Value: "valueA"},
			{Type: logen.TokenTypeLogicalOperator, Value: "AND"},
			{Type: logen.TokenTypeEOF, Value: ""},
		}

		parser := logen.NewParser(tokens)
		_, err := parser.Parse()

		assert.EqualError(t, err, "failed to get filter: failed to get field: failed to consume field: EOF")
	})
}
