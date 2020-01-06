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
	// pre-check to see if the column even exists in the schema
	columnIndex := -1
	for j := range t.Schema {
		if strings.ToUpper(t.Schema[j].Name) == strings.ToUpper(s) {
			columnIndex = j
		}
	}
	if columnIndex == -1 {
		return nil
	}
	columnName := t.Schema[columnIndex].Name
	columnType := t.Schema[columnIndex].ColumnType

	var rows []map[string]string
	for i := range t.Rows {
		row := map[string]string{
			columnName: t.Rows[i].GetValue(columnType, columnIndex),
		}
		rows = append(rows, row)
	}
	return rows
}

func (t *Table) GetMultipleColumns(columns map[string]struct{}) []map[string]string {
	// pre-check to see if the column even exists in the schema
	var cols []Column
	for i := range t.Schema {
		if _, ok := columns[strings.ToUpper(t.Schema[i].Name)]; ok {
			cols = append(cols, t.Schema[i])
		}
	}
	if len(cols) < 1 {
		return nil
	}

	var rows []map[string]string
	for i := range t.Rows {
		row := map[string]string{}
		for j := range cols {
			row[cols[j].Name] = t.Rows[i].GetValue(cols[j].ColumnType, cols[j].Index)
		}
		rows = append(rows, row)
	}
	return rows
}
