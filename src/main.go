package main

import (
	"./configuration"
	"./filereader"
	"./lexer"
	"./parser"
	"./runtime"
	"fmt"
	"os"
)

func main() { os.Exit(mainReturnWithCode()) }

func mainReturnWithCode() int{
	filename, workdir := configuration.GetConfig(os.Args[1:])

	if filename == "" || workdir == "" {
		fmt.Println("Invalid parameters!")
		return 1
	}

	fmt.Println("Working directory: " + workdir)
	fmt.Println("File: " + filename)

	code := filereader.ReadFile(filename)
	commands := lexer.Lex(code)
	root := parser.Parse(commands)
	return int(runtime.Run(root))
}