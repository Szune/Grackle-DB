package querying

import (
	"fmt"
	"grackle/types"
	"grackle/utils"
	"strings"
)

type parserState struct {
	buf    [2]Token
	pos    int
	eot    bool
	tokens []Token
}

func parse(tokens []Token) (instructions []utils.Instruction, err error) {
	state := &parserState{
		buf:    [2]Token{{}, {}},
		pos:    0,
		eot:    false,
		tokens: tokens,
	}
	state.eat()
	state.eat()
	for !state.eot {
		t := state.buf[0]
		switch t.Type {
		case Select:
			instr, err := parseSelect(state)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, *instr)
			break
		case Delete:
			// TODO: parse delete and delete where query
			state.eat()
			break
		case Update:
			// TODO: parse update and update where query
			state.eat()
			break
		case Insert:
			instr, err := parseInsert(state)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, *instr)
			break
		case Pipe:
			state.eat()
		default:
			return nil, fmt.Errorf("unknown token %v", state.buf[0].Type)
		}
	}
	return instructions, err
}

func parseSelect(s *parserState) (*utils.Instruction, error) {
	s.eat()
	column := s.buf[0]
	if column.Type != Identifier && column.Type != Asterisk {
		return nil, fmt.Errorf("failed to parse select")
	}
	s.eat()
	var selectInstruction *utils.Instruction
	if s.buf[0].Type == Comma {
		// get multiple columns, e.g. select x, y from z
		multiColumns := []string{
			strings.ToUpper(column.String),
		}
		for s.buf[0].Type == Comma {
			s.eat()

			if s.buf[0].Type != Identifier {
				return nil, fmt.Errorf("failed to parse select: expected column name(s)")
			}
			multiColumns = append(multiColumns, strings.ToUpper(s.buf[0].String))
			s.eat()
		}

		if s.buf[0].Type != From {
			return nil, fmt.Errorf("failed to parse select: expected 'from'")
		}
		s.eat()

		if s.buf[0].Type != Identifier {
			return nil, fmt.Errorf("failed to parse select: expected table name")
		}
		table := s.buf[0].String
		selectInstruction = &utils.Instruction{
			Op:                    utils.SelectOp,
			Table:                 table,
			SelectOrInsertColumns: multiColumns,
		}
	} else {
		if s.buf[0].Type != From {
			return nil, fmt.Errorf("failed to parse select: expected 'from'")
		}
		s.eat()
		if s.buf[0].Type != Identifier {
			return nil, fmt.Errorf("failed to parse select: expected table name")
		}
		table := s.buf[0].String
		var columns []string

		if column.Type == Asterisk {
			columns = []string{"*"}
		} else {
			columns = []string{column.String}
		}
		selectInstruction = &utils.Instruction{
			Op:                    utils.SelectOp,
			Table:                 table,
			SelectOrInsertColumns: columns,
		}
	}

	s.eat()
	if s.buf[0].Type != Where {
		return selectInstruction, nil
	} else {
		selectWhere, err := parseWhere(s, selectInstruction)
		return selectWhere, err
	}
}

func parseInsert(s *parserState) (*utils.Instruction, error) {
	s.eat()
	if s.buf[0].Type != Into {
		return nil, fmt.Errorf("failed to parse insert: missing 'into'")
	}
	s.eat()
	table := s.buf[0]
	if table.Type != Identifier {
		return nil, fmt.Errorf("failed to parse insert: expected table name after 'insert into'")
	}
	s.eat()
	if s.buf[0].Type != LeftParenthesis {
		return nil, fmt.Errorf("failed to parse insert: expected '(' after table name")
	}
	s.eat()

	// parse columns
	columns := make([]string, 0, 2)

	col := s.buf[0]
	if col.Type != Identifier {
		return nil, fmt.Errorf("failed to parse insert: expected column name after 'insert into tablename('")
	}
	columns = append(columns, col.String)
	s.eat()

	for s.buf[0].Type == Comma {
		s.eat()
		col = s.buf[0]
		if s.buf[0].Type != Identifier {
			return nil, fmt.Errorf("failed to parse insert: expected column name after 'insert into tablename('")
		}
		columns = append(columns, col.String)
		s.eat()
	}

	if s.buf[0].Type != RightParenthesis {
		return nil, fmt.Errorf("failed to parse insert: expected ')' after 'insert into tablename([...]'")
	}
	s.eat()

	if s.buf[0].Type != Values {
		return nil, fmt.Errorf("failed to parse insert: expected 'values' after 'insert into tablename([...])'")
	}
	s.eat()

	if s.buf[0].Type != LeftParenthesis {
		return nil, fmt.Errorf("failed to parse insert: expected '(' after 'insert into tablename([...]) values'")
	}
	s.eat()

	// parse values
	values := make([]utils.QueryValue, 0, 2)

	value := s.buf[0]
	switch value.Type {
	case String:
		values = append(values, utils.QueryValue{Value: utils.StrToBytes(value.String), Type: types.String})
		break
	case Number:
		values = append(values, utils.QueryValue{Value: utils.Int64ToBytes(value.Number), Type: types.Int64})
		break
	case Parameter:
		values = append(values, utils.QueryValue{Value: utils.StrToBytes("@" + value.String)})
		break
	default:
		return nil, fmt.Errorf("failed to parse insert: expected value after 'insert into tablename([...]) values('")
	}
	s.eat()

	for s.buf[0].Type == Comma {
		s.eat()
		value = s.buf[0]
		switch value.Type {
		case String:
			values = append(values, utils.QueryValue{Value: utils.StrToBytes(value.String), Type: types.String})
			break
		case Number:
			values = append(values, utils.QueryValue{Value: utils.Int64ToBytes(value.Number), Type: types.Int64})
			break
		case Parameter:
			values = append(values, utils.QueryValue{Value: utils.StrToBytes("@" + value.String)})
			break
		default:
			return nil, fmt.Errorf("failed to parse insert: expected value after 'insert into tablename([...]) values('")
		}
		s.eat()
	}

	if s.buf[0].Type != RightParenthesis {
		return nil, fmt.Errorf("failed to parse insert: expected ')' after 'insert into tablename([...]) values([...]'")
	}
	s.eat()

	insertInstruction := &utils.Instruction{
		Op:                    utils.InsertOp,
		Table:                 table.String,
		SelectOrInsertColumns: columns,
		InsertValues:          values,
	}
	s.eat()
	return insertInstruction, nil
}

func parseWhere(s *parserState, selectInstr *utils.Instruction) (*utils.Instruction, error) {
	if s.buf[0].Type != Where {
		return nil, fmt.Errorf("expected 'where'")
	}
	s.eat()
	var filters []utils.Filter
	for s.buf[0].Type == Identifier {
		filt, err := parseComparison(s, filters)
		if err != nil {
			return nil, err
		}
		filters = filt

		if s.buf[0].Type == Comma {
			s.eat()
		}
	}

	whereInstruction := &utils.Instruction{
		Op:                    utils.SelectWhereOp,
		Table:                 selectInstr.Table,
		SelectOrInsertColumns: selectInstr.SelectOrInsertColumns,
		Filters:               filters,
	}

	return whereInstruction, nil
}

func parseComparison(s *parserState, filters []utils.Filter) ([]utils.Filter, error) {
	if s.buf[0].Type != Identifier {
		return nil, fmt.Errorf("expected identifier")
	}
	column := s.buf[0].String
	s.eat()

	equals := false
	if s.buf[0].Type == Equals {
		equals = true
	} else {
		return nil, fmt.Errorf("expected '='")
	}
	s.eat()

	var value []byte
	if s.buf[0].Type == String {
		value = utils.StrToBytes(s.buf[0].String)
	} else if s.buf[0].Type == Number {
		// TODO: support int64 as well
		value = utils.Int32ToBytes(int32(s.buf[0].Number))
	}
	s.eat()

	filter := utils.Filter{
		Column: column,
		Equals: equals,
		Value:  value,
	}
	filters = append(filters, filter)
	return filters, nil
}

func (s *parserState) eat() {
	s.buf[0] = s.buf[1]

	if s.pos < len(s.tokens) {
		s.buf[1] = s.tokens[s.pos]
	}
	if s.pos > len(s.tokens) {
		s.eot = true
	}
	s.pos++
}
