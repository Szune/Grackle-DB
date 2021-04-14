package types

type ColumnType byte

const (
	Int32 ColumnType = iota
	Int64
	String
)

func (c ColumnType) String() string {
	switch c {
	case Int32:
		return "Int32"
	case Int64:
		return "Int64"
	case String:
		return "String"
	default:
		return "Unknown"
	}
}
