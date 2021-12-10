package tools

import (
	"strconv"
	"strings"
)

type Code struct {
	original, asmCode, mCode []byte
}

var vCount int64 = 16

// serachVal searches if given key exists in vmap. If not, the variable will be assigned
func (c *Code) searchVal(i, j int) int64 {
	v, found := vmap[string(c.asmCode[i:j])]
	if !found {
		v = vCount
		vmap[string(c.asmCode[i:j])] = vCount
		vCount++
	}
	return v
}

// parseA translates A instruction into binary code
func (c *Code) parseA(i, j int) {
	cmd := []byte{'0','0','0','0','0','0','0','0','0','0','0','0','0','0','0','0', '\n'}
	v, err := strconv.ParseInt(string(c.asmCode[i:j]), 10, 16)
	if err != nil {
		v = c.searchVal(i, j)
	}
	bc := strconv.FormatInt(v, 2)
	for k, el := range bc {
		cmd[len(cmd) - len(bc) + k - 1] = byte(el)
	}
	c.mCode = append(c.mCode, cmd...)
}

// parseC translates C instruction into binary code
func (c *Code) parseC(i, j int) {
	dest := "null"
	comp := ""
	jump := "null"

	for k := i; k < j; k++ {
		switch c.asmCode[k] {
		case '=':
			dest = string(c.asmCode[i:k])
			i = k + 1
			comp = string(c.asmCode[i:j])
		case ';': comp = string(c.asmCode[i:k])
		case 'J': jump = string(c.asmCode[k:j])
		}
	}

	c.mCode = append(c.mCode, '1', '1', '1')
	if strings.Contains(comp, "M") {
		c.mCode = append(c.mCode, '1')
	} else {
		c.mCode = append(c.mCode, '0')
	}

	c.mCode = append(c.mCode, cmap[comp]...)
	c.mCode = append(c.mCode, dmap[dest]...)
	c.mCode = append(c.mCode, jmap[jump]...)
	c.mCode = append(c.mCode, '\n')
}

func Decode(asm []byte) []byte {
	c := Code{
		original: asm,
		asmCode: []byte{},
		mCode: []byte{},
	}

	c.organize()

	for i := 0; i < len(c.asmCode); i++ {
		var j int
		for j = i; j < len(c.asmCode); j++ {
			if c.asmCode[j] == '\n' {
				if c.asmCode[i] == '@' {
					c.parseA(i + 1, j)
				} else {
					c.parseC(i, j)
				}
				break
			}
		}
		i = j
	}

	return c.mCode[:len(c.mCode) - 1]
}