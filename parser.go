package main

type node struct {
	node  []*node
	class string
	name  string
	value string
	cond  []*node
}

// Parser holds all the functions for processing tokens
type Parser struct {
	cursor int
	tokens *Tokens
	AST    []*node
}

func isOperand(token Token) bool {
	switch token.class {
	case "string":
		fallthrough
	case "number":
		fallthrough
	case "operator":
		return true
	default:
		return false
	}
}

func isExpression(token Token) bool {
	switch token.class {
	case "keyword":
		fallthrough
	case "identifier":
		return true
	default:
		return false
	}
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

// Todo fix this to have a better structure for storing elseif and else

func (p *Parser) parseIf(currentNode *node) error {
	var err error
	for {
		condNode := node{}
		nodeNode := node{}
		currentNode.class = "expression"
		currentNode.cond = append(currentNode.cond, &condNode)
		currentNode.name = "if"
		currentNode.node = append(currentNode.node, &nodeNode)

		err = p.parseOperand(&condNode)
		if err != nil {
			return err
		}
		err = p.parseNext(&nodeNode)
		if err != nil {
			return err
		}
		token, err := p.tokens.Peek()
		if err != nil {
			break
		}
		if token.class == "keyword" {
			if token.value == "else" {
				_, _ = p.tokens.Next()
				newNode := node{}
				currentNode.node = append(currentNode.node, &newNode)
				err = p.parseNext(&newNode)
				if err != nil {
					return err
				}
			} else if token.value == "elseif" {
				_, _ = p.tokens.Next()
				continue
			}
			break
		}
	}
	return err
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
func (p *Parser) parsePrint(currentNode *node) error {
	token, err := p.tokens.Next()
	if err != nil {
		return err
	}
	if token.class != "string" {
		return p.tokens.Croak("expected string at " + string(p.cursor))
	}
	currentNode.class = "expression"
	currentNode.name = "print"
	currentNode.value = token.value

	return nil
}
func (p *Parser) parseVar(currentNode *node) error {
	token, err := p.tokens.Next()
	if err != nil {
		return err
	}
	if token.class != "identifier" {
		return p.tokens.Croak("expected identifier at " + string(p.cursor))
	}
	currentNode.class = "expression"
	currentNode.name = "var"
	currentNode.value = token.value

	return nil
}

func (p *Parser) parseNext(currentNode *node) error {
	_, err := p.tokens.Next()
	if err != nil {
		return err
	}
	return p.parseCurrentExpression(currentNode)
}
func (p *Parser) parseCurrentExpression(currentNode *node) error {
	token := p.tokens.Current()
	if !isExpression(token) {
		return p.tokens.Croak("Not expression")
	}
	switch token.value {
	case "if":
		return p.parseIf(currentNode)
	case "else":
		return p.tokens.Croak("Else should not be here!")
	case "elseif":
		return p.tokens.Croak("Elseif should not be here!")
	case "print":
		return p.parsePrint(currentNode)
	case "var":
		return p.parseVar(currentNode)
	default:
		break
	}
	//todo
	// if token.class == "for" {
	// 	return p.parseFor()
	// }

	return p.tokens.Croak("No Expression")
}

func (p *Parser) processOperands(currentNode *node, operands []Token) error {
	// search for the different levels of operands
	// join them up in the array (and/slice/slice) and build node from there
	// join elements that would combine LAST - i.e. first do plus minus
	// then / *
	// then &|
	// then ! =

	if len(operands) == 0 {
		return p.tokens.Croak("somehow ran out of operators")
	}

	if len(operands) == 1 {
		if operands[0].class == "number" || operands[0].class == "string" {
			currentNode.value = operands[0].value
			currentNode.name = operands[0].class
			currentNode.class = "operator"
			return nil
		}
		return p.tokens.Croak("Last operator not number or string")
	}

	items := []string{"=", "|", "&", "/", "*", "-", "+"}

	for _, item := range items {
		for index, nodeI := range operands {
			if nodeI.value == item {
				currentNode.value = nodeI.value
				currentNode.name = nodeI.class
				currentNode.class = "operator"

				beforeNode := node{}
				beforeNode.name = nodeI.value
				beforeNode.class = "operator"
				currentNode.node = append(currentNode.node, &beforeNode)
				//todo add in new node details here
				afterNode := node{}
				afterNode.name = nodeI.value
				afterNode.class = "operator"
				currentNode.node = append(currentNode.node, &afterNode)

				err := p.processOperands(&beforeNode, operands[:index])
				if err != nil {
					return err
				}
				err = p.processOperands(&afterNode, operands[index+1:])
				if err != nil {
					return err
				}
				// We are spawning a new node for each guy so just break the loop here
				return nil
			}
		}
	}

	return p.tokens.Croak("Somehow no operators present")
}

func (p *Parser) parseOperand(currentNode *node) error {
	var operands []Token
	for {
		token, err := p.tokens.Next()
		if err != nil {
			return err
		}
		if isOperand(token) {
			// push into array
			operands = append(operands, token)
		} else {
			return p.tokens.Croak("No Operator")
		}
		if err != nil {
			return err
		}
		token, err = p.tokens.Peek()
		if err != nil {
			return err
		}
		if !isOperand(token) {
			// Validate we have all the requried elements of an operand node
			err = p.processOperands(currentNode, operands)
			if err != nil {
				return err
			}
			return nil
		}
	}
}

// GenerateAST orchestrates the AST generation
func (p *Parser) GenerateAST(text string) error {
	lexer := Lexer{}
	tokens, err := lexer.LexText(text)
	if err != nil {
		panic(err)
	}
	p.tokens = &tokens

	p.cursor = 0
	currentNode := node{}
	p.AST = append(p.AST, &currentNode)

	err = p.parseCurrentExpression(&currentNode)
	if err != nil {
		if err == ErrorMap["EOF"] {
			return nil
		}
		return err
	}
	for {
		currentNode := node{}

		err = p.parseNext(&currentNode)
		if err != nil {
			if err == ErrorMap["EOF"] {
				return nil
			}
			return err
		}
		p.AST = append(p.AST, &currentNode)
	}
}
