package logen

import (
	"errors"
	"fmt"
)

// Grammar (v1)
// query -> expression | sort
// expression -> condition (AND|OR condition)*
// condition -> field operator value
// sort -> SORT field direction

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() (ASTNodeQuery, error) {
	return p.parseQuery()
}

func (p *Parser) parseQuery() (ASTNodeQuery, error) {
	var result ASTNodeQuery

	token := p.peek()

	// TODO refactor to parse sort clause once
	// Check if only sort is exist
	if token.Type == TokenTypeSort {
		sort, err := p.parseSortClause()

		if err != nil {
			return result, fmt.Errorf("failed to get sort clause: %w", err)
		}

		result.Sort = sort
	}

	logicalOperatorTokenFound := false
	filters := make([]ASTNodeFilter, 0)

	for {
		filter, err := p.parseFilterClause()

		if err != nil {
			if errors.Is(err, ErrEOF) && logicalOperatorTokenFound {
				return result, fmt.Errorf("expected a filter after logical operator, got EOF instead")
			} else if !errors.Is(err, ErrEOF) {
				return result, fmt.Errorf("failed to get filter clause: %w", err)
			}

			break
		}

		filters = append(filters, filter)
		nextToken := p.peek()

		if nextToken.Type == TokenTypeLogicalOperator {
			// TODO handle AND/OR grouping
			_, err = p.parseLogicalOperator()

			if err != nil {
				return result, fmt.Errorf("failed to get logical operator: %w", err)
			}

			logicalOperatorTokenFound = true

			continue
		}

		break
	}

	result.Filter = ASTNodeMultiple[ASTNodeFilter]{
		Type:  ASTTypeAnd,
		Items: filters,
	}

	token = p.peek()

	if token.Type == TokenTypeSort {
		sort, err := p.parseSortClause()

		if err != nil {
			return result, fmt.Errorf("failed to get sort clause: %w", err)
		}

		result.Sort = sort
	} else if token.Type != TokenTypeEOF {
		return result, fmt.Errorf("expected EOF or sort clause after filter clauses")
	}

	return result, nil
}

func (p *Parser) parseFilterClause() (ASTNodeFilter, error) {
	var result ASTNodeFilter
	field, err := p.parseField()

	if err != nil {
		return result, fmt.Errorf("failed to get field: %w", err)
	}

	operator, err := p.parseOperator()

	if err != nil {
		return result, fmt.Errorf("failed to get operator: %w", err)
	}

	value, err := p.parseValue()

	if err != nil {
		return result, fmt.Errorf("failed to get value: %w", err)
	}

	result = ASTNodeFilter{
		Type:     ASTTypeFilter,
		Field:    field.Value,
		Operator: operator.Value,
		Value:    value.Value,
	}

	return result, nil
}

func (p *Parser) parseSortClause() (ASTNodeSort, error) {
	var result ASTNodeSort

	_, err := p.parseSort()

	if err != nil {
		return result, fmt.Errorf("failed to get sort: %w", err)
	}

	field, err := p.parseField()

	if err != nil {
		return result, fmt.Errorf("failed to get field: %w", err)
	}

	direction, err := p.parseSortDirection()

	if err != nil {
		return result, fmt.Errorf("failed to get sort direction: %w", err)
	}

	token := p.peek()

	if token.Type != TokenTypeEOF {
		return result, fmt.Errorf("expected no tokens after sort clause")
	}

	result.Type = ASTTypeSort
	result.Field = field.Value
	result.Direction = direction.Value

	return result, nil
}

func (p *Parser) parseField() (Token, error) {
	result, err := p.consume(TokenTypeField)

	if err != nil {
		return result, fmt.Errorf("failed to consume field: %w", err)
	}

	return result, nil
}

func (p *Parser) parseOperator() (Token, error) {
	result, err := p.consume(TokenTypeOperator)

	if err != nil {
		return result, fmt.Errorf("failed to consume operator: %w", err)
	}

	return result, nil
}

func (p *Parser) parseValue() (Token, error) {
	result, err := p.consume(TokenTypeValue)

	if err != nil {
		return result, fmt.Errorf("failed to consume value: %w", err)
	}

	return result, nil
}

func (p *Parser) parseLogicalOperator() (Token, error) {
	result, err := p.consume(TokenTypeLogicalOperator)

	if err != nil {
		return result, fmt.Errorf("failed to consume logical operator: %w", err)
	}

	return result, nil
}

func (p *Parser) parseSort() (Token, error) {
	result, err := p.consume(TokenTypeSort)

	if err != nil {
		return result, fmt.Errorf("failed to consume sort: %w", err)
	}

	return result, nil
}

func (p *Parser) parseSortDirection() (Token, error) {
	result, err := p.consume(TokenTypeSortDirection)

	if err != nil {
		return result, fmt.Errorf("failed to consume sort direction: %w", err)
	}

	return result, nil
}

func (p *Parser) consume(expectedTokenType string) (Token, error) {
	token := p.tokens[p.pos]

	if token.Type == TokenTypeEOF {
		return Token{}, ErrEOF
	}

	if token.Type != expectedTokenType {
		return token, fmt.Errorf("expected token type %s got `%s` (%s) instead at pos %d", expectedTokenType, token.Value, token.Type, p.pos)
	}

	p.pos++

	return token, nil
}

func (p *Parser) peek() Token {
	token := p.tokens[p.pos]

	return token
}

type ASTNodeQuery struct {
	Filter ASTNodeMultiple[ASTNodeFilter]
	Sort   ASTNodeSort
}

type ASTNodeMultiple[T any] struct {
	Type  string
	Items []T
}

type ASTNodeFilter struct {
	Type     string
	Field    string
	Operator string
	Value    string
}

type ASTNodeSort struct {
	Type      string
	Field     string
	Direction string
}

const (
	ASTTypeAnd    = "AND"
	ASTTypeOr     = "OR"
	ASTTypeFilter = "CONDITION"
	ASTTypeSort   = "SORT"
)

var (
	ErrEOF = errors.New("EOF")
)
