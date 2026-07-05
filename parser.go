package logen

import (
	"errors"
	"fmt"
)

// Grammar (v1)
// query -> expression sort
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
	filters := make([]ASTNodeFilter, 0)

	for {
		filter, err := p.parseFilter()

		if err != nil {
			return result, fmt.Errorf("failed to get filter: %w", err)
		}

		filters = append(filters, filter)
		_, err = p.parseLogicalOperator()

		if err != nil {
			if errors.Is(err, ErrEOF) {
				break
			}

			return result, fmt.Errorf("failed to get logical operator: %w", err)
		}
	}

	result.Filter = ASTNodeMultiple[ASTNodeFilter]{
		Type:  ASTTypeAnd,
		Items: filters,
	}

	return result, nil
}

func (p *Parser) parseFilter() (ASTNodeFilter, error) {
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
