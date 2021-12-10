package tools

import (
	"fmt"
	"strconv"
)

type Code struct {
	Acode, Vcmd []byte
	Name string
}

func (c *Code) Extract(line []byte) {
	var i int

	// ignore comments in the source vm code
	for i = 0; i < len(line); i++ {
		if line[i] == '/' && line[i+1] == '/' {
			break
		}
	}

	// remove extra white spaces
	for i > 0 && (line[i-1] == ' ' || line[i-1] == '\t') {
		i--
	}

	c.Vcmd = line[:i]
}

// AppendComment appends the original vm code as comments in the translated assembler code.
func (c *Code) AppendComment() {
	comment := fmt.Sprintf("// %s\n", string(c.Vcmd))
	c.Acode = append(c.Acode, comment...)
}

// ParseCommand determines which command type is required.
func (c *Code) ParseCommand() (string, int) {
	var i int

	for i < len(c.Vcmd) && c.Vcmd[i] != ' ' {
		i++
	}
	
	return CommandType[string(c.Vcmd[:i])], i
}

// ParseArgs parses and returns the arguments of current command.
func (c *Code) ParseArgs(i int) (string, int) {	
	j := i

	for j < len(c.Vcmd) && c.Vcmd[j] != ' ' {
		j++
	}

	segment := string(c.Vcmd[i:j])
	i, _ = strconv.Atoi(string(c.Vcmd[j+1:]))

	return segment, i
}

// NewTranslator initializes a new instance of Code strict with given file name.
func NewTranslator(path string) Code {
	var i, j int
	j = len(path) - 1

	for i = len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			j = i
		} else if path[i] == '/' {
			break
		}
	}

	return Code{
		Acode: []byte{},
		Vcmd: []byte{},
		Name: path[i+1:j],
	}
}