package main

import (
	"fmt"
	"os"
	"strconv"
	"unicode"
)

var Keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Scanner struct {
	source   string
	tokens   []Token
	start    int
	current  int
	line     int
	column   int
	keywords map[string]TokenType
}

func NewScanner(source string) Scanner {
	return Scanner{
		source:   source,
		start:    0,
		current:  0,
		line:     1,
		column:   0,
		keywords: Keywords,
	}
}

// TODO: Lox adds this to a top-level Lox class. See
// http://www.craftinginterpreters.com/scanning.html#error-handling
// Also note that they implement a Lox.hadError variable to prevent
// code with syntax errors from being evaluated.
func SyntaxError(line int, column int, message string) {
	// Not pretty, but scanner.current is always 1 step ahead of the real "current"
	// rune. See Scanner.advance().
	column -= 1
	
	fmt.Fprintf(os.Stderr, "SyntaxError (%d:%d): %s\n", line, column, message)
	// panic("Hang the DJ!")
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	s.current++
	s.column++
	return rune(s.source[s.current-1])
}

// addToken accepts zero or one literals
func (s *Scanner) addToken(tt TokenType, literal ...interface{}) {
	text := string(s.source[s.start:s.current])
	s.tokens = append(s.tokens, Token{tt, text, literal, s.line})
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.ScanToken()
	}

	s.tokens = append(s.tokens, Token{EOF, "", nil, s.line})
	return s.tokens
}

// ScanToken consumes the scanner's source stream one rune at a time and
// coalesces them into tokens. Each token is added to the scanner's internal
// token list.
func (s *Scanner) ScanToken() {
	c := s.advance()

	switch c {
	// Single-rune tokens are easy
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)

	// For some 2-rune operators, we must peek at the next rune in order
	// to determine the token type
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}

	// Comments get discarded by the scanner!
	case '/':
		if s.match('/') {
			// Here the "//" has been consumed and we consume the
			// whole comment, leaving newlines to keep the scanner's
			// line cursor accurate.
			for s.Peek() /*peek() */ != '\n' && !s.isAtEnd() {
				s.advance()
			}

		} else {
			s.addToken(SLASH)
		}
	case '"':
		s.string()

	// Whitespace is consumed on ScanToken entry, but it has no production.
	// It does have side effects for the scanner state, though.
	case ' ', '\t', '\r':
		break
	case '\n':
		s.line++
		s.column = 0

	default:
		// Since we're using equality checks rune-by-rune, identifiers
		// and numbers are tail-ended in the default case. Not exactly
		// elegant, I guess.
		if unicode.IsDigit(c) {
			s.number()
		} else if unicode.IsLetter(c) || c == '_' {
			s.identifier()
		} else {
			SyntaxError(s.line, s.column, fmt.Sprintf("Illegal character: `%c'", c))
		}
	}
}

func (s *Scanner) identifier() {
	for c := s.Peek(); !s.isAtEnd() && unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_'; c = s.Peek() {
		s.advance()
	}

	value := s.source[s.start:s.current]

	if tokentype, ok := Keywords[value]; ok {
		s.addToken(tokentype)
		return
	}

	s.addToken(IDENTIFIER, value)
}

// number identifies only numbers without trailing or leading decimal points.
func (s *Scanner) number() {
	// Consume integral part
	for unicode.IsDigit(s.Peek()) {
		s.advance()
	}

	// Only consume decimal if there is a fractional part next.
	if s.Peek() == '.' && unicode.IsDigit(s.PeekNext()) {
		// Consume decimal
		s.advance()
		// Consume fractional
		for unicode.IsDigit(s.Peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		SyntaxError(s.line, s.column, err.Error())
	}

	s.addToken(NUMBER, value)
}

// string consumes an entire string as a single token. It may contain newlines.
func (s *Scanner) string() {
	// TODO: Lox does not support escaped characters. You'd add unescaping
	// in this method.

	for s.Peek() != '"' && !s.isAtEnd() {
		if s.Peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		SyntaxError(s.line, s.column, "Unterminated string")
	}

	// Consume the closing `"'.
	s.advance()

	// String token does not contain quote runes.
	value := s.source[s.start+1 : s.current-1]

	s.addToken(STRING, value)
}

func (s *Scanner) Peek() rune {
	if s.isAtEnd() {
		return '\x00'
	}

	return rune(s.source[s.current])
}

func (s *Scanner) PeekNext() rune {
	// I suppose we could feasibly implement a peeking method to accept some
	// level of lookahead (currently it's only two). Oh hey, there's a
	// sidebar comment suggesting that having two 0-ary functions informs
	// the reader that our grammar has only two characters of lookahead.
	ahead := 1

	if s.current+ahead >= len(s.source) {
		return '\x00'
	}

	return rune(s.source[s.current+ahead])
}

// match asserts the next character in the source, with a side effect
// of consuming it if the assertion is true.
func (s *Scanner) match(expected rune) bool {
	matches := s.Peek() == expected

	if matches {
		s.advance()
	}

	return matches
}
