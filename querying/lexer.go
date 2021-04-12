package querying

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type lexerState struct {
	buf [2]rune
	pos int
	eot bool
}

func GetTokens(query string) ([]Token, error) {
	runes := []rune(query)
	var tokens []Token
	kw := map[string]TokenType{
		"SELECT": Select,
		"WHERE":  Where,
		"FROM":   From,
		"INSERT": Insert,
		"INTO":   Into,
		"VALUES": Values,
		"DELETE": Delete,
	}
	state := &lexerState{
		buf: [2]rune{' ', ' '},
		pos: 0,
		eot: false,
	}
	consumeRune(state, runes)
	consumeRune(state, runes)

	for !state.eot {
		r := state.buf[0]
		consumeRune(state, runes)
		switch r {
		case ' ':
			break
		case '@':
			consumeRune(state, runes)
			ident := getIdentifier(r, state, runes)
			tokens = append(tokens, Token{Type: Parameter, String: ident})
			break
		case ',':
			tokens = append(tokens, Token{Type: Comma})
		case '(':
			tokens = append(tokens, Token{Type: LeftParenthesis})
		case ')':
			tokens = append(tokens, Token{Type: RightParenthesis})
		case '|':
			tokens = append(tokens, Token{Type: Pipe})
		case '*':
			tokens = append(tokens, Token{Type: Asterisk})
		case '=':
			tokens = append(tokens, Token{Type: Equals})
		case '\'':
			str := getString(state, runes)
			tokens = append(tokens, Token{Type: String, String: str})
		default:
			if isNumber(r) || r == '-' {
				num, err := getNumber(r, state, runes)
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, Token{Type: Number, Number: int64(num)})
			} else if isIdentifier(r) {
				ident := getIdentifier(r, state, runes)
				token, ok := kw[strings.ToUpper(ident)]
				if ok {
					tokens = append(tokens, Token{Type: token})
				} else {
					tokens = append(tokens, Token{Type: Identifier, String: ident})
				}
			} else {
				return nil, fmt.Errorf("failed to tokenize query, unknown token: %c", r)
			}
		}
	}
	return tokens, nil
}

func getString(s *lexerState, runes []rune) string {
	var sb strings.Builder
	for !s.eot && s.buf[0] != '\'' {
		sb.WriteRune(s.buf[0])
		consumeRune(s, runes)
	}
	consumeRune(s, runes) // consume last '
	return sb.String()
}

func getNumber(r rune, s *lexerState, runes []rune) (int, error) {
	var sb strings.Builder
	sb.WriteRune(r)
	for !s.eot && isNumber(s.buf[0]) {
		sb.WriteRune(s.buf[0])
		consumeRune(s, runes)
	}
	return strconv.Atoi(sb.String())
}

func isIdentifier(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isNumber(r rune) bool {
	return unicode.IsDigit(r)
}

func getIdentifier(r rune, s *lexerState, runes []rune) string {
	var sb strings.Builder
	sb.WriteRune(r)
	for !s.eot && isIdentifier(s.buf[0]) {
		sb.WriteRune(s.buf[0])
		consumeRune(s, runes)
	}
	return sb.String()
}

func consumeRune(s *lexerState, runes []rune) {
	s.buf[0] = s.buf[1]

	if s.pos < len(runes) {
		s.buf[1] = runes[s.pos]
	}
	if s.pos > len(runes) {
		s.eot = true
	}
	s.pos++
}
