package lexer

import (
	"alg/config"
	"alg/lexer/token"
	"strings"
	"testing"
)

func loadConfig(t *testing.T, path string) *config.Config {
	t.Helper()
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("could not load config %q: %v", path, err)
	}
	return cfg
}

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
	toks, _ := l.ScanTokens()
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
	l := New(FileContext{source: []rune(input)})
	toks, _ := l.ScanTokens()
	if len(toks) == 0 {
		t.Fatalf("expected tokens, got none")
	}
	for i, tt := range tests {
		tok := toks[i]
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] type: expected=%q got=%q (lexeme %q)",
				i, tt.expectedType, tok.Type, tok.Lexeme)
		}
		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("test[%d] lexeme: expected=%q got=%q",
				i, tt.expectedLexeme, tok.Lexeme)
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
	l := New(FileContext{source: []rune(input)})
	toks, _ := l.ScanTokens()
	if len(toks) == 0 {
		t.Fatalf("expected tokens, got none")
	}
	for i, tt := range tests {
		tok := toks[i]
		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] type: expected=%q got=%q (lexeme %q)",
				i, tt.expectedType, tok.Type, tok.Lexeme)
		}
		if tok.Lexeme != tt.expectedLexeme {
			t.Fatalf("test[%d] lexeme: expected=%q got=%q",
				i, tt.expectedLexeme, tok.Lexeme)
		}
	}
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
		{token.ALGORITHM, "algorithm"},
		{token.VARIABLE, "var"},
		{token.CONSTANT, "const"},
		{token.TYPE, "type"},
		{token.BEGIN, "begin"},
		{token.END, "end"},
		{token.FUNCTION, "function"},
		{token.END_FUNCTION, "endfunction"},
		{token.METHOD, "method"},
		{token.END_METHOD, "endmethod"},
		{token.RETURN, "return"},
		{token.STRUCTURE, "struct"},
		{token.END_STRUCT, "endstruct"},
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
		{token.INTEGER_TYPE, "integer"},
		{token.DOUBLE_TYPE, "real"},
		{token.STRING_TYPE, "string"},
		{token.CHARACTER_TYPE, "char"},
		{token.BOOLEAN_TYPE, "bool"},
		{token.TABLE, "array"},
		{token.OF, "of"},
		{token.AND, "and"},
		{token.OR, "or"},
		{token.NOT, "not"},
		{token.MOD, "mod"},
		{token.WRITE, "write"},
		{token.READ, "read"},
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
	toks, _ := l.ScanTokens()
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

func TestIntegerEdgeCases(t *testing.T) {
	tests := []struct {
		input   string
		lexeme  string
		literal int64
	}{
		{"0", "0", 0},
		{"0XFF", "0XFF", 255},
		{"0x0", "0x0", 0},
		{"0B1111", "0B1111", 15},
		{"0b0", "0b0", 0},
		{"0O77", "0O77", 63},
		{"0o0", "0o0", 0},
		{"1000000", "1000000", 1_000_000},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(FileContext{source: []rune(tt.input)})
			toks, _ := l.ScanTokens()
			if len(toks) != 1 {
				t.Fatalf("expected 1 token, got %d", len(toks))
			}
			tok := toks[0]
			if tok.Type != token.INTEGER {
				t.Errorf("type: expected INTEGER, got %q", tok.Type)
			}
			if tok.Lexeme != tt.lexeme {
				t.Errorf("lexeme: expected %q, got %q", tt.lexeme, tok.Lexeme)
			}
			if tok.Literal.(int64) != tt.literal {
				t.Errorf("literal: expected %d, got %v", tt.literal, tok.Literal)
			}
		})
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
	toks, _ := l.ScanTokens()
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

func TestRealEdgeCases(t *testing.T) {
	tests := []struct {
		input   string
		lexeme  string
		literal float64
	}{
		{"0.0", "0.0", 0.0},
		{"1.0", "1.0", 1.0},
		{"1.0E10", "1.0E10", 1.0e10},
		{"2.5e3", "2.5e3", 2500.0},
		{"1.0e100", "1.0e100", 1.0e100},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(FileContext{source: []rune(tt.input)})
			toks, _ := l.ScanTokens()
			if len(toks) != 1 {
				t.Fatalf("expected 1 token, got %d: %v", len(toks), toks)
			}
			tok := toks[0]
			if tok.Type != token.DOUBLE {
				t.Errorf("type: expected DOUBLE, got %q", tok.Type)
			}
			if tok.Lexeme != tt.lexeme {
				t.Errorf("lexeme: expected %q, got %q", tt.lexeme, tok.Lexeme)
			}
			if tok.Literal.(float64) != tt.literal {
				t.Errorf("literal: expected %v, got %v", tt.literal, tok.Literal)
			}
		})
	}
}

func TestFrenchDecimalSeparator(t *testing.T) {
	cfg := loadConfig(t, "../lang/fr.toml")
	input := `3,14 0,5`
	tests := []struct {
		lexeme  string
		literal float64
	}{
		{"3,14", 3.14},
		{"0,5", 0.5},
	}
	l := New(FileContext{source: []rune(input)}, cfg)
	toks, _ := l.ScanTokens()
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

func TestFrenchZeroCommaNotReal(t *testing.T) {
	cfg := loadConfig(t, "../lang/fr.toml")
	l := New(FileContext{source: []rune("0,abc")}, cfg)
	toks, _ := l.ScanTokens()
	if len(toks) != 3 {
		t.Fatalf("expected 3 tokens, got %d: %v", len(toks), toks)
	}
	if toks[0].Type != token.INTEGER {
		t.Errorf("[0] expected INTEGER, got %q", toks[0].Type)
	}
	if toks[1].Type != token.COMMA {
		t.Errorf("[1] expected COMMA, got %q", toks[1].Type)
	}
	if toks[2].Type != token.IDENTIFIER {
		t.Errorf("[2] expected IDENTIFIER, got %q", toks[2].Type)
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
	toks, _ := l.ScanTokens()
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

func TestCharAllEscapes(t *testing.T) {
	tests := []struct {
		input   string
		literal rune
	}{
		{`'\n'`, '\n'},
		{`'\t'`, '\t'},
		{`'\r'`, '\r'},
		{`'\\'`, '\\'},
		{`'\''`, '\''},
		{`'\0'`, 0},
		{`'\x41'`, 'A'},
		{`'\u0041'`, 'A'},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(FileContext{source: []rune(tt.input)})
			toks, _ := l.ScanTokens()
			if len(toks) != 1 {
				t.Fatalf("expected 1 token, got %d", len(toks))
			}
			if toks[0].Type != token.CHARACTER {
				t.Errorf("expected CHARACTER, got %q", toks[0].Type)
			}
			if toks[0].Literal.(rune) != tt.literal {
				t.Errorf("expected %q (%d), got %q (%d)",
					tt.literal, tt.literal, toks[0].Literal, toks[0].Literal)
			}
		})
	}
}

func TestUnicodeCharLiteral(t *testing.T) {
	l := New(FileContext{source: []rune("'à'")})
	toks, _ := l.ScanTokens()
	if len(toks) != 1 {
		t.Fatalf("expected 1 token, got %d", len(toks))
	}
	if toks[0].Type != token.CHARACTER {
		t.Errorf("expected CHARACTER, got %q", toks[0].Type)
	}
	if toks[0].Literal.(rune) != 'à' {
		t.Errorf("expected 'à', got %q", toks[0].Literal)
	}
}

func TestStringEscapes(t *testing.T) {
	input := `"hello\nworld" "tab\there"`
	tests := []struct{ literal string }{
		{"hello\nworld"},
		{"tab\there"},
	}
	l := New(FileContext{source: []rune(input)})
	toks, _ := l.ScanTokens()
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

func TestStringAllEscapes(t *testing.T) {
	tests := []struct {
		input   string
		literal string
	}{
		{`"\n"`, "\n"},
		{`"\t"`, "\t"},
		{`"\r"`, "\r"},
		{`"\\"`, "\\"},
		{`"\""`, "\""},
		{`"\0"`, "\x00"},
		{`"\x41"`, "A"},
		{`"\u0041"`, "A"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := New(FileContext{source: []rune(tt.input)})
			toks, _ := l.ScanTokens()
			if len(toks) != 1 {
				t.Fatalf("expected 1 token, got %d", len(toks))
			}
			if toks[0].Type != token.STRING {
				t.Errorf("expected STRING, got %q", toks[0].Type)
			}
			if toks[0].Literal.(string) != tt.literal {
				t.Errorf("expected %q, got %q", tt.literal, toks[0].Literal)
			}
		})
	}
}

func TestComments(t *testing.T) {
	input := "42 // this is a comment\n/* block\n   comment */\n99"
	l := New(FileContext{source: []rune(input)})
	toks, _ := l.ScanTokens()
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
	l := New(FileContext{source: []rune("1..10")})
	toks, _ := l.ScanTokens()
	if len(toks) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(toks))
	}
	if toks[1].Type != token.DOT_DOT {
		t.Errorf("expected DOT_DOT, got %q", toks[1].Type)
	}
}

func TestDotDisambiguation(t *testing.T) {
	t.Run("dot between identifiers", func(t *testing.T) {
		l := New(FileContext{source: []rune("foo.bar")})
		toks, _ := l.ScanTokens()
		if len(toks) != 3 {
			t.Fatalf("expected 3 tokens, got %d: %v", len(toks), toks)
		}
		if toks[1].Type != token.DOT {
			t.Errorf("expected DOT, got %q", toks[1].Type)
		}
	})
	t.Run("dot dot in range with reals", func(t *testing.T) {
		l := New(FileContext{source: []rune("1.0..2.0")})
		toks, _ := l.ScanTokens()
		if len(toks) != 3 {
			t.Fatalf("expected 3 tokens, got %d: %v", len(toks), toks)
		}
		if toks[0].Type != token.DOUBLE {
			t.Errorf("[0] expected DOUBLE, got %q", toks[0].Type)
		}
		if toks[1].Type != token.DOT_DOT {
			t.Errorf("[1] expected DOT_DOT, got %q", toks[1].Type)
		}
		if toks[2].Type != token.DOUBLE {
			t.Errorf("[2] expected DOUBLE, got %q", toks[2].Type)
		}
	})
	t.Run("integer followed by dot dot", func(t *testing.T) {
		l := New(FileContext{source: []rune("5..15")})
		toks, _ := l.ScanTokens()
		if len(toks) != 3 {
			t.Fatalf("expected 3 tokens, got %d", len(toks))
		}
		if toks[1].Type != token.DOT_DOT {
			t.Errorf("expected DOT_DOT, got %q", toks[1].Type)
		}
	})
}

func TestZeroFollowedByIdentifier(t *testing.T) {
	l := New(FileContext{source: []rune("0abc")})
	toks, _ := l.ScanTokens()
	if len(toks) != 2 {
		t.Fatalf("expected 2 tokens, got %d: %v", len(toks), toks)
	}
	if toks[0].Type != token.INTEGER || toks[0].Literal.(int64) != 0 {
		t.Errorf("expected INTEGER(0), got %q %v", toks[0].Type, toks[0].Literal)
	}
	if toks[1].Type != token.IDENTIFIER || toks[1].Lexeme != "abc" {
		t.Errorf("expected IDENTIFIER(abc), got %q %q", toks[1].Type, toks[1].Lexeme)
	}
}

func TestIntegerAdjacentToIdentifier(t *testing.T) {
	l := New(FileContext{source: []rune("42abc")})
	toks, _ := l.ScanTokens()
	if len(toks) != 2 {
		t.Fatalf("expected 2 tokens, got %d: %v", len(toks), toks)
	}
	if toks[0].Type != token.INTEGER || toks[0].Literal.(int64) != 42 {
		t.Errorf("expected INTEGER(42), got %q %v", toks[0].Type, toks[0].Literal)
	}
	if toks[1].Type != token.IDENTIFIER || toks[1].Lexeme != "abc" {
		t.Errorf("expected IDENTIFIER(abc), got %q %q", toks[1].Type, toks[1].Lexeme)
	}
}

func TestOperatorsNoWhitespace(t *testing.T) {
	input := `x<-42+y*2==z`
	l := New(FileContext{source: []rune(input)})
	toks, _ := l.ScanTokens()
	expected := []token.TokenType{
		token.IDENTIFIER,
		token.ASSIGN,
		token.INTEGER,
		token.PLUS,
		token.IDENTIFIER,
		token.STAR,
		token.INTEGER,
		token.EQUAL_EQUAL,
		token.IDENTIFIER,
	}
	if len(toks) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %v", len(expected), len(toks), toks)
	}
	for i, tt := range expected {
		if toks[i].Type != tt {
			t.Errorf("[%d] expected %q, got %q (lexeme %q)", i, tt, toks[i].Type, toks[i].Lexeme)
		}
	}
}

func TestKeywordCaseInsensitive(t *testing.T) {
	cfg := loadConfig(t, "../lang/en.toml")
	input := `IF If iF`
	l := New(FileContext{source: []rune(input)}, cfg)
	toks, _ := l.ScanTokens()
	if len(toks) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(toks))
	}
	for i, tok := range toks {
		if tok.Type != token.IF {
			t.Errorf("[%d] expected IF, got %q (lexeme %q)", i, tok.Type, tok.Lexeme)
		}
	}
}

func TestLineAndColumnTracking(t *testing.T) {
	input := "foo\nbar\nbaz"
	l := New(FileContext{source: []rune(input)})
	toks, _ := l.ScanTokens()
	if len(toks) != 3 {
		t.Fatalf("expected 3 tokens, got %d", len(toks))
	}
	cases := []struct{ line, col int }{
		{1, 3},
		{2, 3},
		{3, 3},
	}
	for i, c := range cases {
		if toks[i].Line != c.line {
			t.Errorf("[%d] line: expected %d, got %d", i, c.line, toks[i].Line)
		}
		if toks[i].Column != c.col {
			t.Errorf("[%d] col: expected %d, got %d", i, c.col, toks[i].Column)
		}
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty source", func(t *testing.T) {
		l := New(FileContext{source: []rune("")})
		toks, _ := l.ScanTokens()
		if len(toks) != 0 {
			t.Errorf("expected 0 tokens, got %d", len(toks))
		}
	})
	t.Run("whitespace only", func(t *testing.T) {
		l := New(FileContext{source: []rune("   \t\r\n   ")})
		toks, _ := l.ScanTokens()
		if len(toks) != 0 {
			t.Errorf("expected 0 tokens, got %d", len(toks))
		}
	})
	t.Run("comment only", func(t *testing.T) {
		l := New(FileContext{source: []rune("// nothing\n/* also nothing */")})
		toks, _ := l.ScanTokens()
		if len(toks) != 0 {
			t.Errorf("expected 0 tokens, got %d", len(toks))
		}
	})
}

func TestUnterminatedStringRecovery(t *testing.T) {
	// newline inside string → error recorded, scan resumes on next line
	input := "\"unterminated\n99"
	l := New(FileContext{source: []rune(input)})
	toks, diags := l.ScanTokens()

	if !diags.HasErrors() {
		t.Error("expected at least one diagnostic error")
	}
	found := false
	for _, d := range diags.Errors() {
		if strings.Contains(d.Message, "unterminated string") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'unterminated string' diagnostic, got:\n%s", diags.Format())
	}
	if len(toks) != 1 {
		t.Fatalf("expected 1 token after recovery, got %d: %v", len(toks), toks)
	}
	if toks[0].Literal.(int64) != 99 {
		t.Errorf("expected INTEGER(99), got %v", toks[0].Literal)
	}
}

func TestUnterminatedBlockCommentRecovery(t *testing.T) {
	input := "42 /* this never closes"
	l := New(FileContext{source: []rune(input)})
	toks, diags := l.ScanTokens()

	if !diags.HasErrors() {
		t.Error("expected unterminated block comment diagnostic")
	}
	found := false
	for _, d := range diags.Errors() {
		if strings.Contains(d.Message, "unterminated block comment") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 'unterminated block comment' diagnostic, got:\n%s", diags.Format())
	}
	if len(toks) != 1 || toks[0].Literal.(int64) != 42 {
		t.Errorf("expected INTEGER(42), got %v", toks)
	}
}

func TestUnknownCharacterRecovery(t *testing.T) {
	input := `42 @ 99 $ 1`
	l := New(FileContext{source: []rune(input)})
	toks, diags := l.ScanTokens()

	if len(diags.Errors()) != 2 {
		t.Errorf("expected 2 error diagnostics, got %d:\n%s",
			len(diags.Errors()), diags.Format())
	}
	if len(toks) != 3 {
		t.Fatalf("expected 3 tokens, got %d: %v", len(toks), toks)
	}
	for i, want := range []int64{42, 99, 1} {
		if toks[i].Literal.(int64) != want {
			t.Errorf("[%d] expected %d, got %v", i, want, toks[i].Literal)
		}
	}
}

func TestFullProgram(t *testing.T) {
	input := `
algorithm: test;
var
    x : integer;
    y : real;
begin
    x <- 10;
    y <- 3.14;
    if (x == 10) then:
        write("hello");
    endif
end
`
	cfg := loadConfig(t, "../lang/en.toml")
	l := New(FileContext{source: []rune(input)}, cfg)
	toks, diags := l.ScanTokens()

	if diags.HasErrors() {
		t.Fatalf("unexpected errors in full program scan:\n%s", diags.Format())
	}

	seen := make(map[token.TokenType]bool)
	for _, tok := range toks {
		seen[tok.Type] = true
	}
	mustContain := []token.TokenType{
		token.ALGORITHM, token.VARIABLE, token.INTEGER_TYPE, token.DOUBLE_TYPE,
		token.BEGIN, token.ASSIGN, token.INTEGER, token.DOUBLE,
		token.IF, token.EQUAL_EQUAL, token.THEN, token.WRITE,
		token.STRING, token.ENDIF, token.END,
	}
	for _, want := range mustContain {
		if !seen[want] {
			t.Errorf("token %q missing from full-program scan", want)
		}
	}
}
