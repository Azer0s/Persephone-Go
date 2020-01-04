package runtime

import (
	"../datatypes"
	"../types"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eiannone/keyboard"
)

func replaceAtIndex(in string, r rune, i uint64) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

type stack []datatypes.Data
type intStack []int

var varAddress = 0
var addresses []string
var revAddresses map[string]int
var cbase map[int]int
var constants []datatypes.Data

type nop struct {
}

func (n nop) op() {
}

func (s intStack) Push(v int) intStack {
	return append(s, v)
}

func (s intStack) Pop() (intStack, int) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func (s stack) Push(v datatypes.Data) stack {
	return append(s, v)
}

func (s stack) Pop() (stack, datatypes.Data) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func getBoolFromValue(val datatypes.Data) bool {
	if (val.Type >= datatypes.Int8 && val.Type <= datatypes.Int64) || (val.Type >= datatypes.Ptr && val.Type <= datatypes.Uint64) {
		return getUint64(val) > 0
	} else if val.Type == datatypes.Bit {
		return val.Value.(bool)
	}

	panic("Variable is neither truthy nor falsy!")
}

func pushIntVar(val int64, d datatypes.DataType, s stack) stack {
	switch d {
	case datatypes.Int8:
		s = s.Push(datatypes.Data{Value: int8(val), Type: datatypes.Int8})
	case datatypes.Int16:
		s = s.Push(datatypes.Data{Value: int16(val), Type: datatypes.Int16})
	case datatypes.Int32:
		s = s.Push(datatypes.Data{Value: int32(val), Type: datatypes.Int32})
	case datatypes.Int64:
		s = s.Push(datatypes.Data{Value: val, Type: datatypes.Int64})
	}

	return s
}

func pushUintVar(val uint64, d datatypes.DataType, s stack) stack {
	switch d {
	case datatypes.Ptr:
		s = s.Push(datatypes.Data{Value: uint32(val), Type: datatypes.Ptr})
	case datatypes.Uint8:
		s = s.Push(datatypes.Data{Value: uint8(val), Type: datatypes.Uint8})
	case datatypes.Uint16:
		s = s.Push(datatypes.Data{Value: uint16(val), Type: datatypes.Uint16})
	case datatypes.Uint32:
		s = s.Push(datatypes.Data{Value: uint32(val), Type: datatypes.Uint32})
	case datatypes.Uint64:
		s = s.Push(datatypes.Data{Value: val, Type: datatypes.Uint64})
	}

	return s
}

func pushIntVarMem(val int64, d datatypes.DataType, s string, v map[string]datatypes.Data) map[string]datatypes.Data {
	switch d {
	case datatypes.Int8:
		v[s] = datatypes.Data{Value: int8(val), Type: datatypes.Int8}
	case datatypes.Int16:
		v[s] = datatypes.Data{Value: int16(val), Type: datatypes.Int16}
	case datatypes.Int32:
		v[s] = datatypes.Data{Value: int32(val), Type: datatypes.Int32}
	case datatypes.Int64:
		v[s] = datatypes.Data{Value: val, Type: datatypes.Int64}
	}

	return v
}

func pushUintVarMem(val uint64, d datatypes.DataType, s string, v map[string]datatypes.Data) map[string]datatypes.Data {
	switch d {
	case datatypes.Ptr:
		v[s] = datatypes.Data{Value: uint32(val), Type: datatypes.Ptr}
	case datatypes.Uint8:
		v[s] = datatypes.Data{Value: uint8(val), Type: datatypes.Uint8}
	case datatypes.Uint16:
		v[s] = datatypes.Data{Value: uint16(val), Type: datatypes.Uint16}
	case datatypes.Uint32:
		v[s] = datatypes.Data{Value: uint32(val), Type: datatypes.Uint32}
	case datatypes.Uint64:
		v[s] = datatypes.Data{Value: val, Type: datatypes.Uint64}
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
	case datatypes.Int8:
		return int64(data.Value.(int8))
	case datatypes.Int16:
		return int64(data.Value.(int16))
	case datatypes.Int32:
		return int64(data.Value.(int32))
	case datatypes.Int64:
		return data.Value.(int64)
	case datatypes.Ptr:
		return int64(data.Value.(uint32))
	case datatypes.Uint8:
		return int64(data.Value.(uint8))
	case datatypes.Uint16:
		return int64(data.Value.(uint16))
	case datatypes.Uint32:
		return int64(data.Value.(uint32))
	case datatypes.Uint64:
		return int64(data.Value.(uint64))
	default:
		return 0
	}
}

func getUint64(data datatypes.Data) uint64 {
	switch data.Type {
	case datatypes.Int8:
		return uint64(data.Value.(int8))
	case datatypes.Int16:
		return uint64(data.Value.(int16))
	case datatypes.Int32:
		return uint64(data.Value.(int32))
	case datatypes.Int64:
		return uint64(data.Value.(int64))
	case datatypes.Ptr:
		return uint64(data.Value.(uint32))
	case datatypes.Uint8:
		return uint64(data.Value.(uint8))
	case datatypes.Uint16:
		return uint64(data.Value.(uint16))
	case datatypes.Uint32:
		return uint64(data.Value.(uint32))
	case datatypes.Uint64:
		return data.Value.(uint64)
	default:
		return 0
	}
}

func intOp(s stack, op types.Op) stack {
	var a1 datatypes.Data
	var a2 datatypes.Data

	s, a2 = s.Pop()
	s, a1 = s.Pop()

	if (a1.Type >= datatypes.Ptr && a1.Type <= datatypes.Uint64) && (a2.Type >= datatypes.Ptr && a2.Type <= datatypes.Uint64) {
		min := a1.Type

		if a2.Type > min {
			min = a2.Type
		}

		left := getUint64(a1)
		right := getUint64(a2)

		var result uint64

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
			if op == types.Shl {
				result = left << right
			} else {
				result = left >> right
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

		return pushUintVar(result, min, s)
	}

	if !(((a1.Type >= datatypes.Int8 && a1.Type <= datatypes.Int64) || (a1.Type >= datatypes.Ptr && a1.Type <= datatypes.Uint64)) && ((a2.Type >= datatypes.Int8 && a2.Type <= datatypes.Int64) || (a2.Type >= datatypes.Ptr && a2.Type <= datatypes.Uint64))) {
		panic("One of the values is not of type int or ptr!")
	}

	a1Type := a1.Type
	a2Type := a2.Type

	//Convert uint types to int types for a1
	switch a1Type {
	case datatypes.Uint8:
		a1.Type = datatypes.Int8
	case datatypes.Uint16:
		a1.Type = datatypes.Int16
	case datatypes.Uint32, datatypes.Ptr:
		a1.Type = datatypes.Int32
	case datatypes.Uint64:
		a1.Type = datatypes.Int64
	}

	//Convert uint types to int types for a2
	switch a2Type {
	case datatypes.Uint8:
		a2.Type = datatypes.Int8
	case datatypes.Uint16:
		a2.Type = datatypes.Int16
	case datatypes.Uint32, datatypes.Ptr:
		a2.Type = datatypes.Int32
	case datatypes.Uint64:
		a2.Type = datatypes.Int64
	}

	min := a1.Type

	if a2.Type > min {
		min = a2.Type
	}

	a1.Type = a1Type
	a2.Type = a2Type
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

	return pushIntVar(result, min, s)
}

func intSingleOp(s stack, op types.Op) stack {
	var opv datatypes.Data
	s, opv = s.Pop()

	if opv.Type >= datatypes.Ptr && opv.Type <= datatypes.Uint64 {
		opInt := getUint64(opv)

		switch op {
		case types.Not:
			opInt = ^opInt
		case types.Inc:
			opInt++
		case types.Dec:
			opInt--
		}

		return pushUintVar(opInt, opv.Type, s)
	}

	if !(opv.Type >= datatypes.Int8 && opv.Type <= datatypes.Int64) {
		panic("Value is not of type int or bit!")
	}

	opInt := getInt64(opv)

	switch op {
	case types.Not:
		opInt = ^opInt
	case types.Inc:
		opInt++
	case types.Dec:
		opInt--
	}

	return pushIntVar(opInt, opv.Type, s)
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
		panic("Value is not of type float!")
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
Logical operations
*/

func bitOp(s stack, op types.Op) stack {
	var left datatypes.Data
	var right datatypes.Data

	s, right = s.Pop()

	if right.Type != datatypes.Bit {
		panic("Value is not of type bit!")
	}

	if op == types.Not {
		s = s.Push(datatypes.Data{Value: !right.Value.(bool), Type: datatypes.Bit})
		return s
	}

	s, left = s.Pop()
	if left.Type != datatypes.Bit {
		panic("Value is not of type bit!")
	}

	var result bool

	switch op {
	case types.And:
		result = left.Value.(bool) && right.Value.(bool)
	case types.Or:
		result = left.Value.(bool) || right.Value.(bool)
	case types.Xor:
		result = left.Value.(bool) != right.Value.(bool)
	}

	s = s.Push(datatypes.Data{Value: result, Type: datatypes.Bit})
	return s
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
		c = append(c, datatypes.Data{Value: num, Type: datatypes.Int64})
	}

	return c
}

func declareUintConst(command types.Command, c []datatypes.Data, d datatypes.DataType) []datatypes.Data {
	var num uint64
	num, _ = strconv.ParseUint(command.Param.Text, 0, 64)

	switch d {
	case datatypes.Uint8:
		c = append(c, datatypes.Data{Value: uint8(num), Type: datatypes.Uint8})
	case datatypes.Uint16:
		c = append(c, datatypes.Data{Value: uint16(num), Type: datatypes.Uint16})
	case datatypes.Uint32:
		c = append(c, datatypes.Data{Value: uint32(num), Type: datatypes.Uint32})
	case datatypes.Uint64:
		c = append(c, datatypes.Data{Value: num, Type: datatypes.Uint64})
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
		stringtype = datatypes.StringASCII
	} else {
		stringtype = datatypes.StringUnicode
	}

	c = append(c, datatypes.Data{Value: val, Type: stringtype})
	return c
}

func declareBitConstant(command types.Command, c []datatypes.Data) []datatypes.Data {
	val := false

	switch command.Param.Text {
	case "false":
		val = false
	case "true":
		val = true
	default:
		panic("Value is not of type bit!")
	}

	c = append(c, datatypes.Data{Value: val, Type: datatypes.Bit})
	return c
}

/*
Variable declaration
*/

func declareVar(command types.Command, d datatypes.DataType, v map[string]datatypes.Data) map[string]datatypes.Data {

	addresses = append(addresses, command.Param.Text)
	revAddresses[command.Param.Text] = varAddress
	varAddress++

	switch d {
	case datatypes.Bit:
		v[command.Param.Text] = datatypes.Data{Value: false, Type: datatypes.Bit}
	case datatypes.Ptr:
		v[command.Param.Text] = datatypes.Data{Value: uint32(0x0), Type: datatypes.Ptr}
	case datatypes.Int8:
		v[command.Param.Text] = datatypes.Data{Value: int8(0), Type: datatypes.Int8}
	case datatypes.Int16:
		v[command.Param.Text] = datatypes.Data{Value: int16(0), Type: datatypes.Int16}
	case datatypes.Int32:
		v[command.Param.Text] = datatypes.Data{Value: int32(0), Type: datatypes.Int32}
	case datatypes.Int64:
		v[command.Param.Text] = datatypes.Data{Value: int64(0), Type: datatypes.Int64}
	case datatypes.Uint8:
		v[command.Param.Text] = datatypes.Data{Value: uint8(0), Type: datatypes.Uint8}
	case datatypes.Uint16:
		v[command.Param.Text] = datatypes.Data{Value: uint16(0), Type: datatypes.Uint16}
	case datatypes.Uint32:
		v[command.Param.Text] = datatypes.Data{Value: uint32(0), Type: datatypes.Uint32}
	case datatypes.Uint64:
		v[command.Param.Text] = datatypes.Data{Value: uint64(0), Type: datatypes.Uint64}
	case datatypes.Float32:
		v[command.Param.Text] = datatypes.Data{Value: float32(0), Type: datatypes.Float32}
	case datatypes.Float64:
		v[command.Param.Text] = datatypes.Data{Value: float64(0), Type: datatypes.Float64}
	case datatypes.StringASCII:
		v[command.Param.Text] = datatypes.Data{Value: "", Type: datatypes.StringASCII}
	case datatypes.StringUnicode:
		v[command.Param.Text] = datatypes.Data{Value: "", Type: datatypes.StringUnicode}
	}

	return v
}

/*
Load value on stack
*/

func loadVar(command types.Command, d datatypes.DataType, v map[string]datatypes.Data, s stack) stack {
	name := getByPtr(command, v)
	t := v[name].Type

	//Load uint8 as string
	if d >= datatypes.StringASCII && d <= datatypes.StringUnicode && t == datatypes.Uint8 {
		s = s.Push(datatypes.Data{
			Value: string(v[name].Value.(uint8)),
			Type:  d,
		})
		return s
	}

	//Load string as uint8
	if d == datatypes.Uint8 && t >= datatypes.StringASCII && t <= datatypes.StringUnicode && len(v[name].Value.(string)) == 1 {
		s = s.Push(datatypes.Data{
			Value: uint8(v[name].Value.(string)[0]),
			Type:  t,
		})
		return s
	}

	if d >= datatypes.StringASCII && d <= datatypes.StringUnicode && t >= datatypes.StringASCII && t <= datatypes.StringUnicode {
		s = s.Push(v[name])
		return s
	}

	if d >= datatypes.Float32 && d <= datatypes.Float64 && t >= datatypes.Float32 && t <= datatypes.Float64 {
		return pushFloatVar(getFloat64(v[name]), d, s)
	}

	if d >= datatypes.Int8 && d <= datatypes.Int64 && ((t >= datatypes.Int8 && t <= datatypes.Int64) || (t >= datatypes.Ptr && t <= datatypes.Uint64)) {
		return pushIntVar(getInt64(v[name]), d, s)
	}

	if d >= datatypes.Ptr && d <= datatypes.Uint64 && ((t >= datatypes.Int8 && t <= datatypes.Int64) || (t >= datatypes.Ptr && t <= datatypes.Uint64)) {
		return pushUintVar(getUint64(v[name]), d, s)
	}

	if d == datatypes.Bit && t == datatypes.Bit {
		s = s.Push(v[name])
		return s
	}

	panic("Type mismatch!")
}

func loadConst(command types.Command, d datatypes.DataType, c []datatypes.Data, s stack, line int) stack {
	index, _ := strconv.Atoi(command.Param.Text)

	min := 0
	for e := range cbase {
		if e < line {
			min = e
		}
	}

	index += cbase[min]

	//Load string as uint8
	if c[index].Type >= datatypes.StringASCII && c[index].Type <= datatypes.StringUnicode && d == datatypes.Uint8 && len(c[index].Value.(string)) == 1 {
		s = s.Push(datatypes.Data{
			Value: uint8(c[index].Value.(string)[0]),
			Type:  d,
		})
		return s
	}

	//Load uint8 as string
	if d >= datatypes.StringASCII && d <= datatypes.StringUnicode && c[index].Type == datatypes.Uint8 {
		s = s.Push(datatypes.Data{
			Value: string(c[index].Value.(uint8)),
			Type:  d,
		})
		return s
	}

	if d >= datatypes.StringASCII && d <= datatypes.StringUnicode && c[index].Type >= datatypes.StringASCII && c[index].Type <= datatypes.StringUnicode {
		s = s.Push(c[index])
		return s
	}

	if d >= datatypes.Float32 && d <= datatypes.Float64 && c[index].Type >= datatypes.Float32 && c[index].Type <= datatypes.Float64 {
		return pushFloatVar(getFloat64(c[index]), d, s)
	}

	if d >= datatypes.Int8 && d <= datatypes.Int64 && ((c[index].Type >= datatypes.Int8 && c[index].Type <= datatypes.Int64) || (c[index].Type >= datatypes.Uint8 && c[index].Type <= datatypes.Uint64)) {
		return pushIntVar(getInt64(c[index]), d, s)
	}

	if d >= datatypes.Uint8 && d <= datatypes.Uint64 && ((c[index].Type >= datatypes.Int8 && c[index].Type <= datatypes.Int64) || (c[index].Type >= datatypes.Uint8 && c[index].Type <= datatypes.Uint64)) {
		return pushUintVar(getUint64(c[index]), d, s)
	}

	if d == datatypes.Bit && c[index].Type == datatypes.Bit {
		s = s.Push(c[index])
		return s
	}

	panic("Type mismatch!")
}

/*
Syscall
*/

func syscall(command types.Command, s stack, vars map[string]datatypes.Data) stack {
	var num uint64

	switch command.Param.Kind {
	case types.HexNumber:
		num, _ = strconv.ParseUint(command.Param.Text, 0, 8)
	case types.Number:
		num, _ = strconv.ParseUint(command.Param.Text, 10, 8)
	case types.Pointer:
		val := vars[strings.Trim(strings.Trim(command.Param.Text, "]"), "[")]

		if val.Type != datatypes.Ptr {
			panic("Variable is not of type ptr!")
		}

		num = uint64(val.Value.(uint32))
	}

	switch types.Op(num) {
	case types.Print:
		var v datatypes.Data
		s, v = s.Pop()
		fmt.Print(v.Value)
	case types.Read:
		ch, _, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}
		val := string(ch)
		fmt.Print(val)
		s = s.Push(datatypes.Data{Value: val, Type: datatypes.StringASCII})
	case types.Println:
		var v datatypes.Data
		s, v = s.Pop()
		fmt.Println(v.Value)
	case types.Readln:
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadBytes('\n')
		s = s.Push(datatypes.Data{Value: string(text), Type: datatypes.StringASCII})
	}

	return s
}

/*
Store
*/

func store(command types.Command, s stack, v map[string]datatypes.Data) (stack, map[string]datatypes.Data) {
	name := getByPtr(command, v)
	d := v[name].Type
	var t datatypes.Data
	s, t = s.Pop()

	//Store string as uint8
	if d == datatypes.Uint8 && t.Type >= datatypes.StringASCII && t.Type <= datatypes.StringUnicode && len(t.Value.(string)) == 1 {
		v[name] = datatypes.Data{
			Value: uint8(t.Value.(string)[0]),
			Type:  d,
		}
		return s, v
	}

	if d >= datatypes.StringASCII && d <= datatypes.StringUnicode && t.Type >= datatypes.StringASCII && t.Type <= datatypes.StringUnicode {
		v[name] = t
		return s, v
	}

	//Store uint8 as string
	if d >= datatypes.StringASCII && d <= datatypes.StringUnicode {
		if t.Type == datatypes.Uint8 {
			v[name] = datatypes.Data{Value: string(t.Value.(uint8)), Type: d}
			return s, v
		}
	}

	if d >= datatypes.Float32 && d <= datatypes.Float64 && t.Type >= datatypes.Float32 && t.Type <= datatypes.Float64 {
		return s, pushFloatVarMem(getFloat64(t), d, name, v)
	}

	if d >= datatypes.Int8 && d <= datatypes.Int64 && ((t.Type >= datatypes.Int8 && t.Type <= datatypes.Int64) || (t.Type >= datatypes.Ptr && t.Type <= datatypes.Uint64)) {
		return s, pushIntVarMem(getInt64(t), d, name, v)
	}

	if d >= datatypes.Ptr && d <= datatypes.Uint64 && ((t.Type >= datatypes.Int8 && t.Type <= datatypes.Int64) || (t.Type >= datatypes.Ptr && t.Type <= datatypes.Uint64)) {
		return s, pushUintVarMem(getUint64(t), d, name, v)
	}

	if d == datatypes.Bit && t.Type == datatypes.Bit {
		v[name] = t
		return s, v
	}

	panic("Type mismatch!")
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

func getByPtr(command types.Command, v map[string]datatypes.Data) string {
	switch command.Param.Kind {
	case types.Pointer:
		ptr := v[strings.Trim(strings.Trim(command.Param.Text, "]"), "[")]

		if ptr.Type != datatypes.Ptr {
			panic("Variable is not of type ptr!")
		}

		return addresses[ptr.Value.(uint32)]

	case types.Name:
		return command.Param.Text
	}

	return ""
}

//Run ...Walks the AST and basically executes the program "line-by-line"
func Run(root types.Root) int8 {
	//Global var initialization
	constants = make([]datatypes.Data, 0)
	cbase = make(map[int]int)
	addresses = make([]string, 0)
	revAddresses = make(map[string]int)

	s := make(stack, 0)
	r := make(intStack, 0)
	v := make(map[string]datatypes.Data)

	for k := range root.Labels {
		addresses = append(addresses, k)
		revAddresses[k] = varAddress
		varAddress++
	}

	for e := 0; e < len(root.Commands); e++ {
		if root.Commands[e].Single {
			switch root.Commands[e].Command.Text {
			case "pop":
				s, _ = s.Pop()

			case "nop":
				n := &nop{}
				n.op()

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
				s = intSingleOp(s, types.Not)
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
			case "inc":
				s = intSingleOp(s, types.Inc)
			case "dec":
				s = intSingleOp(s, types.Dec)

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

			case "and":
				s = bitOp(s, types.And)
			case "or":
				s = bitOp(s, types.Or)
			case "xor":
				s = bitOp(s, types.Xor)
			case "not":
				s = bitOp(s, types.Not)

			case "eq":
				var a, b datatypes.Data
				s, a = s.Pop()
				s, b = s.Pop()
				s = s.Push(datatypes.Data{
					Value: a.Value == b.Value,
					Type:  datatypes.Bit,
				})

			/*
				Cbase
			*/
			case "cbase":
				cbase[e] = len(constants)

			/*
				Return
			*/
			case "ret":
				r, e = r.Pop()

			/*
				Concatenate
				Can also be used to convert values to string (by concatenating a value with empty string)
			*/
			case "conc":
				var left datatypes.Data
				var right datatypes.Data

				s, left = s.Pop()
				s, right = s.Pop()

				s = s.Push(datatypes.Data{Value: fmt.Sprintf("%v", left.Value) + fmt.Sprintf("%v", right.Value), Type: datatypes.StringUnicode})
			}
		} else {
			switch root.Commands[e].Command.Text {
			/*
				Syscall
			*/
			case "syscall":
				s = syscall(root.Commands[e], s, v)

			/*
				Store
			*/
			case "store":
				s, v = store(root.Commands[e], s, v)

			/*
				Delete var
			*/
			case "delete":
				varName := getByPtr(root.Commands[e], v)
				_, ok := v[varName]
				if ok {
					delete(v, varName)
				}

			/*
				String functions
			*/
			case "len":
				a := v[getByPtr(root.Commands[e], v)]

				if a.Type == datatypes.StringUnicode || a.Type == datatypes.StringASCII {
					s = s.Push(datatypes.Data{Value: uint32(len(a.Value.(string))), Type: datatypes.Uint32})
				} else {
					switch a.Type {
					case datatypes.Int8, datatypes.Uint8:
						s = s.Push(datatypes.Data{Value: uint32(8), Type: datatypes.Uint32})
					case datatypes.Int16, datatypes.Uint16:
						s = s.Push(datatypes.Data{Value: uint32(16), Type: datatypes.Uint32})
					case datatypes.Int32, datatypes.Uint32, datatypes.Ptr, datatypes.Float32:
						s = s.Push(datatypes.Data{Value: uint32(32), Type: datatypes.Uint32})
					case datatypes.Int64, datatypes.Uint64, datatypes.Float64:
						s = s.Push(datatypes.Data{Value: uint32(64), Type: datatypes.Uint32})
					case datatypes.Bit:
						s = s.Push(datatypes.Data{Value: uint32(1), Type: datatypes.Uint32})
					default:
						panic("Value is not of type stringa or stringu!")
					}
				}

			case "getc":
				a1 := v[getByPtr(root.Commands[e], v)]
				if a1.Type == datatypes.StringUnicode || a1.Type == datatypes.StringASCII {
					var val datatypes.Data
					s, val = s.Pop()
					s = s.Push(datatypes.Data{Value: uint8(a1.Value.(string)[getUint64(val)]), Type: datatypes.Uint8})
				} else {
					panic("Value is not of type stringa or stringu!")
				}

			case "setc":
				a1 := v[getByPtr(root.Commands[e], v)]
				if a1.Type == datatypes.StringUnicode || a1.Type == datatypes.StringASCII {
					var char datatypes.Data
					s, char = s.Pop()

					var pos datatypes.Data
					s, pos = s.Pop()

					if char.Type != datatypes.Uint8 {
						panic("Value is not of type int8!")
					}

					tempVal := a1.Value.(string)
					tempVal = replaceAtIndex(tempVal, rune(char.Value.(uint8)), getUint64(pos))
					a1.Value = tempVal

					s = s.Push(datatypes.Data{Value: tempVal, Type: datatypes.StringUnicode})
				} else {
					panic("Value is not of type stringa or stringu!")
				}

			/*
				Jump
			*/
			case "jmp":
				lbl := getByPtr(root.Commands[e], v)
				e = root.Labels[lbl] - 1
			case "jmpt":
				var val datatypes.Data
				s, val = s.Pop()
				if getBoolFromValue(val) {
					lbl := getByPtr(root.Commands[e], v)
					e = root.Labels[lbl] - 1
				}
			case "jmpf":
				var val datatypes.Data
				s, val = s.Pop()
				if !getBoolFromValue(val) {
					lbl := getByPtr(root.Commands[e], v)
					e = root.Labels[lbl] - 1
				}

			/*
				Call
			*/
			case "call":
				r = r.Push(e)
				lbl := getByPtr(root.Commands[e], v)
				e = root.Labels[lbl] - 1

			/*
				Extern
			*/
			case "extern":
				v = extern(root.Commands[e], v)

			/*
				Declare int constant
			*/
			case "dci8":
				constants = declareIntConst(root.Commands[e], constants, datatypes.Int8)
			case "dci16":
				constants = declareIntConst(root.Commands[e], constants, datatypes.Int16)
			case "dci32", "dci":
				constants = declareIntConst(root.Commands[e], constants, datatypes.Int32)
			case "dci64":
				constants = declareIntConst(root.Commands[e], constants, datatypes.Int64)

			case "dcu8":
				constants = declareUintConst(root.Commands[e], constants, datatypes.Uint8)
			case "dcu16":
				constants = declareUintConst(root.Commands[e], constants, datatypes.Uint16)
			case "dcu32", "dcu":
				constants = declareUintConst(root.Commands[e], constants, datatypes.Uint32)
			case "dcu64":
				constants = declareUintConst(root.Commands[e], constants, datatypes.Uint64)

			/*
				Declare float constant
			*/
			case "dcf32", "dcf":
				constants = declareFloatConst(root.Commands[e], constants, datatypes.Float32)
			case "dcf64", "dcd":
				constants = declareFloatConst(root.Commands[e], constants, datatypes.Float64)

			/*
				Declare string constant
				This implementation of Persephone doesn't differentiate between ASCII and Unicode
			*/
			case "dcsa", "dcsu":
				constants = declareStringConstant(root.Commands[e], constants)

			/*
				Declare bit constant
			*/
			case "dcb":
				constants = declareBitConstant(root.Commands[e], constants)

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
			case "v_uint8":
				v = declareVar(root.Commands[e], datatypes.Uint8, v)
			case "v_uint16":
				v = declareVar(root.Commands[e], datatypes.Uint16, v)
			case "v_uint32", "v_uint":
				v = declareVar(root.Commands[e], datatypes.Uint32, v)
			case "v_uint64":
				v = declareVar(root.Commands[e], datatypes.Uint64, v)
			case "v_float32", "v_float":
				v = declareVar(root.Commands[e], datatypes.Float32, v)
			case "v_float64", "v_double":
				v = declareVar(root.Commands[e], datatypes.Float64, v)
			case "v_stringa":
				v = declareVar(root.Commands[e], datatypes.StringASCII, v)
			case "v_stringu":
				v = declareVar(root.Commands[e], datatypes.StringUnicode, v)
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
			case "ldu8v":
				s = loadVar(root.Commands[e], datatypes.Uint8, v, s)
			case "ldu16v":
				s = loadVar(root.Commands[e], datatypes.Uint16, v, s)
			case "ldu32v", "lduv":
				s = loadVar(root.Commands[e], datatypes.Uint32, v, s)
			case "ldu64v":
				s = loadVar(root.Commands[e], datatypes.Uint64, v, s)
			case "ldf32v", "ldfv":
				s = loadVar(root.Commands[e], datatypes.Float32, v, s)
			case "ldf64v", "lddv":
				s = loadVar(root.Commands[e], datatypes.Float64, v, s)
			case "ldsav":
				s = loadVar(root.Commands[e], datatypes.StringASCII, v, s)
			case "ldsuv":
				s = loadVar(root.Commands[e], datatypes.StringUnicode, v, s)
			case "ldbv":
				s = loadVar(root.Commands[e], datatypes.Bit, v, s)
			case "ldptrv":
				s = loadVar(root.Commands[e], datatypes.Ptr, v, s)

			/*
				Load constant onto stack
			*/
			case "ldi8c":
				s = loadConst(root.Commands[e], datatypes.Int8, constants, s, e)
			case "ldi16c":
				s = loadConst(root.Commands[e], datatypes.Int16, constants, s, e)
			case "ldi32c", "ldic":
				s = loadConst(root.Commands[e], datatypes.Int32, constants, s, e)
			case "ldi64c":
				s = loadConst(root.Commands[e], datatypes.Int64, constants, s, e)
			case "ldu8c":
				s = loadConst(root.Commands[e], datatypes.Uint8, constants, s, e)
			case "ldu16c":
				s = loadConst(root.Commands[e], datatypes.Uint16, constants, s, e)
			case "ldu32c", "lduc":
				s = loadConst(root.Commands[e], datatypes.Uint32, constants, s, e)
			case "ldu64c":
				s = loadConst(root.Commands[e], datatypes.Uint64, constants, s, e)
			case "ldf32c", "ldfc":
				s = loadConst(root.Commands[e], datatypes.Float32, constants, s, e)
			case "ldf64c", "lddc":
				s = loadConst(root.Commands[e], datatypes.Float64, constants, s, e)
			case "ldsac":
				s = loadConst(root.Commands[e], datatypes.StringASCII, constants, s, e)
			case "ldsuc":
				s = loadConst(root.Commands[e], datatypes.StringUnicode, constants, s, e)
			case "ldbc":
				s = loadConst(root.Commands[e], datatypes.Bit, constants, s, e)

			case "ldptr":
				s = s.Push(datatypes.Data{Value: uint32(revAddresses[root.Commands[e].Param.Text]), Type: datatypes.Ptr})
			}
		}
	}

	if val, ok := v["RETURN_CODE"]; ok {
		return val.Value.(int8)
	}

	return 0
}
