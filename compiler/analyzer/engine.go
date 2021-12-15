package analyzer

import "fmt"

var op = []string{"+", "-", "*", "/", "&", "|", "<", ">", "="}

func (c *Code) CompileClass() {
	var i int
	var t Token

	c.XML = append(c.XML, "<class>")

	for i = 0; i < len(c.Tokenized); i++ {
		t = c.Tokenized[i]
		if t.start == "<keyword>" {
			switch t.content {
			case "class": c.appendTerminal(t)
			case "static", "field": i = c.compileClassVarDec(i)
			case "function", "method", "constructor": i = c.compileSubroutine(i)
			}
		} else {
			c.appendTerminal(t)
		}
	}

	c.XML = append(c.XML, "</class>")
}

func (c *Code) compileSubroutine(i int) int {
	var j int
	var t Token

	c.XML = append(c.XML, "<subroutineDec>")
	for j = i; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		if t.content == "{" {
			j = c.compileSubroutineBody(j)
			break
		} else if t.content == "(" {
			j = c.compileParameterList(j)
		} else {
			c.appendTerminal(t)
		}
	}
	c.XML = append(c.XML, "</subroutineDec>")

	return j
}

func (c *Code) compileParameterList(i int) int {
	var j int
	var t = c.Tokenized[i]

	// append "("
	c.appendTerminal(t)

	c.XML = append(c.XML, "<parameterList>")

	for j = i + 1; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		if t.content == ")" {
			break
		} else {
			c.appendTerminal(t)
		}
	}

	c.XML = append(c.XML, "</parameterList>")

	// append ")"
	c.appendTerminal(t)
	
	return j
}

func (c *Code) compileSubroutineBody(i int) int {
	var j int
	var t Token

	c.XML = append(c.XML, "<subroutineBody>")

	// append "{"
	c.appendTerminal(c.Tokenized[i])

	for j = i + 1; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		if t.content == "}" {
			c.appendTerminal(t)
			break
		} else if t.start == "<keyword>" {
			switch t.content {
			case "var": j = c.compileVarDec(j)
			case "let", "if", "while", "do", "return": j = c.compileStatements(j)
			}
		}
	}

	c.XML = append(c.XML, "</subroutineBody>")

	return j
}

func (c *Code) compileStatements(i int) int {
	var j int
	var t Token

	c.XML = append(c.XML, "<statements>")

	for j = i; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		if t.content == "}" {
			break
		}
		switch t.content {
		case "let": j = c.compileLet(j)
		case "if": j = c.compileIf(j)
		case "while": j = c.compileWhile(j)
		case "do": j = c.compileDo(j)
		case "return": j = c.compileReturn(j)
		}
	}

	c.XML = append(c.XML, "</statements>")

	return j - 1
}

func (c *Code) compileExpression(i int, endSymbol string) int {
	var j, numExpLayer int
	var t Token

	c.XML = append(c.XML, "<expression>")

	out:
	for j = i; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		if numExpLayer == 0 {
			switch {
			case t.content == endSymbol:
				c.compileTerm(i, j)
				break out
			case t.content == "(" || t.content == "[":
				numExpLayer++
			case isOp(t.content) && i != j:
				c.compileTerm(i, j)
				c.appendTerminal(t)
				i = j + 1
			}
		} else {
			switch t.content {
			case "(", "[": numExpLayer++
			case ")", "]": numExpLayer--
			}
		}
	}

	c.XML = append(c.XML, "</expression>")

	return j - 1
}

func (c *Code) compileTerm(i, j int) int {
	var t Token
	var nt Token

	c.XML = append(c.XML, "<term>")

	for i < j {
		t = c.Tokenized[i]
		c.appendTerminal(t)
		if t.start == "<identifier>" {
			nt = c.Tokenized[i+1] 
			switch nt.content {
			case "[":
				c.appendTerminal(nt)
				i = c.compileExpression(i+2, "]")
			case "(":
				c.appendTerminal(nt)
				i = c.compileExpressionList(i+2)
			}
		} else if t.content == "(" {
			i = c.compileExpression(i+1, ")")
		} else if t.content == "-" || t.content == "~" {
			i = c.compileTerm(i+1, j)
		}
		i++
	}

	c.XML = append(c.XML, "</term>")

	return i
}

func (c *Code) compileReturn(i int) int {
	var j int
	var t Token

	c.XML = append(c.XML, "<returnStatement>")

	// append "return"
	c.appendTerminal(c.Tokenized[i])
	
	for j = i + 1; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		if t.content == ";" {
			c.appendTerminal(t)
			break
		} else {
			j = c.compileExpression(j, ";")
		}
	}
	
	c.XML = append(c.XML, "</returnStatement>")

	return j
}

func (c *Code) compileDo(i int) int {
	var j int
	var t Token

	c.XML = append(c.XML, "<doStatement>")

	out:
	for j = i; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		c.appendTerminal(t)
		switch t.content {
		case ";": break out
		case "(": j = c.compileExpressionList(j+1)
		}
	}

	c.XML = append(c.XML, "</doStatement>")

	return j
}

func (c *Code) compileExpressionList(i int) int {
	var j int
	var t Token

	c.XML = append(c.XML, "<expressionList>")

	if c.Tokenized[i].content != ")" {
		out:
		for j = i; j < len(c.Tokenized); j++ {
			t = c.Tokenized[j]
			switch t.content {
			case ",":
				j = c.compileExpression(i, ",") + 1
				c.appendTerminal(t)
				i = j + 1
			case ")":
				if c.Tokenized[j+1].content == ";" {
					j = c.compileExpression(i, ")")
					break out
				}
			}
		}
	} else {
		j = i - 1
	}

	c.XML = append(c.XML, "</expressionList>")

	return j
}

func (c *Code) compileWhile(i int) int {
	var j int
	var t Token

	c.XML = append(c.XML, "<whileStatement>")

	out:
	for j = i; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		c.appendTerminal(t)
		switch t.content {
		case "}": break out
		case "(": j = c.compileExpression(j+1, ")")
		case "{": j = c.compileStatements(j+1)
		}
	}

	c.XML = append(c.XML, "</whileStatement>")

	return j
}

func (c *Code) compileIf(i int) int {
	var j int
	var t Token

	c.XML = append(c.XML, "<ifStatement>")

	for j = i; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		c.appendTerminal(t)
		if t.content == "}" && c.Tokenized[j+1].content != "else" {
			break
		}
		switch t.content {
		case "(": j = c.compileExpression(j+1, ")")
		case "{": j = c.compileStatements(j+1)
		}
	}

	c.XML = append(c.XML, "</ifStatement>")

	return j
}

func (c *Code) compileLet(i int) int {
	var j int
	var t Token

	c.XML = append(c.XML, "<letStatement>")

	// left part
	for j = i; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		c.appendTerminal(t)
		if t.content == "[" {
			j = c.compileExpression(j+1, "]")
		} else if t.content == "=" {
			break
		}
	}

	// right part
	j = c.compileExpression(j+1, ";") + 1

	// append ";"
	c.appendTerminal(c.Tokenized[j])

	c.XML = append(c.XML, "</letStatement>")

	return j
}

func (c *Code) compileVarDec(i int) int {
	var j int
	var t Token

	c.XML = append(c.XML, "<varDec>")
	for j = i; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		c.appendTerminal(t)
		if t.content == ";" {
			break
		}
	}
	c.XML = append(c.XML, "</varDec>")

	return j
}

func (c *Code) compileClassVarDec(i int) int {
	var j int
	var t Token
	
	c.XML = append(c.XML, "<classVarDec>")
	for j = i; j < len(c.Tokenized); j++ {
		t = c.Tokenized[j]
		c.appendTerminal(t)
		if t.content == ";" {
			break
		}
	}
	c.XML = append(c.XML, "</classVarDec>")

	return j
}

func (c *Code) appendTerminal(t Token) {
	s := fmt.Sprintf("%s %s %s", t.start, t.content, t.end)
	c.XML = append(c.XML, s)
}

func isOp(c string) bool {
	for i := 0; i < len(op); i++ {
		if c == op[i] {
			return true
		}
	}
	return false
}