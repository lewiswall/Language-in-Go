package main

import (
	"bufio"
	"fmt"
	"language/Format"
	"language/Global"
	"language/ShuntingYard"
	"language/interpreter"
	"language/parser"
	"language/tokenizer"
	"log"
	"os"
	"strconv"
)

func main() {
	var tokens []tokenizer.Token

	f, err := os.Open("testfiles/controlTest")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	myTokenizer := tokenizer.New()

	for scanner.Scan() {
		myTokenizer.NewLine(scanner.Text())
		token, err := myTokenizer.Get()
		if err != nil {
			fmt.Println(err)
			break
		}
		tokens = append(tokens, token)
		for token.Kind != tokenizer.EndOfStatment {
			token, err = myTokenizer.Get()

			if err != nil {
				fmt.Println(err)
				break
			} else {
				tokens = append(tokens, token)
			}
		}
	}
	tokens = append(tokens, tokenizer.CreateToken("END", tokenizer.End, 0, 0))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	FormatChecker := Format.NewFormatChecker(tokens)
	FormatChecker.FormatTokens()
	tokens = FormatChecker.Tokens
	//for _, token := range tokens {
	//	fmt.Println(token.Text)
	//}

	handleStrings(&tokens)
	postFixTokens := ShuntingYard.ShuntingY{Tokens: tokens, Index: 0}
	postFixTokens.ToPostFix()
	post := postFixTokens.Result

	parseTree := parser.NewParser(post)
	parseTree.EvaluateToken()
	tree := parseTree.ParsedLines

	interpreter.Interpret(tree)

}

func handleStrings(tokens *[]tokenizer.Token) {
	for index, token := range *tokens {
		if token.Kind == tokenizer.String {
			Global.Strings = append(Global.Strings, token.Text)
			address := len(Global.Strings) - 1
			(*tokens)[index].Text = strconv.Itoa(address)
		} else if token.Kind == tokenizer.Identifier {
			if !contains(Global.GlobalVarNames, token.Text) {
				Global.GlobalVarNames = append(Global.GlobalVarNames, token.Text)
				address := len(Global.GlobalVarNames) - 1
				(*tokens)[index].Text = strconv.Itoa(address)
				continue
			}
			address := findIndex(token.Text)
			(*tokens)[index].Text = strconv.Itoa(address)
		}
	}
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func findIndex(varName string) int {
	for p, v := range Global.GlobalVarNames {
		if v == varName {
			return p
		}
	}
	return -1
}
