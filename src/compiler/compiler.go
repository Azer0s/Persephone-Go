package compiler

import (
	"encoding/binary"
	"math"
	"os"
	"strconv"
	"strings"

	"../types"
)

var opcodes = map[string]uint16{
	"type":      uint16(0xEEEE),
	"nop":       uint16(0x1000),
	"store":     uint16(0x0000),
	"v_dyn":     uint16(0x0101),
	"v_int8":    uint16(0x0108),
	"v_int16":   uint16(0x0110),
	"v_int32":   uint16(0x0120),
	"v_int":     uint16(0x0120),
	"v_int64":   uint16(0x0140),
	"v_uint8":   uint16(0x1108),
	"v_uint16":  uint16(0x1110),
	"v_uint32":  uint16(0x1120),
	"v_uint":    uint16(0x1120),
	"v_uint64":  uint16(0x1140),
	"v_float32": uint16(0x0121),
	"v_float":   uint16(0x0121),
	"v_float64": uint16(0x0141),
	"v_double":  uint16(0x0141),
	"v_stringa": uint16(0x0131),
	"v_stringu": uint16(0x0132),
	"v_ptr":     uint16(0x0150),
	"v_bit":     uint16(0x0100),
	"delete":    uint16(0x0099),
	"dci8":      uint16(0x0218),
	"dci16":     uint16(0x0210),
	"dci32":     uint16(0x0220),
	"dci":       uint16(0x0220),
	"dci64":     uint16(0x0240),
	"dcu8":      uint16(0x1218),
	"dcu16":     uint16(0x1210),
	"dcu32":     uint16(0x1220),
	"dcu":       uint16(0x1220),
	"dcu64":     uint16(0x1240),
	"dcf32":     uint16(0x0221),
	"dcf":       uint16(0x0221),
	"dcf64":     uint16(0x0241),
	"dcd":       uint16(0x0241),
	"dcsa":      uint16(0x0231),
	"dcsu":      uint16(0x0232),
	"dcb":       uint16(0x0200),
	"lddynv":    uint16(0x0301),
	"ldi8v":     uint16(0x0318),
	"ldi16v":    uint16(0x0310),
	"ldi32v":    uint16(0x0320),
	"ldiv":      uint16(0x0320),
	"ldi64v":    uint16(0x0340),
	"ldu8v":     uint16(0x1318),
	"ldu16v":    uint16(0x1310),
	"ldu32v":    uint16(0x1320),
	"lduv":      uint16(0x1320),
	"ldu64v":    uint16(0x1340),
	"ldf32v":    uint16(0x0321),
	"ldfv":      uint16(0x0321),
	"ldf64v":    uint16(0x0341),
	"lddv":      uint16(0x0341),
	"ldsav":     uint16(0x0331),
	"ldsuv":     uint16(0x0332),
	"ldptrv":    uint16(0x0350),
	"ldbv":      uint16(0x0300),
	"ldi8c":     uint16(0x0418),
	"ldi16c":    uint16(0x0410),
	"ldi32c":    uint16(0x0420),
	"ldic":      uint16(0x0420),
	"ldi64c":    uint16(0x0440),
	"ldu8c":     uint16(0x1418),
	"ldu16c":    uint16(0x1410),
	"ldu32c":    uint16(0x1420),
	"lduc":      uint16(0x1420),
	"ldu64c":    uint16(0x1440),
	"ldf32c":    uint16(0x0421),
	"ldfc":      uint16(0x0421),
	"ldf64c":    uint16(0x0441),
	"lddc":      uint16(0x0441),
	"ldsac":     uint16(0x0431),
	"ldsuc":     uint16(0x0432),
	"ldbc":      uint16(0x0400),
	"cbase":     uint16(0xFFFF),
	"pop":       uint16(0x0001),
	"ret":       uint16(0x0002),
	"eg":        uint16(0x0030),
	"add":       uint16(0x0003),
	"sub":       uint16(0x0004),
	"mul":       uint16(0x0005),
	"div":       uint16(0x0006),
	"mod":       uint16(0x0007),
	"ge":        uint16(0x0008),
	"le":        uint16(0x0009),
	"gt":        uint16(0x000A),
	"lt":        uint16(0x000B),
	"andi":      uint16(0x000C),
	"ori":       uint16(0x000D),
	"xori":      uint16(0x000E),
	"noti":      uint16(0x000F),
	"shl":       uint16(0x0010),
	"shr":       uint16(0x0011),
	"inc":       uint16(0x0012),
	"dec":       uint16(0x0013),
	"addf":      uint16(0x0014),
	"subf":      uint16(0x0015),
	"mulf":      uint16(0x0016),
	"divf":      uint16(0x0017),
	"gef":       uint16(0x0018),
	"lef":       uint16(0x0019),
	"gtf":       uint16(0x001A),
	"ltf":       uint16(0x001B),
	"and":       uint16(0x001C),
	"or":        uint16(0x001D),
	"xor":       uint16(0x001E),
	"not":       uint16(0x001F),
	"conc":      uint16(0x0020),
	"len":       uint16(0x0021),
	"getc":      uint16(0x0022),
	"setc":      uint16(0x0023),
	"syscall":   uint16(0x0024),
	"extern":    uint16(0x0025),
	"call":      uint16(0xF000),
	"jmp":       uint16(0xF001),
	"jmpt":      uint16(0xF002),
	"jmpf":      uint16(0xF003),
}

//Value prefixes for parameters
const (
	Uint     = byte(uint8(0x0))
	Int      = byte(uint8(0x1))
	Float    = byte(uint8(0x2))
	StringA  = byte(uint8(0x3))
	StringU  = byte(uint8(0x4))
	Bit      = byte(uint8(0x5))
	Ptr      = byte(uint8(0x6))
	Label    = byte(uint8(0xE))
	Variable = byte(uint8(0xF))
)

//Size prefixes for int and float for parameters
const (
	Int8  = byte(uint8(0x8))
	Int16 = byte(uint8(0x10))
	Int32 = byte(uint8(0x20))
	Int64 = byte(uint8(0x40))

	Float32 = byte(uint8(0x21))
	Float64 = byte(uint8(0x41))
)

var labels = make(map[string]uint64)
var variables = map[string]uint64{
	"RETURN_CODE": uint64(0x0),
}
var currentVariable = uint64(0xF) //first 15 variables are reserved for external values
var currentLabel = uint64(0x0)

func getUint64Bytes(val uint64) []byte {
	return []byte{
		byte((val & 0xFF00000000000000) >> 56),
		byte((val & 0x00FF000000000000) >> 48),
		byte((val & 0x0000FF0000000000) >> 40),
		byte((val & 0x000000FF00000000) >> 32),
		byte((val & 0x00000000FF000000) >> 24),
		byte((val & 0x0000000000FF0000) >> 16),
		byte((val & 0x000000000000FF00) >> 8),
		byte(val & 0x00000000000000FF),
	}
}

func getUint32Bytes(val uint32) []byte {
	return []byte{
		byte((val & 0xFF000000) >> 24),
		byte((val & 0x00FF0000) >> 16),
		byte((val & 0x0000FF00) >> 8),
		byte(val & 0x000000FF),
	}
}

func getUint16Bytes(val uint16) []byte {
	return []byte{
		byte((val & 0xFF00) >> 8), //Get upper 8 bits
		byte(val & 0x00FF),        //Get lower 8 bits
	}
}

func isASCII(s string) bool {
	f := func(r rune) bool {
		return r < 'A' || r > 'z'
	}
	return !(strings.IndexFunc(s, f) != -1)
}

// Compile compiles an AST to Persephone bytecode
func Compile(root types.Root, outname string) int {
	for e := range root.Labels {
		labels[e] = currentLabel
		currentLabel += uint64(0x1)
	}

	labelToByte := make(map[string]int)

	f, err := os.Create(outname)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	fileBytes := make([]byte, 0)

	write := func(bytes []byte) {
		fileBytes = append(fileBytes, bytes...)
	}

	flush := func(f *os.File) {
		_, err := f.Write(fileBytes)

		if err != nil {
			panic(err)
		}
	}

	for e := 0; e < len(root.Commands); e++ {
		for k, v := range root.Labels { //check if a label points to the current line
			if e == v {
				labelToByte[k] = len(fileBytes)
			}
		}

		isJmp := false
		isUint := false

		switch root.Commands[e].Command.Text {
		case "call", "jmp", "jmpt", "jmpf":
			isJmp = true
		case "dcu8", "dcu16", "dcu32", "dcu", "dcu64":
			isUint = true
		}

		write(getUint16Bytes(opcodes[root.Commands[e].Command.Text]))

		if isJmp && root.Commands[e].Param.Kind == types.Name {
			write([]byte{Label})
			write(getUint64Bytes(labels[root.Commands[e].Param.Text]))
		} else {
			switch root.Commands[e].Param.Kind {
			case types.Bit:
				write([]byte{Bit})
				switch root.Commands[e].Param.Text {
				case "true":
					write([]byte{0x1})
				case "false":
					write([]byte{0x0})
				}

			case types.Name:
				write([]byte{Variable})

				if _, ok := variables[root.Commands[e].Param.Text]; !ok {
					variables[root.Commands[e].Param.Text] = currentVariable
					currentVariable += uint64(0x1)
				}

				write(getUint64Bytes(variables[root.Commands[e].Param.Text]))

			case types.HexNumber, types.Number:
				if isUint {
					write([]byte{Uint})
					var num uint64
					num, _ = strconv.ParseUint(root.Commands[e].Param.Text, 0, 64)

					if num <= 255 {
						write([]byte{Int8})
						write([]byte{byte(int8(num))})
					} else if num <= 65535 {
						write([]byte{Int16})
						write(getUint16Bytes(uint16(num)))
					} else if num <= 4294967295 {
						write([]byte{Int32})
						write(getUint32Bytes(uint32(num)))
					} else {
						write([]byte{Int64})
						write(getUint64Bytes(num))
					}
				} else {
					write([]byte{Int})
					var num int64
					num, _ = strconv.ParseInt(root.Commands[e].Param.Text, 0, 64)

					if num >= -128 && num <= 127 {
						write([]byte{Int8})
						write([]byte{byte(int8(num))})
					} else if num >= -32768 && num <= 32767 {
						write([]byte{Int16})
						write(getUint16Bytes(uint16(num)))
					} else if num >= -2147483648 && num <= 2147483647 {
						write([]byte{Int32})
						write(getUint32Bytes(uint32(num)))
					} else {
						write([]byte{Int64})
						write(getUint64Bytes(uint64(num)))
					}
				}

			case types.Pointer:
				write([]byte{Ptr})

				ptrName := strings.Trim(strings.Trim(root.Commands[e].Param.Text, "]"), "[")

				if _, ok := variables[ptrName]; ok {
					write(getUint64Bytes(variables[ptrName]))
				} else {
					panic("Couldn't find variable: " + ptrName + "!")
				}

			case types.String:
				rawString := strings.Trim(root.Commands[e].Param.Text, "\"")

				if isASCII(rawString) {
					write([]byte{StringA})
					stringBytes := []byte(rawString)
					write(getUint64Bytes(uint64(len(stringBytes))))
					write(stringBytes)
				} else {
					rawRunes := []rune(rawString)
					runeBytes := []byte(string(rawRunes))

					write([]byte{StringU})
					write(getUint64Bytes(uint64(len(runeBytes))))
					write(runeBytes)
				}

			//You'll float too
			case types.Float:
				write([]byte{Float})
				switch root.Commands[e].Param.Size {
				case "64":
					write([]byte{Float64})

					var num float64
					num, _ = strconv.ParseFloat(root.Commands[e].Param.Text, 64)

					bits := math.Float64bits(num)
					bytes := make([]byte, 8)
					binary.LittleEndian.PutUint64(bytes, bits)

					write(bytes[:])
				default: //if size isn't stated, use 32 bit
					write([]byte{Float32})

					var num float64
					num, _ = strconv.ParseFloat(root.Commands[e].Param.Text, 32)

					bits := math.Float32bits(float32(num))
					bytes := make([]byte, 4)
					binary.LittleEndian.PutUint32(bytes, bits)

					write(bytes[:])
				}
			}
		}
	}

	header := make([]byte, 0)
	header = append(header, getUint16Bytes(uint16(len(labelToByte)))...)
	for k, v := range labelToByte {
		header = append(header, getUint64Bytes(uint64(labels[k]))...)
		header = append(header, getUint64Bytes(uint64(v))...)
	}

	fileBytes = append(header, fileBytes...)
	flush(f)
	return 0
}
