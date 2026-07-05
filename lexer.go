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
	if l.isDone() {
		return Token{Type: TokenTypeEOF}, nil
	}

	l.skipWhiteSpace()

	char := l.text[l.pos]

	switch {
	case char == '"':
		token, err := l.readQuotedValue()

		if err != nil {
			return token, fmt.Errorf("failed to read quoted value: %w", err)
		}

		return token, nil
	case char == '=':
		token, err := l.readOperator()

		if err != nil {
			return token, fmt.Errorf("failed to read operator %w", err)
		}

		return token, nil
	case l.isEnglishAlphabet(char):
		token, err := l.readLetter()

		if err != nil {
			return token, fmt.Errorf("failed to read letter: %w", err)
		}

		return token, nil
	}

	return Token{}, fmt.Errorf("failed to get next token: invalid char at pos %d", l.pos)
}

func (l *Lexer) skipWhiteSpace() {
	for l.text[l.pos] == ' ' {
		l.pos++
		continue
	}
}

func (l *Lexer) readQuotedValue() (Token, error) {
	started := false
	var value strings.Builder

	for !l.isDone() {
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

func (l *Lexer) readLetter() (Token, error) {
	var builder strings.Builder

	for !l.isDone() {
		char := l.text[l.pos]

		if l.isEnglishAlphabet(char) {
			builder.WriteByte(char)
			l.pos++
			continue
		}

		break
	}

	value := builder.String()
	token := Token{
		Type:  TokenTypeField,
		Value: value,
	}

	if strings.ToLower(value) == "and" {
		token = Token{
			Type:  TokenTypeLogicalOperator,
			Value: "AND",
		}
	} else if strings.ToLower(value) == "or" {
		token = Token{
			Type:  TokenTypeLogicalOperator,
			Value: "OR",
		}
	}

	return token, nil
}

func (l *Lexer) isEnglishAlphabet(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func (l *Lexer) isDone() bool {
	return l.pos >= len(l.text)
}

type Token struct {
	Type  string
	Value string
}

const (
	TokenTypeEOF             = "EOF"
	TokenTypeField           = "FIELD"
	TokenTypeOperator        = "OPERATOR"
	TokenTypeValue           = "VALUE"
	TokenTypeLogicalOperator = "LOGICAL_OPERATOR"
)
