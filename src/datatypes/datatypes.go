package datatypes

//DataType ...Type for enum for Persephone variable types
type DataType uint16

// Enum for Persephone variable types
const (
	Label DataType = 0x0
	Bit   DataType = 0x01

	Int8  DataType = 0x08
	Int16 DataType = 0x16
	Int32 DataType = 0x32
	Int64 DataType = 0x64

	Ptr    DataType = 0x80
	Uint8  DataType = 0x081
	Uint16 DataType = 0x161
	Uint32 DataType = 0x132
	Uint64 DataType = 0x164

	Float32 DataType = 0x332
	Float64 DataType = 0x364

	StringASCII   DataType = 0x201
	StringUnicode DataType = 0x202

	Dynamic  DataType = 0xFFFE
	Variable DataType = 0xFFFF
)

//Data ...Type for storing variables
type Data struct {
	Value   interface{}
	Type    DataType
	Dynamic bool
}
