package runtime

import (
	"../datatypes"
	"../types"
	"fmt"
	"strconv"
	"strings"
)

type stack []datatypes.Data

var varAddresses []string

func (s stack) Push(v datatypes.Data) stack {
	return append(s, v)
}

func (s stack) Pop() (stack, datatypes.Data) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func pushIntVar(val int64, d datatypes.DataType, s stack) stack {
	switch d {
	case datatypes.Ptr:
		s = s.Push(datatypes.Data{Value: int32(val), Type: datatypes.Ptr})
	case datatypes.Int8:
		s = s.Push(datatypes.Data{Value: int8(val), Type: datatypes.Int8})
	case datatypes.Int16:
		s = s.Push(datatypes.Data{Value: int16(val), Type: datatypes.Int16})
	case datatypes.Int32:
		s = s.Push(datatypes.Data{Value: int32(val), Type: datatypes.Int32})
	case datatypes.Int64:
		s = s.Push(datatypes.Data{Value: int64(val), Type: datatypes.Int64})
	}

	return s
}

func pushIntVarMem(val int64, d datatypes.DataType, s string, v map[string]datatypes.Data) map[string]datatypes.Data {
	switch d {
	case datatypes.Ptr:
		v[s] = datatypes.Data{Value: int32(val), Type: datatypes.Ptr}
	case datatypes.Int8:
		v[s] = datatypes.Data{Value: int8(val), Type: datatypes.Int8}
	case datatypes.Int16:
		v[s] = datatypes.Data{Value: int16(val), Type: datatypes.Int16}
	case datatypes.Int32:
		v[s] = datatypes.Data{Value: int32(val), Type: datatypes.Int32}
	case datatypes.Int64:
		v[s] = datatypes.Data{Value: int64(val), Type: datatypes.Int64}
	}

	return v
}

func pushFloatVar(val float64, d datatypes.DataType, s stack) stack {
	switch d {
	case datatypes.Float32:
		s = s.Push(datatypes.Data{Value: float32(val), Type: datatypes.Float32})
	case datatypes.Float64:
		s = s.Push(datatypes.Data{Value: float64(val), Type: datatypes.Float64})
	}

	return s
}

func pushFloatVarMem(val float64, d datatypes.DataType, s string, v map[string]datatypes.Data) map[string]datatypes.Data {
	switch d {
	case datatypes.Float32:
		v[s] = datatypes.Data{Value: float32(val), Type: datatypes.Float32}
	case datatypes.Float64:
		v[s] = datatypes.Data{Value: float64(val), Type: datatypes.Float64}
	}

	return v
}

/*
Arithmetic operations
*/

func getInt64(data datatypes.Data) int64 {
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

func intOp(s stack, op types.Op) stack {
	var a1 datatypes.Data
	var a2 datatypes.Data

	s, a2 = s.Pop()
	s, a1 = s.Pop()

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

	isPtr := false
	if a1.Type == datatypes.Ptr || a2.Type == datatypes.Ptr {
		isPtr = true
	}

	left := getInt64(a1)
	right := getInt64(a2)

	var result int64

	switch op {
	case types.Add:
		result = left + right
	case types.Sub:
		result = left - right
	case types.Mul:
		result = left * right
	case types.Div:
		result = left / right
	case types.Mod:
		result = left % right
	case types.Shr, types.Shl:
		leftu := uint64(left)
		rightu := uint64(right)

		if op == types.Shl {
			result = int64(leftu << rightu)
		} else {
			result = int64(leftu >> rightu)
		}
	case types.And:
		result = left & right
	case types.Or:
		result = left | right
	case types.Xor:
		result = left ^ right
	case types.Le:
		s = s.Push(datatypes.Data{Value: left <= right, Type: datatypes.Bit})
		return s
	case types.Ge:
		s = s.Push(datatypes.Data{Value: left >= right, Type: datatypes.Bit})
		return s
	case types.L:
		s = s.Push(datatypes.Data{Value: left < right, Type: datatypes.Bit})
		return s
	case types.G:
		s = s.Push(datatypes.Data{Value: left > right, Type: datatypes.Bit})
		return s
	}

	if isPtr {
		s = s.Push(datatypes.Data{Value: int32(result), Type: datatypes.Ptr})
		return s
	}

	return pushIntVar(result, min, s)
}

func negateInt(s stack) stack {
	var op datatypes.Data
	s, op = s.Pop()

	if !(op.Type >= datatypes.Int8 && op.Type <= datatypes.Int64) {
		fmt.Println("Only int or bit allowed in negate operation!")
		return nil
	}

	opInt := getInt64(op)
	opInt = ^opInt

	return pushIntVar(opInt, op.Type, s)
}

func getFloat64(data datatypes.Data) float64 {
	switch data.Type {
	case datatypes.Float32:
		return float64(data.Value.(float32))
	case datatypes.Float64:
		return float64(data.Value.(float64))
	default:
		return 0
	}
}

func floatOp(s stack, op types.Op) stack {
	var a1 datatypes.Data
	var a2 datatypes.Data

	s, a2 = s.Pop()
	s, a1 = s.Pop()

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
	case types.Add:
		result = left + right
	case types.Sub:
		result = left - right
	case types.Mul:
		result = left * right
	case types.Div:
		result = left / right
	case types.Le:
		s = s.Push(datatypes.Data{Value: left <= right, Type: datatypes.Bit})
		return s
	case types.Ge:
		s = s.Push(datatypes.Data{Value: left >= right, Type: datatypes.Bit})
		return s
	case types.L:
		s = s.Push(datatypes.Data{Value: left < right, Type: datatypes.Bit})
		return s
	case types.G:
		s = s.Push(datatypes.Data{Value: left > right, Type: datatypes.Bit})
		return s
	}

	return pushFloatVar(result, min, s)
}

/*
Constant declarations
*/

func declareIntConst(command types.Command, c []datatypes.Data, d datatypes.DataType) []datatypes.Data {
	var num int64
	num, _ = strconv.ParseInt(command.Param.Text, 0, 64)

	switch d {
	case datatypes.Int8:
		c = append(c, datatypes.Data{Value: int8(num), Type: datatypes.Int8})
	case datatypes.Int16:
		c = append(c, datatypes.Data{Value: int16(num), Type: datatypes.Int16})
	case datatypes.Int32:
		c = append(c, datatypes.Data{Value: int32(num), Type: datatypes.Int32})
	case datatypes.Int64:
		c = append(c, datatypes.Data{Value: int64(num), Type: datatypes.Int64})
	}

	return c
}

func declareFloatConst(command types.Command, c []datatypes.Data, d datatypes.DataType) []datatypes.Data {
	var num float64
	num, _ = strconv.ParseFloat(command.Param.Text, 64)

	switch d {
	case datatypes.Float32:
		c = append(c, datatypes.Data{Value: float32(num), Type: datatypes.Float32})
	case datatypes.Float64:
		c = append(c, datatypes.Data{Value: float64(num), Type: datatypes.Float64})
	}

	return c
}

func declareStringConstant(command types.Command, c []datatypes.Data) []datatypes.Data {
	val := strings.Trim(command.Param.Text, "\"")
	var stringtype datatypes.DataType

	if command.Command.Text == "dcsa" {
		stringtype = datatypes.String_ASCII
	} else {
		stringtype = datatypes.String_Unicode
	}

	c = append(c, datatypes.Data{Value: val, Type: stringtype})
	return c
}

func declareBitConstant(command types.Command, c []datatypes.Data) []datatypes.Data {
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

	c = append(c, datatypes.Data{Value: val, Type: datatypes.Bit})
	return c
}

/*
Variable declaration
*/

func declareVar(command types.Command, d datatypes.DataType, v map[string]datatypes.Data) map[string]datatypes.Data {

	varAddresses = append(varAddresses, command.Param.Text)

	switch d {
	case datatypes.Bit:
		v[command.Param.Text] = datatypes.Data{Value: false, Type: datatypes.Bit}
	case datatypes.Ptr:
		v[command.Param.Text] = datatypes.Data{Value: int32(0x0), Type: datatypes.Ptr}
	case datatypes.Int8:
		v[command.Param.Text] = datatypes.Data{Value: int8(0), Type: datatypes.Int8}
	case datatypes.Int16:
		v[command.Param.Text] = datatypes.Data{Value: int16(0), Type: datatypes.Int16}
	case datatypes.Int32:
		v[command.Param.Text] = datatypes.Data{Value: int32(0), Type: datatypes.Int32}
	case datatypes.Int64:
		v[command.Param.Text] = datatypes.Data{Value: int64(0), Type: datatypes.Int64}
	case datatypes.Float32:
		v[command.Param.Text] = datatypes.Data{Value: float32(0), Type: datatypes.Float32}
	case datatypes.Float64:
		v[command.Param.Text] = datatypes.Data{Value: float64(0), Type: datatypes.Float64}
	case datatypes.String_ASCII:
		v[command.Param.Text] = datatypes.Data{Value: "", Type: datatypes.String_ASCII}
	case datatypes.String_Unicode:
		v[command.Param.Text] = datatypes.Data{Value: "", Type: datatypes.String_Unicode}
	}

	return v
}

/*
Load value on stack
*/

func loadVar(command types.Command, d datatypes.DataType, v map[string]datatypes.Data, s stack) stack {
	if d >= datatypes.String_ASCII && d <= datatypes.String_Unicode && v[command.Param.Text].Type >= datatypes.String_ASCII && v[command.Param.Text].Type <= datatypes.String_Unicode {
		s = s.Push(v[command.Param.Text])
		return s
	}

	if d >= datatypes.Float32 && d <= datatypes.Float64 && v[command.Param.Text].Type >= datatypes.Float32 && v[command.Param.Text].Type <= datatypes.Float64 {
		return pushFloatVar(getFloat64(v[command.Param.Text]), d, s)
	}

	if d >= datatypes.Ptr && d <= datatypes.Int64 && v[command.Param.Text].Type >= datatypes.Ptr && v[command.Param.Text].Type <= datatypes.Int64 {
		return pushIntVar(getInt64(v[command.Param.Text]), d, s)
	}

	if d == datatypes.Bit && v[command.Param.Text].Type == datatypes.Bit {
		s = s.Push(v[command.Param.Text])
		return s
	}

	fmt.Println("Type mismatch!")
	return nil
}

func loadConst(command types.Command, d datatypes.DataType, c []datatypes.Data, s stack) stack {
	index, _ := strconv.Atoi(command.Param.Text)
	if d >= datatypes.String_ASCII && d <= datatypes.String_Unicode && c[index].Type >= datatypes.String_ASCII && c[index].Type <= datatypes.String_Unicode {
		s = s.Push(c[index])
		return s
	}

	if d >= datatypes.Float32 && d <= datatypes.Float64 && c[index].Type >= datatypes.Float32 && c[index].Type <= datatypes.Float64 {
		return pushFloatVar(getFloat64(c[index]), d, s)
	}

	if d >= datatypes.Ptr && d <= datatypes.Int64 && c[index].Type >= datatypes.Ptr && c[index].Type <= datatypes.Int64 {
		return pushIntVar(getInt64(c[index]), d, s)
	}

	if d == datatypes.Bit && c[index].Type == datatypes.Bit {
		s = s.Push(c[index])
		return s
	}

	fmt.Println("Type mismatch!")
	return nil
}

/*
Call
*/

func call(command types.Command, s stack) stack {
	var num int64
	num, _ = strconv.ParseInt(command.Param.Text, 0, 8)

	var v datatypes.Data
	s, v = s.Pop()

	switch types.Op(num) {
	case types.Print:
		fmt.Println(v.Value)
	}

	return s
}

/*
Store
*/

func store(command types.Command, s stack, v map[string]datatypes.Data) (stack, map[string]datatypes.Data) {
	d := v[command.Param.Text].Type
	var t datatypes.Data
	s, t = s.Pop()

	if d >= datatypes.String_ASCII && d <= datatypes.String_Unicode && t.Type >= datatypes.String_ASCII && t.Type <= datatypes.String_Unicode {
		v[command.Param.Text] = t
		return s, v
	}

	if d >= datatypes.Float32 && d <= datatypes.Float64 && t.Type >= datatypes.Float32 && t.Type <= datatypes.Float64 {
		return s, pushFloatVarMem(getFloat64(t), d, command.Param.Text, v)
	}

	if d >= datatypes.Ptr && d <= datatypes.Int64 && t.Type >= datatypes.Ptr && t.Type <= datatypes.Int64 {
		return s, pushIntVarMem(getInt64(t), d, command.Param.Text, v)
	}

	if d == datatypes.Bit && t.Type == datatypes.Bit {
		v[command.Param.Text] = t
		return s, v
	}

	fmt.Println("Type mismatch!")
	return nil, nil
}

/*
Extern
*/

func extern(command types.Command, v map[string]datatypes.Data) map[string]datatypes.Data {
	switch command.Param.Text {
	case "RETURN_CODE":
		v["RETURN_CODE"] = datatypes.Data{Value: int8(0), Type: datatypes.Int8}
	}
	return v
}

func Run(root types.Root) int8 {
	s := make(stack, 0)
	varAddresses = make([]string, 0)
	var c []datatypes.Data
	v := make(map[string]datatypes.Data)

	for e := 0; e < len(root.Commands); e++ {
		if root.Commands[e].Single {
			switch root.Commands[e].Command.Text {
			case "pop":
				s, _ = s.Pop()

			/*
				Arithmetic int operations
			*/
			case "add":
				s = intOp(s, types.Add)
			case "sub":
				s = intOp(s, types.Sub)
			case "mul":
				s = intOp(s, types.Mul)
			case "div":
				s = intOp(s, types.Div)
			case "mod":
				s = intOp(s, types.Mod)
			case "andi":
				s = intOp(s, types.And)
			case "ori":
				s = intOp(s, types.Or)
			case "xori":
				s = intOp(s, types.Xor)
			case "noti":
				s = negateInt(s)
			case "shl":
				s = intOp(s, types.Shl)
			case "shr":
				s = intOp(s, types.Shr)
			case "ge":
				s = intOp(s, types.Ge)
			case "le":
				s = intOp(s, types.Le)
			case "gt":
				s = intOp(s, types.G)
			case "lt":
				s = intOp(s, types.L)

			/*
				Arithmetic float operations
			*/
			case "addf":
				s = floatOp(s, types.Add)
			case "subf":
				s = floatOp(s, types.Sub)
			case "mulf":
				s = floatOp(s, types.Mul)
			case "divf":
				s = floatOp(s, types.Div)
			case "gef":
				s = floatOp(s, types.Ge)
			case "lef":
				s = floatOp(s, types.Le)
			case "gtf":
				s = floatOp(s, types.G)
			case "ltf":
				s = floatOp(s, types.L)
			}
		} else {
			switch root.Commands[e].Command.Text {
			/*
				Call
				TODO: Pointers, num values, functions
			*/
			case "call":
				s = call(root.Commands[e], s)

			/*
				Store
				TODO: Pointers, num values
			*/
			case "store":
				s, v = store(root.Commands[e], s, v)

			/*
				Jump
				TODO: Pointers, num values
			*/
			case "jmp":
				lbl := root.Commands[e].Param.Text
				e = root.Labels[lbl] - 1
			case "jmpt":
				var val datatypes.Data
				s, val = s.Pop()

				if val.Value.(bool) {
					lbl := root.Commands[e].Param.Text
					e = root.Labels[lbl] - 1
				}
			case "jmpf":
				var val datatypes.Data
				s, val = s.Pop()

				if !val.Value.(bool) {
					lbl := root.Commands[e].Param.Text
					e = root.Labels[lbl] - 1
				}

			/*
				Extern
			*/
			case "extern":
				v = extern(root.Commands[e], v)

			/*
				Declare int constant
			*/
			case "dci8":
				c = declareIntConst(root.Commands[e], c, datatypes.Int8)
			case "dci16":
				c = declareIntConst(root.Commands[e], c, datatypes.Int16)
			case "dci32", "dci":
				c = declareIntConst(root.Commands[e], c, datatypes.Int32)
			case "dci64":
				c = declareIntConst(root.Commands[e], c, datatypes.Int64)

			/*
				Declare float constant
			*/
			case "dcf32", "dcf":
				c = declareFloatConst(root.Commands[e], c, datatypes.Float32)
			case "dcf64", "dcd":
				c = declareFloatConst(root.Commands[e], c, datatypes.Float64)

			/*
				Declare string constant
				This implementation of Persephone doesn't differentiate between ASCII and Unicode
			*/
			case "dcsa", "dcsu":
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
				s = loadVar(root.Commands[e], datatypes.Int8, v, s)
			case "ldi16v":
				s = loadVar(root.Commands[e], datatypes.Int16, v, s)
			case "ldi32v", "ldiv":
				s = loadVar(root.Commands[e], datatypes.Int32, v, s)
			case "ldi64v":
				s = loadVar(root.Commands[e], datatypes.Int64, v, s)
			case "ldf32v", "ldfv":
				s = loadVar(root.Commands[e], datatypes.Float32, v, s)
			case "ldf64v", "lddv":
				s = loadVar(root.Commands[e], datatypes.Float64, v, s)
			case "ldsav":
				s = loadVar(root.Commands[e], datatypes.String_ASCII, v, s)
			case "ldsuv":
				s = loadVar(root.Commands[e], datatypes.String_Unicode, v, s)
			case "ldbv":
				s = loadVar(root.Commands[e], datatypes.Bit, v, s)
			case "ldptrv":
				s = loadVar(root.Commands[e], datatypes.Ptr, v, s)

			/*
				Load constant onto stack
			*/
			case "ldi8c":
				s = loadConst(root.Commands[e], datatypes.Int8, c, s)
			case "ldi16c":
				s = loadConst(root.Commands[e], datatypes.Int16, c, s)
			case "ldi32c", "ldic":
				s = loadConst(root.Commands[e], datatypes.Int32, c, s)
			case "ldi64c":
				s = loadConst(root.Commands[e], datatypes.Int64, c, s)
			case "ldf32c", "ldfc":
				s = loadConst(root.Commands[e], datatypes.Float32, c, s)
			case "ldf64c", "lddc":
				s = loadConst(root.Commands[e], datatypes.Float64, c, s)
			case "ldsac":
				s = loadConst(root.Commands[e], datatypes.String_ASCII, c, s)
			case "ldsuc":
				s = loadConst(root.Commands[e], datatypes.String_Unicode, c, s)
			case "ldbc":
				s = loadConst(root.Commands[e], datatypes.Bit, c, s)

				//TODO: Pointers, string functions, logical operators, prepare variables
			}
		}
	}

	if val, ok := v["RETURN_CODE"]; ok {
		return val.Value.(int8)
	}

	return 0
}
