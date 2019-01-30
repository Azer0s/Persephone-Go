package compiler

import (
	"../types"
	"os"
)

var opcodes = map[string]uint16{
	"v_int8"	: uint16(0x0118),
	"v_int16"	: uint16(0x0110),
	"v_int32"	: uint16(0x0120),
	"v_int"		: uint16(0x0120),
	"v_int64"	: uint16(0x0140),
	"v_float32"	: uint16(0x0121),
	"v_float"	: uint16(0x0121),
	"v_float64"	: uint16(0x0141),
	"v_double"	: uint16(0x0141),
	"v_stringa"	: uint16(0x0131),
	"v_stringu"	: uint16(0x0132),
	"v_ptr"		: uint16(0x0150),
	"v_bit"		: uint16(0x0100),
	"dci8"		: uint16(0x0218),
	"dci16"		: uint16(0x0210),
	"dci32"		: uint16(0x0220),
	"dci"		: uint16(0x0220),
	"dci64"		: uint16(0x0240),
	"dcf32"		: uint16(0x0221),
	"dcf"		: uint16(0x0221),
	"dcf64"		: uint16(0x0241),
	"dcd"		: uint16(0x0241),
	"dcsa"		: uint16(0x0231),
	"dcsu"		: uint16(0x0232),
	"dcb"		: uint16(0x0200),
	"ldi8v"		: uint16(0x0318),
	"ldi16v"	: uint16(0x0310),
	"ldi32v"	: uint16(0x0320),
	"ldiv"		: uint16(0x0320),
	"ldi64v"	: uint16(0x0340),
	"ldf32v"	: uint16(0x0321),
	"ldfv"		: uint16(0x0321),
	"ldf64v"	: uint16(0x0341),
	"lddv"		: uint16(0x0341),
	"ldsav"		: uint16(0x0331),
	"ldsuv"		: uint16(0x0332),
	"ldptrv"	: uint16(0x0350),
	"ldbv"		: uint16(0x0300),
	"ldi8c"		: uint16(0x0418),
	"ldi16c"	: uint16(0x0410),
	"ldi32c"	: uint16(0x0420),
	"ldic"		: uint16(0x0420),
	"ldi64c"	: uint16(0x0440),
	"ldf32c"	: uint16(0x0421),
	"ldfc"		: uint16(0x0421),
	"ldf64c"	: uint16(0x0441),
	"lddc"		: uint16(0x0441),
	"ldsac"		: uint16(0x0431),
	"ldsuc"		: uint16(0x0432),
	"ldbc"		: uint16(0x0400),
	"pop"		: uint16(0x0001),
	"ret"		: uint16(0x0002),
	"add"		: uint16(0x0003),
	"sub"		: uint16(0x0004),
	"mul"		: uint16(0x0005),
	"div"		: uint16(0x0006),
	"mod"		: uint16(0x0007),
	"ge"		: uint16(0x0008),
	"le"		: uint16(0x0009),
	"gt"		: uint16(0x000A),
	"lt"		: uint16(0x000B),
	"andi"		: uint16(0x000C),
	"ori"		: uint16(0x000D),
	"xori"		: uint16(0x000E),
	"noti"		: uint16(0x000F),
	"shl"		: uint16(0x0010),
	"shr"		: uint16(0x0011),
	"inc"		: uint16(0x0012),
	"dec"		: uint16(0x0013),
	"addf"		: uint16(0x0014),
	"subf"		: uint16(0x0015),
	"mulf"		: uint16(0x0016),
	"divf"		: uint16(0x0017),
	"gef"		: uint16(0x0018),
	"lef"		: uint16(0x0019),
	"gtf"		: uint16(0x001A),
	"ltf"		: uint16(0x001B),
	"and"		: uint16(0x001C),
	"or"		: uint16(0x001D),
	"xor"		: uint16(0x001E),
	"not"		: uint16(0x001F),
	"conc"		: uint16(0x0020),
	"len"		: uint16(0x0021),
	"getc"		: uint16(0x0022),
	"setc"		: uint16(0x0023),
	"syscall"	: uint16(0x0024),
	"call"		: uint16(0xF000),
	"jmp"		: uint16(0xF001),
	"jmpt"		: uint16(0xF002),
	"jmpf"		: uint16(0xF003),
}

func Compile(root types.Root, outname string) int{

	f, err := os.Create(outname)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	for e := 0; e < len(root.Commands); e++ {
		isJmp := false

		switch root.Commands[e].Command.Text {
		case "call","jmp","jmpt","jmpf":
			isJmp = true
		}

		_,err := f.Write([]byte{
			byte(opcodes[root.Commands[e].Command.Text] & 0xFF00), //Get upper 8 bits
			byte(opcodes[root.Commands[e].Command.Text] & 0x00FF)}) //Get lower 8 bits

		if err != nil {
			panic(err)
		}

		if isJmp && root.Commands[e].Param.Kind == types.Name{
			//Label name -> int64
		}else{
			switch root.Commands[e].Param.Kind {

			}
		}
	}

	return 0
}
