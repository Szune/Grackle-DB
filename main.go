package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"grackle/db"
	ql "grackle/querying"
	"grackle/types"
	"grackle/utils"
	"os"
)

func main() {
	// TODO: add update query
	// TODO: add delete query
	// TODO: select top(10), e.g. add top parameter to select instructions
	// TODO: io -> handle writing to files on a separate thread or goroutine than the query engine
	// TODO: net -> handle HTTP and other network requests
	fmt.Printf(".--------------------.\n")
	fmt.Printf("|     GrackleDB      |\n")
	fmt.Printf("'--------------------'\n")
	table := &db.Table{
		Name: "DaysOfWeek",
		Schema: []db.Column{
			{
				Index:      0,
				Name:       "Number",
				ColumnType: types.Int32,
			},
			{
				Index:      1,
				Name:       "Name",
				ColumnType: types.String,
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
				ColumnType: types.Int32,
			},
			{
				Index:      1,
				Name:       "Name",
				ColumnType: types.String,
			},
			{
				Index:      2,
				Name:       "Season",
				ColumnType: types.String,
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

	/*
		db.Print(database)

		table.Update(2, &db.Row{
			Values: [][]byte{
				utils.Int32ToBytes(10),
				utils.StrToBytes("Scrambleday"),
			},
		})
	*/

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
				fmt.Printf("[Error] query failed:\n\"%v\"\nError message: %s\n", text, err)
				continue
			}
			if len(resultSets) > 0 {
				fmt.Printf("Returned %v result sets:\n", len(resultSets))
				serialized, err := json.MarshalIndent(resultSets, "", "\t")
				if err == nil {
					fmt.Println(string(serialized))
				}
			} else {
				fmt.Printf("Returned 0 result sets.\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("GrackleDB exited")
}
