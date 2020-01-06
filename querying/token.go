package querying

type TokenType int8

const (
	Identifier TokenType = iota
	Select
	String
	Number
	From
	Equals
	Asterisk
	Comma
)

func (t TokenType) String() string {
	switch t {
	case Identifier:
		return "Identifier"
	case Select:
		return "Select"
	case String:
		return "String"
	case Number:
		return "Number"
	case From:
		return "From"
	case Equals:
		return "Equals"
	case Asterisk:
		return "Asterisk"
	case Comma:
		return "Comma"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type TokenType
	String string
	Number int64
}
