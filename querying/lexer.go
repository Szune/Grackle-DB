package querying

import (
	"strings"
)

type lexerState struct {
	buf []rune
	pos int
	eof bool
}

func GetTokens(query string) []Token {
	runes := []rune(query)
	var tokens []Token
	kw := map[string]TokenType{
		"SELECT": Select,
		"FROM": From,
	}
	state := &lexerState{
		buf: []rune{' ', ' '},
		pos: 0,
		eof: false,
	}
	consume(state, runes)
	consume(state, runes)

	for !state.eof {
		r := state.buf[0]
		consume(state, runes)
		switch r {
		case ' ':
			break
		case ',':
			tokens = append(tokens, Token{Type:Comma})
		case '*':
			tokens = append(tokens, Token{Type:Asterisk})
		case '=':
			tokens = append(tokens, Token{Type:Equals})
		default:
			if isIdentifier(r) {
				ident := getIdentifier(r, state, runes)
				token, ok := kw[strings.ToUpper(ident)]
				if ok {
					tokens = append(tokens, Token{Type: token})
				} else {
					tokens = append(tokens, Token{Type: Identifier, String: ident})
				}
			}
		}
	}
	return tokens
}

func isIdentifier(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func getIdentifier(r rune, s *lexerState, runes []rune) string {
	var sb strings.Builder
	sb.WriteRune(r)
	for !s.eof && isIdentifier(s.buf[0]) {
		sb.WriteRune(s.buf[0])
		consume(s, runes)
	}
	return sb.String()
}

func consume(s *lexerState, runes []rune) {
	s.buf[0] = s.buf[1]

	if s.pos < len(runes) {
		s.buf[1] = runes[s.pos]
	}
	if s.pos > len(runes) {
		s.eof = true
	}
	s.pos++
}
