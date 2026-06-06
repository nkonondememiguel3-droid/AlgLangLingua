package token

import "strconv"

type TokenType int

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL TokenType = iota
	EOF

	IDENTIFIER
	INT

	ASSIGN
	PLUS
	MINUS
	BANG
	STAR
	SLASH

	LT
	GT
	EQ
	NOT_EQ

	COMMA
	SEMICOLON

	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE

	FUNCTION
	LET
	TRUE
	FALSE
	IF
	ELSE
	RETURN
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdentifer(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}

var tokens = [...]string{
	ILLEGAL: "illegal",
	EOF:     "eof",

	IDENTIFIER: "identifier",
	INT:        "int",

	ASSIGN: "assign",
	PLUS:   "plus",
	MINUS:  "minus",
	BANG:   "!",
	STAR:   "*",
	SLASH:  "/",

	LT:     "less",
	GT:     "greater",
	EQ:     "==",
	NOT_EQ: "!=",

	COMMA:     "comma",
	SEMICOLON: "semicolon",

	LEFT_PAREN:  "left_parenthese",
	RIGHT_PAREN: "right_parenthese",
	LEFT_BRACE:  "left_brace",
	RIGHT_BRACE: "right_brace",

	FUNCTION: "function",
	LET:      "let",
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
