package main

import (
	"errors"
	"io/ioutil"
	"os"
)

// Compiler holds the new file to be written and the compilation methods
type Compiler struct {
	content string
}

func isSingleOperator(currNode node) bool {
	switch currNode.value {
	case "!":
		return true
	default:
		return false
	}
}

func isMultiOperator(currNode node) bool {
	switch currNode.value {
	case "=":
		fallthrough
	case "+":
		fallthrough
	case "-":
		fallthrough
	case "/":
		fallthrough
	case "*":
		fallthrough
	case "&":
		fallthrough
	case "|":
		return true
	default:
		return false
	}
}

func (c *Compiler) processOperators(currNode *node) (string, error) {
	if isSingleOperator(*currNode) {
		switch currNode.value {
		case "!":
			result := "!("
			content, err := c.processOperators(currNode.node[0])
			if err != nil {
				return "", nil
			}
			return result + content + ")", nil
		default:
			return "", ErrorMap["Expected !"]
		}
	} else if isMultiOperator(*currNode) {
		switch currNode.value {
		case "&":
			fallthrough
		case "|":
			fallthrough
		case "=":
			beforeContent, err := c.processOperators(currNode.node[0])
			if err != nil {
				return "", err
			}
			afterContent, err := c.processOperators(currNode.node[1])
			if err != nil {
				return "", err
			}
			return "(" + beforeContent + ") " + currNode.value + currNode.value + " (" + afterContent + ")", nil
		case "+":
			fallthrough
		case "-":
			fallthrough
		case "/":
			fallthrough
		case "*":
			beforeContent, err := c.processOperators(currNode.node[0])
			if err != nil {
				return "", err
			}
			afterContent, err := c.processOperators(currNode.node[1])
			if err != nil {
				return "", err
			}
			return "(" + beforeContent + ") " + currNode.value + " (" + afterContent + ")", nil
		default:
			return "", ErrorMap["Expected multi node operator"]
		}
	} else if currNode.name == "number" {
		return currNode.value, nil
	}
	return "", ErrorMap["Expected number or operator"]
}

func (c *Compiler) processNode(currNode *node) error {
	switch currNode.name {
	case "if":
		for index := range currNode.node {
			if index == 0 {
				// handle initial if
				c.content = c.content + "if "
				operators, err := c.processOperators(currNode.cond[index])
				if err != nil {
					return err
				}
				c.content = c.content + operators + " {\n"
				c.processNode(currNode.node[index])
				c.content = c.content + "}"
			} else if index >= len(currNode.cond) {
				//handle else
				c.content = c.content + " else {\n"
				c.processNode(currNode.node[index])
				c.content = c.content + "}"
			} else {
				//handle elseif
				c.content = c.content + " elseif "
				operators, err := c.processOperators(currNode.cond[index])
				if err != nil {
					return err
				}
				c.content = c.content + operators + " {\n"
				c.processNode(currNode.node[index])
				c.content = c.content + "}"
			}
			if len(currNode.node)-1 == index {
				c.content = c.content + "\n"
			}
		}
		break
	case "print":
		c.content = c.content + `fmt.Println("` + currNode.name + `")` + "\n"
		break
	case "var":
		c.content = c.content + `fmt.Println("` + currNode.name + `")` + "\n"
		break
	default:
		return errors.New("Invalid AST")
	}
	// for index, astNode := range currNode.node {

	// }
	return nil
}

// Start begins the compilation process
func (c *Compiler) Start(AST []*node) error {
	c.content = `
package main

import (
    "fmt"
)

func main () {
`
	for _, astNode := range AST {
		err := c.processNode(astNode)
		if err != nil {
			return err
		}
	}
	c.content = c.content + "\n}"

	ioutil.WriteFile("./built/compiled.go", []byte(c.content), os.FileMode(777))
	return nil
}
