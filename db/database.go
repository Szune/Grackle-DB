package db

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Database struct {
	Name   string
	Tables []*Table
}

func (d Database) GetTable(name string) *Table {
	for i := range d.Tables {
		if strings.ToUpper(d.Tables[i].Name) == strings.ToUpper(name) {
			return d.Tables[i]
		}
	}
	return nil
}

func Print(db *Database) {
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
