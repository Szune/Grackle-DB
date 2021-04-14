package utils

import (
	"grackle/types"
)

type Operation int8

const (
	SelectOp Operation = iota
	SelectWhereOp
	DeleteOp
	DeleteWhereOp
	InsertOp
	UpdateOp
	UpdateWhereOp
)

type Filter struct {
	Column string
	Equals bool
	Value  []byte
}

type QueryValue struct {
	Value []byte
	Type  types.ColumnType
}

type Instruction struct {
	Op                    Operation
	Table                 string
	SelectOrInsertColumns []string
	Filters               []Filter
	InsertValues          []QueryValue
}
