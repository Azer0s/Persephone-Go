package datatypes

//DataType ...Type for enum for Persephone variable types
type DataType uint16

// Enum for Persephone variable types
const (
	Label DataType = 0x0
	Bit   DataType = 0x01
	Ptr   DataType = 0x02

	Int8  DataType = 0x08
	Int16 DataType = 0x16
	Int32 DataType = 0x32
	Int64 DataType = 0x64

	Float32 DataType = 0x132
	Float64 DataType = 0x164

	StringASCII   DataType = 0x201
	StringUnicode DataType = 0x202

	Variable DataType = 0xFFFF
)

//Data ...Type for storing variables
type Data struct {
	Value interface{}
	Type  DataType
}
