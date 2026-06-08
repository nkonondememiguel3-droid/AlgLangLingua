package lexer

import (
	"alg/config"
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
	fileStack   []*FileContext // context stack.
	currentFile *FileContext   // active compiling context.
	decimalSep  rune           // '.' (default) or ','
}

// helper function mimicing the ternary operator.
func If[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func New(fileContext FileContext, cfg ...*config.Config) *Lexer {

	l := &Lexer{
		currentFile: &fileContext,
		decimalSep:  '.',
	}
	if len(cfg) > 0 && cfg[0] != nil {
		// Build the keyword map from the config before scanning anything.
		token.RegisterKeywords(buildKeywordMap(cfg[0]))

		if cfg[0].Meta.DecimalSep == "," {
			l.decimalSep = ','
		}
	}

	if fileContext.Filename != "" {
		if err := l.PushFile(fileContext.Filename); err != nil {
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

// buildKeywordMap converts a config.Keywords struct into the
// surface-word -> TokenType map that the lexer uses for lookup.
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

	// we append the current file only if it is not the first file being scanned.
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

	// we get the current file.
	l.currentFile = l.fileStack[n-1]
	// we shrink the fileStack by one.
	l.fileStack = l.fileStack[:n-1]

	return true
}

// advance reads the next rune and moves the read position forward.
// Returns 0 when the source is exhausted.
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

// peek returns the next rune without consuming it.
func (l *Lexer) peek() rune {
	f := l.currentFile
	if f.readPos >= len(f.source) {
		return 0
	}
	return f.source[f.readPos]
}

// peekAt returns the rune n positions ahead without consuming anything.
func (l *Lexer) peekAt(n int) rune {
	f := l.currentFile
	pos := f.readPos + n
	if pos >= len(f.source) {
		return 0
	}
	return f.source[pos]
}

// newToken constructs a Token from the current file context.
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

func (l *Lexer) ScanTokens() []token.Token {
	// Close real files on exit; tests use nil File, so guard first.
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
		case ' ', '\t', '\r':
			// skip whitespace (newlines already handled in advance)
			continue
		case '\n':
			continue

			// one character token.
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

			// one or two character tokens.
		case '.':
			tokens = append(tokens, l.matchToken('.', token.DOT, ".", token.DOT_DOT, ".."))
		case '-':
			tokens = append(tokens, l.matchToken('-', token.MINUS, "-", token.MINUS_MINUS, "--"))
		case '+':
			tokens = append(tokens, l.matchToken('+', token.PLUS, "+", token.PLUS_PLUS, "++"))
		case '/':
			if l.match('/') {
				for l.peek() != '\n' && !l.isAtEnd() {
					l.advance()
				}
			} else if l.match('*') {
				if err := l.skipBlockComment(); err != nil {
					fmt.Fprintf(os.Stderr, "%v\n", err)
				}
			} else {
				tokens = append(tokens, l.newToken(token.SLASH, "/", nil))
			}
		case '*':
			tokens = append(tokens, l.matchToken('*', token.STAR, "*", token.STAR_STAR, "**"))
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
				fmt.Fprintf(os.Stderr, "%s:%d:%d: unexpected character %q\n",
					f.Filename, f.line, f.column, ch)
			}
		case '"':
			value, err := l.scanString()

			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}

			tokens = append(tokens,
				l.newToken(token.STRING, value, value))

		case '\'':
			tok, err := l.scanChar()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			tokens = append(tokens, tok)

		default:
			if isDigit(ch) {
				tokens = append(tokens, l.scanNumber(ch))
			} else if isLetter(ch) {
				word := l.scanIdentifier()
				tt := token.LookupKeyword(word)
				tokens = append(tokens, l.newToken(tt, word, nil))
			} else {
				fmt.Fprintf(os.Stderr, "%s:%d:%d: unexpected character %q\n",
					f.Filename, f.line, f.column, ch)
			}
		}
	}

	return tokens
}

// skip comments
func (l *Lexer) skipBlockComment() error {
	for {
		ch := l.advance()
		if ch == 0 {
			return fmt.Errorf("%s:%d:%d: unterminated block comment",
				l.currentFile.Filename, l.currentFile.line, l.currentFile.column)
		}
		if ch == '*' && l.peek() == '/' {
			l.advance() // consume '/'
			return nil
		}
	}
}

// match consumes the next rune only if it equals expected.
func (l *Lexer) match(expected rune) bool {
	if l.peek() != expected {
		return false
	}
	l.advance() // consume it
	return true
}

// isAtEnd reports whether all source runes have been consumed.
func (l *Lexer) isAtEnd() bool {
	return l.currentFile.readPos >= len(l.currentFile.source)
}

// matchToken tries to match the next rune against next.
// On success it emits the two-char token (doubleType/doubleLexeme).
// On failure it emits the single-char token (singleType/singleLexeme).
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

// scanIdentifier reads a full identifier or keyword starting from the
// current character (already consumed by advance into ch).
func (l *Lexer) scanIdentifier() string {
	f := l.currentFile
	// start is the position of the first character (already consumed).
	start := f.position
	for isLetter(l.peek()) || isDigit(l.peek()) {
		l.advance()
	}
	return string(f.source[start : f.position+1])
}

func (l *Lexer) scanString() (string, error) {
	var buf strings.Builder
	for {
		ch := l.advance()
		switch ch {
		case 0:
			return "", fmt.Errorf("%s:%d:%d: unterminated string",
				l.currentFile.Filename, l.currentFile.line, l.currentFile.column)
		case '"':
			return buf.String(), nil
		case '\\':
			r, err := l.scanEscape()
			if err != nil {
				return "", err
			}
			buf.WriteRune(r)
		case '\n':
			// a raw newline inside a string is an error - strings must be
			// closed on the same line (use \n escape for a newline value)
			return "",
				fmt.Errorf("%s:%d:%d: unterminated string (newline before closing quote)", l.currentFile.Filename, l.currentFile.line, l.currentFile.column)
		default:
			buf.WriteRune(ch)
		}
	}
}

// scanChar scans a character literal after the opening ' has been consumed.
func (l *Lexer) scanChar() (token.Token, error) {
	ch := l.advance()
	if ch == 0 {
		return token.Token{}, fmt.Errorf("%s:%d:%d: unterminated char literal",
			l.currentFile.Filename, l.currentFile.line, l.currentFile.column)
	}

	var value rune
	if ch == '\\' {
		escaped, err := l.scanEscape()
		if err != nil {
			return token.Token{}, err
		}
		value = escaped
	} else {
		value = ch
	}

	if l.advance() != '\'' {
		return token.Token{}, fmt.Errorf("%s:%d:%d: char literal not terminated",
			l.currentFile.Filename, l.currentFile.line, l.currentFile.column)
	}

	lexeme := "'" + string(value) + "'"
	return l.newToken(token.CHARACTER, lexeme, value), nil
}

// scanEscape processes a backslash escape sequence.
// The leading '\' has already been consumed.
func (l *Lexer) scanEscape() (rune, error) {
	ch := l.advance()
	switch ch {
	case 'n':
		return '\n', nil
	case 't':
		return '\t', nil
	case 'r':
		return '\r', nil
	case '\\':
		return '\\', nil
	case '"':
		return '"', nil
	case '\'':
		return '\'', nil
	case '0':
		return 0, nil
	case 'x':
		// \xHH
		h1 := l.advance()
		h2 := l.advance()
		val, err := strconv.ParseInt(string([]rune{h1, h2}), 16, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid hex escape \\x%c%c", h1, h2)
		}
		return rune(val), nil
	case 'u':
		// \uHHHH
		buf := make([]rune, 4)
		for i := range buf {
			buf[i] = l.advance()
		}
		val, err := strconv.ParseInt(string(buf), 16, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid unicode escape \\u%s", string(buf))
		}
		return rune(val), nil
	default:
		return 0, fmt.Errorf("%s:%d:%d: unknown escape sequence '\\%c'",
			l.currentFile.Filename, l.currentFile.line, l.currentFile.column, ch)
	}
}

func FileExist(filePath string) (bool, error) {

	// check if the files exists.
	_, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("File '%s' does not exists.\n", filePath)
			return false, err
		} else {
			fmt.Printf("Error checking file: %v\n", err)
			return false, err
		}
	}

	// check the extension of the file to see if it matches the appropriated extension.
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

// scanNumber scans an integer or real literal.
// ch is the first digit, already consumed.
func (l *Lexer) scanNumber(first rune) token.Token {
	f := l.currentFile
	start := f.position

	// prefixed integer bases
	if first == '0' {
		switch l.peek() {
		case 'x', 'X':
			l.advance() // consume 'x'
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

	// decimal integer or real
	for isDigit(l.peek()) {
		l.advance()
	}

	// check for decimal separator followed by more digits → real
	if l.peek() == l.decimalSep && isDigit(l.peekAt(1)) {
		l.advance() // consume separator
		for isDigit(l.peek()) {
			l.advance()
		}
		// optional exponent: e / E followed by optional '-' then digits
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
		// normalise: replace ',' with '.' so strconv.ParseFloat works
		val, _ := strconv.ParseFloat(strings.ReplaceAll(lexeme, ",", "."), 64)
		return l.newToken(token.DOUBLE, lexeme, val)
	}

	lexeme := string(f.source[start : f.position+1])
	val, _ := strconv.ParseInt(lexeme, 10, 64)
	return l.newToken(token.INTEGER, lexeme, val)
}

func isHexDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

func isOctalDigit(ch rune) bool {
	return ch >= '0' && ch <= '7'
}
