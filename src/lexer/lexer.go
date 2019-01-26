package lexer

import (
	"../types"
	"fmt"
	"regexp"
)

var current = types.Token{Kind: "", Text:""}
var code string
var index = 0

func consume(){
	current.Text += letter()
	index++
}

func discard()  {
	index++
}

func letter() string{
	return string(code[index])
}

func Lex(lines []string) (tokens []types.Token){
	add := func(){tokens = append(tokens, current); current = types.Token{Kind:"", Text:""}}
	trimComments := regexp.MustCompile("(.*)#.*")

	for e := range lines {
		code += trimComments.ReplaceAllString(lines[e], "$1") + " "
	}

	isLetter := regexp.MustCompile("[a-zA-Z]")
	isNumber := regexp.MustCompile("[0-9]")

	for index = 0; index < len(code); index++ {
		if letter() == " "{
			continue
		}

		if isLetter.MatchString(letter()) {
			for isLetter.MatchString(letter()) || isNumber.MatchString(letter()) || letter() == "_"{
				consume()
			}

			current.Kind = types.Name

			if letter() == ":" {
				current.Kind = types.Label
				discard()
			}
		}else if isNumber.MatchString(letter()) {
			consume()

			if letter() == "x" {
				current.Kind = types.HexNumber
				consume()
			}else{
				current.Kind = types.Number
			}

			for isNumber.MatchString(letter()){
				consume()
			}

			if letter() == "." {
				if current.Kind == types.HexNumber {
					fmt.Println("Hexnumber can't have a decimal point!")
					return nil
				}
				consume()
				current.Kind = types.Float
				for isNumber.MatchString(letter()){
					consume()
				}
			}
		}else if letter() == "{"{
			current.Kind = types.Lbrace
			current.Text = "{"
		}else if letter() == "}"{
			current.Kind = types.Rbrace
			current.Text = "}"
		}else if letter() == "["{
			current.Kind = types.Pointer
			consume()

			for isLetter.MatchString(letter()) || isNumber.MatchString(letter()) || letter() == "_"{
				consume()
			}

			if letter() != "]" {
				fmt.Println("Expected a closing ], got: " + letter() + "!")
				return nil
			}else{
				consume()
			}
		}else if letter() == "\""{
			current.Kind = types.String
			consume()
			for letter() != "\"" {
				consume()
			}
			consume()
		}else{
			fmt.Println("Unknown token: " + letter() + "!")
			return nil
		}

		add()
	}

	return
}
