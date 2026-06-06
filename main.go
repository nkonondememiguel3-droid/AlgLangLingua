package main

import (
	"alg/lexer"
	"log"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatalf("Usage: alg <scrip>.al")
		os.Exit(64)
	}
	filePath := os.Args[1]

	l := lexer.New()
	_ = l.ScanTokens(filePath)

	// for _, filePath := range filesPath {
	// 	// check if the files exists.
	// 	fileInfo, err := os.Stat(filePath)
	// 	if err != nil {
	// 		if errors.Is(err, os.ErrNotExist) {
	// 			log.Fatalf("File '%s' does not exists.\n", filePath)
	// 		} else {
	// 			log.Fatalf("Error checking file: %v\n", err)
	// 		}
	// 		os.Exit(64)
	// 	}

	// 	// check the extension of the file to see if it matches the appropriated extension.
	// 	ext := filepath.Ext(filePath)
	// 	if ext == "" {
	// 		log.Fatalf("The file has no extension.")
	// 	} else if ext != ".al" {
	// 		log.Fatalf("File extension '%s' no recognize.\n", ext)
	// 	}

	// 	readFIle(fileInfo.Name())

	// }

}

func readFIle(name string) {}
