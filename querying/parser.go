package querying

import (
	"fmt"
	"grackle/utils"
	"strings"
)

type parserState struct {
	buf [2]Token
	pos int
	eot bool
}

func parse(tokens []Token) (instructions []utils.Instruction, err error) {
	state := &parserState{
		buf: [2]Token{{}, {}},
		pos: 0,
		eot: false,
	}
	consumeToken(state, tokens)
	consumeToken(state, tokens)
	for !state.eot {
		t := state.buf[0]
		switch t.Type {
		case Select:
			instr, err := parseSelect(state, tokens)
			if err != nil {
				return nil, err
			}
			instructions = append(instructions, *instr)
			break
		case Pipe:
			consumeToken(state, tokens)
		default:
			return nil, fmt.Errorf("unknown token %v", state.buf[0].Type)
		}
	}
	return instructions, err
}

func parseSelect(s *parserState, tokens []Token) (*utils.Instruction, error) {
	consumeToken(s, tokens)
	column := s.buf[0]
	if column.Type != Identifier && column.Type != Asterisk {
		return nil, fmt.Errorf("failed to parse select")
	}
	consumeToken(s, tokens)
	if s.buf[0].Type == Comma {
		// get multiple columns, e.g. select x, y from z
		multiColumns := []string{
			strings.ToUpper(column.String),
		}
		for s.buf[0].Type == Comma {
			consumeToken(s, tokens)

			if s.buf[0].Type != Identifier {
				return nil, fmt.Errorf("failed to parse select: expected column name(s)")
			}
			multiColumns = append(multiColumns, strings.ToUpper(s.buf[0].String))
			consumeToken(s, tokens)
		}

		if s.buf[0].Type != From {
			return nil, fmt.Errorf("failed to parse select: expected 'from'")
		}
		consumeToken(s, tokens)

		if s.buf[0].Type != Identifier {
			return nil, fmt.Errorf("failed to parse select: expected table name")
		}
		table := s.buf[0].String
		selectInstruction := &utils.Instruction{
			Op:            utils.SelectOp,
			Table:         table,
			SelectColumns: multiColumns,
		}
		consumeToken(s, tokens)
		if s.buf[0].Type != Where {
			return selectInstruction, nil
		} else {
			selectWhere, err := parseWhere(s, tokens, selectInstruction)
			return selectWhere, err
		}
	} else {
		if s.buf[0].Type != From {
			return nil, fmt.Errorf("failed to parse select: expected 'from'")
		}
		consumeToken(s, tokens)
		if s.buf[0].Type != Identifier {
			return nil, fmt.Errorf("failed to parse select: expected table name")
		}
		table := s.buf[0].String
		var selectInstruction *utils.Instruction

		if column.Type == Asterisk {
			selectInstruction = &utils.Instruction{
				Op:            utils.SelectOp,
				Table:         table,
				SelectColumns: []string{"*"},
			}
		} else {
			selectInstruction = &utils.Instruction{
				Op:            utils.SelectOp,
				Table:         table,
				SelectColumns: []string{column.String},
			}
		}

		consumeToken(s, tokens)
		if s.buf[0].Type != Where {
			return selectInstruction, nil
		} else {
			selectWhere, err := parseWhere(s, tokens, selectInstruction)
			return selectWhere, err
		}

	}
}

func parseWhere(s *parserState, tokens []Token, selectInstr *utils.Instruction) (*utils.Instruction, error) {
	if s.buf[0].Type != Where {
		return nil, fmt.Errorf("expected 'where'")
	}
	consumeToken(s, tokens)
	var filters []utils.Filter
	for s.buf[0].Type == Identifier {
		filt, err := parseComparison(s, tokens, filters)
		if err != nil {
			return nil, err
		}
		filters = filt

		if s.buf[0].Type == Comma {
			consumeToken(s, tokens)
		}
	}

	whereInstruction := &utils.Instruction{
		Op:            utils.SelectWhereOp,
		Table:         selectInstr.Table,
		SelectColumns: selectInstr.SelectColumns,
		Filters:       filters,
	}

	return whereInstruction, nil
}

func parseComparison(s *parserState, tokens []Token, filters []utils.Filter) ([]utils.Filter, error) {
	if s.buf[0].Type != Identifier {
		return nil, fmt.Errorf("expected identifier")
	}
	column := s.buf[0].String
	consumeToken(s, tokens)

	equals := false
	if s.buf[0].Type == Equals {
		equals = true
	} else {
		return nil, fmt.Errorf("expected '='")
	}
	consumeToken(s, tokens)

	value := []byte{}
	if s.buf[0].Type == String {
		value = utils.StrToBytes(s.buf[0].String)
	} else if s.buf[0].Type == Number {
		// TODO: support int64 as well
		value = utils.Int32ToBytes(int32(s.buf[0].Number))
	}
	consumeToken(s, tokens)

	filter := utils.Filter{
		Column: column,
		Equals: equals,
		Value:  value,
	}
	filters = append(filters, filter)
	return filters, nil
}

func consumeToken(s *parserState, tokens []Token) {
	s.buf[0] = s.buf[1]

	if s.pos < len(tokens) {
		s.buf[1] = tokens[s.pos]
	}
	if s.pos > len(tokens) {
		s.eot = true
	}
	s.pos++
}
