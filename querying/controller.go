package querying

import (
	"../db"
	"errors"
)


// change the name of this entire thing


func ExecuteQuery(tokens []Token, database *db.Database) ([]map[string]string, error) {
	var result []map[string]string
	switch tokens[0].Type {
	case Select:
		result = executeSelect(tokens, database)

	default:
		result = nil
	}

	if result == nil {
		return nil, errors.New("query failed")
	}

	return result, nil
}
