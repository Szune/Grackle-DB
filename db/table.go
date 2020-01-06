package db

import "strings"

type Table struct {
	Rows      []Row
	Name      string
	Schema    []Column
	LastRowId int64
}


func (t *Table) Insert(row *Row) {
	t.LastRowId += 1
	row.Id = t.LastRowId
	t.Rows = append(t.Rows, *row)
}

func (t *Table) Update(id int64, row *Row) {
	for i := range t.Rows {
		if t.Rows[i].Id == id {
			row.Id = t.Rows[i].Id // can't change id
			t.Rows[i] = *row
		}
	}
}

func (t *Table) GetAll() []map[string]string {
	var rows []map[string]string
	idx := 0
	for i := range t.Rows {
		rows = append(rows, map[string]string{})
		for j := range t.Schema {
			rows[idx][t.Schema[j].Name] = t.Rows[i].GetValue(t.Schema[j].ColumnType, j)
		}
		idx++
	}
	return rows
}

func (t *Table) GetColumn(s string) []map[string]string {
	var rows []map[string]string
	idx := 0
	for i := range t.Rows {
		rows = append(rows, map[string]string{})
		for j := range t.Schema {
			if strings.ToUpper(t.Schema[j].Name) != strings.ToUpper(s) { // very not optimized, as is most of the code atm
				continue
			}
			rows[idx][t.Schema[j].Name] = t.Rows[i].GetValue(t.Schema[j].ColumnType, j)
		}
		idx++
	}
	return rows
}
