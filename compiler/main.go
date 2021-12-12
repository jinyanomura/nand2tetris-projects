package main

import (
	"bufio"
	"compiler/analyzer"
	"flag"
	"fmt"
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
		c := analyzer.NewCompiler(jack)
		c.Tokenize()
		writeFile(path[:len(path) - 5], &c)
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
				c := analyzer.NewCompiler(jack)
				c.Tokenize()
				c.CompileClass()
				writeFile(fmt.Sprintf("%s/%s", path, fileName[:len(fileName) - 5]), &c)
			}
		}
	}
}

func writeFile(path string, c *analyzer.Code) {
	xml, err := os.Create(fmt.Sprintf("%sJ.xml", path))
	if err != nil {
		log.Fatal(err)
	}
	defer xml.Close()
	fw := bufio.NewWriter(xml)
	x := strings.Join(c.XML, "\n")
	if _, err = fw.Write([]byte(x)); err != nil {
		log.Fatal(err)
	}
	if err = fw.Flush(); err != nil {
		log.Fatal(err)
	}
}