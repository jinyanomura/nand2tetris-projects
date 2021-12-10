package analyzer

var Symbols = []byte{'{', '}', '(', ')', '[', ']', '.', ',', ';', '+', '-', '*', '/', '&', '|', '<', '>', '=', '~'}

var Keywords = map[string]Token{
	"class":       {start: "<keyword>", content: "class", end: "</keyword>"},
	"constructor": {start: "<keyword>", content: "constructor", end: "</keyword>"},
	"function":    {start: "<keyword>", content: "function", end: "</keyword>"},
	"method":      {start: "<keyword>", content: "method", end: "</keyword>"},
	"field":       {start: "<keyword>", content: "field", end: "</keyword>"},
	"static":      {start: "<keyword>", content: "static", end: "</keyword>"},
	"var":         {start: "<keyword>", content: "var", end: "</keyword>"},
	"int":         {start: "<keyword>", content: "int", end: "</keyword>"},
	"char":        {start: "<keyword>", content: "char", end: "</keyword>"},
	"boolean":     {start: "<keyword>", content: "boolean", end: "</keyword>"},
	"void":        {start: "<keyword>", content: "void", end: "</keyword>"},
	"true":        {start: "<keyword>", content: "true", end: "</keyword>"},
	"false":       {start: "<keyword>", content: "false", end: "</keyword>"},
	"null":        {start: "<keyword>", content: "null", end: "</keyword>"},
	"this":        {start: "<keyword>", content: "this", end: "</keyword>"},
	"let":         {start: "<keyword>", content: "let", end: "</keyword>"},
	"do":          {start: "<keyword>", content: "do", end: "</keyword>"},
	"if":          {start: "<keyword>", content: "if", end: "</keyword>"},
	"else":        {start: "<keyword>", content: "else", end: "</keyword>"},
	"while":       {start: "<keyword>", content: "while", end: "</keyword>"},
	"return":      {start: "<keyword>", content: "return", end: "</keyword>"},
}
