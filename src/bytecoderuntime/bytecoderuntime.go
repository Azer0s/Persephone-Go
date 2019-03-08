package bytecoderuntime

import (
	"../compiler"
	"../datatypes"
	"encoding/binary"
	"math"
)

type command struct {
	opcode uint16
	param  datatypes.Data
}

func getNextUint64(code []byte, e *int) uint64 {
	val := getUint64FromBytes(code[*e], code[(*e)+1], code[(*e)+2], code[(*e)+3], code[(*e)+4], code[(*e)+5], code[(*e)+6], code[(*e)+7])
	*e += 8
	return val
}

func getUint64FromBytes(a, b, c, d, e, f, g, h byte) uint64 {
	return (uint64(a) << 56) + (uint64(b) << 48) + (uint64(c) << 40) + (uint64(d) << 32) + (uint64(e) << 24) + (uint64(f) << 16) + (uint64(g) << 8) + uint64(h)
}

func getUint32FromBytes(a, b, c, d byte) uint32 {
	return (uint32(a) << 24) + (uint32(b) << 16) + (uint32(c) << 8) + uint32(d)
}

func getUint16FromBytes(a, b byte) uint16 {
	return (uint16(a) << 8) + uint16(b)
}

func getNextUint64(code []byte, e *int) uint64 {
	val := getUint64FromBytes(code[*e], code[(*e)+1], code[(*e)+2], code[(*e)+3], code[(*e)+4], code[(*e)+5], code[(*e)+6], code[(*e)+7])
	*e += 8
	return val
}

//Run ...Runs a compiled Persephone file
func Run(bytes []byte) int8 {
	labels := make(map[uint64]uint64)
	labelCount := getUint16FromBytes(bytes[0], bytes[1])

	labelPtr := 2
	for e := uint16(0); e < labelCount; e++ {
		labels[getUint64FromBytes(bytes[labelPtr],
			bytes[labelPtr+1],
			bytes[labelPtr+2],
			bytes[labelPtr+3],
			bytes[labelPtr+4],
			bytes[labelPtr+5],
			bytes[labelPtr+6],
			bytes[labelPtr+7])] =
			getUint64FromBytes(bytes[labelPtr+8],
				bytes[labelPtr+9],
				bytes[labelPtr+10],
				bytes[labelPtr+11],
				bytes[labelPtr+12],
				bytes[labelPtr+13],
				bytes[labelPtr+14],
				bytes[labelPtr+15])

		labelPtr += 16
	}

	code := bytes[labelPtr:]

	statements := make([]command, 0)

	for e := 0; e < len(code); {
		opcode := getUint16FromBytes(code[e], code[e+1])
		e += 2

		parameter := false

		switch opcode {
		case uint16(0x0003),
			uint16(0x0004),
			uint16(0x0005),
			uint16(0x0006),
			uint16(0x0007),
			uint16(0x000C),
			uint16(0x000D),
			uint16(0x000E),
			uint16(0x000F),
			uint16(0x0010),
			uint16(0x0011),
			uint16(0x0014),
			uint16(0x0015),
			uint16(0x0016),
			uint16(0x0017),
			uint16(0x0001),
			uint16(0x0008),
			uint16(0x0009),
			uint16(0x000A),
			uint16(0x000B),
			uint16(0x0018),
			uint16(0x0019),
			uint16(0x001A),
			uint16(0x001B),
			uint16(0x0012),
			uint16(0x0013),
			uint16(0xFFFF),
			uint16(0x001C),
			uint16(0x001D),
			uint16(0x001E),
			uint16(0x001F),
			uint16(0x0002),
			uint16(0x0020):
			statements = append(statements, command{opcode, datatypes.Data{}})
			continue
		default:
			//Opcode has a parameter
			parameter = true
		}

		if parameter {
			paramType := code[e]
			e++

			param := datatypes.Data{}

			switch paramType {
			case compiler.Int:
				intSize := code[e]
				e++

				switch intSize {
				case compiler.Int8:
					param.Type = datatypes.Int8
					param.Value = int8(code[e])
					e++

				case compiler.Int16:
					param.Type = datatypes.Int16
					param.Value = int16(getUint16FromBytes(code[e], code[e+1]))
					e += 2

				case compiler.Int32:
					param.Type = datatypes.Int32
					param.Value = int32(getUint32FromBytes(code[e], code[e+1], code[e+2], code[e+3]))
					e += 4

				case compiler.Int64:
					param.Type = datatypes.Int64
					param.Value = int64(getUint64FromBytes(code[e], code[e+1], code[e+2], code[e+3], code[e+4], code[e+5], code[e+6], code[e+7]))
					e += 8
				}
			case compiler.Float:
				floatSize := code[e]
				e++

				switch floatSize {
				case compiler.Float32:
					bytes := []byte{code[e], code[e+1], code[e+2], code[e+3]}
					e += 4

					bits := binary.LittleEndian.Uint32(bytes)

					param.Type = datatypes.Float32
					param.Value = math.Float32frombits(bits)

				case compiler.Float64:
					bytes := []byte{code[e], code[e+1], code[e+2], code[e+3], code[e+4], code[e+5], code[e+6], code[e+7]}
					e += 8

					bits := binary.LittleEndian.Uint64(bytes)

					param.Type = datatypes.Float64
					param.Value = math.Float64frombits(bits)
				}
			case compiler.StringA:
			case compiler.StringU:
			case compiler.Bit:
				param.Type = datatypes.Bit
				val := code[e]
				e++

				if val == 0x0 {
					param.Value = false
				} else {
					param.Value = true
				}
			case compiler.Ptr:
				param.Type = datatypes.Ptr
				param.Value = getNextUint64(code, &e)
			case compiler.Label:
				param.Type = datatypes.Label
				param.Value = getNextUint64(code, &e)
			case compiler.Variable:
				param.Type = datatypes.Variable
				param.Value = getNextUint64(code, &e)
			}

			statements = append(statements, command{opcode, param})
		}
	}

	return 0
}
