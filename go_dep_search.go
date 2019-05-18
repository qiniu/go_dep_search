package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
)

const usage = `Usage:

go list -json ./... | %s package_names

Args:
`

type DepInfo struct {
	ImportPath string   `json:"ImportPath"`
	Name       string   `json:"Name"`
	Deps       []string `json:"Deps"`
}

func (d DepInfo) HasDep(packageName string) bool {
	for _, v := range d.Deps {
		if v == packageName {
			return true
		}
	}
	return false
}

func main() {
	onlyMain := flag.Bool("main", false, "only show main package")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}
	dec := json.NewDecoder(os.Stdin)
	var once sync.Once
	for {
		var di DepInfo
		err := dec.Decode(&di)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if *onlyMain && di.Name != "main" {
			continue
		}
		for _, dep := range flag.Args() {
			if di.HasDep(dep) {
				once.Do(func() {
					fmt.Fprintln(os.Stderr, "Name\t->\tImportPath\t->\tdep_package")
				})
				fmt.Printf("%s\t->\t%s\t->\t%s\n", di.Name, di.ImportPath, dep)
			}
		}
	}
	once.Do(func() {
		fmt.Fprintf(os.Stderr, "package not found: %v\n", flag.Args())
	})
}
