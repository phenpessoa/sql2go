package parser

import (
	"strings"
)

type lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

func newLexer(input string) *lexer {
	l := &lexer{input: strings.ReplaceAll(input, "\r\n", "\n")}
	l.readChar()
	return l
}

func (l *lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *lexer) moveBack() {
	if l.pos > 0 && l.pos < len(l.input) {
		l.readPos = l.pos
		l.pos--
		l.ch = l.input[l.readPos]
	} else {
		l.ch = 0
	}
}

func (l *lexer) readLine() string {
	pos := l.pos
	for l.ch != '\n' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

const (
	nameBytes = " name: "
)

// isName detects if we are in a name token,
// if true it will consume the bytes
func (l *lexer) isName() bool {
	counter := 0
	for counter < len(nameBytes) &&
		l.input[l.readPos+counter] == nameBytes[counter] {
		counter++
	}
	if counter == len(nameBytes) {
		for i := 0; i < len(nameBytes); i++ {
			l.readChar()
		}
		return true
	}
	return false
}

func (l *lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *lexer) readRawInput() string {
	pos := l.pos
outer:
	for {
		switch l.ch {
		case '-':
			if l.peekChar() == '-' {
				break outer
			}
		case ';', '\'', '"', '`':
			break outer
		case '\n', 0:
			break outer
		}
		l.readChar()
	}
	data := l.input[pos:l.pos]
	l.moveBack()
	return data
}

func (l *lexer) nextToken() token {
	var t token

	l.skipWhitespace()

	switch l.ch {
	case '-':
		if l.peekChar() == '-' {
			l.readChar()
			if l.isName() {
				t.literal = "-- name: "
				t.typ = tokenTypeName
			} else {
				t.literal = "--"
				t.typ = tokenTypeComment
			}
		} else {
			t.typ = tokenTypeUndefined
		}
	case '\n':
		t.typ = tokenTypeNewLine
		t.literal = "\n"
	case '"', '`':
		t.typ = tokenTypeIdentifier
		t.literal = string(l.ch)
	case ';':
		t.typ = tokenTypeSemicolon
		t.literal = ";"
	case 0:
		t.typ = tokenTypeEOF
	case '\'':
		t.typ = tokenTypeStringLiteral
		t.literal = "'"
	default:
		t.typ = tokenTypeRawInput
		t.literal = l.readRawInput()
	}

	l.readChar()
	return t
}
