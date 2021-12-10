package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"myasm/tools"
)

func main() {
	decoded := tools.Decode(readFile())
	
	name := giveName()

	writeFile(name, decoded)
}

// readFile reads and returns the contents of the file specified by command line input
func readFile() []byte {
	flag.Parse()
	file := flag.Arg(0)
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

// writeFile creates a new file with given name and content
func writeFile(name string, content []byte) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fw :=bufio.NewWriter(f)
	_, err = fw.Write(content)
	if err != nil {
		log.Fatal(err)
	}

	err = fw.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

// giveName replaces .asm with .hack, and returns a new file name
func giveName() string {
	cdl := -1
	name := flag.Arg(0)

	for i, el := range name {
		if el == '/' {
			cdl = i
		} else if el == '.' {
			name = name[cdl+1:i] + ".hack"
			break
		}
	}

	return name
}