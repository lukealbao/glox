package main

import "testing"

func TTypes(tokens []Token) []TokenType {
	v := make([]TokenType, len(tokens))
	for i, t := range tokens {
		v[i] = t.Type
	}

	return v
}

func TestScannerTokenTypes(t *testing.T) {
	cases := map[string][]TokenType{
		// LL(1)
		"(){},.-+;!*=></": []TokenType{
			LEFT_PAREN,
			RIGHT_PAREN,
			LEFT_BRACE,
			RIGHT_BRACE,
			COMMA,
			DOT,
			MINUS,
			PLUS,
			SEMICOLON,
			BANG,
			STAR,
			EQUAL,
			GREATER,
			LESS,
			SLASH,
			EOF,
		},
		// LL(2)
		"!= == <= >=": []TokenType{
			BANG_EQUAL,
			EQUAL_EQUAL,
			LESS_EQUAL,
			GREATER_EQUAL,
			EOF,
		},
		// Comments
		"// comment\n": []TokenType{
			EOF,
		},
		// Numbers
		"123456 1 1.23": []TokenType{NUMBER, NUMBER, NUMBER, EOF},
		// Maximum munch
		"orthogonal andy and or": []TokenType{IDENTIFIER, IDENTIFIER, AND, OR, EOF},
	}

	for source, types := range cases {
		scanner := NewScanner(source)
		scanner.ScanTokens()

		for i, tt := range types {
			if scanner.tokens[i].Type != tt {
				t.Errorf("\n```\n%s\n```\nExpected: %s, Got: %s\n",
					source, tt, scanner.tokens[i].Type)
			}
		}
	}
}
