package tree

type ChildrenNodes struct {
	LeftChild  Node
	RightChild Node
	LeftVal    Value
	RightVal   Value
}

func (children *ChildrenNodes) RetrieveChildrensNodeValues() error {
	var err error
	children.LeftVal, err = children.LeftChild.Evaluate()
	if err != nil {
		return err
	}
	if children.RightVal, err = children.RightChild.Evaluate(); err != nil {
		return err
	}

	if children.oneOrBothOfChildrenAreIdentifiers() {
		if err = children.retrieveValuesFromIdentifiers(); err != nil {
			return err
		}
	}
	return nil
}

func (children ChildrenNodes) oneOrBothOfChildrenAreIdentifiers() bool {
	if children.LeftVal.ValueType == Identifier ||
		children.RightVal.ValueType == Identifier {
		return true
	}
	return false
}

func (children *ChildrenNodes) retrieveValuesFromIdentifiers() error {
	var err error
	if children.LeftVal.ValueType == Identifier {
		if children.LeftVal, err = fetchVar(children.LeftVal.Value); err != nil {
			return err
		}
	}
	if children.RightVal.ValueType == Identifier {
		if children.RightVal, err = fetchVar(children.RightVal.Value); err != nil {
			return err
		}
	}
	return nil
}
