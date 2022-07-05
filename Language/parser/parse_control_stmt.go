package parser

import (
	tree "language/syntax_tree"
	"language/tokenizer"
)

func (t *TreeBuilder) handleControlStatement() {
	controlType := t.tokens[t.index]
	t.index += 1
	switch controlType.Kind {
	case tokenizer.If:
		t.handleIfStatement(controlType)
	case tokenizer.While:
		t.handleWhileStatement(controlType)
	}
	t.index += 1
	t.EvaluateToken()
}

func (t *TreeBuilder) handleIfStatement(ifToken tokenizer.Token) {
	expression := t.getExpressionInFormSyntaxTree()
	statements := t.getControlsStatementsInFormSyntaxTree()
	t.Stack = append(t.Stack, tree.IfNode{ifToken, expression, statements})
}

func (t *TreeBuilder) handleWhileStatement(whileToken tokenizer.Token) {
	expression := t.getExpressionInFormSyntaxTree()
	statements := t.getControlsStatementsInFormSyntaxTree()
	t.Stack = append(t.Stack, tree.WhileNode{whileToken, expression, statements})
}

func (t *TreeBuilder) getControlsStatementsInFormSyntaxTree() []tree.Node {
	statementsStart := t.index + 2
	t.index += 1
	statementsEnd := t.getStatementsEnd()
	return retrieveParsedExpressions(t.tokens[statementsStart:statementsEnd], statementsEnd-statementsStart)
}

func (t *TreeBuilder) getExpressionInFormSyntaxTree() tree.Node {
	expressionStart := t.index
	expressionEnd := t.getExpressionEnd()
	return retrieveParsedExpressions(t.tokens[expressionStart:expressionEnd], expressionEnd-expressionStart)[0]
}

func (t *TreeBuilder) getStatementsEnd() int {
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
	return t.index
}

func (t *TreeBuilder) getExpressionEnd() int {
	for t.tokens[t.index].Kind != tokenizer.BlockStart {
		t.index += 1
	}
	return t.index
}

func retrieveParsedExpressions(tokens []tokenizer.Token, numOfTokens int) []tree.Node {
	var expression = make([]tokenizer.Token, numOfTokens)
	copy(expression, tokens)
	addEndToken(&expression)

	parsedExpression := parseTokensAndReturnTopNode(expression)
	return parsedExpression
}

func parseTokensAndReturnTopNode(tokens []tokenizer.Token) []tree.Node {
	expressionParser := NewParser(tokens)
	expressionParser.EvaluateToken()

	return expressionParser.ParsedLines
}

func addEndToken(tokens *[]tokenizer.Token) {
	end := tokenizer.CreateToken("END", tokenizer.EndOfStatment, 0, 0)
	*tokens = append((*tokens), end)
}
