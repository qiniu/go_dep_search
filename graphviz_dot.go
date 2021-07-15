// +build !graphviz

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func write(f *os.File, data string) {
	_, err := f.WriteString(data)
	if err != nil {
		log.Panicln(err)
	}
}

func createHeader(f *os.File) {
	write(f, fmt.Sprintln("digraph \"\" {"))
	write(f, fmt.Sprintln("\tgraph [ rankdir=LR ];"))
}

func createNode(name string, f *os.File, nodes map[string]struct{}) {
	if _, ok := nodes[name]; ok {
		return
	}
	write(f, fmt.Sprintf("\t\"%s\"\t[ style=bold ];\n", name))
	nodes[name] = struct{}{}
}

func createLine(from, to string, f *os.File) {
	format := "\t\"%s\" -> \"%s\" [ key=\"%s\", URL=\"#%s->%s\", style=bold ];\n"
	write(f, fmt.Sprintf(format, from, to, from, from, to))
}

func createFooter(f *os.File) {
	write(f, fmt.Sprintln("}"))
}

func resultToSvg(result map[string][]string) {
	fn := *graphResultFile
	if !strings.HasSuffix(fn, ".dot") {
		fn += ".dot"
	}
	f, err := os.Create(fn)
	if err != nil {
		log.Panicln("open result file failed", result, err)
	}
	defer f.Close()

	createHeader(f)
	nodes := make(map[string]struct{})
	for from, tos := range result {
		createNode(from, f, nodes)
		for _, to := range tos {
			createNode(to, f, nodes)
			createLine(from, to, f)
		}
	}
	createFooter(f)

	err = f.Close()
	if err != nil {
		log.Panicln(err)
	}

	if *graphResultFile == fn {
		fmt.Println("result saved to " + *graphResultFile)
	} else {
		fmt.Println("\ninstall graphviz and run the following command to generate result file:")
		format := strings.Trim(filepath.Ext(*graphResultFile), ".")
		fmt.Printf("dot -T%v %v -o %v\n\n", format, fn, *graphResultFile)
	}
}
