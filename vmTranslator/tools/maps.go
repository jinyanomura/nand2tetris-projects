package tools

var CommandType = map[string]string{
	"add": "C_ARITHMETIC",
	"sub": "C_ARITHMETIC",
	"neg": "C_ARITHMETIC",
	"eq": "C_ARITHMETIC",
	"gt": "C_ARITHMETIC",
	"lt": "C_ARITHMETIC",
	"and": "C_ARITHMETIC",
	"or": "C_ARITHMETIC",
	"not": "C_ARITHMETIC",
	"push": "C_PUSH",
	"pop": "C_POP",
	"label": "C_BRANCH",
	"goto": "C_BRANCH",
	"if-goto": "C_BRANCH",
	"function": "C_FUNCTION",
	"call": "C_CALL",
	"return": "C_RETURN",
}

var Segments = map[string]string{
	"local": "LCL",
	"argument": "ARG",
	"this": "THIS",
	"that": "THAT",
	"temp": "TEMP",
	"R5": "R5",
}