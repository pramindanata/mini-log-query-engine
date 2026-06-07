package logen

import "fmt"

// Grammar (v1)
// (TODO) query -> expression sort
// expression -> condition (AND|OR condition)*
// condition -> field operator value
// (TODO) sort -> SORT field direction

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() (ASTNode[ASTNodeValueCondition], error) {
	return p.parseCondition()
}

// func (p *Parser) parseExpression() (ASTNode[any], error) {
// 	var result ASTNode[any]

// 	condition, err := p.parseCondition()

// 	if err != nil {
// 		return result, fmt.Errorf("failed to parse condition: %w", err)
// 	}

// 	result.Type = ASTNodeTypeBinary
// 	result.Value = ASTNodeValueBinary[ASTNodeValueCondition, any]{
// 		Left: condition,
// 	}

// 	return result, nil
// }

func (p *Parser) parseCondition() (ASTNode[ASTNodeValueCondition], error) {
	var result ASTNode[ASTNodeValueCondition]
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

	result.Type = ASTNodeTypeCondition
	result.Value = ASTNodeValueCondition{
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

func (p *Parser) consume(expectedTokenType string) (Token, error) {
	if p.pos >= len(p.tokens) {
		return Token{}, fmt.Errorf("EOF")
	}

	token := p.tokens[p.pos]

	if token.Type != expectedTokenType {
		return token, fmt.Errorf("expected token type %s got %s instead at pos %d", expectedTokenType, token.Type, p.pos)
	}

	p.pos++

	return token, nil
}

type ASTNode[T any] struct {
	Type  string
	Value T
}

type ASTNodeValueBinary[L any, R any] struct {
	Left  ASTNode[L]
	Right *ASTNode[R]
}

type ASTNodeValueCondition struct {
	Field    string
	Operator string
	Value    string
}

const (
	ASTNodeTypeBinary    = "Binary"
	ASTNodeTypeCondition = "Condition"
)
