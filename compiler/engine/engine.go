package engine

import (
	"fmt"
	"jackcompiler/analyzer"
	"jackcompiler/symboltable"
	"jackcompiler/writer"
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

// CompileExpressionList compiles arguments of subroutine calls, and returns its number.
func (e *Engine) CompileExpressionList() int {
	if e.NextToken().Content == ")" {
		e.Forward(1)
		return 0
	}

	numExp := 0
	for e.Current.Content != ")" {
		e.CompileExpression()
		numExp++
	}

	return numExp
}

// CompileTerm compiles terms.
func (e *Engine) CompileTerm() {
	if e.Current.Key == "identifier" {
		name := e.Current.Content
		switch e.NextToken().Content {
		case "[":
			e.WritePush(e.KindOf(e.Current.Content), e.IndexOf(e.Current.Content))
			e.Forward(1)
			e.CompileExpression()
			e.WriteArithmetic("+")
			e.WritePop("pointer", 1)
			e.WritePush("that", 0)
		case "(":
			e.Forward(1)
			name = fmt.Sprintf("%s.%s", e.ClassName, name)
			e.WritePush("pointer", 0)
			e.WriteCall(name, e.CompileExpressionList() + 1)
		case ".":
			e.Forward(2)
			// if method call
			if e.IndexOf(name) > -1 {
				e.WritePush(e.KindOf(name), e.IndexOf(name))
				name = fmt.Sprintf("%s.%s", e.TypeOf(name), e.Current.Content)
				e.Forward(1)
				e.WriteCall(name, e.CompileExpressionList() + 1)
			} else {
				name = fmt.Sprintf("%s.%s", name, e.Current.Content)
				e.Forward(1)
				e.WriteCall(name, e.CompileExpressionList())
			}
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
			case "false", "null":
				e.WritePush("constant", 0)
			case "this":
				e.WritePush("pointer", 0)
			}
		case "stringConstant":
			e.WritePush("constant", len(e.Current.Content))
			e.WriteCall("String.new", 1)
			for _, char := range e.Current.Content {
				e.WritePush("constant", int(char))
				e.WriteCall("String.appendChar", 2)
			}
		}
	}
}

// CmpileExpression compiles expressions.
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

// CompileReturn compiles return statements.
func (e *Engine) CompileReturn() {
	if e.NextToken().Content == ";" {
		e.WritePush("constant", 0)
		e.Forward(1)
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
		log.Fatal("Do statement must end with ';'.")
	}
}

// CompileWhile compiles while statements.
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

// CompileIf compiles if statements.
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

// CompileLet compiles let statements.
func (e *Engine) CompileLet() {
	// save variable name
	e.Forward(1)
	varName := e.Current.Content

	e.Forward(1)
	switch e.Current.Content {
	case "=":
		e.CompileExpression()
		e.WritePop(e.KindOf(varName), e.IndexOf(varName))
	case "[":
		e.WritePush(e.KindOf(varName), e.IndexOf(varName))
		e.CompileExpression()
		e.WriteArithmetic("+")
		e.Forward(1)
		if e.Current.Content != "=" {
			log.Fatal(fmt.Sprintf("Expected '=' but found %s.", e.Current.Content))
		}
		e.CompileExpression()
		e.WritePop("temp", 0)
		e.WritePop("pointer", 1)
		e.WritePush("temp", 0)
		e.WritePop("that", 0)
	default:
		log.Fatal(fmt.Sprintf("Expected '=' or '[' but found %s.", e.Current.Content))
	}

	if e.Current.Content != ";" {
		log.Fatal("let statement must end with ';'.")
	}
}

// CompileStatements progresses compiling body of subroutines and if/while statements.
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

// CompileVarDec registers local variables declared in subroutines to symbol table.
func (e *Engine) CompileVarDec() {
	e.Forward(1)
	varType := e.Current.Content

	e.Forward(1)
	for e.Current.Content != ";" {
		e.Table.Define(e.Current.Content, varType, "var")
		e.Forward(1)
		if e.Current.Content == "," {
			e.Forward(1)
		}
	}
}

// CompileSubroutineBody compiles body of subroutines.
func (e *Engine) CompileSubroutineBody(funcType string) {
	e.Forward(1)

	// compile var declaration.
	for e.Current.Content == "var" {
		e.CompileVarDec()
		e.Forward(1)
	}

	// write function code with total number of local variables.
	e.WriteFunction(e.SubroutineName, e.Count.Var)

	// do object manipulation if funcType is either 'method' or 'constructor'
	switch funcType {
	case "method":
		e.WritePush("argument", 0)
		e.WritePop("pointer", 0)
	case "constructor":
		e.WritePush("constant", e.Count.Field)
		e.WriteCall("Memory.alloc", 1)
		e.WritePop("pointer", 0)
	}

	// compile statements.
	e.CompileStatements()
}

// CompileParameterList registers arguments to symbol table.
func (e *Engine) CompileParameterList(funcType string) {
	e.Forward(1)
	argType, c := "", ""

	// set 'this' as first argument to symbol table if funcType is 'method'
	if funcType == "method" {
		e.Table.Define("this", e.ClassName, "argument")
	}

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

// CompileSubroutine compiles subroutines(constructor, method, and function).
func (e *Engine) CompileSubroutine() {
	// save the type of function(function, constructor, method)
	funcType := e.Current.Content

	// reset the local symbol table
	e.Table.Reset()

	// write function and its arguments declaration
	e.Forward(2)
	e.SubroutineName = fmt.Sprintf("%s.%s", e.ClassName, e.Current.Content)
	e.Forward(1)
	e.CompileParameterList(funcType)

	// compile subroutine body
	e.Forward(1)
	if e.Current.Content == "{" {
		e.CompileSubroutineBody(funcType)
	} else {
		log.Fatal("Subroutine body must start with '{'.")
	}
}

// CompileClassVarDec registers declared variables to Global symbol map.
func (e *Engine) CompileClassVarDec() {
	// stores kind(field, static) in k
	k := e.Current.Content

	e.Forward(1)
	
	// stores variable type in t
	t := e.Current.Content

	e.Forward(1)
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
	e.Current = e.Tokenized[e.Index]
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
func New(t *analyzer.Tokenizer) *Engine {
	return &Engine{
		Tokenizer:      t,
		Table:          symboltable.New(),
		Writer:         writer.New(),
		Current:        t.Tokenized[3],
		Index:          3,
		ClassName:      t.Tokenized[1].Content,
		SubroutineName: "",
		LabelCount:     0,
	}
}
