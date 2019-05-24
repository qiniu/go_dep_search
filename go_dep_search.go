package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/ma6174/go_dep_search/depgraph"
)

const usage = `Usage:

go list -json all | %s package_names

Args:
`

func main() {
	onlyMain := flag.Bool("main", false, "only show main package")
	chain := flag.Bool("chain", false, "show dep chained")
	unused := flag.Bool("unused", false, "list unused packages")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() == 0 && !*unused {
		flag.Usage()
		return
	}
	if *chain {
		*onlyMain = true
	}
	dg, err := depgraph.LoadDeps(os.Stdin)
	if err != nil {
		log.Fatalln("LoadDeps failed", err)
	}
	log.Printf("successfuly load %d packages (%d main packages)", dg.CountAll(), dg.CountMain())
	if *unused {
		log.Println("unused packages:")
		fmt.Println(strings.Join(dg.ListUnUsed(), "\n"))
	}
	for _, dep := range flag.Args() {
		if *chain {
			chains := dg.SearchChain(dep)
			if len(chains) == 0 {
				log.Printf("%v not found", dep)
			}
			for _, chain := range chains {
				fmt.Println(strings.Join(chain, " -> "))
			}
		} else if *onlyMain {
			packages := dg.SearchMain(dep)
			if len(packages) == 0 {
				log.Printf("%v not found", dep)
			}
			for _, p := range packages {
				fmt.Println(strings.Join([]string{"main", p, dep}, " -> "))
			}
		} else {
			packages := dg.SearchAll(dep)
			if len(packages) == 0 {
				log.Printf("%v not found", dep)
			}
			for _, p := range packages {
				name := "main"
				if !dg.IsMainPackage(p) {
					name = path.Base(p)
				}
				fmt.Println(strings.Join([]string{name, p, dep}, " -> "))
			}
		}
	}
}
