package logen

import (
	"fmt"
	"strings"
)

type Lexer struct {
	text string
	pos  int
}

func NewLexer(text string) *Lexer {
	return &Lexer{text: text, pos: 0}
}

func (l *Lexer) GetNextToken() (Token, error) {
	if l.pos >= len(l.text) {
		return Token{Type: TokenTypeEOF}, nil
	}

	l.removeWhiteScape()

	char := l.text[l.pos]

	switch char {
	case '"':
		token, err := l.readQuotedValue()

		if err != nil {
			return token, fmt.Errorf("failed to read quoted value: %w", err)
		}

		return token, nil
	case '=':
		token, err := l.readOperator()

		if err != nil {
			return token, fmt.Errorf("failed to read operator %w", err)
		}

		return token, nil
	}

	return Token{}, nil
}

func (l *Lexer) removeWhiteScape() {
	for l.text[l.pos] == ' ' {
		l.pos++
		continue
	}
}

func (l *Lexer) readQuotedValue() (Token, error) {
	started := false
	var value strings.Builder

	for l.pos < len(l.text) {
		char := l.text[l.pos]

		if !started && char == '"' {
			started = true
			l.pos++
			continue
		} else if started && char == '"' {
			started = false
			l.pos++
			break
		}

		value.WriteByte(char)
		l.pos++
	}

	if started {
		return Token{}, fmt.Errorf("no end quote found")
	}

	token := Token{
		Type:  TokenTypeValue,
		Value: strings.TrimSpace(value.String()),
	}

	return token, nil
}

func (l *Lexer) readOperator() (Token, error) {
	l.pos++
	token := Token{
		Type:  TokenTypeOperator,
		Value: "=",
	}

	return token, nil
}

type Token struct {
	Type  string
	Value string
}

const (
	TokenTypeEOF      = "EOF"
	TokenTypeField    = "FIELD"
	TokenTypeOperator = "OPERATOR"
	TokenTypeValue    = "VALUE"
)
