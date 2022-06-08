package LanErrs

import (
	"language/tokenizer"
	"strconv"
)

type MultipleErrors struct {
	Errors []error
}

func (e MultipleErrors) Error() {

}

type IncompatibleTypeError struct {
	Token tokenizer.Token
}

func (e IncompatibleTypeError) Error() string {
	return "ERROR: Incompatible types on each side of operation : " + tokenizer.TKString(e.Token.Kind) + " at : line number " +
		strconv.Itoa(e.Token.LineNum) + "; cursor number " + strconv.Itoa(e.Token.Cursor)
}

type NoIdentifierAvailableError struct {
	Identifier string
}

func (e NoIdentifierAvailableError) Error() string {
	return "ERROR: Cannot find identifier \"" + e.Identifier + "\""
}

type ExpectedBoolError struct {
	Token tokenizer.Token
}

func (e ExpectedBoolError) Error() string {
	return "ERROR: Expected return type \"Bool\" on either side of token at line num :" + strconv.Itoa(e.Token.LineNum) +
		", cursor :" + strconv.Itoa(e.Token.Cursor)
}

type ExpectedBoolWithControlError struct {
	Token tokenizer.Token
}

func (e ExpectedBoolWithControlError) Error() string {
	return "ERROR: Expected Bool type after control statement at line num :" + strconv.Itoa(e.Token.LineNum) +
		", cursor :" + strconv.Itoa(e.Token.Cursor)
}

type MustBeNumWithComparisonOp struct {
	Token tokenizer.Token
}

func (e MustBeNumWithComparisonOp) Error() string {
	return "ERROR: Must use numbers on either side of Comparison operator at line num :" + strconv.Itoa(e.Token.LineNum) +
		", cursor :" + strconv.Itoa(e.Token.Cursor)
}

type WrongTypeUsedWithBinOpError struct {
	Token tokenizer.Token
}

func (e WrongTypeUsedWithBinOpError) Error() string {
	return "ERROR: Wrong Types used with Binary Operation at line num :" + strconv.Itoa(e.Token.LineNum) +
		", cursor :" + strconv.Itoa(e.Token.Cursor)
}

type ExpectedIdentifierError struct {
	Token tokenizer.Token
}

func (e ExpectedIdentifierError) Error() string {
	return "ERROR: Expected identifier next to Delete Operation at line num :" + strconv.Itoa(e.Token.LineNum) +
		", cursor :" + strconv.Itoa(e.Token.Cursor)
}
