package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"grackle/db"
	ql "grackle/querying"
	"grackle/utils"
	"os"
)

func main() {
	// TODO: select top(10), e.g. add top parameter to select instructions
	fmt.Printf(".--------------------.\n")
	fmt.Printf("|     Grackle DB     |\n")
	fmt.Printf("'--------------------'\n")
	table := &db.Table{
		Name: "DaysOfWeek",
		Schema: []db.Column{
			{
				Index:      0,
				Name:       "Number",
				ColumnType: db.Int32,
			},
			{
				Index:      1,
				Name:       "Name",
				ColumnType: db.String,
			},
		},
		Rows:      []db.Row{},
		LastRowId: 0,
	}

	table2 := &db.Table{
		Name: "Months",
		Schema: []db.Column{
			{
				Index:      0,
				Name:       "Number",
				ColumnType: db.Int32,
			},
			{
				Index:      1,
				Name:       "Name",
				ColumnType: db.String,
			},
			{
				Index:      2,
				Name:       "Season",
				ColumnType: db.String,
			},
		},
		Rows:      []db.Row{},
		LastRowId: 0,
	}

	database := &db.Database{
		Name: "Schedule",
		Tables: []*db.Table{
			table,
			table2,
		},
	}

	table.Insert(&db.Row{
		Values: [][]byte{
			utils.Int32ToBytes(1),
			utils.StrToBytes("Monday"),
		},
	})
	table.Insert(&db.Row{
		Values: [][]byte{
			utils.Int32ToBytes(2),
			utils.StrToBytes("Tuesday"),
		},
	})
	table.Insert(&db.Row{
		Values: [][]byte{
			utils.Int32ToBytes(3),
			utils.StrToBytes("Wednesday"),
		},
	})

	table2.Insert(&db.Row{
		Values: [][]byte{
			utils.Int32ToBytes(1),
			utils.StrToBytes("January"),
			utils.StrToBytes("Winter"),
		},
	})
	table2.Insert(&db.Row{
		Values: [][]byte{
			utils.Int32ToBytes(2),
			utils.StrToBytes("February"),
			utils.StrToBytes("Winter"),
		},
	})
	table2.Insert(&db.Row{
		Values: [][]byte{
			utils.Int32ToBytes(3),
			utils.StrToBytes("March"),
			utils.StrToBytes("Spring"),
		},
	})

	db.Print(database)

	table.Update(2, &db.Row{
		Values: [][]byte{
			utils.Int32ToBytes(10),
			utils.StrToBytes("Scrambleday"),
		},
	})

	db.Print(database)

	fmt.Printf("Waiting for command...\n")
	fmt.Printf("Enter query (quit with q):\n")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if text := scanner.Text(); text == "q" {
			break
		} else {
			fmt.Printf("Input: %v\n", text)
			fmt.Printf("Executing query...\n")
			resultSets, err := ql.ExecuteQuery(text, database)
			if err != nil {
				fmt.Printf("[Error] query failed: %v\n", text)
				continue
			}
			fmt.Printf("Returned %v result sets:\n", len(resultSets))
			serialized, err := json.MarshalIndent(resultSets, "", "\t")
			if err == nil {
				fmt.Println(string(serialized))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Grackle DB exited")
}
