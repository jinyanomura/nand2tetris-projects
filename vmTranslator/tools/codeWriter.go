package tools

import "fmt"

var labelNo int
var asm, l1, l2, jump string
var spf string = "@SP\nA=M\nM=D\n@SP\nM=M+1"

func (c *Code) WritePush(seg string, i int) {
	switch seg {
	case "local", "argument", "this", "that": pushWithIndex(seg, i)
	case "temp": asm = fmt.Sprintf("@%d\nD=A\n@R5\nA=A+D\nD=M\n%s\n", i, spf)
	case "constant": asm = fmt.Sprintf("@%d\nD=A\n%s\n", i, spf)
	case "pointer":
		if i == 0 {
			pushNoIndex("THIS")
		} else {
			pushNoIndex("THAT")
		}
	case "static":
		name := fmt.Sprintf("%s.%d", c.Name, i)
		pushNoIndex(name)
	}

	c.Acode = append(c.Acode, asm...)
}

func (c *Code) WritePop(seg string, i int) {
	switch seg {
	case "local", "argument", "this", "that": popWithIndex(seg, i)
	case "temp": asm = fmt.Sprintf("@R5\nD=A\n@%d\nD=A+D\n@SP\nAM=M-1\nD=D+M\nA=D-M\nM=D-A\n", i)
	case "pointer":
		if i == 0 {
			popNoIndex("THIS")
		} else {
			popNoIndex("THAT")
		}
	case "static":
		name := fmt.Sprintf("%s.%d", c.Name, i)
		popNoIndex(name)
	}

	c.Acode = append(c.Acode, asm...)
}

func (c *Code) WriteArithmetic(cmd []byte) {
	switch string(cmd) {
	case "add": asm = "@SP\nAM=M-1\nD=M\nA=A-1\nM=D+M\n"
	case "sub": asm = "@SP\nAM=M-1\nD=M\nA=A-1\nM=M-D\n"
	case "neg": asm = "@SP\nA=M-1\nM=-M\n"
	case "eq", "gt", "lt": asm = translateComp(string(cmd))
	case "and": asm = "@SP\nAM=M-1\nD=M\nA=A-1\nM=D&M\n"
	case "or": asm = "@SP\nAM=M-1\nD=M\nA=A-1\nM=D|M\n"
	case "not": asm = "@SP\nA=M-1\nM=!M\n"
	}

	c.Acode = append(c.Acode, asm...)
}

func (c *Code) WriteBranch(cmd, l []byte) {
	label := string(l)
	switch string(cmd) {
	case "label": asm = fmt.Sprintf("(%s)\n", label)
	case "goto": asm = fmt.Sprintf("@%s\n0;JMP\n", label)
	case "if-goto": asm = fmt.Sprintf("@SP\nAM=M-1\nD=M\n@%s\nD;JNE\n", label)
	}

	c.Acode = append(c.Acode, asm...)
}

func (c *Code) WriteFunction(funcName string, numVar int) {
	c.Acode = append(c.Acode, fmt.Sprintf("(%s)\n", funcName)...)
	for i := 0; i < numVar; i++ {
		c.WritePush("constant", 0)
	}
}

func (c *Code) WriteCall(funcName string, numArgs int) {
	ret := fmt.Sprintf("%s$ret.%d", funcName, labelNo)
	labelNo++

	pushRetAdd := fmt.Sprintf("@%s\nD=A\n%s", ret, spf)
	pushSeg := fmt.Sprintf("@LCL\nD=M\n%s\n@ARG\nD=M\n%s\n@THIS\nD=M\n%s\n@THAT\nD=M\n%s", spf, spf, spf, spf)
	setARG := fmt.Sprintf("@SP\nD=M\n@5\nD=D-A\n@%d\nD=D-A\n@ARG\nM=D", numArgs)
	setLCL := "@SP\nD=M\n@LCL\nM=D"
	gotoF := fmt.Sprintf("@%s\n0;JMP", funcName)
	setRetLabel := fmt.Sprintf("(%s)", ret)

	asm = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n", pushRetAdd, pushSeg, setARG, setLCL, gotoF, setRetLabel)
	c.Acode = append(c.Acode, asm...)
}

func (c *Code) WriteReturn() {
	asm = "@LCL\nD=M\n@frame\nM=D\n@5\nA=D-A\nD=M\n@retAddr\nM=D\n@SP\nAM=M-1\nD=M\n@ARG\nA=M\nM=D\n@ARG\nD=M+1\n@SP\nM=D\n@frame\nAM=M-1\nD=M\n@THAT\nM=D\n@frame\nAM=M-1\nD=M\n@THIS\nM=D\n@frame\nAM=M-1\nD=M\n@ARG\nM=D\n@frame\nAM=M-1\nD=M\n@LCL\nM=D\n@retAddr\nA=M\n0;JMP\n"
	c.Acode = append(c.Acode, asm...)
}

func (c *Code) WriteBootStrap() {
	asm = "// bootstrap\n@256\nD=A\n@SP\nM=D\n"
	c.Acode = append(c.Acode, asm...)
	c.WriteCall("Sys.init", 0)
}

func pushWithIndex(seg string, i int) {
	seg = Segments[seg]
	asm = fmt.Sprintf("@%d\nD=A\n@%s\nA=D+M\nD=M\n%s\n", i, seg, spf)
}

func pushNoIndex(v string) {
	asm = fmt.Sprintf("@%s\nD=M\n%s\n", v, spf)
}

func popWithIndex(seg string, i int) {
	seg = Segments[seg]
	asm = fmt.Sprintf("@%s\nD=M\n@%d\nD=A+D\n@SP\nAM=M-1\nD=D+M\nA=D-M\nM=D-A\n", seg, i)
}

func popNoIndex(v string) {
	asm = fmt.Sprintf("@SP\nAM=M-1\nD=M\n@%s\nM=D\n", v)
}

func translateComp(condition string) string {
	l1 = fmt.Sprintf("COMP_TRUE_%d", labelNo)
	l2 = fmt.Sprintf("SKIP_%d", labelNo)
	labelNo++

	switch condition {
	case "eq": jump = "JEQ"
	case "gt": jump = "JGT"
	case "lt": jump = "JLT"
	}

	return fmt.Sprintf("@SP\nAM=M-1\nD=M\nA=A-1\nD=M-D\n@%s\nD;%s\n@SP\nA=M-1\nM=0\n@%s\n0;JMP\n(%s)\n@SP\nA=M-1\nM=-1\n(%s)\n", l1, jump, l2, l1, l2)
}