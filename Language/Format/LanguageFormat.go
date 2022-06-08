package Format

import (
	"language/tokenizer"
)

func NewFormatChecker(tokens []tokenizer.Token) FormatChecker {
	return FormatChecker{tokens, 0}
}

type FormatChecker struct {
	Tokens []tokenizer.Token
	Index  int
}

func (f *FormatChecker) FormatTokens() {
	for f.Index < len(f.Tokens) {
		switch f.Tokens[f.Index].Kind {
		case tokenizer.EndOfStatment:
			f.HandleNewLine()
		}
		f.Index += 1
	}
}

func (f *FormatChecker) HandleNewLine() {
	for f.Tokens[f.Index+1].Kind == tokenizer.EndOfStatment {
		f.RemoveToken()
	}
	f.Index += 1
	f.FormatTokens()
}

func (f *FormatChecker) RemoveToken() {
	newTokens := make([]tokenizer.Token, 0)
	newTokens = append(newTokens, f.Tokens[:f.Index]...)
	newTokens = append(newTokens, f.Tokens[f.Index+1:]...)
	f.Tokens = newTokens
}
