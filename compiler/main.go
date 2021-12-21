package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"jackcompiler/analyzer"
	"jackcompiler/engine"
	"log"
	"os"
	"strings"
)

var dirPath string

func main() {
	// Parse the file given from command line.
	flag.Parse()
	path := flag.Arg(0)
	dirPath = path
	
	// Opens given directory and compiles its contents.
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		readFile(file, path)
	}

	// Compile OS libraries.
	path = "../os"
	files, err = os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		readFile(file, path)
	}
}

func readFile(file fs.DirEntry, path string) {
	fileName := file.Name()
	if fileName[len(fileName) - 5:] == ".jack" {
		jack, err := os.ReadFile(fmt.Sprintf("%s/%s", path, fileName))
		if err != nil {
			log.Fatal(err)
		}

		// tokenize
		a := analyzer.New(jack)
		a.Tokenize()

		// compile
		e := engine.New(&a)
		e.CompileClass()

		writeFile(fileName[:len(fileName) - 5], e.VM)
	}
}

func writeFile(fileName string, code []string) {
	vm, err := os.Create(fmt.Sprintf("%s/%s.vm", dirPath, fileName))
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