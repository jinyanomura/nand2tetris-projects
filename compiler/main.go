package main

import (
	"bufio"
	"flag"
	"fmt"
	"jackcompiler/analyzer"
	"jackcompiler/engine"
	"log"
	"os"
	"strings"
)

func main() {
	// Parse the file given from command line.
	flag.Parse()
	path := flag.Arg(0)
	
	// Opens the given directory or file and pass the contents to readFile.
	var jack []byte
	files, err := os.ReadDir(path)
	// if a single .jack file
	if err != nil {
		if jack, err = os.ReadFile(path); err != nil {
			log.Fatal(err)
		}
		a := analyzer.New(jack)
		a.Tokenize()
		e := engine.New(&a)
		e.CompileClass()
		writeFile(path[:len(path) - 5], e.VM)
	// if directory containing several .jack files
	} else {
		i := len(path) - 1
		for i >= 0 && path[i] != '/' {
			i--
		}
		for _, file := range files {
			fileName := file.Name()
			if fileName[len(fileName) - 5:] == ".jack" {
				jack, err = os.ReadFile(fmt.Sprintf("%s/%s", path, fileName))
				if err != nil {
					log.Fatal(err)
				}
				// tokenize
				a := analyzer.New(jack)
				a.Tokenize()

				// compile
				e := engine.New(&a)
				e.CompileClass()

				writeFile(fmt.Sprintf("%s/%s", path, fileName[:len(fileName) - 5]), e.VM)
			}
		}
	}
}

func writeFile(path string, code []string) {
	vm, err := os.Create(fmt.Sprintf("%s.vm", path))
	if err != nil {
		log.Fatal(err)
	}
	defer vm.Close()
	fw := bufio.NewWriter(vm)
	x := strings.Join(code, "\n")
	if _, err = fw.Write([]byte(x)); err != nil {
		log.Fatal(err)
	}
	if err = fw.Flush(); err != nil {
		log.Fatal(err)
	}
}