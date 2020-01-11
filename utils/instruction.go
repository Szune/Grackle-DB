package utils

type Operation int8

const (
	SelectOp Operation = iota
	SelectWhereOp
)

type Filter struct {
	Column string
	Equals bool
	Value  []byte
}

type Instruction struct {
	Op            Operation
	Table         string
	SelectColumns []string
	Filters       []Filter
}
