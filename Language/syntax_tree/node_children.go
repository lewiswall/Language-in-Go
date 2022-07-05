package tree

type childrenNodes struct {
	leftChild  Node
	rightChild Node
	leftVal    Value
	rightVal   Value
}

func (children childrenNodes) RetrieveChildrensNodeValues() error {
	var err error
	if children.leftVal, err = children.leftChild.Evaluate(); err != nil {
		return err
	}
	if children.rightVal, err = children.rightChild.Evaluate(); err != nil {
		return err
	}

	if children.oneOrBothOfChildrenAreIdentifiers() {
		children.retrieveValuesFromIdentifiers()
	}

	return nil
}

func (children childrenNodes) oneOrBothOfChildrenAreIdentifiers() bool {
	if children.leftVal.ValueType == Identifier ||
		children.rightVal.ValueType == Identifier {
		return true
	}
	return false
}

func (children childrenNodes) retrieveValuesFromIdentifiers() error {
	var err error
	if children.leftVal.ValueType == Identifier {
		children.leftVal, err = fetchVar(children.leftVal.Value)
		if err != nil {
			return err
		}
	}
	if children.rightVal.ValueType == Identifier {
		children.rightVal, err = fetchVar(children.rightVal.Value)
		if err != nil {
			return err
		}
	}
	return nil
}
