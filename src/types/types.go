package types

type Op int8

const(
	Add Op = 0x0
	Sub Op = 0x1
	Div Op = 0x2
	Mul Op = 0x3
	Mod Op = 0x4
	Shl Op = 0x5
	Shr Op = 0x6

	And Op = 0x10
	Or Op = 0x11
	Xor Op = 0x12
	Not Op = 0x13
)

const(
	Print Op = 0x1
)

type Kind string

const(
	Name Kind = "name"
	Number Kind = "number"
	String Kind = "string"
	HexNumber Kind = "hexnum"
	Float Kind = "float"
	Pointer Kind = "pointer"
	Lbrace Kind = "lbrace"
	Rbrace Kind = "rbrace"
	Label Kind = "label"
	Unknown Kind = "unknown"
)

type Token struct {
	Kind Kind
	Text string
}

type Root struct{
	Functions []Function
	Commands []Command
	Labels map[string]int
}

type Function struct{
	Name Token
	Commands []Command
}

type Command struct {
	Single bool
	Command, Param Token
}
