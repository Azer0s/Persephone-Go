package parser

import (
	"../types"
	"fmt"
)

var noArgCommands = []string{"shl","shr","add","sub","mul","div","mod","pop"}
var functionCommands = []string{"put","ret"}

var tokens []types.Token
var index int

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func parseStatement() (command types.Command){
	if tokens[index].Kind == types.Name {
		if contains(noArgCommands,tokens[index].Text) {
			command.Single = true
			command.Command = tokens[index]
			return
		}

		command.Single = false
		command.Command = tokens[index]
		index++

		if tokens[index].Kind == types.Name || tokens[index].Kind == types.Number || tokens[index].Kind == types.HexNumber || tokens[index].Kind == types.String || tokens[index].Kind == types.Pointer {
			command.Param = tokens[index]
			return
		}else{
			fmt.Println("Expected name, number, hexnumber, string or pointer, got: " + tokens[index].Kind)
			return types.Command{}
		}
	}else {
		fmt.Println("Expected name, got: " + tokens[index].Kind)
		return types.Command{}
	}
}

func Parse (tks []types.Token) (root types.Root){
	tokens = tks
	index = 0

	for index = 0; index < len(tokens); index++{
		if tokens[index].Kind == types.Name && tokens[index+1].Kind == types.Lbrace {
			fn := types.Function{}
			fn.Name = tokens[index]
			index++
			index++

			for tokens[index].Kind != types.Rbrace {
				fn.Commands = append(fn.Commands, parseStatement())
				index++
			}

			root.Functions = append(root.Functions, fn)
		}else{
			root.Commands = append(root.Commands, parseStatement())
		}
	}

	return
}
