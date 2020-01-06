package db

import "strings"

type Database struct {
	Name string
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
