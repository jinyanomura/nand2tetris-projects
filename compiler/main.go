package main

import (
	"bufio"
	"compiler/analyzer"
	"flag"
	"fmt"
	"log"
	"os"
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
				// writeFile(fmt.Sprintf("%s/%s", path, fileName[:len(fileName) - 5]), &c)
				for _, el := range c.Tokenized {
					fmt.Println(el)
				}
			}
		}
	}
}

func writeFile(path string, c *analyzer.Code) {
	xml, err := os.Create(fmt.Sprintf("%sT1.xml", path))
	if err != nil {
		log.Fatal(err)
	}
	defer xml.Close()
	fw := bufio.NewWriter(xml)
	if _, err = fw.Write([]byte(c.XML)); err != nil {
		log.Fatal(err)
	}
	if err = fw.Flush(); err != nil {
		log.Fatal(err)
	}
}