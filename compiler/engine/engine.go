package engine

import (
	"compiler/analyzer"
	"compiler/symboltable"
	"compiler/writer"
	"fmt"
	"log"
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
	LabelCount     int
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
	if e.Current.Key == "identifier" {
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
		op := e.Current.Content
		e.Forward(1)
		e.CompileTerm()
		if op == "-" {
			op = "neg"
		}
		e.WriteArithmetic(op)
	} else {
		switch e.Current.Key {
		case "integerConstant":
			i, _ := strconv.Atoi(e.Current.Content)
			e.WritePush("constant", i)
		case "keyword":
			switch e.Current.Content {
			case "true":
				e.WritePush("constant", 0)
				e.WriteArithmetic("~")
			case "false":
				e.WritePush("constant", 0)
			case "null":
			case "this":
			}
		case "stringConstant":
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
	if e.NextToken().Content == ";" {
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
	// create 2 labels and increment LabelCount by 1
	l1 := fmt.Sprintf("WHILE_EXP_%d", e.LabelCount)
	l2 := fmt.Sprintf("WHILE_END_%d", e.LabelCount)
	e.LabelCount++

	// write label L1
	e.WriteLabel(l1)

	// evaluate expression, and reverse the result
	e.Forward(1)
	e.CompileExpression()
	e.WriteArithmetic("~")

	// write if-goto L2 and increment LabelCount by 1
	e.WriteIf(l2)

	// compile statements
	e.Forward(1)
	if e.Current.Content != "{" {
		log.Fatal("While statements must start with '{'.")
	} else {
		e.Forward(1)
		e.CompileStatements()
	}

	// write goto L1
	e.WriteGoto(l1)

	// write label L2
	e.WriteLabel(l2)
}

func (e *Engine) CompileIf() {
	// create 2 labels and increment LabelCount by 1
	l1 := fmt.Sprintf("IF_FALSE_%d", e.LabelCount)
	l2 := fmt.Sprintf("IF_TRUE_%d", e.LabelCount)
	e.LabelCount++

	// evaluate expression and reverse the result
	e.Forward(1)
	e.CompileExpression()
	e.WriteArithmetic("~")

	// write if-goto L1.
	e.WriteIf(l1)

	// compile statements
	e.Forward(1)
	if e.Current.Content != "{" {
		log.Fatal("If statements must start with '{'.")
	} else {
		e.Forward(1)
		e.CompileStatements()
	}

	// write goto L2
	e.WriteGoto(l2)

	// write label L1
	e.WriteLabel(l1)

	// compile statements if 'else' keyword is found
	if e.NextToken().Content == "else" {
		e.Forward(2)
		if e.Current.Content != "{" {
			log.Fatal("Else statements must start with '{'.")
		} else {
			e.Forward(1)
			e.CompileStatements()
		}
	}

	// write label L2
	e.WriteLabel(l2)
}

func (e *Engine) CompileLet() {
	// save variable name
	e.Forward(1)
	varName := e.Current.Content

	e.Forward(1)
	if e.Current.Content != "=" {
		// should handle array element here.
	} else {
		e.CompileExpression()
	}

	s, ok := e.Table.Local[varName]
	if !ok {
		s, ok = e.Table.Global[varName]
		if !ok {
			log.Fatal("Variable not found in both local and global symbol tables.")
		}
	}
	
	e.WritePop(s.Kind, s.Index)
}

func (e *Engine) CompileStatements() {
	for e.Current.Content != "}" {
		switch e.Current.Content {
		case "let": e.CompileLet()
		case "if": e.CompileIf()
		case "while": e.CompileWhile()
		case "do": e.CompileDo()
		case "return": e.CompileReturn()
		}
		e.Forward(1)
	}
}

func (e *Engine) CompileVarDec() {
	e.Forward(1)
	varType := e.Current.Content

	e.Forward(1)
	for e.Current.Content != ";" {
		e.Table.Define(e.Current.Content, varType, "local")
		e.Forward(1)
		if e.Current.Content == "," {
			e.Forward(1)
		}
	}
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
	e.CompileStatements()
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
			e.Table.Define(c, argType, "argument")
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
		Tokenizer:      c,
		Table:          symboltable.New(),
		Writer:         writer.New(),
		Current:        c.Tokenized[3],
		Index:          3,
		ClassName:      c.Tokenized[1].Content,
		SubroutineName: "",
		LabelCount:     0,
	}
}
