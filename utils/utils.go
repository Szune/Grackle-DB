package utils

import (
	"encoding/binary"
)

type ResultSet []map[string]string

func Int32ToBytes(num int32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(num))
	return bytes
}

func Int64ToBytes(num int64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, uint64(num))
	return bytes
}

func StrToBytes(str string) []byte {
	return []byte(str)
}

func BytesToStr(bytes []byte) string {
	return string(bytes)
}
