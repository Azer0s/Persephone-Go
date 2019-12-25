package filereader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

//ReadFile reads file and does preprocessing
func ReadFile(filename, workdir string) (code []string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	code = []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		code = append(code, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	region := regexp.MustCompile("^ *%(end)?region *.*$")
	include := regexp.MustCompile("^ *%include +(.+) *")

	for e := 0; e < len(code); e++ {
		if region.MatchString(code[e]) {
			code = append(code[:e], code[e+1:]...)
		}

		if include.MatchString(code[e]) {
			tempLines := code[e+1:]
			newFileLines := ReadFile(filepath.Join(workdir, include.ReplaceAllString(code[e], "$1")), workdir)
			code = append(code[:e], append(newFileLines, tempLines...)...)
		}
	}

	return
}
