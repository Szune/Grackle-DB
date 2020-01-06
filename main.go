package main

import (
	"./db"
	ql "./querying"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	fmt.Printf(".--------------------.\n")
	fmt.Printf("|     Grackle DB     |\n")
	fmt.Printf("'--------------------'\n")
	table := &db.Table{
		Name: "DaysOfWeek",
		Schema: []db.Column{
			{
				Name:       "Number",
				ColumnType: db.Int32,
			},
			{
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
				Name: "Number",
				ColumnType: db.Int32,
			},
			{
				Name: "Name",
				ColumnType: db.String,
			},
		},
		Rows: []db.Row{},
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
			db.Int32ToBytes(1),
			db.StrToBytes("Monday"),
		},
	})
	table.Insert(&db.Row{
		Values: [][]byte{
			db.Int32ToBytes(2),
			db.StrToBytes("Tuesday"),
		},
	})
	table.Insert(&db.Row{
		Values: [][]byte{
			db.Int32ToBytes(3),
			db.StrToBytes("Wednesday"),
		},
	})

	table2.Insert(&db.Row{
		Values:[][]byte{
			db.Int32ToBytes(1),
			db.StrToBytes("January"),
		},
	})
	table2.Insert(&db.Row{
		Values:[][]byte{
			db.Int32ToBytes(2),
			db.StrToBytes("February"),
		},
	})

	db.PrintDb(database)

	table.Update(2, &db.Row{
		Values: [][]byte{
			db.Int32ToBytes(10),
			db.StrToBytes("Scrambleday"),
		},
	})

	db.PrintDb(database)

	fmt.Printf("Waiting for command...\n")
	fmt.Printf("Enter query (quit with q):\n")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if text := scanner.Text(); text == "q" {
			break
		} else {
			fmt.Printf("Input: %v\n", text)
			fmt.Printf("Executing query...\n")
			tokens := ql.GetTokens(text)
			rows := ql.Interpret(tokens, database)
			fmt.Printf("Returned %v rows:\n", len(rows))
			serialized, err := json.MarshalIndent(rows, "", "\t")
			if err == nil {
				fmt.Println(string(serialized))
			}
			/*for i := range rows {
				fmt.Printf("Row %v\n", i)
				for k, v := range rows[i] {
					fmt.Printf("\t%v: %v\n", k, v)
				}
			}*/
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Grackle DB exited")
}

