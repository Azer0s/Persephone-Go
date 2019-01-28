package runtime

import (
	"../datatypes"
	"../types"
	"fmt"
	"strconv"
	"strings"
)

type stack []datatypes.Data

func (s stack) Push(v datatypes.Data) stack {
	return append(s, v)
}

func (s stack) Pop() (stack, datatypes.Data) {
	l := len(s)
	return  s[:l-1], s[l-1]
}

/*
Arithmetic operations
 */
func getInt64(data datatypes.Data) int64{
	switch data.Type {
	case datatypes.Ptr:
		return int64(data.Value.(int32))
	case datatypes.Int8:
		return int64(data.Value.(int8))
	case datatypes.Int16:
		return int64(data.Value.(int16))
	case datatypes.Int32:
		return int64(data.Value.(int32))
	case datatypes.Int64:
		return int64(data.Value.(int64))
	default:
		return 0
	}
}

func intOp(s stack, op datatypes.Op) stack{
	var a1 datatypes.Data
	var a2 datatypes.Data

	s,a1 = s.Pop()
	s,a2 = s.Pop()

	if !(a1.Type >= datatypes.Ptr && a1.Type <= datatypes.Int64 && a2.Type >= datatypes.Ptr && a2.Type <= datatypes.Int64) {
		fmt.Println("Only int or ptr allowed in int operations!")
		return nil
	}

	min := a1.Type

	if a2.Type > min {
		min = a2.Type
	}

	if min == datatypes.Ptr {
		min = datatypes.Int32
	}

	left := getInt64(a1)
	right := getInt64(a2)

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
	case datatypes.Shr, datatypes.Shl:
		leftu := uint64(left)
		rightu := uint64(right)

		if op == datatypes.Shl {
			result = int64(leftu << rightu)
		}else {
			result = int64(leftu >> rightu)
		}
	case datatypes.And:
		result = left & right
	case datatypes.Or:
		result = left | right
	case datatypes.Xor:
		result = left ^ right
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

func getFloat64(data datatypes.Data) float64{
	switch data.Type {
	case datatypes.Float32:
		return float64(data.Value.(float32))
	case datatypes.Float64:
		return float64(data.Value.(float64))
	default:
		return 0
	}
}

func floatOp(s stack, op datatypes.Op) stack{
	var a1 datatypes.Data
	var a2 datatypes.Data

	s,a1 = s.Pop()
	s,a2 = s.Pop()

	if !(a1.Type >= datatypes.Float32 && a1.Type <= datatypes.Float64 && a2.Type >= datatypes.Float32 && a2.Type <= datatypes.Float64) {
		fmt.Println("Only float allowed in float operations!")
		return nil
	}

	min := a1.Type

	if a2.Type > min {
		min = a2.Type
	}

	left := getFloat64(a1)
	right := getFloat64(a2)

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

/*
Constant declarations
 */

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

func declareStringConstant(command types.Command, c []datatypes.Data) []datatypes.Data{
	val := strings.Trim(command.Param.Text,"\"")
	var stringtype datatypes.DataType

	if command.Command.Text == "dcsa"  {
		stringtype = datatypes.String_ASCII
	}else {
		stringtype = datatypes.String_Unicode
	}

	c = append(c, datatypes.Data{Value:val,Type:stringtype})
	return c
}

func declareBitConstant(command types.Command, c []datatypes.Data) []datatypes.Data{
	val := false

	switch command.Param.Text {
	case "0":
		val = false
	case "1":
		val = true
	default:
		fmt.Println("Error, expected a bit value, received: " + command.Param.Text + "!")
		return nil
	}

	c = append(c, datatypes.Data{Value:val, Type:datatypes.Bit})
	return c
}

/*
Variable declaration
 */
func declareVar(command types.Command, d datatypes.DataType, v map[string]datatypes.Data) map[string]datatypes.Data{

	switch d {
	case datatypes.Bit:
		v[command.Param.Text] = datatypes.Data{Value:false,Type:datatypes.Bit}
	case datatypes.Ptr:
		v[command.Param.Text] = datatypes.Data{Value:int32(0x0),Type:datatypes.Ptr}
	case datatypes.Int8:
		v[command.Param.Text] = datatypes.Data{Value:int8(0),Type:datatypes.Int8}
	case datatypes.Int16:
		v[command.Param.Text] = datatypes.Data{Value:int16(0),Type:datatypes.Int16}
	case datatypes.Int32:
		v[command.Param.Text] = datatypes.Data{Value:int32(0),Type:datatypes.Int32}
	case datatypes.Int64:
		v[command.Param.Text] = datatypes.Data{Value:int64(0),Type:datatypes.Int64}
	case datatypes.Float32:
		v[command.Param.Text] = datatypes.Data{Value:float32(0),Type:datatypes.Float32}
	case datatypes.Float64:
		v[command.Param.Text] = datatypes.Data{Value:float64(0),Type:datatypes.Float64}
	case datatypes.String_ASCII:
		v[command.Param.Text] = datatypes.Data{Value:"",Type:datatypes.String_ASCII}
	case datatypes.String_Unicode:
		v[command.Param.Text] = datatypes.Data{Value:"",Type:datatypes.String_Unicode}
	}

	return v
}

func Run (root types.Root){
	s := make(stack,0)
	var c []datatypes.Data
	v := make(map[string]datatypes.Data)

	for e := range root.Commands {
		if root.Commands[e].Single {
			switch root.Commands[e].Command.Text {
			case "pop":
				s,_ = s.Pop()

			/*
			Arithmetic int operations
			 */
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
			case "andi":
				s = intOp(s, datatypes.And)
			case "ori":
				s = intOp(s, datatypes.Or)
			case "xori":
				s = intOp(s, datatypes.Xor)
			case "noti":
				//TODO
			case "shl":
				s = intOp(s, datatypes.Shl)
			case "shr":
				s = intOp(s, datatypes.Shr)

			/*
			Arithmetic float operations
			 */
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
			/*
			Declare int constant
			 */
			case "dci8":
				c = declareIntConst(root.Commands[e], c, datatypes.Int8)
			case "dci16":
				c = declareIntConst(root.Commands[e], c, datatypes.Int16)
			case "dci32","dci":
				c = declareIntConst(root.Commands[e], c, datatypes.Int32)
			case "dci64":
				c = declareIntConst(root.Commands[e], c, datatypes.Int64)

			/*
			Declare float constant
			 */
			case "dcf32","dcf":
				c = declareFloatConst(root.Commands[e], c, datatypes.Float32)
			case "dcf64","dcd":
				c = declareFloatConst(root.Commands[e], c, datatypes.Float64)

			/*
			Declare string constant
			This implementation of Persephone doesn't differentiate between ASCII and Unicode
			 */
			case "dcsa","dcsu":
				c = declareStringConstant(root.Commands[e], c)

			/*
			Declare bit constant
			 */
			case "dcb":
				c = declareBitConstant(root.Commands[e], c)

			/*
			Variable creation
			 */
			case "v_int8":
				v = declareVar(root.Commands[e], datatypes.Int8, v)
			case "v_int16":
				v = declareVar(root.Commands[e], datatypes.Int16, v)
			case "v_int32", "v_int":
				v = declareVar(root.Commands[e], datatypes.Int32, v)
			case "v_int64":
				v = declareVar(root.Commands[e], datatypes.Int64, v)
			case "v_float32", "v_float":
				v = declareVar(root.Commands[e], datatypes.Float32, v)
			case "v_float64", "v_double":
				v = declareVar(root.Commands[e], datatypes.Float64, v)
			case "v_stringa":
				v = declareVar(root.Commands[e], datatypes.String_ASCII, v)
			case "v_stringu":
				v = declareVar(root.Commands[e], datatypes.String_Unicode, v)
			case "v_bit":
				v = declareVar(root.Commands[e], datatypes.Bit, v)
			case "v_ptr":
				v = declareVar(root.Commands[e], datatypes.Ptr, v)

			/*
			Load variable onto stack
			 */
			case "ldi8v":
			case "ldi16v":
			case "ldi32v", "ldiv":
			case "ldi64v":
			case "ldf32v", "ldfv":
			case "ldf64v", "lddv":
			case "ldsav", "ldsuv":
			case "ldbv":
			case "ldptrv":

			}
		}
	}
}
