package tools

var cmdCount int64

// organize refactors the original assembler code
func (c *Code) organize() {
	comment := false
	for i := 0; i < len(c.original); i++ {
		if !comment {
			switch c.original[i] {
			case '/': comment = true
			case '\n':
				if len(c.asmCode) != 0 && c.asmCode[len(c.asmCode)-1] != '\n' {
					c.asmCode = append(c.asmCode, '\n')
					cmdCount++
				}
			case '(': i = c.registerLabel(i)
			case ' ':
			case '\r':
			default: c.asmCode = append(c.asmCode, c.original[i])
			}
		} else if c.original[i] == '\n' {
			comment = false
			if len(c.asmCode) != 0 && c.asmCode[len(c.asmCode)-1] != '\n' {
				c.asmCode = append(c.asmCode, '\n')
				cmdCount++
			}
		}
	}
	if c.asmCode[len(c.asmCode) - 1] != '\n' {
		c.asmCode = append(c.asmCode, '\n')
	}
}

// registerLabel registers a label and its line number to vmap
func (c *Code) registerLabel(i int) int {
	var k int
	for k = i + 1; k < len(c.original); k++ {
		if c.original[k] == ')' {
			break
		}
	}
	vmap[string(c.original[i+1:k])] = cmdCount
	return k
}