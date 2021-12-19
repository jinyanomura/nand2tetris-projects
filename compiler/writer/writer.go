package writer

import "fmt"

type Writer struct {
	VM          []string
	LabelNumber int
}

var command = map[string]string{
	"+": "add",
	"-": "sub",
	"=": "eq",
	">": "gt",
	"<": "lt",
	"&": "and",
	"|": "or",
	"~": "not",
}

// WritePush writes a VM push function.
func (w *Writer) WritePush(seg string, index int) {
	w.VM = append(w.VM, fmt.Sprintf("push %s %d", seg, index))
}

// WritePop writes a VM pop function.
func (w *Writer) WritePop(seg string, index int) {
	w.VM = append(w.VM, fmt.Sprintf("pop %s %d", seg, index))
}

// WriteArithmetic writes a VM arithmetic-logical command.
func (w *Writer) WriteArithmetic(cmd string) {
	c, ok := command[cmd]
	if ok {
		w.VM = append(w.VM, c)
	} else {
		switch cmd {
		case "neg": // unary negation command.
		case "*":
			w.VM = append(w.VM, "call Math.multiply 2")
		case "/": // should call Math.devide function
		}
	}
}

func (w *Writer) WriteLabel(label string) {

}

func (w *Writer) WriteGoto(label string) {

}

func (w *Writer) WriteIf(label string) {

}

// WriteCall writes a VM call command.
func (w *Writer) WriteCall(name string, numArgs int) {
	w.VM = append(w.VM, fmt.Sprintf("call %s %d", name, numArgs))
}

// WriteFunction writes a VM function command.
func (w *Writer) WriteFunction(name string, numVars int) {
	w.VM = append(w.VM, fmt.Sprintf("function %s %d", name, numVars))
}

// WriteReturn writes a VM return command.
func (w *Writer) WriteReturn() {
	w.VM = append(w.VM, "return")
}

func New() *Writer {
	return &Writer{
		VM: make([]string, 0),
		LabelNumber: 0,
	}
}