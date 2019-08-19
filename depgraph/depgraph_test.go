package depgraph

import (
	"os"
	"testing"
)

func TestDepGraph(t *testing.T) {
	f, err := os.Open("../testdata/go1.12.5_deps.json")
	if err != nil {
		panic(err)
	}
	dg, err := LoadDeps(f)
	if err != nil {
		panic(err)
	}
	if !dg.Exists("fmt") {
		t.Error("fmt should exists")
	}
	if dg.Exists("fmtxxxxxxx") {
		t.Error("fmtxxxxxxx should not exists")
	}
	mains := map[string]bool{
		"cmd/vet": true,
		"fmt":     false,
	}
	for k, v := range mains {
		if dg.IsMainPackage(k) != v {
			t.Error(k, v)
		}
	}
	chains := dg.SearchChain("fmt")
	if len(chains) <= 0 {
		t.Error("chains empty")
	}
	for _, chain := range chains {
		if chain[0] != "main" || chain[len(chain)-1] != "fmt" {
			t.Error("result error", chain)
		}
		for _, v := range chain {
			if v == "..." {
				t.Error("should not contains ...")
			}
		}
	}
	if dg.CountAll() != 371 || dg.CountMain() != 20 {
		t.Error(dg.CountAll(), dg.CountMain())
	}
	all := dg.SearchAll("net/url")
	if !sliceContains(all, "net/http") {
		t.Error("not found")
	}
	dg.Add(DepInfo{
		ImportPath: "x",
		Name:       "x",
		Deps:       []string{"fmt"},
		Imports:    []string{"fmt"},
	})
	unused := dg.ListUnUsed()
	if !sliceContains(unused, "x") {
		t.Error("unused not contains x")
	}
}

func sliceContains(s []string, t string) bool {
	for _, v := range s {
		if v == t {
			return true
		}
	}
	return false
}
