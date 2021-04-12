package querying

import (
	"fmt"
	"grackle/db"
	"grackle/utils"
	"strings"
)

type InstructionAndTable struct {
	Instruction utils.Instruction
	Table       *db.Table
}

func ExecuteQuery(query string, database *db.Database) ([]utils.ResultSet, error) {
	tokens, err := GetTokens(query)
	if err != nil {
		return nil, err
	}
	instructions, err := parse(tokens)
	if err != nil {
		return nil, err
	}
	instructionsWithTables, err := getTables(instructions, database)
	if err != nil {
		return nil, err
	}

	return executeInstructions(instructionsWithTables)
}

func getTables(instructions []utils.Instruction, database *db.Database) ([]InstructionAndTable, error) {
	var tables []InstructionAndTable
	for i := range instructions {
		table := database.GetTable(instructions[i].Table)
		if table != nil {
			tables = append(tables, InstructionAndTable{Instruction: instructions[i], Table: table})
		} else {
			return nil, fmt.Errorf("unknown table '%v'", instructions[i].Table)
		}
	}
	return tables, nil
}

func merge(a map[string]string, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}

	return a
}

func join(a []map[string]string, b []map[string]string, joiner string) []map[string]string {
	for x := range a {
		for y := range b {
			if a[x][joiner] == b[y][joiner] {
				for k, v := range b[y] {
					a[x][k] = v
				}
				break
			}
		}
	}

	return a
}

func stringListToSet(columns []string) (set map[string]struct{}) {
	set = map[string]struct{}{}
	for i := range columns {
		set[strings.ToUpper(columns[i])] = struct{}{}
	}
	return set
}

func executeInstructions(instructions []InstructionAndTable) (rows []utils.ResultSet, err error) {
	for i := range instructions {
		instr := instructions[i]
		// TODO: clean this up
		switch instr.Instruction.Op {
		case utils.SelectOp:
			if instr.Instruction.SelectOrInsertColumns[0] == "*" {
				rows = append(rows, instr.Table.GetAll())
			} else {
				rows = append(rows, instr.Table.GetMultipleColumns(stringListToSet(instr.Instruction.SelectOrInsertColumns)))
			}
			break
		case utils.SelectWhereOp:
			if instr.Instruction.SelectOrInsertColumns[0] == "*" {
				rows = append(rows, instr.Table.GetAllWhere(instr.Instruction.Filters))
			} else {
				rows = append(rows, instr.Table.GetMultipleColumnsWhere(instr.Instruction.Filters, stringListToSet(instr.Instruction.SelectOrInsertColumns)))
			}
			break
		case utils.InsertOp:
			// TODO: add validation by checking against the schema
			instr.Table.Insert(&db.Row{
				Values: instr.Instruction.InsertValues,
			})
			break
		}
	}
	return rows, nil
}
