package parser

import (
	"../types"
	"fmt"
)

var noArgCommands = []string{"add", "sub", "mul", "div", "mod", "andi", "ori", "xori", "noti", "shl", "shr", "addf", "subf", "mulf", "divf", "pop"}
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

func parseStatement(function bool) (command types.Command){
	if tokens[index].Kind == types.Name {

		hasNoArgsCommands := contains(noArgCommands,tokens[index].Text)
		hasFunctionCommands := contains(functionCommands,tokens[index].Text)

		if hasNoArgsCommands || hasFunctionCommands {
			command.Single = true
			command.Command = tokens[index]
			if hasNoArgsCommands {
				return
			}

			if hasFunctionCommands {
				if function {
					return
				}else {
					fmt.Println("put or ret can only be used in functions!")
					return types.Command{}
				}
			}
		}

		command.Single = false
		command.Command = tokens[index]
		index++

		if tokens[index].Kind == types.Name || tokens[index].Kind == types.Number || tokens[index].Kind == types.HexNumber || tokens[index].Kind == types.Float || tokens[index].Kind == types.String || tokens[index].Kind == types.Pointer {
			command.Param = tokens[index]
			return
		}else{
			fmt.Println("Expected name, number, hexnumber, flaot, string or pointer, got: " + tokens[index].Kind)
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
	root.Labels = make(map[string]int)

	for index = 0; index < len(tokens); index++{
		if tokens[index].Kind == types.Name && tokens[index+1].Kind == types.Lbrace {
			fn := types.Function{}
			fn.Name = tokens[index]
			index++
			index++

			for tokens[index].Kind != types.Rbrace {
				fn.Commands = append(fn.Commands, parseStatement(true))
				index++
			}

			root.Functions = append(root.Functions, fn)
		}else if tokens[index].Kind == types.Label{
			root.Labels[tokens[index].Text] = len(root.Commands)
		}else{
			root.Commands = append(root.Commands, parseStatement(false))
		}
	}

	return
}
