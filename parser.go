package main

// Parser holds all the functions for processing tokens
type Parser struct {
}

// ProcessText orchestrates the parsing of a string
func (p Parser) ProcessText(text string) string {
	lexer := Lexer{}
	_, _ = lexer.LexText(text)
	return ""
}
