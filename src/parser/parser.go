package parser

import (
	"../types"
	"fmt"
)

var noArgCommands = []string{"add", "sub", "mul", "div", "mod", "andi", "ori", "xori", "noti", "shl", "shr", "addf", "subf", "mulf", "divf", "pop", "ge", "le", "gt", "lt", "gef", "lef", "gtf", "ltf", "inc", "dec", "cbase", "and", "or", "xor", "not", "ret", "conc"}

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

func parseStatement() (command types.Command) {
	if tokens[index].Kind == types.Name {

		hasNoArgsCommands := contains(noArgCommands, tokens[index].Text)

		if hasNoArgsCommands {
			command.Single = true
			command.Command = tokens[index]
			if hasNoArgsCommands {
				return
			}
		}

		command.Single = false
		command.Command = tokens[index]
		index++

		if tokens[index].Kind == types.Name || tokens[index].Kind == types.Number || tokens[index].Kind == types.HexNumber || tokens[index].Kind == types.Float || tokens[index].Kind == types.String || tokens[index].Kind == types.Pointer {
			command.Param = tokens[index]
			return
		}

		fmt.Println("Expected name, number, hexnumber, flaot, string or pointer, got: " + tokens[index].Kind)
		return types.Command{}
	} else {
		fmt.Println("Expected name, got: " + tokens[index].Kind)
		return types.Command{}
	}
}

//Parse ...Parses a list of tokens into an AST
func Parse(tks []types.Token) (root types.Root) {
	tokens = tks
	index = 0
	root.Labels = make(map[string]int)

	for index = 0; index < len(tokens); index++ {
		if tokens[index].Kind == types.Label {
			root.Labels[tokens[index].Text] = len(root.Commands)
		} else {
			root.Commands = append(root.Commands, parseStatement())
		}
	}

	return
}
