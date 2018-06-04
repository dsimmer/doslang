package main

import "errors"

type node struct {
	operator string
	node     []node
	name     string
	cond     []node
}

// Parser holds all the functions for processing tokens
type Parser struct {
	cursor int
	tokens Tokens
	AST    []node
}

func isOperand(token Token) bool {
	switch token.class {
	case "string":
	case "number":
	case "operator":
		return true
	default:
		return false
	}
	return false
}

func isExpression(token Token) bool {
	switch token.class {
	case "keyword":
	case "identifier":
		return true
	default:
		return false
	}
	return false
}

// function parse_if() {
//     skip_kw("if");
//     var cond = parse_expression();
//     if (!is_punc("{")) skip_kw("then");
//     var then = parse_expression();
//     var ret = { type: "if", cond: cond, then: then };
//     if (is_kw("else")) {
//         input.next();
//         ret.else = parse_expression();
//     }
//     return ret;
// }

func (p Parser) parseIf(currentNode *node) (bool, error) {
	condNode := node{}
	nodeNode := node{}
	currentNode.name = "if"
	currentNode.cond = append(currentNode.cond, condNode)
	currentNode.operator = "if"
	currentNode.node = append(currentNode.node, nodeNode)

	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	done, err := p.parseOperand(&condNode)
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	done, err = p.parseExpression(&condNode)
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}

	return done, err
}
func (p Parser) parseElseif(currentNode *node) (bool, error) {
	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func (p Parser) parseElse(currentNode *node) (bool, error) {
	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

// todo
// func (p Parser) parseFor(currentNode *node) (bool, error) {
// 	token, err := p.tokens.Next()
// 	if err != nil {
// 		if err == ErrorMap["EOF"] {
// 			return true, nil
// 		}
// 		return false, err
// 	}
// 	return false, nil
// }
func (p Parser) parsePrint(currentNode *node) (bool, error) {
	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func (p Parser) parseVar(currentNode *node) (bool, error) {
	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	if token.class != "identifier" {
		return false, errors.New("expected identifier at " + string(p.cursor))
	}
	p.AST = append(p.AST, node{
		operator: "var",
		name:     token.value,
	})
	return false, nil
}

func (p Parser) parseExpression(currentNode *node) (bool, error) {
	token, err := p.tokens.Next()
	if err != nil {
		return false, err
	}
	if !isExpression(token) {
		return false, p.tokens.Croak("Not expression")
	}
	if token.value == "if" {
		return p.parseIf(currentNode)
	}
	//todo
	// if token.class == "for" {
	// 	return p.parseFor()
	// }
	if token.value == "else" {
		return p.parseElse(currentNode)
	}
	if token.value == "elseif" {
		return p.parseElseif(currentNode)
	}
	if token.value == "print" {
		return p.parsePrint(currentNode)
	}
	if token.value == "var" {
		return p.parseVar(currentNode)
	}
	return false, p.tokens.Croak("No Expression")
}

func (p Parser) processOperands(currentNode *node, operands []Token) error {
	// search for the different levels of operands
	// join them up in the array (and/slice/slice) and build node from there
	// join elements that would combine LAST - i.e. first do plus minus
	// then / *
	// then &|
	// then ! =
	items := []string{"-", "+", "/", "*"}

	for _, item := range items {
		for index, nodeI := range operands {
			if nodeI.value == item {
				beforeNode := node{}
				beforeNode.name = nodeI.class
				//todo add in new node details here
				afterNode := node{}
				afterNode.name = nodeI.class
				//todo add in new node details here
				err := p.processOperands(&beforeNode, operands[:index])
				if err != nil {
					return err
				}
				err = p.processOperands(&afterNode, operands[index+1:])
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (p Parser) parseOperand(currentNode *node) (bool, error) {
	var operands []Token
	for {
		token, err := p.tokens.Next()
		if err != nil {
			if err == ErrorMap["EOF"] {
				return true, nil
			}
			return false, err
		}
		if isOperand(token) {
			// push into array
			operands = append(operands, token)
		} else {
			return false, p.tokens.Croak("No Operator")
		}
		if err != nil {
			return false, err
		}
		token, err = p.tokens.Peek()
		if err != nil {
			return false, err
		}
		if !isOperand(token) {
			// Validate we have all the requried elements of an operand node
			err = p.processOperands(currentNode, operands)
			if err != nil {
				return false, err
			}
			return false, nil
		}
	}
}

// GenerateAST orchestrates the AST generation
func (p Parser) GenerateAST(text string) string {
	lexer := Lexer{}
	tokens, err := lexer.LexText(text)
	if err != nil {
		panic(err)
	}
	p.tokens = tokens

	p.cursor = 0
	for {
		currentNode := node{}
		done, err := p.parseExpression(&currentNode)
		if err != nil {
			return ""
		}
		if done {
			return ""
		}
	}
}
