package tree

import (
	"language/Global"
	"language/LanErrs"
	"language/tokenizer"
	"math"
	"strconv"
)

type MultiplyNode struct {
	Token tokenizer.Token
	ChildrenNodes
}

func (node MultiplyNode) Evaluate() (Value, error) {
	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

	if left.ValueType != right.ValueType && !(isNum(left.ValueType) && isNum(right.ValueType)) {
		return Value{}, LanErrs.IncompatibleTypeError{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			return intValue(intUncast(left.Value) * intUncast(right.Value)), nil
		}
		return DecimalValue(float64(intUncast(left.Value)) * DecimalUncast(right.Value)), nil
	case Decimal:
		if right.ValueType == Decimal {
			return DecimalValue(DecimalUncast(left.Value) * DecimalUncast(right.Value)), nil
		}
		return DecimalValue(DecimalUncast(left.Value) * float64(intUncast(right.Value))), nil
	}

	//Error
	return Value{}, LanErrs.WrongTypeUsedWithBinOpError{node.Token}
}

type AddNode struct {
	Token tokenizer.Token
	ChildrenNodes
}

func (node AddNode) Evaluate() (Value, error) {

	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

	if left.ValueType != right.ValueType && !(isNum(left.ValueType) && isNum(right.ValueType)) {
		return Value{}, LanErrs.IncompatibleTypeError{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			return intValue(intUncast(left.Value) + intUncast(right.Value)), nil
		}
		return DecimalValue(float64(intUncast(left.Value)) + DecimalUncast(right.Value)), nil
	case Decimal:
		if right.ValueType == Decimal {
			return DecimalValue(DecimalUncast(left.Value) + DecimalUncast(right.Value)), nil
		}
		return DecimalValue(DecimalUncast(left.Value) + float64(intUncast(right.Value))), nil
	case str:
		val, err := addStrings(Global.Strings[left.Value], Global.Strings[right.Value])
		if err != nil {
			return Value{}, err
		}
		return val, nil
	}

	//return error at the end - the node types cannot be used with '+'
	return Value{}, LanErrs.WrongTypeUsedWithBinOpError{node.Token}
}

func addStrings(l string, r string) (Value, error) {
	Global.Strings = append(Global.Strings, l+r)
	index := strconv.Itoa(len(Global.Strings) - 1)

	val, err := stringValue(index)
	if err != nil {
		return Value{}, err
	}
	return val, nil
}

type DivideNode struct {
	Token tokenizer.Token
	ChildrenNodes
}

func (node DivideNode) Evaluate() (Value, error) {

	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

	if left.ValueType != right.ValueType && !(isNum(left.ValueType) && isNum(right.ValueType)) {
		return Value{}, LanErrs.IncompatibleTypeError{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			return intValue(intUncast(left.Value) / intUncast(right.Value)), nil
		}
		return DecimalValue(float64(intUncast(left.Value)) / DecimalUncast(right.Value)), nil
	case Decimal:
		if right.ValueType == Decimal {
			return DecimalValue(DecimalUncast(left.Value) / DecimalUncast(right.Value)), nil
		}
		return DecimalValue(DecimalUncast(left.Value) / float64(intUncast(right.Value))), nil
	}

	//return error
	return Value{}, LanErrs.WrongTypeUsedWithBinOpError{node.Token}
}

type SubtractNode struct {
	Token tokenizer.Token
	ChildrenNodes
}

func (node SubtractNode) Evaluate() (Value, error) {
	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

	if left.ValueType != right.ValueType && !(isNum(left.ValueType) && isNum(right.ValueType)) {
		return Value{}, LanErrs.IncompatibleTypeError{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			return intValue(intUncast(left.Value) - intUncast(right.Value)), nil
		}
		return DecimalValue(float64(intUncast(left.Value)) - DecimalUncast(right.Value)), nil
	case Decimal:
		if right.ValueType == Decimal {
			return DecimalValue(DecimalUncast(left.Value) - DecimalUncast(right.Value)), nil
		}
		return DecimalValue(DecimalUncast(left.Value) - float64(intUncast(right.Value))), nil
	}

	return Value{}, LanErrs.WrongTypeUsedWithBinOpError{node.Token}
}

type ExpoNode struct {
	Token tokenizer.Token
	ChildrenNodes
}

func (node ExpoNode) Evaluate() (Value, error) {

	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

	if left.ValueType != right.ValueType && !(isNum(left.ValueType) && isNum(right.ValueType)) {
		return Value{}, LanErrs.IncompatibleTypeError{node.Token}
	}

	switch left.ValueType {
	case Integer:
		if right.ValueType == Integer {
			return intValue(int(math.Pow(float64(intUncast(left.Value)), float64(intUncast(right.Value))))), nil
		}
		return DecimalValue(math.Pow(float64(intUncast(left.Value)), DecimalUncast(right.Value))), nil
	case Decimal:
		if right.ValueType == Decimal {
			return DecimalValue(math.Pow(DecimalUncast(left.Value), DecimalUncast(right.Value))), nil
		}
		return DecimalValue(math.Pow(DecimalUncast(left.Value), float64(intUncast(right.Value)))), nil
	}

	//Return error because node types cannot be used with exspo
	return Value{}, LanErrs.WrongTypeUsedWithBinOpError{node.Token}
}
