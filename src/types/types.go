package types

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
