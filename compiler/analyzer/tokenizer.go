package analyzer

import "strconv"

type Code struct {
	Jack []byte
	Tokenized []Token
	XML string
}

type Token struct {
	start, content, end string
}

// Tokenize organizes and transfers given .jack source code into .xml code.
func (c *Code) Tokenize() {
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
		}
	}
}

// regSymbol appends given symbol token.
func (c *Code) regSymbol(s byte) {
	c.Tokenized = append(c.Tokenized, Token{start: "<symbol>", content: string(s), end: "</symbol>"})
}

// regIntConst reads .jack code and appends its integer constant starts from given index.
func (c *Code) regIntConst(i int) int {
	var j int
	for j = i + 1; j < len(c.Jack); j++ {
		_, err := strconv.ParseInt(string(c.Jack[j]), 10, 16)
		if err != nil {
			break
		}
	}
	t := Token{start: "<integerConstant>", content: string(c.Jack[i:j]), end: "</integerConstant>"}
	c.Tokenized = append(c.Tokenized, t)
	return j - 1
}

// regStrConst reads .jack code and appends its string constant starts from given index.
func (c *Code) regStrConst(i int) int {
	var j = i
	for j < len(c.Jack) && c.Jack[j] != '"' {
		j++
	}
	t := Token{start: "<stringConstant>", content: string(c.Jack[i:j]), end: "</stringConstant>"}
	c.Tokenized = append(c.Tokenized, t)
	return j
}

// regKeyVar classifies a token starts from given index into Keyword or Identifier.
func (c *Code) regKeyVar(i int) int {
	var j = i
	for c.Jack[j] != ' ' && !isSymbol(c.Jack[j]) {
		j++
	}
	word := string(c.Jack[i:j])
	t, isKey := Keywords[word]
	if !isKey {
		t = Token{start: "<identifier>", content: word, end: "</identifier>"}
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
	_, err := strconv.ParseInt(string(b), 10, 4)
	return err == nil
}

// NewCompiler initializes and returns Code struct with given .jack file.
func NewCompiler(jack []byte) Code {
	return Code{
		Jack: jack,
		Tokenized: make([]Token, 0),
		XML: "",
	}
}