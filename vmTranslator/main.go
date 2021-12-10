package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"myvmt/tools"
	"os"
)

func main() {
	// Parse the file given from command line.
	flag.Parse()
	path := flag.Arg(0)
	c := tools.NewTranslator(path)
	
	// Opens the given directory or file and pass the contents to readFile.
	var vm *os.File
	var dirName string
	files, err := os.ReadDir(path)
	// if a single .vm file
	if err != nil {
		if vm, err = os.Open(path); err != nil {
			log.Fatal(err)
		}
		readFile(vm, &c)
		defer vm.Close()
		path = path[:len(path)-3]
	// if directory containing several .vm files
	} else {
		i := len(path) - 1
		for i >= 0 && path[i] != '/' {
			i--
		}
		dirName = path[i+1:]
		c.WriteBootStrap()
		for _, file := range files {
			fileName := file.Name()
			if fileName[len(fileName) - 3:] == ".vm" {
				c.Name = fileName[:len(fileName) - 3]
				vm, err = os.Open(fmt.Sprintf("%s/%s", path, fileName))
				if err != nil {
					log.Fatal(err)
				}
				readFile(vm, &c)
				defer vm.Close()
			}
		}
		path = fmt.Sprintf("%s/%s", path, dirName)
	}

	// Creates a new file to write the translated code.
	asm, err := os.Create(fmt.Sprintf("%s.asm", path))
	if err != nil {
		log.Fatal(err)
	}
	defer asm.Close()
	fw := bufio.NewWriter(asm)
	if _, err = fw.Write(c.Acode); err != nil {
		log.Fatal(err)
	}
	if err = fw.Flush(); err != nil {
		log.Fatal(err)
	}
}

func readFile(f *os.File, c *tools.Code) {
	s := bufio.NewScanner(f)
	for s.Scan() {
		c.Extract(s.Bytes())
		if len(c.Vcmd) > 0 {
			c.AppendComment()
			t, i := c.ParseCommand()
			switch t {
			case "C_ARITHMETIC": c.WriteArithmetic(c.Vcmd[:i])
			case "C_PUSH": c.WritePush(c.ParseArgs(i+1))
			case "C_POP": c.WritePop(c.ParseArgs(i+1))
			case "C_BRANCH": c.WriteBranch(c.Vcmd[:i], c.Vcmd[i+1:])
			case "C_FUNCTION": c.WriteFunction(c.ParseArgs(i+1))
			case "C_CALL": c.WriteCall(c.ParseArgs(i+1))
			case "C_RETURN": c.WriteReturn()
			}
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
}