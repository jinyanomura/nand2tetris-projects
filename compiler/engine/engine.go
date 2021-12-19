package engine

import (
	"compiler/analyzer"
	"compiler/symboltable"
	"compiler/writer"
	"fmt"
	"strconv"
)

type Engine struct {
	*analyzer.Tokenizer
	*symboltable.Table
	*writer.Writer
	Current        analyzer.Token
	Index          int
	ClassName      string
	SubroutineName string
}

var opcode = []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}

func (e *Engine) CompileExpressionList() int {
	if e.NextToken().Content == ")" {
		return 0
	}

	numExp := 0
	for e.Current.Content != ")" {
		e.CompileExpression()
		numExp++
	}

	return numExp
}

func (e *Engine) CompileTerm() {
	if e.Current.Key == "<identifier>" {
		name := e.Current.Content
		nt := e.Tokenizer.Tokenized[e.Index + 1]
		switch nt.Content {
		case "[":
			// compilation for an array element
		case "(":
			e.WriteCall(name, e.CompileExpressionList())
		case ".":
			e.Forward(2)
			name = fmt.Sprintf("%s.%s", name, e.Current.Content)
			e.Forward(1)
			e.WriteCall(name, e.CompileExpressionList())
		default:
			kind, index := e.Table.KindOf(e.Current.Content), e.Table.IndexOf(e.Current.Content)
			if index > -1 {
				e.WritePush(kind, index)
			}
		}
	} else if e.Current.Content == "(" {
		e.CompileExpression()
	} else if e.Current.Content == "-" || e.Current.Content == "~" {
		// compilation for unary term
	} else {
		// compilation for constants
		switch e.Current.Key {
		case "<integerConstant>":
			i, _ := strconv.Atoi(e.Current.Content)
			e.WritePush("constant", i)
		case "<stringConstant>":
		case "<keywordConstant>":
		}
	}
}

func (e *Engine) CompileExpression() {
	// compile the first term
	e.Forward(1)
	e.CompileTerm()

	// compile other terms as long as the current token is an opcode.
	e.Forward(1)
	for isOp(e.Current) {
		op := e.Current.Content
		e.Forward(1)
		e.CompileTerm()
		e.WriteArithmetic(op)
		e.Forward(1)
	}
}

func (e *Engine) CompileReturn() {
	e.Forward(1)

	if e.Current.Content == ";" {
		e.WritePush("constant", 0)
	} else {
		e.CompileExpression()
	}

	e.WriteReturn()
}

func (e *Engine) CompileDo() {
	e.CompileExpression()

	if e.Current.Content == ";" {
		e.WritePop("temp", 0)
	} else {
		fmt.Println("Do statement must end with ';'.")
	}
}

func (e *Engine) CompileWhile() {

}

func (e *Engine) CompileIf() {

}

func (e *Engine) CompileLet() {

}

func (e *Engine) CompileStatements() {
	switch e.Current.Content {
	case "let":
	case "if":
	case "while":
	case "do": e.CompileDo()
	case "return": e.CompileReturn()
	}
}

func (e *Engine) CompileVarDec() {

}

func (e *Engine) CompileSubroutineBody() {
	e.Forward(1)

	// compile var declaration.
	for e.Current.Content == "var" {
		e.CompileVarDec()
		e.Forward(1)
	}

	// write function code with total number of local variables.
	e.WriteFunction(e.SubroutineName, e.Table.Count.Var)

	// compile statements.
	for e.Current.Content != "}" {
		e.CompileStatements()
		e.Forward(1)
	}
}

func (e *Engine) CompileParameterList() {
	e.Forward(1)
	argType, c := "", ""

	for e.Current.Content != ")" {
		c = e.Current.Content
		if argType == "" {
			argType = c
		} else if c == "," {
			argType = ""
		} else {
			e.Table.Define(c, argType, "arg")
		}
		e.Forward(1)
	}
}

func (e *Engine) CompileSubroutine() {
	e.Table.Reset()

	switch e.Current.Content {
	case "function":
		e.Forward(2)
		e.SubroutineName = fmt.Sprintf("%s.%s", e.ClassName, e.Current.Content)
		e.Forward(1)
		e.CompileParameterList()
	case "constructor":
	case "method":
	}

	e.Forward(1)
	if e.Current.Content == "{" {
		e.CompileSubroutineBody()
	} else {
		fmt.Println("Subroutine body must start with '{'.")
	}
}

// CompileClassVarDec registers declared variables to Global symbol map.
func (e *Engine) CompileClassVarDec() {
	k := e.Current.Content

	e.Forward(1)
	t := e.Current.Content

	for e.Current.Content != ";" {
		if e.Current.Key == "identifier" {
			e.Table.Define(e.Current.Content, t, k)
		}
		e.Forward(1)
	}
}

// CompileClass starts compiling its contents. Should be called at first.
func (e *Engine) CompileClass() {
	for e.Current.Content != "}" {
		switch e.Current.Content {
		case "static", "field": e.CompileClassVarDec()
		case "constructor", "function", "method": e.CompileSubroutine()
		}
		e.Forward(1)
	}
}

// NextToken returns the next token.
func (e *Engine) NextToken() analyzer.Token {
	return e.Tokenized[e.Index + 1]
}

// Forward progresses the current token and its index for given number.
func (e *Engine) Forward(i int) {
	e.Index += i
	e.Current = e.Tokenizer.Tokenized[e.Index]
}

// isOp checks if the content of given token is an opcode.
func isOp(t analyzer.Token) bool {
	for _, el := range opcode {
		if t.Content == el {
			return true
		}
	}
	return false
}

// New initializes and returns a new instance of Engine struct.
func New(c *analyzer.Tokenizer) *Engine {
	return &Engine{
		Tokenizer: c,
		Table:     symboltable.New(),
		Writer:    writer.New(),
		Current:   c.Tokenized[3],
		Index:     3,
		ClassName: c.Tokenized[1].Content,
	}
}
