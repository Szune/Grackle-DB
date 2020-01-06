package db

import (
	"encoding/binary"
	"fmt"
)

func Int32ToBytes(num int32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(num))
	return bytes
}

func StrToBytes(str string) []byte {
	return []byte(str)
}
func PrintDb(db *Database) {
	fmt.Printf("Database '%v'\n{\n", db.Name)
	for t := range db.Tables {
		fmt.Printf("\tTable '%v'\n\t[\n", db.Tables[t].Name)
		for i := range db.Tables[t].Rows {
			fmt.Printf("\t\tRow (id %v):\n", db.Tables[t].Rows[i].Id)
			for j := range db.Tables[t].Schema {
				fmt.Printf("\t\t\t%v (%v): ", db.Tables[t].Schema[j].Name, db.Tables[t].Schema[j].ColumnType)
				switch db.Tables[t].Schema[j].ColumnType {
				case Int32:
					fmt.Printf("%v\n", int32(binary.LittleEndian.Uint32(db.Tables[t].Rows[i].Values[j])))
				case Int64:
					fmt.Printf("%v\n", int64(binary.LittleEndian.Uint64(db.Tables[t].Rows[i].Values[j])))
				case String:
					fmt.Printf("%v\n", string(db.Tables[t].Rows[i].Values[j]))
				default:
					fmt.Printf("--CORRUPT VALUE--\n")
				}
			}
		}
		fmt.Printf("\t]\n")
	}
	fmt.Printf("}\n")
}
