package parser

import (
	"language/tokenizer"
	"language/tree"
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
		t.handleControl()

	case isEndOfLine(t.tokens[t.index]):
		node := popFromStack(&t.Stack)
		t.ParsedLines = append(t.ParsedLines, node)
		t.index += 1
		t.EvaluateToken()
	}

}

func (t *TreeBuilder) handleControl() {
	t.controlFlow = true
	contType := t.tokens[t.index]
	t.index += 1

	expressionStart := t.index
	for t.tokens[t.index].Kind != tokenizer.BlockStart {
		t.index += 1
	}
	expressionEnd := t.index

	var expression = make([]tokenizer.Token, expressionEnd-expressionStart)
	copy(expression, t.tokens[expressionStart:expressionEnd])
	end := tokenizer.CreateToken("END", tokenizer.End, 0, 0)
	expression = append(expression, end)
	leftTree := NewParser(expression)
	leftTree.EvaluateToken()
	l := leftTree.Stack[0]

	blockStart := t.index + 2
	t.index += 1
	nestedBlock := 0
	for t.tokens[t.index].Kind != tokenizer.BlockEnd || nestedBlock > 0 {
		if t.tokens[t.index].Kind == tokenizer.BlockStart {
			nestedBlock += 1
		}
		if t.tokens[t.index].Kind == tokenizer.BlockEnd {
			nestedBlock -= 1
		}
		t.index += 1
	}
	blockEnd := t.index

	var block = make([]tokenizer.Token, blockEnd-blockStart)
	copy(block, t.tokens[blockStart:blockEnd])

	block = append(block, tokenizer.CreateToken("END", tokenizer.End, 0, 0))
	rightTree := NewParser(block)
	rightTree.EvaluateToken()
	r := rightTree.ParsedLines

	switch contType.Kind {
	case tokenizer.If:
		t.Stack = append(t.Stack, tree.IfNode{contType, l, r})

	case tokenizer.While:
		t.Stack = append(t.Stack, tree.WhileNode{contType, l, r})
	}

	t.index += 1
	t.EvaluateToken()
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

func (t *TreeBuilder) handleBinOp() {

	token := t.tokens[t.index]

	switch t.tokens[t.index].Kind {
	case tokenizer.Multiply:
		right := popFromStack(&t.Stack)
		left := popFromStack(&t.Stack)
		t.Stack = append(t.Stack, tree.MultiplyNode{Token: token, Right: right, Left: left})

	case tokenizer.Add:
		right := popFromStack(&t.Stack)
		left := popFromStack(&t.Stack)
		t.Stack = append(t.Stack, tree.AddNode{Token: token, Right: right, Left: left})

	case tokenizer.Divide:
		right := popFromStack(&t.Stack)
		left := popFromStack(&t.Stack)
		t.Stack = append(t.Stack, tree.DivideNode{Token: token, Right: right, Left: left})

	case tokenizer.Subtract:
		right := popFromStack(&t.Stack)
		left := popFromStack(&t.Stack)
		t.Stack = append(t.Stack, tree.SubtractNode{Token: token, Right: right, Left: left})

	case tokenizer.Exspo:
		right := popFromStack(&t.Stack)
		left := popFromStack(&t.Stack)
		t.Stack = append(t.Stack, tree.ExpoNode{Token: token, Right: right, Left: left})

	case tokenizer.Unary:
		right := popFromStack(&t.Stack)
		t.Stack = append(t.Stack, tree.UnaryNode{Token: token, Right: right})
	}
	t.index += 1
	t.EvaluateToken()
}

func (t *TreeBuilder) handleLiteral() {
	//if t.tokens[t.index].Kind == tokenizer.Int {
	//	intNode := tree.IntNode{t.tokens[t.index]}
	//	t.Stack = append(t.Stack, intNode)
	//} else if t.tokens[t.index].Kind == tokenizer.Float {
	//	floatNode := tree.FloatNode{t.tokens[t.index]}
	//	t.Stack = append(t.Stack, floatNode)
	//}

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
