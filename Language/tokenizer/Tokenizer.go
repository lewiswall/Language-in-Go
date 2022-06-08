package tokenizer

import (
	"errors"
	"strconv"
	"unicode"
)

type tokenizer struct {
	text       string
	cursor     int
	lineNumber int
}

// Creates a new Tokenizer
func New() tokenizer {
	e := tokenizer{lineNumber: 0, cursor: 0}
	return e
}

func (tokenizer *tokenizer) NewLine(text string) {
	tokenizer.cursor = 0
	tokenizer.text = text
	tokenizer.lineNumber += 1
}

func (tokenizer *tokenizer) Get() (Token, error) {

	for tokenizer.cursor < len(tokenizer.text) {
		char := []rune(tokenizer.text)[tokenizer.cursor]
		switch char {

		case '"':
			tokenizer.cursor += 1
			stringHead := tokenizer.cursor

			for tokenizer.cursor < len(tokenizer.text) {
				if []rune(tokenizer.text)[tokenizer.cursor] == '"' {
					stringTail := tokenizer.cursor
					tokenizer.cursor += 1

					return CreateToken(tokenizer.text[stringHead:stringTail], String, tokenizer.cursor, tokenizer.lineNumber), nil
				}
				tokenizer.cursor += 1
			}
			err := "There is no closing `\"` on string @ Line Int : " + strconv.Itoa(tokenizer.lineNumber) + "; Cursor Int : " + strconv.Itoa(tokenizer.cursor)
			return Token{}, errors.New(err)

		case ' ', '\t':
			tokenizer.cursor += 1

		case '\n':
			//idenStart := tokenizer.cursor
			tokenizer.cursor += 1
			//tokenizer.lineNumber += 1
			//return CreateToken("NL", EndOfStatment, idenStart, tokenizer.lineNumber-1), nil

		case '+', '-', '*', '/', '(', ')', '^':
			opToken := createOperatorToken(char, tokenizer)
			tokenizer.cursor += 1
			return opToken, nil

		case '!':
			if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == '=' {
				tokenizer.cursor += 2
				return CreateToken(tokenizer.text[tokenizer.cursor-2:tokenizer.cursor], BooleanOp, tokenizer.cursor-2, tokenizer.lineNumber), nil
			}
			tokenizer.cursor += 1
			return CreateToken(tokenizer.text[tokenizer.cursor-1:tokenizer.cursor], Unary, tokenizer.cursor, tokenizer.lineNumber), nil

		case '<':
			if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == '=' {
				tokenizer.cursor += 2
				return CreateToken(tokenizer.text[tokenizer.cursor-2:tokenizer.cursor], BooleanOp, tokenizer.cursor-2, tokenizer.lineNumber), nil
			}
			tokenizer.cursor += 1
			return CreateToken(tokenizer.text[tokenizer.cursor-1:tokenizer.cursor], BooleanOp, tokenizer.cursor, tokenizer.lineNumber), nil

		case '>':
			if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == '=' {
				tokenizer.cursor += 2
				return CreateToken(tokenizer.text[tokenizer.cursor-2:tokenizer.cursor], BooleanOp, tokenizer.cursor-2, tokenizer.lineNumber), nil
			}
			tokenizer.cursor += 1
			return CreateToken(tokenizer.text[tokenizer.cursor-1:tokenizer.cursor], BooleanOp, tokenizer.cursor, tokenizer.lineNumber), nil

		case '=':
			tokenizer.cursor += 1
			return CreateToken(tokenizer.text[tokenizer.cursor-1:tokenizer.cursor], BooleanOp, tokenizer.cursor, tokenizer.lineNumber), nil

		case '&':
			tokenizer.cursor += 1
			return CreateToken(tokenizer.text[tokenizer.cursor-1:tokenizer.cursor], BoolConnector, tokenizer.cursor-1, tokenizer.lineNumber), nil

		case '|':
			tokenizer.cursor += 1
			return CreateToken(tokenizer.text[tokenizer.cursor-1:tokenizer.cursor], BoolConnector, tokenizer.cursor-1, tokenizer.lineNumber), nil

		case ':':
			if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == '=' {
				tokenizer.cursor += 2
				return CreateToken(tokenizer.text[tokenizer.cursor-2:tokenizer.cursor], Assign, tokenizer.cursor-2, tokenizer.lineNumber), nil
			}

		case '{':
			tokenizer.cursor += 1
			return CreateToken(tokenizer.text[tokenizer.cursor-1:tokenizer.cursor], BlockStart, tokenizer.cursor, tokenizer.lineNumber), nil

		case '}':
			tokenizer.cursor += 1
			return CreateToken(tokenizer.text[tokenizer.cursor-1:tokenizer.cursor], BlockEnd, tokenizer.cursor, tokenizer.lineNumber), nil

		case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
			identifierStart := tokenizer.cursor
			tokenizer.cursor += 1

			//Checking for a single digit float number
			if tokenizer.cursor < len(tokenizer.text) && []rune(tokenizer.text)[tokenizer.cursor] == '.' {
				tokenizer.cursor += 1
				for tokenizer.cursor < len(tokenizer.text) && unicode.IsDigit([]rune(tokenizer.text)[tokenizer.cursor]) {
					tokenizer.cursor += 1
				}
				return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], Decimal, tokenizer.cursor, tokenizer.lineNumber), nil
			}

			for tokenizer.cursor < len(tokenizer.text) && unicode.IsDigit([]rune(tokenizer.text)[tokenizer.cursor]) {
				tokenizer.cursor += 1

				if tokenizer.cursor < len(tokenizer.text) && []rune(tokenizer.text)[tokenizer.cursor] == '.' {
					tokenizer.cursor += 1
					for tokenizer.cursor < len(tokenizer.text) && unicode.IsDigit([]rune(tokenizer.text)[tokenizer.cursor]) {
						tokenizer.cursor += 1
					}
					return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], Decimal, tokenizer.cursor, tokenizer.lineNumber), nil
				}
			}
			return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], Int, tokenizer.cursor, tokenizer.lineNumber), nil
		default:
			if unicode.IsLetter(char) {
				identifierStart := tokenizer.cursor

				if char == 't' {
					if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == 'r' {
						if tokenizer.cursor+2 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+2] == 'u' {
							if tokenizer.cursor+3 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+3] == 'e' {
								tokenizer.cursor += 4
								return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], Bool, identifierStart, tokenizer.lineNumber), nil
							}
						}
					}
				}

				if char == 'f' {
					if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == 'a' {
						if tokenizer.cursor+2 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+2] == 'l' {
							if tokenizer.cursor+3 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+3] == 's' {
								if tokenizer.cursor+4 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+4] == 'e' {
									tokenizer.cursor += 5
									return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], Bool, identifierStart, tokenizer.lineNumber), nil
								}
							}
						}
					}
				}

				if char == 'p' {
					if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == 'r' {
						if tokenizer.cursor+2 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+2] == 'i' {
							if tokenizer.cursor+3 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+3] == 'n' {
								if tokenizer.cursor+4 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+4] == 't' {
									tokenizer.cursor += 5
									return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], Print, identifierStart, tokenizer.lineNumber), nil
								}
							}
						}
					}
				}

				if char == 'i' {
					if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == 'f' {
						tokenizer.cursor += 2
						return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], If, identifierStart, tokenizer.lineNumber), nil
					}
				}

				if char == 'w' {
					if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == 'h' {
						if tokenizer.cursor+2 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+2] == 'i' {
							if tokenizer.cursor+3 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+3] == 'l' {
								if tokenizer.cursor+4 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+4] == 'e' {
									tokenizer.cursor += 5
									return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], While, identifierStart, tokenizer.lineNumber), nil
								}
							}
						}
					}
				}

				if char == 'i' {
					if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == 'n' {
						if tokenizer.cursor+2 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+2] == 'p' {
							if tokenizer.cursor+3 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+3] == 'u' {
								if tokenizer.cursor+4 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+4] == 't' {
									tokenizer.cursor += 5
									return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], Input, identifierStart, tokenizer.lineNumber), nil
								}
							}
						}
					}
				}

				if char == 'd' {
					if tokenizer.cursor+1 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+1] == 'e' {
						if tokenizer.cursor+2 < len(tokenizer.text) && tokenizer.text[tokenizer.cursor+2] == 'l' {
							tokenizer.cursor += 3
							return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], Del, identifierStart, tokenizer.lineNumber), nil

						}
					}
				}

				tokenizer.cursor += 1

				for tokenizer.cursor < len(tokenizer.text) && !unicode.IsSpace([]rune(tokenizer.text)[tokenizer.cursor]) {
					tokenizer.cursor += 1
				}
				return CreateToken(tokenizer.text[identifierStart:tokenizer.cursor], Identifier, tokenizer.cursor, tokenizer.lineNumber), nil
			}
			err := "Invalid Charecter @ Line Int : " + strconv.Itoa(tokenizer.lineNumber) + "; Cursor Int : " + strconv.Itoa(tokenizer.cursor)
			return Token{}, errors.New(err)
		}
	}
	return CreateToken("NL", EndOfStatment, tokenizer.cursor, tokenizer.lineNumber), nil
}

func createOperatorToken(char rune, tokenizer *tokenizer) Token {

	switch char {
	case '+':
		return CreateToken(string(char), Add, tokenizer.cursor, tokenizer.lineNumber)

	case '-':
		return CreateToken(string(char), Subtract, tokenizer.cursor, tokenizer.lineNumber)

	case '*':
		return CreateToken(string(char), Multiply, tokenizer.cursor, tokenizer.lineNumber)

	case '/':
		return CreateToken(string(char), Divide, tokenizer.cursor, tokenizer.lineNumber)

	case '(':
		return CreateToken(string(char), Openbrack, tokenizer.cursor, tokenizer.lineNumber)

	case ')':
		return CreateToken(string(char), Closebrack, tokenizer.cursor, tokenizer.lineNumber)

	case '^':
		return CreateToken(string(char), Exspo, tokenizer.cursor, tokenizer.lineNumber)
	}
	return Token{}

}
