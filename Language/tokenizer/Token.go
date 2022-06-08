package tokenizer

type Token struct {
	Text    string
	Kind    TokenKind
	Cursor  int
	LineNum int
}

func CreateToken(text string, kind TokenKind, cursor int, lineNum int) Token {
	return Token{Text: text, Kind: kind, Cursor: cursor, LineNum: lineNum}
}
