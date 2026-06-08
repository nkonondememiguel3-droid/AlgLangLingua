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
	tokens, diags := l.ScanTokens()
	if diags.HasErrors() {
		fmt.Fprint(os.Stderr, diags.Format())
		os.Exit(1)
	}
	// warnings don't stop compilation
	if len(diags.Warnings()) > 0 {
		fmt.Fprint(os.Stderr, diags.Format())
	}

	for _, tk := range tokens {
		fmt.Println(tk)
	}
}
