package db

import (
	"encoding/binary"
	"strconv"
)

type Row struct {
	Id     int64
	Values [][]byte // Values[column] = the column value as an array of bytes
}

func (r *Row) GetValue(c ColumnType, i int) string {
	switch c {
	case Int32:
		return strconv.Itoa(int(binary.LittleEndian.Uint32(r.Values[i])))
	case Int64:
		return strconv.Itoa(int(binary.LittleEndian.Uint64(r.Values[i])))
	case String:
		return string(r.Values[i])
	default:
		return "--CORRUPT VALUE--\n"
	}
}
