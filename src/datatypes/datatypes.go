package datatypes

type DataType int16
type PtrType int8

const (
	Bit DataType = 0x01
	Ptr DataType = 0x02

	Int8  DataType = 0x08
	Int16 DataType = 0x16
	Int32 DataType = 0x32
	Int64 DataType = 0x64

	Float32 DataType = 0x132
	Float64 DataType = 0x164

	String_ASCII   DataType = 0x201
	String_Unicode DataType = 0x202
)

type Data struct {
	Value interface{}
	Type  DataType
}
