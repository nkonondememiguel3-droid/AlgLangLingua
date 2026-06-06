package lexer

import (
	"alg/lexer/token"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type FileContext struct {
	filename string
	file     *os.File
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

func New() *Lexer {
	l := &Lexer{}

	return l
}

func (l *Lexer) PushFile(path string) error {

	result, err := fileExist(path)
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
		filename: path,
		file:     f,
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
	l.currentFile.file.Close()

	n := len(l.fileStack)
	if n == 0 {
		return false
	}

	l.currentFile = l.fileStack[n-1]
	l.fileStack = l.fileStack[:n-1]

	return true
}

func (l *Lexer) ScanTokens(filePath string) []token.Token {
	l.PushFile(filePath)
	defer l.PopFile()

	fmt.Println(string(l.currentFile.source))

	return nil
}

func fileExist(filePath string) (bool, error) {

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
