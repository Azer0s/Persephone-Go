package runtime

import (
	"../datatypes"
	"../types"
	"strconv"
)

type stack []datatypes.Data

func (s stack) Push(v datatypes.Data) stack {
	return append(s, v)
}

func (s stack) Pop() (stack, datatypes.Data) {
	l := len(s)
	return  s[:l-1], s[l-1]
}

func intOp(s stack, op datatypes.Op) stack{
	var a1 datatypes.Data
	var a2 datatypes.Data

	s,a1 = s.Pop()
	s,a2 = s.Pop()

	min := a1.Type

	if a2.Type > min {
		min = a2.Type
	}

	var left int64
	var right int64

	switch a1.Type {
	case datatypes.Int8:
		left = int64(a1.Value.(int8))
	case datatypes.Int16:
		left = int64(a1.Value.(int16))
	case datatypes.Int32:
		left = int64(a1.Value.(int32))
	case datatypes.Int64:
		left = int64(a1.Value.(int64))
	}

	switch a2.Type {
	case datatypes.Int8:
		right = int64(a2.Value.(int8))
	case datatypes.Int16:
		right = int64(a2.Value.(int16))
	case datatypes.Int32:
		right = int64(a2.Value.(int32))
	case datatypes.Int64:
		right = int64(a2.Value.(int64))
	}

	var result int64

	switch op {
	case datatypes.Add:
		result = left + right
	case datatypes.Sub:
		result = left - right
	case datatypes.Mul:
		result = left * right
	case datatypes.Div:
		result = left / right
	case datatypes.Mod:
		result = left % right
	}

	switch min {
	case datatypes.Int8:
		s = s.Push(datatypes.Data{Value:int8(result),Type:datatypes.Int8})
	case datatypes.Int16:
		s = s.Push(datatypes.Data{Value:int16(result),Type:datatypes.Int16})
	case datatypes.Int32:
		s = s.Push(datatypes.Data{Value:int32(result),Type:datatypes.Int32})
	case datatypes.Int64:
		s = s.Push(datatypes.Data{Value:int64(result),Type:datatypes.Int64})
	}

	return s
}

func floatOp(s stack, op datatypes.Op) stack{
	var a1 datatypes.Data
	var a2 datatypes.Data

	s,a1 = s.Pop()
	s,a2 = s.Pop()

	min := a1.Type

	if a2.Type > min {
		min = a2.Type
	}

	var left float64
	var right float64

	switch a1.Type {
	case datatypes.Float32:
		left = float64(a1.Value.(float32))
	case datatypes.Float64:
		left = float64(a1.Value.(float64))
	}

	switch a2.Type {
	case datatypes.Float32:
		right = float64(a2.Value.(float32))
	case datatypes.Float64:
		right = float64(a2.Value.(float64))
	}

	var result float64

	switch op {
	case datatypes.Add:
		result = left + right
	case datatypes.Sub:
		result = left - right
	case datatypes.Mul:
		result = left * right
	case datatypes.Div:
		result = left / right
	}

	switch min {
	case datatypes.Float32:
		s = s.Push(datatypes.Data{Value:float32(result),Type:datatypes.Float32})
	case datatypes.Float64:
		s = s.Push(datatypes.Data{Value:float64(result),Type:datatypes.Float64})
	}

	return s
}

func declareIntConst(command types.Command, c []datatypes.Data, d datatypes.DataType) []datatypes.Data{
	var num int64
	num,_ = strconv.ParseInt(command.Param.Text,0,64)

	switch d {
	case datatypes.Int8:
		c = append(c,datatypes.Data{Value:int8(num),Type:datatypes.Int8})
	case datatypes.Int16:
		c = append(c,datatypes.Data{Value:int16(num),Type:datatypes.Int16})
	case datatypes.Int32:
		c = append(c,datatypes.Data{Value:int32(num),Type:datatypes.Int32})
	case datatypes.Int64:
		c = append(c,datatypes.Data{Value:int64(num),Type:datatypes.Int64})
	}

	return c
}

func declareFloatConst(command types.Command, c []datatypes.Data, d datatypes.DataType) []datatypes.Data{
	var num float64
	num,_ = strconv.ParseFloat(command.Param.Text,64)

	switch d {
	case datatypes.Float32:
		c = append(c,datatypes.Data{Value:float32(num),Type:datatypes.Float32})
	case datatypes.Float64:
		c = append(c,datatypes.Data{Value:float64(num),Type:datatypes.Float64})
	}

	return c
}

func Run (root types.Root){
	s := make(stack,0)
	c := []datatypes.Data{}

	for e := range root.Commands {
		if root.Commands[e].Single {
			switch root.Commands[e].Command.Text {
			case "pop":
				s,_ = s.Pop()
			case "add":
				s = intOp(s, datatypes.Add)
			case "sub":
				s = intOp(s, datatypes.Sub)
			case "mul":
				s = intOp(s, datatypes.Mul)
			case "div":
				s = intOp(s, datatypes.Div)
			case "mod":
				s = intOp(s, datatypes.Mod)

			case "addf":
				s = floatOp(s, datatypes.Add)
			case "subf":
				s = floatOp(s, datatypes.Sub)
			case "mulf":
				s = floatOp(s, datatypes.Mul)
			case "divf":
				s = floatOp(s, datatypes.Div)
			}
		}else{
			switch root.Commands[e].Command.Text {
			case "dci8":
				c = declareIntConst(root.Commands[e], c, datatypes.Int8)
			case "dci16":
				c = declareIntConst(root.Commands[e], c, datatypes.Int16)
			case "dci":
			case "dci32":
				c = declareIntConst(root.Commands[e], c, datatypes.Int32)
			case "dci64":
				c = declareIntConst(root.Commands[e], c, datatypes.Int64)

			case "dcf":
			case "dcf32":
				c = declareFloatConst(root.Commands[e], c, datatypes.Float32)
			case "dcd":
			case "dcf64":
				c = declareFloatConst(root.Commands[e], c, datatypes.Float64)
			}
		}
	}
}
