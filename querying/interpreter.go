package querying
import "../db"


// change the name of this entire thing


func Interpret(tokens []Token, database *db.Database) []map[string]string {
	switch tokens[0].Type {
	case Select:
		return executeSelect(tokens, database)

	default:
		return nil
	}
}

func executeSelect(tokens []Token, database *db.Database) []map[string]string {
	if len(tokens) < 4 {
		return nil
	}
	// this is not important for now, I'm more interested in the part where this actually gets the data
	getting := tokens[1]
	if getting.Type != Identifier && getting.Type != Asterisk {
		return nil
	}
	if tokens[2].Type != From {
		return nil
	}
	from := tokens[3]
	if tokens[3].Type != Identifier {
		return nil
	}

	table := database.GetTable(from.String)
	if table == nil {
		return nil
	}
	if getting.Type == Asterisk {
		return table.GetAll()
	} else {
		return table.GetColumn(getting.String)
	}
}
