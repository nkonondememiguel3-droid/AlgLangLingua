package main

import (
	"alg/lexer"
	"fmt"
	"log"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("Usage: alg <scrip>.al")
	}
	filePath := os.Args[1]

	fileContext := lexer.FileContext{
		Filename: filePath,
	}
	l := lexer.New(fileContext)
	tokens := l.ScanTokens()

	for _, tk := range tokens {
		fmt.Println(tk)
	}

}
