package main

import "errors"

type node struct {
	operator     string
	codeBlock    string
	altCodeBlock string
	name         string
}

// Parser holds all the functions for processing tokens
type Parser struct {
	cursor int
	tokens Tokens
	AST    []node
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

func (p Parser) parseIf() (bool, error) {
	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func (p Parser) parseElseif() (bool, error) {
	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func (p Parser) parseElse() (bool, error) {
	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func (p Parser) parseFor() (bool, error) {
	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func (p Parser) parsePrint() (bool, error) {
	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func (p Parser) parseVar() (bool, error) {
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
		operator:     "var",
		codeBlock:    "",
		altCodeBlock: "",
		name:         token.value,
	})
	return false, nil
}

func (p Parser) parseExpression() (bool, error) {
	// construct a full on tree here
	token, err := p.tokens.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	if token.class == "if" {
		return p.parseIf()
	}
	if token.class == "for" {
		return p.parseFor()
	}
	if token.class == "else" {
		return p.parseElse()
	}
	if token.class == "elseif" {
		return p.parseElseif()
	}
	if token.class == "print" {
		return p.parsePrint()
	}
	if token.class == "var" {
		return p.parseVar()
	}
	return false, nil
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
		done, err := p.parseExpression()
		if err != nil {
			return ""
		}
		if done {
			return ""
		}
	}
}
