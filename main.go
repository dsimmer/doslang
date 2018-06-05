package main

import (
	"errors"
	"fmt"
	"io/ioutil"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// ErrorMap holds all the standard errors for identification
var ErrorMap = map[string]error{
	"EOF":          errors.New("End of input"),
	"noExpression": errors.New("No expression"),
	"noValue":      errors.New("No value, expected a non expression in this context"),
}

// if for var + - / * = "string" 1234567890 //comment
func main() {
	content, err := ioutil.ReadFile("./input.dos")
	check(err)
	parser := Parser{}
	err = parser.GenerateAST(string(content))
	check(err)
	compiler := Compiler{}
	err = compiler.Start(parser.AST)
	check(err)
	fmt.Println("done")
}
