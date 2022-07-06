package tree

import (
	"errors"
	"fmt"
	"language/Global"
	"language/LanErrs"
	"language/tokenizer"
	"strconv"
	"unsafe"
)

//Interface used so different nodes can be linked together
type Node interface {
	Evaluate() (Value, error)
}

// Value for storing values
type valueKind int

const (
	Integer valueKind = iota
	Decimal
	str
	Bool
	Identifier
)

type Value struct {
	ValueType valueKind
	Value     uint64
}

//bit casting
func intValue(val int) Value {
	p := unsafe.Pointer(&val)
	v := *(*uint64)(p)
	return Value{ValueType: Integer, Value: v}
}

func DecimalValue(val float64) Value {
	p := unsafe.Pointer(&val)
	f := *(*uint64)(p)
	return Value{ValueType: Decimal, Value: f}
}

func stringValue(val string) (Value, error) {
	u, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return Value{}, err
	}
	return Value{ValueType: str, Value: u}, nil
}

//un-casting
func intUncast(val uint64) int {
	p := unsafe.Pointer(&val)
	v := *(*int)(p)
	return v
}

func DecimalUncast(val uint64) float64 {
	p := unsafe.Pointer(&val)
	v := *(*float64)(p)
	return v
}

//Value nodes - will return Value structs when evaluated
type BoolNode struct {
	Token tokenizer.Token
}

func (node BoolNode) Evaluate() (Value, error) {
	switch node.Token.Text {
	case "false":
		return Value{Bool, 0}, nil
	}
	return Value{Bool, 1}, nil
}

type IntNode struct {
	Token tokenizer.Token
}

func (node IntNode) Evaluate() (Value, error) {
	number, err := strconv.Atoi(node.Token.Text)
	if err != nil {
		return Value{}, err
	}
	return intValue(number), nil
}

type DecimalNode struct {
	Token tokenizer.Token
}

func (node DecimalNode) Evaluate() (Value, error) {
	number, err := strconv.ParseFloat(node.Token.Text, 64)
	if err != nil {
		return Value{}, err
	}
	return DecimalValue(number), nil
}

type StringNode struct {
	Token tokenizer.Token
}

func (node StringNode) Evaluate() (Value, error) {
	val, err := stringValue(node.Token.Text)
	if err != nil {
		return Value{}, err
	}
	return val, nil
}

//Boolean Comparison nodes
type OrNode struct {
	Token tokenizer.Token
	Left  Node
	Right Node
}

func (node OrNode) Evaluate() (Value, error) {
	left, err := node.Left.Evaluate()
	if err != nil {
		return Value{}, err
	}
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		val, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = val
	}

	if left.ValueType == Identifier {
		val, err := fetchVar(left.Value)
		if err != nil {
			return Value{}, err
		}
		left = val
	}

	if left.ValueType != Bool || right.ValueType != Bool {
		return Value{}, LanErrs.ExpectedBoolError{node.Token}
	}

	if left.Value == 1 || right.Value == 1 {
		return Value{Bool, 1}, nil
	}
	return Value{Bool, 0}, nil
}

type AndNode struct {
	Token tokenizer.Token
	Left  Node
	Right Node
}

func (node AndNode) Evaluate() (Value, error) {
	left, err := node.Left.Evaluate()
	if err != nil {
		return Value{}, err
	}
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		val, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = val
	}

	if left.ValueType == Identifier {
		val, err := fetchVar(left.Value)
		if err != nil {
			return Value{}, err
		}
		left = val
	}

	if left.ValueType != Bool || right.ValueType != Bool {
		return Value{}, LanErrs.ExpectedBoolError{node.Token}
	}

	if left.Value == 1 && right.Value == 1 {
		return Value{Bool, 1}, nil
	}
	return Value{Bool, 0}, nil
}

type DoesEqualNode struct {
	Token tokenizer.Token
	Left  Node
	Right Node
}

func (node DoesEqualNode) Evaluate() (Value, error) {
	left, err := node.Left.Evaluate()
	if err != nil {
		return Value{}, err
	}
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		val, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = val
	}

	if left.ValueType == Identifier {
		val, err := fetchVar(left.Value)
		if err != nil {
			return Value{}, err
		}
		left = val
	}

	if left.ValueType != right.ValueType && !(isNum(left.ValueType) && isNum(right.ValueType)) {
		return Value{}, LanErrs.IncompatibleTypeError{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			if intUncast(left.Value) == intUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if float64(intUncast(left.Value)) == DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		}

	case Decimal:
		if right.ValueType == Decimal {
			if DecimalUncast(left.Value) == DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if DecimalUncast(left.Value) == float64(intUncast(right.Value)) {
				return Value{Bool, 1}, nil
			}
		}

	case str:
		l := Global.Strings[left.Value]
		r := Global.Strings[right.Value]
		if l == r {
			return Value{Bool, 1}, nil
		}

	case Bool:
		if left.Value == right.Value {
			return Value{Bool, 1}, nil
		}
	}

	return Value{Bool, 0}, nil
}

type NotEqualNode struct {
	Token tokenizer.Token
	Left  Node
	Right Node
}

func (node NotEqualNode) Evaluate() (Value, error) {
	left, err := node.Left.Evaluate()
	if err != nil {
		return Value{}, err
	}
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		val, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = val
	}

	if left.ValueType == Identifier {
		val, err := fetchVar(left.Value)
		if err != nil {
			return Value{}, err
		}
		left = val
	}

	if left.ValueType != right.ValueType && !(isNum(left.ValueType) && isNum(right.ValueType)) {
		return Value{}, LanErrs.IncompatibleTypeError{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			if intUncast(left.Value) != intUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if float64(intUncast(left.Value)) != DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		}

	case Decimal:
		if right.ValueType == Decimal {
			if DecimalUncast(left.Value) != DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if DecimalUncast(left.Value) != float64(intUncast(right.Value)) {
				return Value{Bool, 1}, nil
			}
		}

	case str:
		l := Global.Strings[left.Value]
		r := Global.Strings[right.Value]
		if l != r {
			return Value{Bool, 1}, nil
		}

	case Bool:
		if left.Value != right.Value {
			return Value{Bool, 1}, nil
		}
	}
	return Value{Bool, 0}, nil
}

func isNum(v valueKind) bool {
	if v == Integer || v == Decimal {
		return true
	}
	return false
}

type BigThanEqualNode struct {
	Token tokenizer.Token
	Left  Node
	Right Node
}

func (node BigThanEqualNode) Evaluate() (Value, error) {
	left, err := node.Left.Evaluate()
	if err != nil {
		return Value{}, err
	}
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		val, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = val
	}

	if left.ValueType == Identifier {
		val, err := fetchVar(left.Value)
		if err != nil {
			return Value{}, err
		}
		left = val
	}

	if !isNum(left.ValueType) || !isNum(right.ValueType) {
		return Value{}, LanErrs.MustBeNumWithComparisonOp{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			if intUncast(left.Value) >= intUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if float64(intUncast(left.Value)) >= DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		}

	case Decimal:
		if right.ValueType == Decimal {
			if DecimalUncast(left.Value) >= DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if DecimalUncast(left.Value) >= float64(intUncast(right.Value)) {
				return Value{Bool, 1}, nil
			}
		}
	}

	return Value{Bool, 0}, nil
}

type BigThanNode struct {
	Token tokenizer.Token
	Left  Node
	Right Node
}

func (node BigThanNode) Evaluate() (Value, error) {
	left, err := node.Left.Evaluate()
	if err != nil {
		return Value{}, err
	}
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		val, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = val
	}

	if left.ValueType == Identifier {
		val, err := fetchVar(left.Value)
		if err != nil {
			return Value{}, err
		}
		left = val
	}

	if !isNum(left.ValueType) || !isNum(right.ValueType) {
		return Value{}, LanErrs.MustBeNumWithComparisonOp{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			if intUncast(left.Value) > intUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if float64(intUncast(left.Value)) > DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		}

	case Decimal:
		if right.ValueType == Decimal {
			if DecimalUncast(left.Value) > DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if DecimalUncast(left.Value) > float64(intUncast(right.Value)) {
				return Value{Bool, 1}, nil
			}
		}
	}
	return Value{Bool, 0}, nil
}

type SmallThanEqualNode struct {
	Token tokenizer.Token
	Left  Node
	Right Node
}

func (node SmallThanEqualNode) Evaluate() (Value, error) {
	left, err := node.Left.Evaluate()
	if err != nil {
		return Value{}, err
	}
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		val, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = val
	}

	if left.ValueType == Identifier {
		val, err := fetchVar(left.Value)
		if err != nil {
			return Value{}, err
		}
		left = val
	}

	if !isNum(left.ValueType) || !isNum(right.ValueType) {
		return Value{}, LanErrs.MustBeNumWithComparisonOp{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			if intUncast(left.Value) <= intUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if float64(intUncast(left.Value)) <= DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		}

	case Decimal:
		if right.ValueType == Decimal {
			if DecimalUncast(left.Value) <= DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if DecimalUncast(left.Value) <= float64(intUncast(right.Value)) {
				return Value{Bool, 1}, nil
			}
		}
	}
	return Value{Bool, 0}, nil
}

type SmallThanNode struct {
	Token tokenizer.Token
	Left  Node
	Right Node
}

func (node SmallThanNode) Evaluate() (Value, error) {
	left, err := node.Left.Evaluate()
	if err != nil {
		return Value{}, err
	}
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		val, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = val
	}

	if left.ValueType == Identifier {
		val, err := fetchVar(left.Value)
		if err != nil {
			return Value{}, err
		}
		left = val
	}

	if !isNum(left.ValueType) || !isNum(right.ValueType) {
		return Value{}, LanErrs.MustBeNumWithComparisonOp{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			if intUncast(left.Value) < intUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if float64(intUncast(left.Value)) < DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		}

	case Decimal:
		if right.ValueType == Decimal {
			if DecimalUncast(left.Value) < DecimalUncast(right.Value) {
				return Value{Bool, 1}, nil
			}
		} else {
			if DecimalUncast(left.Value) < float64(intUncast(right.Value)) {
				return Value{Bool, 1}, nil
			}
		}
	}
	return Value{Bool, 0}, nil
}

type UnaryNode struct {
	Token tokenizer.Token
	Right Node
}

func (node UnaryNode) Evaluate() (Value, error) {
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		val, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = val
	}

	if node.Token.Text == "-" {
		switch right.ValueType {
		case Integer:
			r := intUncast(right.Value)
			r = r * -1
			return intValue(r), nil

		case Decimal:
			r := DecimalUncast(right.Value)
			r = r * -1
			return DecimalValue(r), nil
		}
		return Value{}, errors.New("ERROR: Expected number type with unary negation token at line num :" +
			strconv.Itoa(node.Token.LineNum) + ", cursor :" + strconv.Itoa(node.Token.Cursor))
	}

	if right.ValueType == Bool {
		if right.Value == 0 {
			return Value{Bool, 1}, nil
		}
		return Value{Bool, 0}, nil
	}
	return Value{}, errors.New("ERROR: Expected type Bool with unary NOT token at line num :" +
		strconv.Itoa(node.Token.LineNum) + ", cursor :" + strconv.Itoa(node.Token.Cursor))
}

//For Variables
type IdentifierNode struct {
	Token tokenizer.Token
}

func (node IdentifierNode) Evaluate() (Value, error) {
	if _, ok := globalVars[node.Token.Text]; ok {
		val := globalVars[node.Token.Text]
		return val, nil
	}
	return identifierValue(node.Token.Text), nil
}

//Stores the varibales
var globalVars = make(map[string]Value)

//Gets variable from Global vars
func getVar(name string) (Value, error) {
	if _, ok := globalVars[name]; ok {
		val := globalVars[name]
		return val, nil
	}
	return Value{}, LanErrs.NoIdentifierAvailableError{name}
}

//Sets a var
func setVar(name string, value Value) {
	globalVars[name] = value
}

//used to store an index into the Global array which stores varible names
func identifierValue(val string) Value {
	u, err := strconv.ParseUint(val, 10, 64)
	if err == nil {
		return Value{ValueType: Identifier, Value: u}
	}
	return Value{}
}

//used to do all the work of getting a varible - means less code in each Evaluate node method
func fetchVar(v uint64) (Value, error) {
	s := Global.GlobalVarNames[v]
	val, err := getVar(s)
	if err != nil {
		return Value{}, err
	}
	return val, nil
}

//Used for assigning values to varibles
type AssignmentNode struct {
	Token tokenizer.Token
	Left  Node
	Right Node
}

func (node AssignmentNode) Evaluate() (Value, error) {
	left, err := node.Left.Evaluate()
	if err != nil {
		return Value{}, err
	}
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		val, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = val
	}

	if left.ValueType != Identifier {
		return Value{}, errors.New("ERROR: Expected identifier on the left of assignment at line num :" +
			strconv.Itoa(node.Token.LineNum) + ", cursor : " + strconv.Itoa(node.Token.Cursor))
	}

	identifierStr := Global.GlobalVarNames[left.Value]

	setVar(identifierStr, right)

	return Value{}, nil
}

type PrintNode struct {
	Token tokenizer.Token
	Right Node
}

func (node PrintNode) Evaluate() (Value, error) {
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType == Identifier {
		this, err := fetchVar(right.Value)
		if err != nil {
			return Value{}, err
		}
		right = this
	}

	switch right.ValueType {
	case Integer:
		i := intUncast(right.Value)
		fmt.Println(i)

	case Decimal:
		i := DecimalUncast(right.Value)
		fmt.Println(i)

	case Bool:
		if right.Value == 1 {
			fmt.Println("True")
		} else {
			fmt.Println("False")
		}

	case str:
		s := Global.Strings[right.Value]
		fmt.Println(s)
	}
	return Value{}, nil
}

//Control Flow

type IfNode struct {
	Token      tokenizer.Token
	Expression Node
	Statements []Node
}

func (node IfNode) Evaluate() (Value, error) {
	left, err := node.Expression.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if left.ValueType != Bool {
		return Value{}, LanErrs.ExpectedBoolWithControlError{node.Token}
	}

	if left.Value == 1 {
		for _, statement := range node.Statements {
			_, err := statement.Evaluate()
			if err != nil {
				return Value{}, err
			}
		}
	}

	return Value{}, nil
}

type WhileNode struct {
	Token      tokenizer.Token
	Expression Node
	Statements []Node
}

func (node WhileNode) Evaluate() (Value, error) {
	left, err := node.Expression.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if left.ValueType != Bool {
		return Value{}, LanErrs.ExpectedBoolWithControlError{node.Token}
	}

	for left.Value == 1 {
		for _, statement := range node.Statements {
			_, err := statement.Evaluate()
			if err != nil {
				return Value{}, err
			}
		}
		left, err = node.Expression.Evaluate()
		if err != nil {
			return Value{}, err
		}
		if left.ValueType != Bool {
			return Value{}, LanErrs.ExpectedBoolWithControlError{node.Token}
		}
	}

	return Value{}, nil
}

type InputNode struct {
	Token tokenizer.Token
	Right Node
}

func (node InputNode) Evaluate() (Value, error) {
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}
	r := Global.Strings[right.Value]
	fmt.Println(r)

	var input string
	fmt.Scanln(&input)
	if len(input) == 0 {
		input = ""
	}
	Global.Strings = append(Global.Strings, input)
	address := len(Global.Strings) - 1
	val, err := stringValue(strconv.Itoa(address))
	if err != nil {
		return Value{}, err
	}
	return val, nil
}

type DelNode struct {
	Token tokenizer.Token
	Right Node
}

func (node DelNode) Evaluate() (Value, error) {
	right, err := node.Right.Evaluate()
	if err != nil {
		return Value{}, err
	}

	if right.ValueType != Identifier {
		return Value{}, LanErrs.ExpectedIdentifierError{node.Token}
	}
	index := Global.GlobalVarNames[right.Value]
	delete(globalVars, index)
	return Value{}, nil
}
