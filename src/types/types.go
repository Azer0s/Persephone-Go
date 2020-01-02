package types

//Op operation type
type Op int8

//Operation types
const (
	Add Op = 0x0
	Sub Op = 0x1
	Div Op = 0x2
	Mul Op = 0x3
	Mod Op = 0x4
	Shl Op = 0x5
	Shr Op = 0x6

	And Op = 0x10
	Or  Op = 0x11
	Xor Op = 0x12
	Not Op = 0x13

	Le Op = 0x20
	Ge Op = 0x21
	L  Op = 0x22
	G  Op = 0x23

	Inc Op = 0x30
	Dec Op = 0x31
)

//Print syscall
const (
	Fork    Op = 0x0
	Print   Op = 0x1
	Read    Op = 0x2
	Println Op = 0x10
	Readln  Op = 0x20
)

//Kind AST kinds
type Kind string

//AST kinds
const (
	Name      Kind = "name"
	Number    Kind = "number"
	String    Kind = "string"
	HexNumber Kind = "hexnum"
	Float     Kind = "float"
	Pointer   Kind = "pointer"
	Label     Kind = "label"
	Unknown   Kind = "unknown"
	Bit       Kind = "bit"
)

//Token Lexeme
type Token struct {
	Kind Kind
	Size string
	Text string
}

//Root AST root node
type Root struct {
	Commands []Command
	Labels   map[string]int
}

//Function function node
type Function struct {
	Name     Token
	Commands []Command
}

//Command Command node
type Command struct {
	Single         bool
	Command, Param Token
}
