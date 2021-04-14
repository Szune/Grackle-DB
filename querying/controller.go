package querying

import (
	"fmt"
	"grackle/db"
	"grackle/types"
	"grackle/utils"
	"math"
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
			if len(instr.Instruction.InsertValues) != len(instr.Instruction.SelectOrInsertColumns) {
				return nil, fmt.Errorf("different amount of columns compared to values on insert command: %v, %v", len(instr.Instruction.SelectOrInsertColumns), len(instr.Instruction.InsertValues))
			}
			values := make([][]byte, len(instr.Table.Schema))
			for i1, v1 := range instr.Instruction.SelectOrInsertColumns {
				set := false
				for i2, v2 := range instr.Table.Schema {
					if strings.EqualFold(v2.Name, v1) {
						converted, err := convertTypes(instr.Instruction.InsertValues[i1].Value, v2.Name, instr.Instruction.InsertValues[i1].Type, v2.ColumnType)
						if err != nil {
							return nil, err
						} else {
							values[i2] = converted
							set = true
						}
					}
				}
				if !set {
					return nil, fmt.Errorf("column '%s' does not exist in table '%s'", v1, instr.Table.Name)
				}
			}

			instr.Table.Insert(&db.Row{
				Values: values,
			})
			break
		}
	}
	return rows, nil
}

func convertTypes(value []byte, columnName string, from types.ColumnType, to types.ColumnType) ([]byte, error) {
	if from == to {
		return value, nil
	}
	if from == types.Int32 && to == types.Int64 {
		before := utils.BytesToInt32(value)
		return utils.Int64ToBytes(int64(before)), nil
	} else if from == types.Int64 && to == types.Int32 {
		before := utils.BytesToInt64(value)
		if before > math.MaxInt32 || before < math.MinInt32 {
			return nil, fmt.Errorf("value out of range on insert, column '%s' expected type '%s', received '%s' (%v)", columnName, to.String(), from.String(), before)
		}
		return utils.Int32ToBytes(int32(before)), nil
	}
	return nil, fmt.Errorf("wrong type on insert, column '%s' expected type '%s', received '%s'", columnName, to.String(), from.String())
}
