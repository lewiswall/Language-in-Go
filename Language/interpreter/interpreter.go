package interpreter

import (
	"fmt"
	"language/tree"
)

func Interpret(treee []tree.Node) {

	for _, node := range treee {
		_, err := node.Evaluate()
		if err != nil {
			fmt.Println(err)
		}
	}

}
