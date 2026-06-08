package token

import "strconv"

type TokenType int

type Token struct {
	Type    TokenType // type of the token.
	Lexeme  string    // the lexeme of the token. eg: "morning" -> string
	Literal any       // the literal of the token. eg: "morning" -> morning
	Line    int
	Column  int
	File    string // the file where the token occured
}

const (
	_ TokenType = iota

	// Single character tokens
	COMMA
	SEMICOLON
	COLON
	DOT
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACKET
	RIGHT_BRACKET
	LEFT_BRACE
	RIGHT_BRACE
	BANG
)

var keywords = map[string]TokenType{}

var tokens = [...]string{
	// Single character tokens
	COMMA:         "comma",
	SEMICOLON:     "semicolon",
	COLON:         "colon",
	DOT:           "dot",
	LEFT_PAREN:    "left_paren",
	RIGHT_PAREN:   "right_paren",
	LEFT_BRACKET:  "left_bracket",
	RIGHT_BRACKET: "right_bracket",
	LEFT_BRACE:    "left_brace",
	RIGHT_BRACE:   "right_brace",
	BANG:          "bang",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token [ADD], the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token [IDENT], the string is "IDENT").
func (tok TokenType) String() string {
	s := ""
	if 0 <= tok && tok < TokenType(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}
