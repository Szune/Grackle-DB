package querying

type TokenType int8

const (
	Identifier TokenType = iota
	Select
	Where
	Insert
	Into
	Values
	Update
	Delete
	Parameter
	String
	Number
	From
	Equals
	Asterisk
	Comma
	Pipe
	LeftParenthesis
	RightParenthesis
)

func (t TokenType) String() string {
	switch t {
	case Identifier:
		return "Identifier"
	case Select:
		return "Select"
	case Where:
		return "Where"
	case Insert:
		return "Insert"
	case Into:
		return "Into"
	case Values:
		return "Values"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	case Parameter:
		return "Parameter"
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
	case Pipe:
		return "Pipe"
	case LeftParenthesis:
		return "LeftParenthesis"
	case RightParenthesis:
		return "RightParenthesis"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type   TokenType
	String string
	Number int64
}
