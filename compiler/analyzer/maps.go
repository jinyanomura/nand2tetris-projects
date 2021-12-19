package analyzer

var Symbols = []byte{'{', '}', '(', ')', '[', ']', '.', ',', ';', '+', '-', '*', '/', '&', '|', '<', '>', '=', '~'}

var Keywords = map[string]Token{
	"class":       {Key: "<keyword>", Content: "class", end: "</keyword>"},
	"constructor": {Key: "<keyword>", Content: "constructor", end: "</keyword>"},
	"function":    {Key: "<keyword>", Content: "function", end: "</keyword>"},
	"method":      {Key: "<keyword>", Content: "method", end: "</keyword>"},
	"field":       {Key: "<keyword>", Content: "field", end: "</keyword>"},
	"static":      {Key: "<keyword>", Content: "static", end: "</keyword>"},
	"var":         {Key: "<keyword>", Content: "var", end: "</keyword>"},
	"int":         {Key: "<keyword>", Content: "int", end: "</keyword>"},
	"char":        {Key: "<keyword>", Content: "char", end: "</keyword>"},
	"boolean":     {Key: "<keyword>", Content: "boolean", end: "</keyword>"},
	"void":        {Key: "<keyword>", Content: "void", end: "</keyword>"},
	"true":        {Key: "<keyword>", Content: "true", end: "</keyword>"},
	"false":       {Key: "<keyword>", Content: "false", end: "</keyword>"},
	"null":        {Key: "<keyword>", Content: "null", end: "</keyword>"},
	"this":        {Key: "<keyword>", Content: "this", end: "</keyword>"},
	"let":         {Key: "<keyword>", Content: "let", end: "</keyword>"},
	"do":          {Key: "<keyword>", Content: "do", end: "</keyword>"},
	"if":          {Key: "<keyword>", Content: "if", end: "</keyword>"},
	"else":        {Key: "<keyword>", Content: "else", end: "</keyword>"},
	"while":       {Key: "<keyword>", Content: "while", end: "</keyword>"},
	"return":      {Key: "<keyword>", Content: "return", end: "</keyword>"},
}
