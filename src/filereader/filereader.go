package filereader

import (
	"bufio"
	"fmt"
	"os"
)

func ReadFile(filename string) (code []string){
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	code = []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		code = append(code,scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}
	return
}
