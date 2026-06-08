package lexer

import (
	"alg/config"
	errs "alg/err"
	"alg/lexer/token"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type FileContext struct {
	Filename string
	File     *os.File
	source   []rune
	position int
	readPos  int
	ch       rune
	line     int
	column   int
}

type Lexer struct {
	fileStack   []*FileContext
	currentFile *FileContext
	decimalSep  rune
	Diagnostics errs.DiagnosticList
}

func If[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// diagnosticErrorf records an error at the current position into
// Diagnostics and returns it as a Go error for early returns in helpers.
func (l *Lexer) diagnosticErrorf(format string, args ...any) error {
	f := l.currentFile
	msg := fmt.Sprintf(format, args...)
	l.Diagnostics.Errorf(f.Filename, f.line, f.column, "%s", msg)
	return fmt.Errorf("%s", msg)
}

func New(fileContext FileContext, cfg ...*config.Config) *Lexer {
	l := &Lexer{
		currentFile: &fileContext,
		decimalSep:  '.',
	}
	if len(cfg) > 0 && cfg[0] != nil {
		token.RegisterKeywords(buildKeywordMap(cfg[0]))
		if cfg[0].Meta.DecimalSep == "," {
			l.decimalSep = ','
		}
	}
	if fileContext.Filename != "" {
		if err := l.PushFile(fileContext.Filename); err != nil {
			// File-open failures are fatal setup errors, not scan diagnostics.
			fmt.Fprintf(os.Stderr, "lexer: %v\n", err)
		}
	} else {
		if l.currentFile.line == 0 {
			l.currentFile.line = 1
			l.currentFile.column = 0
		}
	}
	return l
}

func buildKeywordMap(cfg *config.Config) map[string]token.TokenType {
	kw := cfg.Keywords
	return map[string]token.TokenType{
		kw.Algorithm: token.ALGORITHM,
		kw.Variable:  token.VARIABLE,
		kw.Constant:  token.CONSTANT,
		kw.Type:      token.TYPE,
		kw.Begin:     token.BEGIN,
		kw.End:       token.END,

		kw.Function:    token.FUNCTION,
		kw.EndFunction: token.END_FUNCTION,
		kw.Method:      token.METHOD,
		kw.EndMethod:   token.END_METHOD,
		kw.Return:      token.RETURN,

		kw.Structure: token.STRUCTURE,
		kw.EndStruct: token.END_STRUCT,

		kw.If:       token.IF,
		kw.Then:     token.THEN,
		kw.Else:     token.ELSE,
		kw.ElseIf:   token.ELSEIF,
		kw.EndIf:    token.ENDIF,
		kw.For:      token.FOR,
		kw.To:       token.TO,
		kw.Step:     token.STEP,
		kw.EndFor:   token.ENDFOR,
		kw.While:    token.WHILE,
		kw.Do:       token.DO,
		kw.EndWhile: token.ENDWHILE,
		kw.Repeat:   token.REPEAT,
		kw.Until:    token.UNTIL,

		kw.IntegerType:   token.INTEGER_TYPE,
		kw.DoubleType:    token.DOUBLE_TYPE,
		kw.StringType:    token.STRING_TYPE,
		kw.CharacterType: token.CHARACTER_TYPE,
		kw.BooleanType:   token.BOOLEAN_TYPE,
		kw.Table:         token.TABLE,
		kw.Of:            token.OF,

		kw.And: token.AND,
		kw.Or:  token.OR,
		kw.Not: token.NOT,
		kw.Mod: token.MOD,

		kw.Write: token.WRITE,
		kw.Read:  token.READ,

		kw.True:  token.TRUE,
		kw.False: token.FALSE,

		kw.Nil:   token.NIL,
		kw.Class: token.CLASS,
	}
}

func (l *Lexer) PushFile(path string) error {
	result, err := FileExist(path)
	if result != true {
		return err
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	stat, _ := f.Stat()
	data := make([]byte, stat.Size())
	f.Read(data)
	cxt := &FileContext{
		Filename: path,
		File:     f,
		source:   []rune(string(data)),
		line:     1,
		column:   1,
	}
	if l.currentFile != nil {
		l.fileStack = append(l.fileStack, l.currentFile)
	}
	l.currentFile = cxt
	return nil
}

func (l *Lexer) PopFile() bool {
	if l.currentFile.Filename != "main.al" && l.currentFile.File != nil {
		l.currentFile.File.Close()
	}
	n := len(l.fileStack)
	if n == 0 {
		return false
	}
	l.currentFile = l.fileStack[n-1]
	l.fileStack = l.fileStack[:n-1]
	return true
}

func (l *Lexer) advance() rune {
	f := l.currentFile
	if f.readPos >= len(f.source) {
		f.ch = 0
		return 0
	}
	f.ch = f.source[f.readPos]
	f.position = f.readPos
	f.readPos++
	if f.ch == '\n' {
		f.line++
		f.column = 0
	} else {
		f.column++
	}
	return f.ch
}

func (l *Lexer) peek() rune {
	f := l.currentFile
	if f.readPos >= len(f.source) {
		return 0
	}
	return f.source[f.readPos]
}

func (l *Lexer) peekAt(n int) rune {
	f := l.currentFile
	pos := f.readPos + n
	if pos >= len(f.source) {
		return 0
	}
	return f.source[pos]
}

func (l *Lexer) newToken(tt token.TokenType, lexeme string, literal any) token.Token {
	f := l.currentFile
	return token.Token{
		Type:    tt,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    f.line,
		Column:  f.column,
		File:    f.Filename,
	}
}

// ScanTokens scans the entire source and returns all tokens together
// with the full list of diagnostics collected during scanning.
// Errors do not stop the scan — all diagnostics accumulate so the
// caller sees every problem in one pass.
func (l *Lexer) ScanTokens() ([]token.Token, *errs.DiagnosticList) {
	if l.currentFile.File != nil {
		defer l.currentFile.File.Close()
	}

	var tokens []token.Token

	for {
		ch := l.advance()
		if ch == 0 {
			break
		}

		f := l.currentFile

		switch ch {
		case ' ', '\t', '\r', '\n':
			continue

		case ',':
			tokens = append(tokens, l.newToken(token.COMMA, ",", nil))
		case ';':
			tokens = append(tokens, l.newToken(token.SEMICOLON, ";", nil))
		case ':':
			tokens = append(tokens, l.newToken(token.COLON, ":", nil))
		case '(':
			tokens = append(tokens, l.newToken(token.LEFT_PAREN, "(", nil))
		case ')':
			tokens = append(tokens, l.newToken(token.RIGHT_PAREN, ")", nil))
		case '[':
			tokens = append(tokens, l.newToken(token.LEFT_BRACKET, "[", nil))
		case ']':
			tokens = append(tokens, l.newToken(token.RIGHT_BRACKET, "]", nil))
		case '{':
			tokens = append(tokens, l.newToken(token.LEFT_BRACE, "{", nil))
		case '}':
			tokens = append(tokens, l.newToken(token.RIGHT_BRACE, "}", nil))
		case '!':
			tokens = append(tokens, l.newToken(token.BANG, "!", nil))

		case '.':
			tokens = append(tokens, l.matchToken('.', token.DOT, ".", token.DOT_DOT, ".."))
		case '-':
			tokens = append(tokens, l.matchToken('-', token.MINUS, "-", token.MINUS_MINUS, "--"))
		case '+':
			tokens = append(tokens, l.matchToken('+', token.PLUS, "+", token.PLUS_PLUS, "++"))
		case '*':
			tokens = append(tokens, l.matchToken('*', token.STAR, "*", token.STAR_STAR, "**"))

		case '/':
			if l.match('/') {
				// single-line comment: consume to end of line
				for l.peek() != '\n' && !l.isAtEnd() {
					l.advance()
				}
			} else if l.match('*') {
				l.skipBlockComment()
			} else {
				tokens = append(tokens, l.newToken(token.SLASH, "/", nil))
			}

		case '<':
			switch {
			case l.match('='):
				tokens = append(tokens, l.newToken(token.LESS_OR_EQUAL, "<=", nil))
			case l.match('>'):
				tokens = append(tokens, l.newToken(token.DIFF, "<>", nil))
			case l.match('-'):
				tokens = append(tokens, l.newToken(token.ASSIGN, "<-", nil))
			default:
				tokens = append(tokens, l.newToken(token.LESS, "<", nil))
			}

		case '>':
			tokens = append(tokens, l.matchToken('=', token.GREATER, ">", token.GREATER_OR_EQUAL, ">="))

		case '=':
			if l.match('=') {
				tokens = append(tokens, l.newToken(token.EQUAL_EQUAL, "==", nil))
			} else {
				l.diagnosticErrorf("unexpected character %q", ch)
			}

		case '"':
			if value, ok := l.scanString(); ok {
				tokens = append(tokens, l.newToken(token.STRING, value, value))
			}
			// on failure diagnosticErrorf already recorded it; continue scanning

		case '\'':
			if tok, ok := l.scanChar(); ok {
				tokens = append(tokens, tok)
			}

		default:
			if isDigit(ch) {
				tokens = append(tokens, l.scanNumber(ch))
			} else if isLetter(ch) {
				word := l.scanIdentifier()
				tt := token.LookupKeyword(word)
				tokens = append(tokens, l.newToken(tt, word, nil))
			} else {
				l.diagnosticErrorf("unexpected character %q", ch)
			}

			// suppress the unused variable warning for f on error-only paths
			_ = f
		}
	}

	return tokens, &l.Diagnostics
}

// skipBlockComment consumes everything up to and including the closing */.
// Any error is recorded in Diagnostics.
func (l *Lexer) skipBlockComment() {
	for {
		ch := l.advance()
		if ch == 0 {
			l.diagnosticErrorf("unterminated block comment")
			return
		}
		if ch == '*' && l.peek() == '/' {
			l.advance() // consume '/'
			return
		}
	}
}

func (l *Lexer) match(expected rune) bool {
	if l.peek() != expected {
		return false
	}
	l.advance()
	return true
}

func (l *Lexer) isAtEnd() bool {
	return l.currentFile.readPos >= len(l.currentFile.source)
}

func (l *Lexer) matchToken(
	next rune,
	singleType token.TokenType, singleLexeme string,
	doubleType token.TokenType, doubleLexeme string,
) token.Token {
	if l.match(next) {
		return l.newToken(doubleType, doubleLexeme, nil)
	}
	return l.newToken(singleType, singleLexeme, nil)
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isHexDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

func isOctalDigit(ch rune) bool {
	return ch >= '0' && ch <= '7'
}

func (l *Lexer) scanIdentifier() string {
	f := l.currentFile
	start := f.position
	for isLetter(l.peek()) || isDigit(l.peek()) {
		l.advance()
	}
	return string(f.source[start : f.position+1])
}

// scanString scans a double-quoted string literal.
// Returns (value, true) on success, ("", false) on error.
// Errors are recorded in Diagnostics so scanning continues.
func (l *Lexer) scanString() (string, bool) {
	var buf strings.Builder
	for {
		ch := l.advance()
		switch ch {
		case 0:
			l.diagnosticErrorf("unterminated string")
			return "", false
		case '\n':
			l.diagnosticErrorf("unterminated string (newline before closing quote)")
			return "", false
		case '"':
			return buf.String(), true
		case '\\':
			r, ok := l.scanEscape()
			if !ok {
				return "", false
			}
			buf.WriteRune(r)
		default:
			buf.WriteRune(ch)
		}
	}
}

// scanChar scans a single-quoted character literal.
// Returns (token, true) on success, (zero, false) on error.
func (l *Lexer) scanChar() (token.Token, bool) {
	ch := l.advance()
	if ch == 0 {
		l.diagnosticErrorf("unterminated char literal")
		return token.Token{}, false
	}

	var value rune
	if ch == '\\' {
		r, ok := l.scanEscape()
		if !ok {
			return token.Token{}, false
		}
		value = r
	} else {
		value = ch
	}

	if l.advance() != '\'' {
		l.diagnosticErrorf("char literal not terminated")
		return token.Token{}, false
	}

	lexeme := "'" + string(value) + "'"
	return l.newToken(token.CHARACTER, lexeme, value), true
}

// scanEscape processes a backslash escape sequence.
// The leading '\' has already been consumed.
// Returns (rune, true) on success, (0, false) on error.
func (l *Lexer) scanEscape() (rune, bool) {
	ch := l.advance()
	switch ch {
	case 'n':
		return '\n', true
	case 't':
		return '\t', true
	case 'r':
		return '\r', true
	case '\\':
		return '\\', true
	case '"':
		return '"', true
	case '\'':
		return '\'', true
	case '0':
		return 0, true
	case 'x':
		h1 := l.advance()
		h2 := l.advance()
		val, err := strconv.ParseInt(string([]rune{h1, h2}), 16, 32)
		if err != nil {
			l.diagnosticErrorf("invalid hex escape \\x%c%c", h1, h2)
			return 0, false
		}
		return rune(val), true
	case 'u':
		buf := make([]rune, 4)
		for i := range buf {
			buf[i] = l.advance()
		}
		val, err := strconv.ParseInt(string(buf), 16, 32)
		if err != nil {
			l.diagnosticErrorf("invalid unicode escape \\u%s", string(buf))
			return 0, false
		}
		return rune(val), true
	default:
		l.diagnosticErrorf("unknown escape sequence '\\%c'", ch)
		return 0, false
	}
}

// scanNumber scans an integer or real literal.
// ch is the first digit, already consumed.
func (l *Lexer) scanNumber(first rune) token.Token {
	f := l.currentFile
	start := f.position

	if first == '0' {
		switch l.peek() {
		case 'x', 'X':
			l.advance()
			for isHexDigit(l.peek()) {
				l.advance()
			}
			lexeme := string(f.source[start : f.position+1])
			val, _ := strconv.ParseInt(lexeme[2:], 16, 64)
			return l.newToken(token.INTEGER, lexeme, val)

		case 'b', 'B':
			l.advance()
			for l.peek() == '0' || l.peek() == '1' {
				l.advance()
			}
			lexeme := string(f.source[start : f.position+1])
			val, _ := strconv.ParseInt(lexeme[2:], 2, 64)
			return l.newToken(token.INTEGER, lexeme, val)

		case 'o', 'O':
			l.advance()
			for isOctalDigit(l.peek()) {
				l.advance()
			}
			lexeme := string(f.source[start : f.position+1])
			val, _ := strconv.ParseInt(lexeme[2:], 8, 64)
			return l.newToken(token.INTEGER, lexeme, val)
		}
	}

	for isDigit(l.peek()) {
		l.advance()
	}

	if l.peek() == l.decimalSep && isDigit(l.peekAt(1)) {
		l.advance()
		for isDigit(l.peek()) {
			l.advance()
		}
		if l.peek() == 'e' || l.peek() == 'E' {
			l.advance()
			if l.peek() == '-' {
				l.advance()
			}
			for isDigit(l.peek()) {
				l.advance()
			}
		}
		lexeme := string(f.source[start : f.position+1])
		val, _ := strconv.ParseFloat(strings.ReplaceAll(lexeme, ",", "."), 64)
		return l.newToken(token.DOUBLE, lexeme, val)
	}

	lexeme := string(f.source[start : f.position+1])
	val, _ := strconv.ParseInt(lexeme, 10, 64)
	return l.newToken(token.INTEGER, lexeme, val)
}

func FileExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("File '%s' does not exists.\n", filePath)
			return false, err
		}
		fmt.Printf("Error checking file: %v\n", err)
		return false, err
	}
	ext := filepath.Ext(filePath)
	if ext == "" {
		fmt.Printf("The file has no extension.")
		return false, err
	} else if ext != ".al" {
		fmt.Printf("File extension '%s' no recognize.\n", ext)
		return false, err
	}
	return true, nil
}
