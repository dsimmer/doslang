package main

import (
	"errors"
	"io/ioutil"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// ErrorMap holds all the standard errors for identification
var ErrorMap = map[string]error{
	"EmptyToken": errors.New("No more tokens"),
	"EOF":        errors.New("End of input"),
}

// if for var + - / * = "string" 1234567890 //comment
func main() {
	content, err := ioutil.ReadFile("./input.dos")
	check(err)
	parser := Parser{}
	_ = parser.ProcessText(string(content))
}
