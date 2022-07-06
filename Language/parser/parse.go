package parser

import (
	tree "language/syntax_tree"
	"language/tokenizer"
)

type TreeBuilder struct {
	tokens      []tokenizer.Token
	Stack       []tree.Node
	index       int
	ParsedLines []tree.Node
	controlFlow bool
}

func NewParser(tokens []tokenizer.Token) TreeBuilder {
	t := TreeBuilder{tokens: tokens, index: 0, controlFlow: false}
	return t
}

func (t *TreeBuilder) EvaluateToken() {
	if t.index >= len(t.tokens) {
		return
	}
	switch true {
	case isLiteral(t.tokens[t.index]):
		t.handleLiteral()

	case isBinOp(t.tokens[t.index]):
		t.handleBinOp()

	case isBoolOp(t.tokens[t.index]):
		t.handleBoolOp()

	case isIdentifier(t.tokens[t.index]):
		t.handleIdentifier()

	case isAssingment(t.tokens[t.index]):
		t.handleAssignment()

	case isFunc(t.tokens[t.index]):
		t.handleFunc()

	case isControl(t.tokens[t.index]):
		t.handleControlStatement()

	case isEndOfLine(t.tokens[t.index]):
		if len(t.Stack) > 0 {
			node := popFromStack(&t.Stack)
			t.ParsedLines = append(t.ParsedLines, node)
		}
		t.index += 1
		t.EvaluateToken()
	}

}

func isControl(token tokenizer.Token) bool {
	if token.Kind == tokenizer.If || token.Kind == tokenizer.While {
		return true
	}
	return false
}

func (t *TreeBuilder) handleFunc() {
	right := popFromStack(&t.Stack)
	token := t.tokens[t.index]
	switch token.Kind {
	case tokenizer.Print:
		t.Stack = append(t.Stack, tree.PrintNode{token, right})

	case tokenizer.Input:
		t.Stack = append(t.Stack, tree.InputNode{token, right})

	case tokenizer.Del:
		t.Stack = append(t.Stack, tree.DelNode{token, right})
	}

	t.index += 1
	t.EvaluateToken()
}

func isFunc(token tokenizer.Token) bool {
	switch token.Kind {
	case tokenizer.Print, tokenizer.Input, tokenizer.Del:
		return true
	}
	return false
}

func isEndOfLine(token tokenizer.Token) bool {
	if token.Kind == tokenizer.EndOfStatment {
		return true
	}
	return false
}

func (t *TreeBuilder) handleAssignment() {
	right := popFromStack(&t.Stack)
	left := popFromStack(&t.Stack)
	token := t.tokens[t.index]

	t.Stack = append(t.Stack, tree.AssignmentNode{token, left, right})
	t.index += 1
	t.EvaluateToken()
}

func isAssingment(token tokenizer.Token) bool {
	if token.Kind == tokenizer.Assign {
		return true
	}
	return false
}

func (t *TreeBuilder) handleIdentifier() {
	token := t.tokens[t.index]
	t.Stack = append(t.Stack, tree.IdentifierNode{token})

	t.index += 1
	t.EvaluateToken()
}

func isIdentifier(token tokenizer.Token) bool {
	if token.Kind == tokenizer.Identifier {
		return true
	}
	return false
}

func (t *TreeBuilder) handleBoolOp() {
	right := popFromStack(&t.Stack)
	left := popFromStack(&t.Stack)
	token := t.tokens[t.index]

	switch t.tokens[t.index].Text {
	case "<":
		t.Stack = append(t.Stack, tree.SmallThanNode{token, left, right})

	case "<=":
		t.Stack = append(t.Stack, tree.SmallThanEqualNode{token, left, right})

	case ">":
		t.Stack = append(t.Stack, tree.BigThanNode{token, left, right})

	case ">=":
		t.Stack = append(t.Stack, tree.BigThanEqualNode{token, left, right})

	case "!=":
		t.Stack = append(t.Stack, tree.NotEqualNode{token, left, right})

	case "=":
		t.Stack = append(t.Stack, tree.DoesEqualNode{token, left, right})

	case "&":
		t.Stack = append(t.Stack, tree.AndNode{token, left, right})

	case "|":
		t.Stack = append(t.Stack, tree.OrNode{token, left, right})
	}
	t.index += 1
	t.EvaluateToken()
}

func isBoolOp(token tokenizer.Token) bool {
	switch token.Kind {
	case tokenizer.BooleanOp, tokenizer.BoolConnector:
		return true

	default:
		return false
	}
}

func (t *TreeBuilder) createChildrenStructure() tree.ChildrenNodes {
	right := popFromStack(&t.Stack)
	left := popFromStack(&t.Stack)

	return tree.ChildrenNodes{
		LeftChild:  left,
		RightChild: right,
		LeftVal:    tree.Value{},
		RightVal:   tree.Value{},
	}
}

func (t *TreeBuilder) handleBinOp() {

	token := t.tokens[t.index]

	switch t.tokens[t.index].Kind {
	case tokenizer.Multiply:
		children := t.createChildrenStructure()
		t.Stack = append(t.Stack, tree.MultiplyNode{token, children})

	case tokenizer.Add:
		children := t.createChildrenStructure()
		t.Stack = append(t.Stack, tree.AddNode{token, children})

	case tokenizer.Divide:
		// children := t.createChildrenStructure()
		//left := popFromStack(&t.Stack)
		//right := popFromStack(&t.Stack)
		//t.Stack = append(t.Stack, tree.DivideNode{token, left, right})
		children := t.createChildrenStructure()
		t.Stack = append(t.Stack, tree.DivideNode{token, children})

	case tokenizer.Subtract:
		// children := t.createChildrenStructure()
		//left := popFromStack(&t.Stack)
		//right := popFromStack(&t.Stack)
		//t.Stack = append(t.Stack, tree.SubtractNode{token, left, right})
		children := t.createChildrenStructure()
		t.Stack = append(t.Stack, tree.SubtractNode{token, children})

	case tokenizer.Exspo:
		// children := t.createChildrenStructure()
		//left := popFromStack(&t.Stack)
		//right := popFromStack(&t.Stack)
		//t.Stack = append(t.Stack, tree.ExpoNode{token, left, right})
		children := t.createChildrenStructure()
		t.Stack = append(t.Stack, tree.ExpoNode{token, children})

	case tokenizer.Unary:
		right := popFromStack(&t.Stack)
		t.Stack = append(t.Stack, tree.UnaryNode{Token: token, Right: right})
	}
	t.index += 1
	t.EvaluateToken()
}

func (t *TreeBuilder) handleLiteral() {

	switch t.tokens[t.index].Kind {
	case tokenizer.Int:
		intNode := tree.IntNode{t.tokens[t.index]}
		t.Stack = append(t.Stack, intNode)

	case tokenizer.Decimal:
		floatNode := tree.DecimalNode{t.tokens[t.index]}
		t.Stack = append(t.Stack, floatNode)

	case tokenizer.String:
		stringNode := tree.StringNode{t.tokens[t.index]}
		t.Stack = append(t.Stack, stringNode)

	case tokenizer.Bool:
		boolNode := tree.BoolNode{t.tokens[t.index]}
		t.Stack = append(t.Stack, boolNode)
	}
	t.index += 1
	t.EvaluateToken()
}

func isLiteral(token tokenizer.Token) bool {
	if token.Kind == tokenizer.Int || token.Kind == tokenizer.Decimal || token.Kind == tokenizer.String || token.Kind == tokenizer.Bool {
		return true
	}

	return false
}

func isBinOp(token tokenizer.Token) bool {
	switch token.Kind {
	case tokenizer.Exspo, tokenizer.Subtract, tokenizer.Add, tokenizer.Divide, tokenizer.Multiply, tokenizer.Unary:
		return true

	default:
		return false
	}
}

func popFromStack(nodes *[]tree.Node) tree.Node {
	node := (*nodes)[len(*nodes)-1]
	*nodes = (*nodes)[:len(*nodes)-1]
	return node
}
