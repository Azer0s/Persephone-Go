package bytecoderuntime

func getUint64FromBytes(a, b, c, d, e, f, g, h byte) uint64 {
	return (uint64(a) << 56) + (uint64(b) << 48) + (uint64(c) << 40) + (uint64(d) << 32) + (uint64(e) << 24) + (uint64(f) << 16) + (uint64(g) << 8) + uint64(h)
}

func getUint16FromBytes(a, b byte) uint16 {
	return (uint16(a) << 8) + uint16(b)
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

	for e := 0; e < len(code); e++ {

	}

	return 0
}
