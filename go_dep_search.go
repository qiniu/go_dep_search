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

var (
	onlyMain        = flag.Bool("main", false, "only show main package")
	onlyTest        = flag.Bool("test", false, "only show test package")
	chain           = flag.Bool("chain", false, "show dep chain, eg: main->package_a->b->c->d")
	unused          = flag.Bool("unused", false, "list unused packages")
	graph           = flag.Bool("graph", false, "show dep graph: -graph <from_package> <to_package>")
	graphResultFile = flag.String("o", "/tmp/dep.svg", "used with -graph, result file path, supported type: svg,png,jpg,dot")
)

func main() {
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
	log.Printf("successfully load %d packages (%d main packages, %d test packages)",
		dg.CountAll(), dg.CountMain(), dg.CountTest())
	if *unused {
		log.Println("unused packages:")
		fmt.Println(strings.Join(dg.ListUnUsed(), "\n"))
		return
	}
	if *graph {
		result := dg.SearchGraph(flag.Arg(0), flag.Arg(1))
		resultToSvg(result)
		return
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
				deps := []string{"main", p}
				if p != dep {
					deps = append(deps, dep)
				}
				fmt.Println(strings.Join(deps, " -> "))
			}
		} else if *onlyTest {
			packages := dg.SearchTest(dep)
			if len(packages) == 0 {
				log.Printf("%v not found", dep)
			}
			for _, p := range packages {
				p = strings.TrimSuffix(p, ".test")
				fmt.Println(strings.Join([]string{"test", p, dep}, " -> "))
			}
		} else {
			if dg.Exists(dep) {
				fmt.Println(strings.Join([]string{"[self]", dep}, " -> "))
			}
			packages := dg.SearchAll(dep)
			if len(packages) == 0 && !dg.Exists(dep) {
				log.Printf("%v not found", dep)
			}
			for _, p := range packages {
				name := path.Base(p)
				if dg.IsMainPackage(p) {
					name = "[main]"
				} else if dg.IsTestPackage(p) {
					name = "[test]"
					p = strings.TrimSuffix(p, ".test")
				}
				fmt.Println(strings.Join([]string{name, p, dep}, " -> "))
			}
		}
	}
}
