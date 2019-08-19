package depgraph

import (
	"encoding/json"
	"io"
	"sort"
	"strings"
)

type DepInfo struct {
	ImportPath string   `json:"ImportPath"`
	Name       string   `json:"Name"`
	Deps       []string `json:"Deps"`
	Imports    []string `json:"Imports"`
}

func (d *DepInfo) ImportsMap() map[string]bool {
	m := make(map[string]bool, len(d.Imports))
	for _, v := range d.Imports {
		m[v] = true
	}
	return m
}

func (d *DepInfo) DepsMap() map[string]bool {
	m := make(map[string]bool, len(d.Deps))
	for _, v := range d.Deps {
		m[v] = true
	}
	return m
}

type DepGraph struct {
	imports      map[string]map[string]bool
	allDeps      map[string]map[string]bool
	mainPackages map[string]bool
	testPackages map[string]bool
}

func (g *DepGraph) Add(d DepInfo) {
	if g.imports == nil {
		g.imports = make(map[string]map[string]bool)
		g.imports["main"] = make(map[string]bool)
	}
	if g.mainPackages == nil {
		g.mainPackages = make(map[string]bool)
	}
	if g.testPackages == nil {
		g.testPackages = make(map[string]bool)
	}
	if g.allDeps == nil {
		g.allDeps = make(map[string]map[string]bool)
	}
	if strings.HasSuffix(d.ImportPath, "]") { // skip test package
		return
	}
	isTestPackage := strings.HasSuffix(d.ImportPath, ".test")
	if d.Name == "main" {
		if isTestPackage {
			g.testPackages[d.ImportPath] = true
		} else {
			g.mainPackages[d.ImportPath] = true
		}
	}
	g.imports[d.ImportPath] = d.ImportsMap()
	g.allDeps[d.ImportPath] = d.DepsMap()
}

func (g *DepGraph) CountAll() int {
	return len(g.imports)
}
func (g *DepGraph) CountMain() int {
	return len(g.mainPackages)
}

func (g *DepGraph) CountTest() int {
	return len(g.testPackages)
}

func reverseSlice(a []string) {
	if len(a) <= 1 {
		return
	}
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func (g *DepGraph) SearchMain(packageName string) (packages []string) {
	for v := range g.mainPackages {
		if g.allDeps[v][packageName] {
			packages = append(packages, v)
		}
	}
	return
}

func (g *DepGraph) SearchTest(packageName string) (packages []string) {
	for v := range g.testPackages {
		if g.allDeps[v][packageName] {
			packages = append(packages, v)
		}
	}
	return
}

func (g *DepGraph) Exists(packageName string) bool {
	_, exists := g.allDeps[packageName]
	return exists
}

func (g *DepGraph) SearchAll(packageName string) (packages []string) {
	for k, v := range g.allDeps {
		if v[packageName] {
			packages = append(packages, k)
		}
	}
	return
}

func (g *DepGraph) ListUnUsed() (packages []string) {
	defer func() {
		sort.Strings(packages)
	}()
	for p := range g.allDeps {
		if g.mainPackages[p] || g.testPackages[p] {
			continue
		}
		found := false
		for m := range g.allDeps {
			if g.allDeps[m][p] {
				found = true
			}
		}
		if !found {
			packages = append(packages, p)
		}
	}
	return
}

func (g *DepGraph) IsMainPackage(packageName string) bool {
	return g.mainPackages[packageName]
}

func (g *DepGraph) IsTestPackage(packageName string) bool {
	return g.testPackages[packageName]
}

func (g *DepGraph) SearchChain(packageName string) (chains [][]string) {
	for _, p := range g.SearchMain(packageName) {
		chain := []string{}
		checked := make(map[string]bool)
		chain, found := g.search(p, packageName, chain, checked)
		if !found {
			// dep存在，但是找不到依赖链，说明依赖关系导入不全，比如缺少标准库
			chain = []string{packageName, "..."}
		}
		chain = append(chain, p)
		chain = append(chain, "main")
		reverseSlice(chain)
		chains = append(chains, chain)
	}
	return
}

func (g *DepGraph) search(start, packageName string, current []string, checked map[string]bool) (after []string, found bool) {
	if checked[start] {
		return
	}
	checked[start] = true
	if g.imports[start][packageName] {
		found = true
		after = append(current, packageName)
		return
	}
	for p := range g.imports[start] {
		if after, ok := g.search(p, packageName, current, checked); ok {
			after = append(after, p)
			return after, true
		}
	}
	return
}

func LoadDeps(r io.Reader) (dg *DepGraph, err error) {
	dec := json.NewDecoder(r)
	dg = &DepGraph{}
	for {
		var di DepInfo
		err = dec.Decode(&di)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		dg.Add(di)
	}
	return
}
