package configuration

import (
	"path/filepath"
)

func GetConfig(args []string) (file, workdir string, compile bool) {

	for k, v := range args {
		switch v {
		case "--workdir":
			workdir = args[k+1]

		case "-i","--input":
			file = args[k+1]

		case "--compile", "-c":
			compile = true
		}
	}

	if workdir == "" {
		dir, err := filepath.Abs(filepath.Dir(file))
		if err != nil {
			panic(err)
		}
		workdir = dir
	}

	if file == "" {
		panic("Input file cannot be empty!")
	}

	return
}
