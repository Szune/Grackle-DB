package utils

import (
	"encoding/binary"
)

type ResultSet []map[string]string

func Int32ToBytes(num int32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(num))
	return bytes
}

func Int64ToBytes(num int64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(num))
	return bytes
}

func StrToBytes(str string) []byte {
	return []byte(str)
}

func BytesToInt32(bytes []byte) int32 {
	return int32(binary.BigEndian.Uint32(bytes))
}

func BytesToInt64(bytes []byte) int64 {
	return int64(binary.BigEndian.Uint64(bytes))
}

func BytesToStr(bytes []byte) string {
	return string(bytes)
}
