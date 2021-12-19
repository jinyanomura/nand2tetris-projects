package analyzer

import (
	"strconv"
)

type Tokenizer struct {
	Jack []byte
	Tokenized []Token
}

type Token struct {
	Key, Content string
}

var Symbols = []byte{'{', '}', '(', ')', '[', ']', '.', ',', ';', '+', '-', '*', '/', '&', '|', '<', '>', '=', '~'}

var Keywords = map[string]Token{
	"class":       {Key: "keyword", Content: "class"},
	"constructor": {Key: "keyword", Content: "constructor"},
	"function":    {Key: "keyword", Content: "function"},
	"method":      {Key: "keyword", Content: "method"},
	"field":       {Key: "keyword", Content: "field"},
	"static":      {Key: "keyword", Content: "static"},
	"var":         {Key: "keyword", Content: "var"},
	"int":         {Key: "keyword", Content: "int"},
	"char":        {Key: "keyword", Content: "char"},
	"boolean":     {Key: "keyword", Content: "boolean"},
	"void":        {Key: "keyword", Content: "void"},
	"true":        {Key: "keyword", Content: "true"},
	"false":       {Key: "keyword", Content: "false"},
	"null":        {Key: "keyword", Content: "null"},
	"this":        {Key: "keyword", Content: "this"},
	"let":         {Key: "keyword", Content: "let"},
	"do":          {Key: "keyword", Content: "do"},
	"if":          {Key: "keyword", Content: "if"},
	"else":        {Key: "keyword", Content: "else"},
	"while":       {Key: "keyword", Content: "while"},
	"return":      {Key: "keyword", Content: "return"},
}

// Tokenize organizes and transfers given .jack source code into .xml code.
func (c *Tokenizer) Tokenize() {
	var el byte
	var isComment, apiComment = false, false
	for i := 0; i < len(c.Jack); i++ {
		el = c.Jack[i]
		if !isComment && !apiComment {
			if el != ' ' && el != '\r' && el != '\n' && el != '\t' {
				switch {
				case el == '/' && c.Jack[i+1] == '/': isComment = true
				case el == '/' && c.Jack[i+1] == '*': apiComment = true
				case isSymbol(el): c.regSymbol(el)
				case isIntConst(el): i = c.regIntConst(i)
				case el == '"': i = c.regStrConst(i+1)
				case el == ' ' || el == '\r' || el == '\n' || el == '\t':
				default: i = c.regKeyVar(i)
				}
			}
		} else if isComment && el == '\n' {
			isComment = false
		} else if apiComment && el == '*' && c.Jack[i+1] == '/' {
			apiComment = false
			i++
		}
	}
}

// regSymbol appends given symbol token.
func (c *Tokenizer) regSymbol(s byte) {
	c.Tokenized = append(c.Tokenized, Token{Key: "symbol", Content: string(s)})
}

// regIntConst reads .jack code and appends its integer constant starts from given index.
func (c *Tokenizer) regIntConst(i int) int {
	var j int
	for j = i + 1; j < len(c.Jack); j++ {
		_, err := strconv.ParseInt(string(c.Jack[j]), 10, 16)
		if err != nil {
			break
		}
	}
	t := Token{Key: "integerConstant", Content: string(c.Jack[i:j])}
	c.Tokenized = append(c.Tokenized, t)
	return j - 1
}

// regStrConst reads .jack code and appends its string constant starts from given index.
func (c *Tokenizer) regStrConst(i int) int {
	var j = i
	for j < len(c.Jack) && c.Jack[j] != '"' {
		j++
	}
	t := Token{Key: "stringConstant", Content: string(c.Jack[i:j])}
	c.Tokenized = append(c.Tokenized, t)
	return j
}

// regKeyVar classifies a token starts from given index into Keyword or Identifier.
func (c *Tokenizer) regKeyVar(i int) int {
	var j = i
	for c.Jack[j] != ' ' && !isSymbol(c.Jack[j]) {
		j++
	}
	word := string(c.Jack[i:j])
	t, isKey := Keywords[word]
	if !isKey {
		t = Token{Key: "identifier", Content: word}
	}
	c.Tokenized = append(c.Tokenized, t)	
	return j - 1
}

// isSymbol chacks if given character is a symbol
func isSymbol(b byte) bool {
	for _, el := range Symbols {
		if b == el {
			return true
		}
	}
	return false
}

// isIntConst checks if given character represents an integer.
func isIntConst(b byte) bool {
	_, err := strconv.ParseInt(string(b), 10, 0)
	return err == nil
}

// New initializes and returns Tokenizer struct with given .jack file.
func New(jack []byte) Tokenizer {
	return Tokenizer{
		Jack:      jack,
		Tokenized: make([]Token, 0),
	}
}