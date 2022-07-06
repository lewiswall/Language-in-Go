package interpreter

import (
	"fmt"
	tree "language/syntax_tree"
)

func Interpret(treee []tree.Node) {
	for _, node := range treee {
		_, err := node.Evaluate()
		if err != nil {
			fmt.Println(err)
		}
	}

}
