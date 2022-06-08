package ShuntingYard

import (
	"language/tokenizer"
)

type ShuntingY struct {
	Tokens []tokenizer.Token
	stack  []tokenizer.Token
	Result []tokenizer.Token
	Index  int
}

var operations = map[tokenizer.TokenKind]struct {
	prec  int
	assoc bool
}{
	tokenizer.Exspo:         {5, true},
	tokenizer.Unary:         {5, true},
	tokenizer.Multiply:      {4, false},
	tokenizer.Divide:        {4, false},
	tokenizer.Add:           {3, false},
	tokenizer.Subtract:      {3, false},
	tokenizer.BooleanOp:     {2, true},
	tokenizer.BoolConnector: {2, false},
	tokenizer.Input:         {2, false},
	tokenizer.Assign:        {1, true},
	tokenizer.Print:         {1, true},
	tokenizer.Del:           {1, true},
	tokenizer.EndOfStatment: {0, false},
	tokenizer.BlockStart:    {0, false},
	tokenizer.BlockEnd:      {0, false},
}

func (s *ShuntingY) ToPostFix() {
	if s.Tokens[s.Index].Kind == tokenizer.End {
		// drain stack to Result
		for len(s.stack) > 0 {
			s.Result = append(s.Result, s.stack[len(s.stack)-1])
			s.stack = s.stack[:len(s.stack)-1]
		}
		s.Result = append(s.Result, s.Tokens[s.Index])
		return
	} else if s.Tokens[s.Index].Kind == tokenizer.EndOfStatment {
		for len(s.stack) > 0 {
			s.Result = append(s.Result, s.stack[len(s.stack)-1])
			s.stack = s.stack[:len(s.stack)-1]
		}
	}

	switch true {
	case isOpenBrack(s.Tokens[s.Index]):
		s.stack = append(s.stack, s.Tokens[s.Index])
		s.Index += 1
		s.ToPostFix()

	case isCloseBrack(s.Tokens[s.Index]):
		s.handleCloseBrack()
		s.Index += 1
		s.ToPostFix()

	case isOp(s.Tokens[s.Index]):
		s.testUnary()
		s.handleOp()
		s.Index += 1
		s.ToPostFix()

	default: // token is an operator
		s.Result = append(s.Result, s.Tokens[s.Index])
		s.Index += 1
		s.ToPostFix()
	}
}

func isEnd(t tokenizer.Token) bool {
	if t.Kind == tokenizer.EndOfStatment {
		return true
	}
	return false
}

func isOpenBrack(t tokenizer.Token) bool {

	if t.Kind == tokenizer.Openbrack {
		return true
	}
	return false
}

func isCloseBrack(t tokenizer.Token) bool {
	if t.Kind == tokenizer.Closebrack {
		return true
	}
	return false
}

func isOp(token tokenizer.Token) bool {
	switch token.Kind {
	case tokenizer.Exspo, tokenizer.Subtract, tokenizer.Add, tokenizer.Divide, tokenizer.Multiply,
		tokenizer.Unary, tokenizer.BooleanOp, tokenizer.BoolConnector, tokenizer.Assign, tokenizer.Print,
		tokenizer.BlockStart, tokenizer.BlockEnd, tokenizer.Input, tokenizer.Del:
		return true

	default:
		return false
	}
}

func (s *ShuntingY) handleCloseBrack() {
	for {
		var op tokenizer.Token

		op, s.stack = s.stack[len(s.stack)-1], s.stack[:len(s.stack)-1] // Pops the last value of stack to op

		if op.Kind == tokenizer.Openbrack {
			break // Will delete the close bracket token
		}
		s.Result = append(s.Result, op)
	}
}

func (s *ShuntingY) handleOp() {
	op1 := operations[s.Tokens[s.Index].Kind]
	for len(s.stack) > 0 {
		topOp := s.stack[len(s.stack)-1] // Token on top of the stack

		if op2, isOperator := operations[topOp.Kind]; !isOperator || op1.prec > op2.prec ||
			op1.prec == op2.prec && op1.assoc {
			break
		}
		// top item is an operator that needs to come off
		s.stack = s.stack[:len(s.stack)-1] // pop it
		s.Result = append(s.Result, topOp) // add it to Result
	}
	// push operator (the new one) to stack
	s.stack = append(s.stack, s.Tokens[s.Index])
}

func (s *ShuntingY) testUnary() {
	if s.Index == 0 || isOpenBrack(s.Tokens[s.Index-1]) || isOp(s.Tokens[s.Index-1]) ||
		s.Tokens[s.Index-1].Kind == tokenizer.If || s.Tokens[s.Index-1].Kind == tokenizer.While {
		if s.Tokens[s.Index].Kind != tokenizer.Print && s.Tokens[s.Index].Kind != tokenizer.Input &&
			s.Tokens[s.Index].Kind != tokenizer.Del {
			s.Tokens[s.Index].Kind = tokenizer.Unary
		}
	}
}
