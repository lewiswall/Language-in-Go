package tokenizer

// enum for Token kind

type TokenKind int

const (
	End TokenKind = iota //enum index 0
	Identifier
	String
	Int
	Decimal
	Add
	Subtract
	Divide
	Multiply
	Openbrack
	Closebrack
	Exspo
	Bool
	Unary
	BooleanOp
	BoolConnector
	Assign
	EndOfStatment
	Print
	If
	While
	BlockStart
	BlockEnd
	Input
	Del
)

func TKString(tK TokenKind) string {
	return [...]string{"End", "Identifier", "String", "Int", "Decimal", "Add", "Subtract", "Divide", "Multiply", "Openbrack",
		"Closebrack", "Exspo", "Bool", "Unary", "BooleanOp", "BoolConnector", "Assign", "EndOfStatement", "Print", "If", "While",
		"BlockStart", "BlockEnd", "Input", "Del"}[tK]
}
