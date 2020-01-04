package main

import (
	"./bytecoderuntime"
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
	filename, workDir, out, compile, isBinary := configuration.GetConfig(os.Args[1:])

	if filename == "" || workDir == "" {
		fmt.Println("Invalid parameters!")
		return 1
	}

	if isBinary {
		f, err := os.Open(filename)

		if err != nil {
			panic(err)
		}

		stat, err := f.Stat()

		if err != nil {
			panic(err)
		}

		buf := make([]byte, stat.Size())
		_, err = f.Read(buf)

		if err != nil {
			panic(err)
		}

		return int(bytecoderuntime.Run(buf))
	}

	fmt.Println("Working directory: " + workDir)
	fmt.Println("File: " + filename)

	code := filereader.ReadFile(filename, workDir)
	commands := lexer.Lex(code)
	root := parser.Parse(commands)

	if compile {
		return int(compiler.Compile(root, out))
	}

	return int(runtime.Run(root))
}
