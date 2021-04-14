package db

import (
	"encoding/binary"
	"grackle/types"
	"strconv"
)

type Row struct {
	Id     int64
	Values [][]byte // Values[column] = the column value as an array of bytes
}

func (r *Row) GetValue(c types.ColumnType, i int) string {
	switch c {
	case types.Int32:
		return strconv.Itoa(int(binary.BigEndian.Uint32(r.Values[i])))
	case types.Int64:
		return strconv.Itoa(int(binary.BigEndian.Uint64(r.Values[i])))
	case types.String:
		return string(r.Values[i])
	default:
		return "--CORRUPT VALUE--\n"
	}
}
