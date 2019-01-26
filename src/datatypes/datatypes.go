package datatypes

type DataType int16
type Op int8

const(
	Add Op = 0x0
	Sub Op = 0x1
	Div Op = 0x2
	Mul Op = 0x3
	Mod Op = 0x4
)

const(
	Int8 DataType = 0x08
	Int16 DataType = 0x16
	Int32 DataType = 0x32
	Int64 DataType = 0x64

	Float32 DataType = 0x132
	Float64 DataType = 0x164
)

type Data struct {
	Value interface{}
	Type DataType
}
