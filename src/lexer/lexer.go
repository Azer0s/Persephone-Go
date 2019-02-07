package lexer

import (
	"../types"
	"fmt"
	"regexp"
	"strconv"
)

var current = types.Token{Kind: "", Text: ""}
var code []rune
var index = 0

func consume() {
	current.Text += string(letter())
	index++
}

func discard() {
	index++
}

func letter() rune {
	return code[index]
}

//Lex ...Lexes Persephone code and returns an array of tokens
func Lex(lines []string) (tokens []types.Token) {
	add := func() { tokens = append(tokens, current); current = types.Token{Kind: "", Text: ""} }
	trimComments := regexp.MustCompile("(.*)#.*")

	for e := range lines {
		code = append(code, []rune(trimComments.ReplaceAllString(lines[e], "$1")+" ")...)
	}

	isLetter := regexp.MustCompile("[a-zA-Z]")
	isNumber := regexp.MustCompile("[0-9]")

	for index = 0; index < len(code); index++ {
		if letter() == rune(' ') || letter() == rune('\t') {
			continue
		}

		if isLetter.MatchString(string(letter())) {
			for isLetter.MatchString(string(letter())) || isNumber.MatchString(string(letter())) || letter() == rune('_') {
				consume()
			}

			current.Kind = types.Name

			if letter() == rune(':') {
				current.Kind = types.Label
				discard()
			}

			if current.Text == "true" || current.Text == "false"{
				current.Kind = types.Bit
			}
		} else if isNumber.MatchString(string(letter())) || (letter() == rune('-') && isNumber.MatchString(string(code[index + 1]))){
			if letter() == rune('-') {
				consume()
			}

			consume()

			if letter() == rune('x') {
				current.Kind = types.HexNumber
				consume()
			} else {
				current.Kind = types.Number
			}

			for isNumber.MatchString(string(letter())) {
				consume()
			}

			if letter() == rune('.') {
				if current.Kind == types.HexNumber {
					fmt.Println("Hexnumber can't have a decimal point!")
					return nil
				}
				consume()
				current.Kind = types.Float
				for isNumber.MatchString(string(letter())) {
					consume()
				}

				if letter() == rune('[') { //size is explicitly stated
					discard()
					if letter() == rune('3') {
						discard()
						if letter() == rune('2') {
							discard()
							current.Size = "32"
						} else {
							panic("Expected either [32] or [64]!")
						}
					} else if letter() == rune('6') {
						discard()
						if letter() == rune('4') {
							discard()
							current.Size = "64"
						} else {
							panic("Expected either [32] or [64]!")
						}
					} else {
						panic("Expected either [32] or [64]!")
					}

					if letter() != rune(']') {
						panic("Expected either [32] or [64]!")
					}
					discard()
				}
			}
		} else if letter() == rune('[') {
			current.Kind = types.Pointer
			consume()

			for isLetter.MatchString(string(letter())) || isNumber.MatchString(string(letter())) || letter() == rune('_') {
				consume()
			}

			if letter() != rune(']') {
				fmt.Println("Expected a closing ], got: " + string(letter()) + "!")
				return nil
			}

			consume()
		} else if letter() == rune('"') {
			current.Kind = types.String
			consume()
			for letter() != rune('"') {
				consume()
			}
			consume()
		} else if letter() == rune('\'') {
			discard()
			consume()
			discard()

			current.Text = strconv.Itoa(int(int8(current.Text[0])))
			current.Kind = types.Number
		} else {
			fmt.Println("Unknown token: " + string(letter()) + "!")
			return nil
		}

		add()
	}

	return
}
