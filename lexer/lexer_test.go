package lexer

import (
	"alg/lexer/token"
	"testing"
)

func TestSingleCharToken(t *testing.T) {
	input := `,;:.()[]{}!`

	tests := []struct {
		expectedType   token.TokenType
		expectedLexeme string
	}{
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
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
