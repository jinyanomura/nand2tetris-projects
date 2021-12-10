package analyzer

type Code struct {
	Name, XML string
}

func (c *Code) Tokenize(jack []byte) {
	c.XML = "If you see this message, then it worked as expected! Congrats :)"
}

func NewCompiler(path string) Code {
	var i, j int
	j = len(path) - 1

	for i = len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			j = i
		} else if path[i] == '/' {
			break
		}
	}

	return Code{
		Name: path[i+1:j],
		XML: "",
	}
}