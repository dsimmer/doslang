package main

import "errors"

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
func (t *Tokens) Peek() (Token, error) {
	if int64(len(t.tokens)) > t.cursor+1 {
		return t.tokens[t.cursor+1], nil
	}
	return Token{}, ErrorMap["EOF"]
}

// LookAhead allows you to view any tokens ahead of the current cursor
func (t *Tokens) LookAhead(step int64) (Token, error) {
	if int64(len(t.tokens)) > t.cursor+step {
		return t.tokens[t.cursor+step], nil
	}
	return Token{}, ErrorMap["EOF"]
}

// Next gives you the token after the cursor
func (t *Tokens) Next() (Token, error) {
	if int64(len(t.tokens)) > t.cursor+1 {
		t.cursor++
		return t.tokens[t.cursor], nil
	}
	return Token{}, ErrorMap["EOF"]
}

// Current gives you the current token
func (t *Tokens) Current() Token {
	return t.tokens[t.cursor]
}

// todo this should map to the cursor error, should return in a channel essentially

// Croak gives you the token after the cursor
func (t *Tokens) Croak(msg string) error {
	return errors.New(msg + " (" + string(t.cursor) + ")")
}
