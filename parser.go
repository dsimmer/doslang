package main

// Parser holds all the functions for processing tokens
type Parser struct {
	cursor int
	tokens Tokens
}

func (p Parser) parseExpression() (bool, error) {
	// construct a full on tree here
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
