package main

// Token stores the data struct for tokens
type Token struct {
	class string
	value string
}

// Tokens stores the class for iterating on the token list
type Tokens struct {
	tokens []Token
	cursor int64
}

// Peek shows you the next token without increasing the cursor
func (t Tokens) Peek() (Token, error) {
	if int64(len(t.tokens)) >= t.cursor+2 {
		return t.tokens[t.cursor+2], nil
	}
	return Token{}, ErrorMap["EmptyToken"]
}

// LookAhead allows you to view any tokens ahead of the current cursor
func (t Tokens) LookAhead(step int64) (Token, error) {
	if int64(len(t.tokens)) >= t.cursor+1+step {
		return t.tokens[t.cursor+1+step], nil
	}
	return Token{}, ErrorMap["EmptyToken"]
}

// Next gives you the token after the cursor
func (t Tokens) Next() (Token, error) {
	if int64(len(t.tokens)) >= t.cursor+2 {
		t.cursor++
		return t.tokens[t.cursor+1], nil
	}
	return Token{}, ErrorMap["EmptyToken"]
}
