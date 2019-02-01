package main

import (
	"./compiler"
	"./configuration"
	"./filereader"
	"./lexer"
	"./parser"
	"./runtime"
	"fmt"
	"os"
)

func main() { os.Exit(mainReturnWithCode()) }

func mainReturnWithCode() int {
	filename, workdir, out, compile := configuration.GetConfig(os.Args[1:])

	if filename == "" || workdir == "" {
		fmt.Println("Invalid parameters!")
		return 1
	}

	fmt.Println("Working directory: " + workdir)
	fmt.Println("File: " + filename)

	code := filereader.ReadFile(filename, workdir)
	commands := lexer.Lex(code)
	root := parser.Parse(commands)

	if compile {
		return int(compiler.Compile(root, out))
	}

	return int(runtime.Run(root))
}
