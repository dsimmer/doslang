package main

import (
	"errors"
	"strings"
)

// InputStream handles char input for us
type InputStream struct {
	input []string
	pos   int
	line  int
	col   int
}

// Init the inputstream
func (i *InputStream) Init(text string) {
	i.input = strings.Split(text, "\n")
	i.pos = 0
	i.line = 0
	i.col = 0
}

//todo must be a better way of doing this

// Next get the next char in the stream
func (i *InputStream) Next() (string, error) {
	i.pos++
	if len(i.input[i.line]) > i.col+1 {
		// inc column
		i.col++
		return string(i.input[i.line][i.col]), nil
	} else if len(i.input) > i.line+1 {
		// inc line, we reached the end of col
		i.line++
		i.col = 0
		return string(i.input[i.line][i.col]), nil
	}
	return "", ErrorMap["EOF"]
}

// Curr get the current char in the stream
func (i *InputStream) Curr() string {
	return string(i.input[i.line][i.col])
}

// Peek but dont increment the counter
func (i *InputStream) Peek() (string, error) {
	if len(i.input[i.line]) > i.col+1 {
		return string(i.input[i.line][i.col+1]), nil
	} else if len(i.input) > i.line+1 {
		return string(i.input[i.line+1][0]), nil
	}
	return "", ErrorMap["EOF"]

}

// EOF we are done
func (i *InputStream) EOF() bool {
	_, err := i.Peek()
	return err != nil
}

// Croak gives the current line of the error
func (i *InputStream) Croak(msg string) error {
	return errors.New(msg + " (" + string(i.line) + ":" + string(i.col) + ")")
}
