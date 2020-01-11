package db

type ColumnType int8

const (
	Int32 ColumnType = iota
	Int64
	String
)

type Column struct {
	Index      int
	Name       string
	ColumnType ColumnType
}

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
