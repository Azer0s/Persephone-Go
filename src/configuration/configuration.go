package configuration

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetConfig(args []string)(file, workdir string){

	for k, v := range args {
		switch v {
		case "--workdir":
			workdir = args[k+1]

		case "--input":
		case "-i":
			file = args[k+1]
		}
	}

	if workdir == ""{
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			fmt.Println(err)
			return "",""
		}
		workdir = dir
	}

	if file == "" {
		fmt.Println("Input file cannot be empty!")
		return "",""
	}

	return
}