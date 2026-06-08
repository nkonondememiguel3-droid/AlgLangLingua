package main

import (
	"alg/config"
	"alg/lexer"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("usage: alg <lang>.toml <script>.al")
	}

	cfg, err := config.Load(os.Args[1])
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	fileContext := lexer.FileContext{Filename: os.Args[2]}
	l := lexer.New(fileContext, cfg)
	tokens := l.ScanTokens()

	for _, tk := range tokens {
		fmt.Println(tk)
	}
}
