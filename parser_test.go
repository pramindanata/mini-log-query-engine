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

	t.Run("should return correct AST for a query with single filter & a sort clause", func(t *testing.T) {
		tokens := []logen.Token{
			{Type: logen.TokenTypeField, Value: "fieldA"},
			{Type: logen.TokenTypeOperator, Value: "="},
			{Type: logen.TokenTypeValue, Value: "valueA"},
			{Type: logen.TokenTypeSort, Value: "SORT"},
			{Type: logen.TokenTypeField, Value: "fieldA"},
			{Type: logen.TokenTypeSortDirection, Value: "ASC"},
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
			Sort: logen.ASTNodeSort{
				Type:      logen.ASTTypeSort,
				Field:     "fieldA",
				Direction: "ASC",
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("should return correct AST for a query with multiple filters & a sort clause", func(t *testing.T) {
		tokens := []logen.Token{
			{Type: logen.TokenTypeField, Value: "fieldA"},
			{Type: logen.TokenTypeOperator, Value: "="},
			{Type: logen.TokenTypeValue, Value: "valueA"},
			{Type: logen.TokenTypeLogicalOperator, Value: "AND"},
			{Type: logen.TokenTypeField, Value: "fieldB"},
			{Type: logen.TokenTypeOperator, Value: "="},
			{Type: logen.TokenTypeValue, Value: "valueB"},
			{Type: logen.TokenTypeSort, Value: "SORT"},
			{Type: logen.TokenTypeField, Value: "fieldA"},
			{Type: logen.TokenTypeSortDirection, Value: "ASC"},
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
				},
			},
			Sort: logen.ASTNodeSort{
				Type:      logen.ASTTypeSort,
				Field:     "fieldA",
				Direction: "ASC",
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

		assert.EqualError(t, err, "unexpected token `invalid` (FIELD) at post 3")
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

		assert.EqualError(t, err, "expected a filter after logical operator, got EOF instead")
	})
}
