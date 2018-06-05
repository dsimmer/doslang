package main

import "unicode"

// Lexer orchestrates the tokenizing of the file
type Lexer struct {
	inputStream *InputStream
	cursor      int64
	result      Tokens
}

// http://blog.leahhanson.us/post/recursecenter2016/recipeparser.html
// { type: "punc", value: "(" }           // punctuation: parens, comma, semicolon etc.
//     { type: "num", value: 5 }              // numbers
//     { type: "str", value: "Hello World!" } // strings
//     { type: "kw", value: "lambda" }        // keywords
//     { type: "var", value: "a" }            // identifiers
//     { type: "op", value: "!=" }            // operators

// First off, skip over whitespace.
// If input.eof() then return null.
// If it's a sharp sign (#), skip comment (retry after the end of line).
// If it's a quote then read a string.
// If it's a digit, then we proceed to read a number.
// If it's a “letter”, then read an identifier or a keyword token.
// If it's one of the punctuation characters, return a punctuation token.
// If it's one of the operator characters, return an operator token.
// If none of the above, error out with input.croak()
func isNumber(char string) bool {
	if len(char) > 0 {
		return unicode.IsNumber(rune(char[0]))
	}
	return false
}
func isString(char string) bool {
	switch char {
	case `"`:
		fallthrough
	case `'`:
		return true
	default:
		return false
	}
}
func isKeywordOrIdentifier(char string) bool {
	if len(char) > 0 {
		return unicode.IsLetter(rune(char[0]))
	}
	return false
}
func isOperator(char string) bool {
	switch char {
	case "!":
		fallthrough
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

func isKeyword(word string) bool {
	switch word {
	case "if":
		fallthrough
	case "for":
		fallthrough
	case "else":
		fallthrough
	case "elseif":
		fallthrough
	case "print":
		fallthrough
	case "var":
		return true
	default:
		return false
	}
}

func isPunct(char string) bool {
	switch char {
	case " ":
		fallthrough
	case "(":
		fallthrough
	case "{":
		fallthrough
	case "}":
		return true
	case ")":
		return true
	default:
		return false
	}
}

func (l *Lexer) consumeNumber() error {
	value := l.inputStream.Curr()
	for {
		char, err := l.inputStream.Peek()
		if err != nil {
			return err
		}
		if isNumber(char) {
			value = value + char
			l.inputStream.Next()
		} else {
			l.result.tokens = append(l.result.tokens, Token{class: "number", value: value})
			return nil
		}
	}
}

func (l *Lexer) consumeString() error {
	endChar := l.inputStream.Curr()
	var value string
	for {
		char, err := l.inputStream.Next()
		if err != nil {
			return err
		}
		if char == endChar {
			l.result.tokens = append(l.result.tokens, Token{class: "string", value: value})
			return nil
		}
		value = value + char
	}
}

func (l *Lexer) consumeOperator() error {
	value := l.inputStream.Curr()
	// todo investigate whether we should consume two at once sometimes? e.g. !=
	// char, err := l.inputStream.Peek()
	// if err != nil {
	// 	return err
	// }
	// if isOperator(char) {
	// 	value = value + char
	// 	l.inputStream.Next()
	// }
	l.result.tokens = append(l.result.tokens, Token{class: "operator", value: value})
	return nil
}

func (l *Lexer) consumeKeywordOrIdentifier() error {
	value := l.inputStream.Curr()
	for {
		char, err := l.inputStream.Peek()
		if err != nil {
			return err
		}
		if isKeywordOrIdentifier(char) {
			value = value + char
			l.inputStream.Next()
		} else {
			break
		}
	}
	if isKeyword(value) {
		l.result.tokens = append(l.result.tokens, Token{class: "keyword", value: value})
	} else {
		l.result.tokens = append(l.result.tokens, Token{class: "identifier", value: value})
	}
	return nil
}

func (l *Lexer) consumePunct() error {
	// _, err = l.inputStream.Next()
	// Just dont do anything and it will skip to the next item
	//todo add in a scope token
	return nil
}

func (l *Lexer) read() (bool, error) {
	var err error
	char := l.inputStream.Curr()

	switch true {
	case isNumber(char):
		err = l.consumeNumber()
		break
	case isString(char):
		err = l.consumeString()
		break
	case isKeywordOrIdentifier(char):
		err = l.consumeKeywordOrIdentifier()
		break
	case isPunct(char):
		err = l.consumePunct()
		break
	case isOperator(char):
		err = l.consumeOperator()
		break
	default:
		return false, l.inputStream.Croak("Syntax error")
	}
	if err != nil {
		if err == ErrorMap["EOF"] {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (l *Lexer) readNext() (bool, error) {
	_, err := l.inputStream.Next()
	if err != nil {
		if err == ErrorMap["EOF"] {
			// we finished!
			return true, nil
		}
		return false, err
	}
	return l.read()
}

func (l *Lexer) readCurr() (bool, error) {
	return l.read()
}

//todo turn lextext into a parallel channel

// LexText tokenizes a given text string
func (l *Lexer) LexText(text string) (Tokens, error) {
	l.inputStream = &InputStream{}
	l.inputStream.Init(text)
	l.cursor = 0
	l.result = Tokens{tokens: []Token{}, cursor: 0}
	done, err := l.readCurr()
	if err != nil {
		return Tokens{}, err
	}
	if done {
		return l.result, nil
	}
	for {
		done, err := l.readNext()
		if err != nil {
			return Tokens{}, err
		}
		if done {
			return l.result, nil
		}
	}
}
