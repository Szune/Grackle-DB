package querying

import (
	"../db"
	"strings"
)

func executeSelect(tokens []Token, database *db.Database) []map[string]string {
	if len(tokens) < 4 { // [1]select [2]* [3]from [4]<table>
		return nil
	}

	tPos := 1
	column := tokens[tPos]
	if column.Type != Identifier && column.Type != Asterisk {
		return nil
	}
	tPos++
	if tokens[tPos].Type == Comma {
		// get multiple columns, e.g. select x, y from z
		multiColumns := map[string]struct{}{} // closest to a set afaik
		multiColumns[strings.ToUpper(column.String)] = struct{}{}
		for tokens[tPos].Type == Comma {
			tPos++

			if tokens[tPos].Type != Identifier {
				return nil
			}
			multiColumns[strings.ToUpper(tokens[tPos].String)] = struct{}{}
			tPos++
		}
		table := GetTable(tokens, database, tPos)
		if table == nil {
			return nil
		}
		return table.GetMultipleColumns(multiColumns)
	} else {
		table := GetTable(tokens, database, 2)
		if table == nil {
			return nil
		}
		if column.Type == Asterisk {
			return table.GetAll()
		} else {
			return table.GetColumn(column.String)
		}
	}
}

func GetTable(tokens []Token, database *db.Database, pos int) *db.Table {
	if tokens[pos].Type != From {
		return nil
	}
	pos++
	tableName := tokens[pos]
	if tableName.Type != Identifier {
		return nil
	}

	table := database.GetTable(tableName.String)
	if table == nil {
		return nil
	}
	return table
}
