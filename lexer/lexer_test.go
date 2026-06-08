package lexer

import (
	"alg/config"
	"alg/lexer/token"
	"testing"
)

func TestOneOrTwoCharToken(t *testing.T) {
	input := `< <= <> <- > >= - -- + ++ / * ** == ..`

	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.LESS, "<"},
		{token.LESS_OR_EQUAL, "<="},
		{token.DIFF, "<>"},
		{token.ASSIGN, "<-"},
		{token.GREATER, ">"},
		{token.GREATER_OR_EQUAL, ">="},
		{token.MINUS, "-"},
		{token.MINUS_MINUS, "--"},
		{token.PLUS, "+"},
		{token.PLUS_PLUS, "++"},
		{token.SLASH, "/"},
		{token.STAR, "*"},
		{token.STAR_STAR, "**"},
		{token.EQUAL_EQUAL, "=="},
		{token.DOT_DOT, ".."},
	}

	fileContext := FileContext{
		source: []rune(input),
	}
	l := New(fileContext)
	tokens := l.ScanTokens()

	if len(tokens) == 0 {
		t.Fatalf("Length of the tokens scanned should be greater than zero. not zero.\n")
	}

	for i, tt := range tests {
		token_ := tokens[i]
		if token_.Type != tt.expectedType {
			t.Fatalf("Test[%d] - token type wrong. expected=%q, got=%q. (%q)", i, tt.expectedType, token_.Type, token_.Lexeme)
		}

		if token_.Lexeme != tt.expectedLexeme {
			t.Fatalf("Test[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedLexeme, token_.Literal)
		}
	}
}

func TestSingleCharToken(t *testing.T) {
	input := `,;:.()[]{}!`

	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.COLON, ":"},
		{token.DOT, "."},
		{token.LEFT_PAREN, "("},
		{token.RIGHT_PAREN, ")"},
		{token.LEFT_BRACKET, "["},
		{token.RIGHT_BRACKET, "]"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},
		{token.BANG, "!"},
	}

	fileContext := FileContext{
		source: []rune(input),
	}
	l := New(fileContext)
	tokens := l.ScanTokens()

	if len(tokens) == 0 {
		t.Fatalf("Length of the tokens scanned should be greater than zero. not zero.\n")
	}

	for i, tt := range tests {
		token_ := tokens[i]
		if token_.Type != tt.expectedType {
			t.Fatalf("Test[%d] - token type wrong. expected=%q, got=%q. (%q)", i, tt.expectedType, token_.Type, token_.Lexeme)
		}

		if token_.Lexeme != tt.expectedLexeme {
			t.Fatalf("Test[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedLexeme, token_.Literal)
		}
	}
}

// loadConfig is a test helper that loads a .toml file and fails the
// test immediately if it cannot be parsed.
func loadConfig(t *testing.T, path string) *config.Config {
	t.Helper()
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("could not load config %q: %v", path, err)
	}
	return cfg
}

func TestKeywordsEnglish(t *testing.T) {
	cfg := loadConfig(t, "../lang/en.toml")

	input := `algorithm var const type begin end
              function endfunction method endmethod return
              struct endstruct
              if then else elseif endif
              for to step endfor
              while do endwhile
              repeat until
              integer real string char bool array of
              and or not mod
              write read
              true false nil class`

	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		// Program structure
		{token.ALGORITHM, "algorithm"},
		{token.VARIABLE, "var"},
		{token.CONSTANT, "const"},
		{token.TYPE, "type"},
		{token.BEGIN, "begin"},
		{token.END, "end"},
		// Functions / Methods
		{token.FUNCTION, "function"},
		{token.END_FUNCTION, "endfunction"},
		{token.METHOD, "method"},
		{token.END_METHOD, "endmethod"},
		{token.RETURN, "return"},
		// Structures
		{token.STRUCTURE, "struct"},
		{token.END_STRUCT, "endstruct"},
		// Control flow
		{token.IF, "if"},
		{token.THEN, "then"},
		{token.ELSE, "else"},
		{token.ELSEIF, "elseif"},
		{token.ENDIF, "endif"},
		{token.FOR, "for"},
		{token.TO, "to"},
		{token.STEP, "step"},
		{token.ENDFOR, "endfor"},
		{token.WHILE, "while"},
		{token.DO, "do"},
		{token.ENDWHILE, "endwhile"},
		{token.REPEAT, "repeat"},
		{token.UNTIL, "until"},
		// Types
		{token.INTEGER_TYPE, "integer"},
		{token.DOUBLE_TYPE, "real"},
		{token.STRING_TYPE, "string"},
		{token.CHARACTER_TYPE, "char"},
		{token.BOOLEAN_TYPE, "bool"},
		{token.TABLE, "array"},
		{token.OF, "of"},
		// Logical
		{token.AND, "and"},
		{token.OR, "or"},
		{token.NOT, "not"},
		{token.MOD, "mod"},
		// I/O
		{token.WRITE, "write"},
		{token.READ, "read"},
		// Literals / other
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.NIL, "nil"},
		{token.CLASS, "class"},
	}

	runKeywordTest(t, cfg, input, tests)
}

func TestKeywordsFrench(t *testing.T) {
	cfg := loadConfig(t, "../lang/fr.toml")

	input := `algorithme variables constantes type debut fin
              fonction finfonction methode finmethode retourne
              structure finstructure
              si alors sinon sinonsi finsi
              pour à pas finpour
              tantque faire fintantque
              repeter jusqua
              entier reel chaine caractere booleen tableau de
              et ou non mod
              ecrire lire
              vrai faux nil classe`

	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.ALGORITHM, "algorithme"},
		{token.VARIABLE, "variables"},
		{token.CONSTANT, "constantes"},
		{token.TYPE, "type"},
		{token.BEGIN, "debut"},
		{token.END, "fin"},
		{token.FUNCTION, "fonction"},
		{token.END_FUNCTION, "finfonction"},
		{token.METHOD, "methode"},
		{token.END_METHOD, "finmethode"},
		{token.RETURN, "retourne"},
		{token.STRUCTURE, "structure"},
		{token.END_STRUCT, "finstructure"},
		{token.IF, "si"},
		{token.THEN, "alors"},
		{token.ELSE, "sinon"},
		{token.ELSEIF, "sinonsi"},
		{token.ENDIF, "finsi"},
		{token.FOR, "pour"},
		{token.TO, "à"},
		{token.STEP, "pas"},
		{token.ENDFOR, "finpour"},
		{token.WHILE, "tantque"},
		{token.DO, "faire"},
		{token.ENDWHILE, "fintantque"},
		{token.REPEAT, "repeter"},
		{token.UNTIL, "jusqua"},
		{token.INTEGER_TYPE, "entier"},
		{token.DOUBLE_TYPE, "reel"},
		{token.STRING_TYPE, "chaine"},
		{token.CHARACTER_TYPE, "caractere"},
		{token.BOOLEAN_TYPE, "booleen"},
		{token.TABLE, "tableau"},
		{token.OF, "de"},
		{token.AND, "et"},
		{token.OR, "ou"},
		{token.NOT, "non"},
		{token.MOD, "mod"},
		{token.WRITE, "ecrire"},
		{token.READ, "lire"},
		{token.TRUE, "vrai"},
		{token.FALSE, "faux"},
		{token.NIL, "nil"},
		{token.CLASS, "classe"},
	}

	runKeywordTest(t, cfg, input, tests)
}

func TestIdentifier(t *testing.T) {
	// Words that are NOT keywords in English should come back as IDENTIFIER.
	cfg := loadConfig(t, "../lang/en.toml")

	input := `foo bar _x myVar hello_world x123`

	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.IDENTIFIER, "foo"},
		{token.IDENTIFIER, "bar"},
		{token.IDENTIFIER, "_x"},
		{token.IDENTIFIER, "myVar"},
		{token.IDENTIFIER, "hello_world"},
		{token.IDENTIFIER, "x123"},
	}

	runKeywordTest(t, cfg, input, tests)
}

// runKeywordTest is shared by all keyword/identifier table tests.
func runKeywordTest(
	t *testing.T,
	cfg *config.Config,
	input string,
	tests []struct {
		expectedType   token.TokenType
		expectedLexeme string
	},
) {
	t.Helper()

	fc := FileContext{source: []rune(input)}
	l := New(fc, cfg)
	toks := l.ScanTokens()

	if len(toks) < len(tests) {
		t.Fatalf("expected at least %d tokens, got %d", len(tests), len(toks))
	}

	for i, tt := range tests {
		got := toks[i]
		if got.Type != tt.expectedType {
			t.Errorf("test[%d] type: expected %q, got %q (lexeme %q)",
				i, tt.expectedType, got.Type, got.Lexeme)
		}
		if got.Lexeme != tt.expectedLexeme {
			t.Errorf("test[%d] lexeme: expected %q, got %q",
				i, tt.expectedLexeme, got.Lexeme)
		}
	}
}

func TestIntegerLiterals(t *testing.T) {
	input := ` 0 42 255 0xFF 0b1010 0o17 `

	tests := []struct {
		lexeme  string
		literal int64
	}{
		{"0", 0},
		{"42", 42},
		{"255", 255},
		{"0xFF", 255},
		{"0b1010", 10},
		{"0o17", 15},
	}

	l := New(FileContext{source: []rune(input)})
	toks := l.ScanTokens()

	if len(toks) < len(tests) {
		t.Fatalf("expected %d tokens, got %d", len(tests), len(toks))
	}
	for i, tt := range tests {
		tok := toks[i]
		if tok.Type != token.INTEGER {
			t.Errorf("[%d] type: expected INTEGER, got %q", i, tok.Type)
		}
		if tok.Lexeme != tt.lexeme {
			t.Errorf("[%d] lexeme: expected %q, got %q", i, tt.lexeme, tok.Lexeme)
		}
		if tok.Literal.(int64) != tt.literal {
			t.Errorf("[%d] literal: expected %d, got %v", i, tt.literal, tok.Literal)
		}
	}
}

func TestRealLiterals(t *testing.T) {
	input := `3.14 0.5 1.0e10 2.5E-3`
	tests := []struct {
		lexeme  string
		literal float64
	}{
		{"3.14", 3.14},
		{"0.5", 0.5},
		{"1.0e10", 1.0e10},
		{"2.5E-3", 2.5e-3},
	}

	l := New(FileContext{source: []rune(input)})
	toks := l.ScanTokens()

	if len(toks) < len(tests) {
		t.Fatalf("expected %d tokens, got %d", len(tests), len(toks))
	}
	for i, tt := range tests {
		tok := toks[i]
		if tok.Type != token.DOUBLE {
			t.Errorf("[%d] type: expected DOUBLE, got %q", i, tok.Type)
		}
		if tok.Lexeme != tt.lexeme {
			t.Errorf("[%d] lexeme: expected %q, got %q", i, tt.lexeme, tok.Lexeme)
		}
		if tok.Literal.(float64) != tt.literal {
			t.Errorf("[%d] literal: expected %v, got %v", i, tt.literal, tok.Literal)
		}
	}
}

func TestFrenchDecimalSeparator(t *testing.T) {
	cfg, err := config.Load("../lang/fr.toml") // comma as separator
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	input := `3,14 0,5`
	tests := []struct {
		lexeme  string
		literal float64
	}{
		{"3,14", 3.14},
		{"0,5", 0.5},
	}

	l := New(FileContext{source: []rune(input)}, cfg)
	toks := l.ScanTokens()

	if len(toks) < len(tests) {
		t.Fatalf("expected %d tokens, got %d", len(tests), len(toks))
	}
	for i, tt := range tests {
		tok := toks[i]
		if tok.Type != token.DOUBLE {
			t.Errorf("[%d] type: expected DOUBLE, got %q", i, tok.Type)
		}
		if tok.Lexeme != tt.lexeme {
			t.Errorf("[%d] lexeme: expected %q, got %q", i, tt.lexeme, tok.Lexeme)
		}
		if tok.Literal.(float64) != tt.literal {
			t.Errorf("[%d] literal: expected %v, got %v", i, tt.literal, tok.Literal)
		}
	}
}

func TestCharLiterals(t *testing.T) {
	input := `'a' '\n' '\t' '\\'`
	tests := []struct {
		lexeme  string
		literal rune
	}{
		{"'a'", 'a'},
		{"'\n'", '\n'},
		{"'\t'", '\t'},
		{"'\\'", '\\'},
	}

	l := New(FileContext{source: []rune(input)})
	toks := l.ScanTokens()

	if len(toks) < len(tests) {
		t.Fatalf("expected %d tokens, got %d", len(tests), len(toks))
	}
	for i, tt := range tests {
		tok := toks[i]
		if tok.Type != token.CHARACTER {
			t.Errorf("[%d] type: expected CHARACTER, got %q", i, tok.Type)
		}
		if tok.Literal.(rune) != tt.literal {
			t.Errorf("[%d] literal: expected %q, got %q", i, tt.literal, tok.Literal)
		}
	}
}

func TestStringEscapes(t *testing.T) {
	// "hello\nworld"  "tab\there"
	input := `"hello\nworld" "tab\there"`
	tests := []struct {
		literal string
	}{
		{"hello\nworld"},
		{"tab\there"},
	}

	l := New(FileContext{source: []rune(input)})
	toks := l.ScanTokens()

	if len(toks) < len(tests) {
		t.Fatalf("expected %d tokens, got %d", len(tests), len(toks))
	}
	for i, tt := range tests {
		tok := toks[i]
		if tok.Type != token.STRING {
			t.Errorf("[%d] type: expected STRING, got %q", i, tok.Type)
		}
		if tok.Literal.(string) != tt.literal {
			t.Errorf("[%d] literal: expected %q, got %q", i, tt.literal, tok.Literal)
		}
	}
}

func TestComments(t *testing.T) {
	input := `
42 // this is a comment
/* block
   comment */
99
`
	l := New(FileContext{source: []rune(input)})
	toks := l.ScanTokens()

	// only 42 and 99 should survive; comments are discarded
	if len(toks) != 2 {
		t.Fatalf("expected 2 tokens, got %d: %v", len(toks), toks)
	}
	if toks[0].Literal.(int64) != 42 {
		t.Errorf("expected 42, got %v", toks[0].Literal)
	}
	if toks[1].Literal.(int64) != 99 {
		t.Errorf("expected 99, got %v", toks[1].Literal)
	}
}

func TestDotVsDotDot(t *testing.T) {
	// A plain dot must not consume the second dot of ".."
	input := `1..10`
	l := New(FileContext{source: []rune(input)})
	toks := l.ScanTokens()

	// expected: INTEGER(1)  DOT_DOT  INTEGER(10)
	if len(toks) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(toks))
	}
	if toks[1].Type != token.DOT_DOT {
		t.Errorf("expected DOT_DOT, got %q", toks[1].Type)
	}
}
