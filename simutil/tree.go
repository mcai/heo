package simutil

import "fmt"

func _printNode(node interface{}, prefix string, tail bool, value func(node interface{}) interface{}, children func(node interface{}) []interface{}) {
	var line string

	if tail {
		line = "└── "
	} else {
		line = "├── "
	}

	fmt.Printf("%s%s%s\n", prefix, line, value(node))

	if len(children(node)) > 0 {
		for i := 0; i < len(children(node)) - 1; i++ {
			var childNode = children(node)[i]

			if tail {
				line = "    "
			} else {
				line = "│   "
			}

			_printNode(childNode, fmt.Sprintf("%s%s", prefix, line), false, value, children)
		}
		if len(children(node)) >= 1 {
			var lastNode = children(node)[len(children(node)) - 1]

			if tail {
				line = "    "
			} else {
				line = "│   "
			}

			_printNode(lastNode, fmt.Sprintf("%s%s", prefix, line), true, value, children)
		}
	}
}

func PrintNode(node interface{}, value func(node interface{}) interface{}, children func(node interface{}) []interface{}) {
	_printNode(node, "", true, value, children)
}
