package lexer

import (
	"alg/lexer/token"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
}

// helper function mimicing the ternary operator.
func If[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

func New(fileContext FileContext) *Lexer {
	l := &Lexer{
		currentFile: &fileContext,
	}
	// Only try to open a file when a real path was provided.
	if fileContext.Filename != "" {
		if err := l.PushFile(fileContext.Filename); err != nil {
			// caller will see an empty scan; you could also return an error
			fmt.Fprintf(os.Stderr, "lexer: %v\n", err)
		}
	} else {
		// In-memory source (tests): initialise line/column if not set.
		if l.currentFile.line == 0 {
			l.currentFile.line = 1
			l.currentFile.column = 0
		}
	}
	return l
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
		f.column = 1
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

		case ',':
			tokens = append(tokens, l.newToken(token.COMMA, ",", nil))
		case ';':
			tokens = append(tokens, l.newToken(token.SEMICOLON, ";", nil))
		case ':':
			tokens = append(tokens, l.newToken(token.COLON, ":", nil))
		case '.':
			tokens = append(tokens, l.newToken(token.DOT, ".", nil))
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

		default:
			fmt.Fprintf(os.Stderr, "%s:%d:%d: unexpected character %q\n",
				f.Filename, f.line, f.column, ch)
		}
	}

	return tokens
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
