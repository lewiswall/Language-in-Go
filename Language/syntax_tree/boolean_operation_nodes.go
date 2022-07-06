package tree

import (
	"language/Global"
	"language/LanErrs"
	"language/tokenizer"
)

type OrNode struct {
	Token tokenizer.Token
	ChildrenNodes
}

func (node OrNode) Evaluate() (Value, error) {
	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

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
	ChildrenNodes
}

func (node AndNode) Evaluate() (Value, error) {
	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

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
	ChildrenNodes
}

func (node DoesEqualNode) Evaluate() (Value, error) {
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
	ChildrenNodes
}

func (node NotEqualNode) Evaluate() (Value, error) {
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

type BigThanEqualNode struct {
	Token tokenizer.Token
	ChildrenNodes
}

func (node BigThanEqualNode) Evaluate() (Value, error) {
	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

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
	ChildrenNodes
}

func (node BigThanNode) Evaluate() (Value, error) {
	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

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
	ChildrenNodes
}

func (node SmallThanEqualNode) Evaluate() (Value, error) {
	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

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
	ChildrenNodes
}

func (node SmallThanNode) Evaluate() (Value, error) {
	if err := node.ChildrenNodes.RetrieveChildrensNodeValues(); err != nil {
		return Value{}, err
	}
	left, right := node.ChildrenNodes.LeftVal, node.ChildrenNodes.RightVal

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
