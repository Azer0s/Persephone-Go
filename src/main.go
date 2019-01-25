package main

import (
	"./configuration"
	"./filereader"
	"./lexer"
	"./parser"
	"fmt"
	"os"
)

func main(){
	filename, workdir := configuration.GetConfig(os.Args[1:])

	if filename == "" || workdir == "" {
		fmt.Println("Invalid parameters!")
		return
	}

	fmt.Println("Working directory: " + workdir)
	fmt.Println("File: " + filename)

	code := filereader.ReadFile(filename)
	commands := lexer.Lex(code)
	root := parser.Parse(commands)
}