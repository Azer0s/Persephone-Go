package runtime

import (
	"../types"
)

type stack []interface{}

func (s stack) Push(v interface{}) stack {
	return append(s, v)
}

func (s stack) Pop() (stack, interface{}) {
	l := len(s)
	return  s[:l-1], s[l-1]
}

//TODO: Add struct that stores data and datatype
func add(s stack) stack{
	var a1 interface{}
	var a2 interface{}

	s,a1 = s.Pop()
	s,a2 = s.Pop()

	var inta int64
	var intb int64

	switch a1.(type) {
	case int8:
		inta = int64(a1.(int8))
	case int16:
		inta = int64(a1.(int16))
	case int32:
		inta = int64(a1.(int32))
	case int64:
		inta = int64(a1.(int64))
	}

	switch a2.(type) {
	case int8:
		inta = int64(a2.(int8))
	case int16:
		inta = int64(a2.(int16))
	case int32:
		inta = int64(a2.(int32))
	case int64:
		inta = int64(a2.(int64))
	}

	s.Push(inta + intb)
	return s
}


func Run (root types.Root){
	s := make(stack,0)

	for e := range root.Commands {
		if root.Commands[e].Single {
			switch root.Commands[e].Command.Text {
				case "pop":
					s.Pop()
				case "add":
					s = add(s)
			}
		}else{
			switch root.Commands[e].Command.Text {

			}
		}
	}
}
