package token

import (
	"strconv"
	"strings"
)

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

	// One or two character tokens
	LESS
	LESS_OR_EQUAL
	DIFF
	ASSIGN
	GREATER
	GREATER_OR_EQUAL
	MINUS
	MINUS_MINUS
	PLUS
	PLUS_PLUS
	SLASH
	STAR
	STAR_STAR
	EQUAL_EQUAL
	DOT_DOT // .. for array ranges

	// Literals
	IDENTIFIER

	// Type keywords
	INTEGER
	INTEGER_TYPE
	DOUBLE
	DOUBLE_TYPE
	STRING
	STRING_TYPE
	CHARACTER
	CHARACTER_TYPE
	BOOLEAN
	BOOLEAN_TYPE

	TABLE // Array/Table

	// Program structure keywords
	ALGORITHM
	VARIABLE
	CONSTANT
	TYPE
	BEGIN
	END

	// Function/Method keywords
	FUNCTION
	END_FUNCTION

	METHOD
	END_METHOD

	RETURN

	// Structure keywords
	STRUCTURE
	END_STRUCT

	// Control flow keywords
	IF
	THEN
	ELSE
	ELSEIF
	ENDIF

	FOR
	TO
	STEP
	ENDFOR

	WHILE
	DO
	ENDWHILE

	REPEAT
	UNTIL

	// Logical operators
	AND
	OR
	NOT

	// I/O keywords
	WRITE
	READ

	// Boolean literals
	TRUE
	FALSE

	// Other keywords
	NIL
	CLASS

	MOD // Modulo operator
	OF  // For array declarations: array[1..10] of integer

	// Ambiguous/Special
	INDENT
	EOF
)

var keywords = map[string]TokenType{}

// RegisterKeywords populates the keyword lookup table from the loaded config.
// Must be called once, before the lexer runs.
func RegisterKeywords(kw map[string]TokenType) {
	for k, v := range kw {
		keywords[strings.ToLower(k)] = v
	}
	// maps.Copy(keywords, kw)
}

// LookupKeyword returns the TokenType for a surface word, or IDENTIFIER if
// the word is not a keyword in the active language config.
func LookupKeyword(word string) TokenType {
	if tt, ok := keywords[strings.ToLower(word)]; ok {
		return tt
	}
	return IDENTIFIER
}

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

	// One or two character tokens
	LESS:             "<",
	LESS_OR_EQUAL:    "<=",
	DIFF:             "<>",
	ASSIGN:           "<-",
	GREATER:          ">",
	GREATER_OR_EQUAL: ">=",
	MINUS:            "-",
	MINUS_MINUS:      "--",
	PLUS:             "+",
	PLUS_PLUS:        "++",
	SLASH:            "/",
	STAR:             "*",
	STAR_STAR:        "**",
	EQUAL_EQUAL:      "==",
	DOT_DOT:          "..",

	// Literals
	IDENTIFIER: "IDENTIFIER",

	// Type keywords  (surface names vary, but the token constant name is fixed)
	INTEGER:        "INTEGER_LITERAL",
	INTEGER_TYPE:   "INTEGER_TYPE",
	DOUBLE:         "REAL_LITERAL",
	DOUBLE_TYPE:    "REAL_TYPE",
	STRING:         "STRING_LITERAL",
	STRING_TYPE:    "STRING_TYPE",
	CHARACTER:      "CHAR_LITERAL",
	CHARACTER_TYPE: "CHAR_TYPE",
	BOOLEAN:        "BOOL_LITERAL",
	BOOLEAN_TYPE:   "BOOL_TYPE",
	TABLE:          "ARRAY",

	// Program structure
	ALGORITHM: "ALGORITHM",
	VARIABLE:  "VAR",
	CONSTANT:  "CONST",
	TYPE:      "TYPE",
	BEGIN:     "BEGIN",
	END:       "END",

	// Functions / Methods
	FUNCTION:     "FUNCTION",
	END_FUNCTION: "ENDFUNCTION",
	METHOD:       "METHOD",
	END_METHOD:   "ENDMETHOD",
	RETURN:       "RETURN",

	// Structures
	STRUCTURE:  "STRUCT",
	END_STRUCT: "ENDSTRUCT",

	// Control flow
	IF:       "IF",
	THEN:     "THEN",
	ELSE:     "ELSE",
	ELSEIF:   "ELSEIF",
	ENDIF:    "ENDIF",
	FOR:      "FOR",
	TO:       "TO",
	STEP:     "STEP",
	ENDFOR:   "ENDFOR",
	WHILE:    "WHILE",
	DO:       "DO",
	ENDWHILE: "ENDWHILE",
	REPEAT:   "REPEAT",
	UNTIL:    "UNTIL",

	// Logical
	AND: "AND",
	OR:  "OR",
	NOT: "NOT",
	MOD: "MOD",

	// I/O
	WRITE: "WRITE",
	READ:  "READ",

	// Boolean literals
	TRUE:  "TRUE",
	FALSE: "FALSE",

	// Other
	NIL:   "NIL",
	CLASS: "CLASS",
	OF:    "OF",

	// Special
	EOF: "EOF",
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
