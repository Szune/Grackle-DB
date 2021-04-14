package db

import "grackle/types"

type Column struct {
	Index      int
	Name       string
	ColumnType types.ColumnType
}
