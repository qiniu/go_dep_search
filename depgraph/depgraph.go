package depgraph

import (
	"container/list"
	"encoding/json"
	"io"
	"sort"
	"strings"
)

type DepInfo struct {
	ImportPath string   `json:"ImportPath"` // import path of package in dir
	Name       string   `json:"Name"`       // package name
	Deps       []string `json:"Deps"`       // all (recursively) imported dependencies
	Imports    []string `json:"Imports"`    // import paths used by this package
}

func (d *DepInfo) ImportsMap() packages {
	m := make(packages, len(d.Imports))
	for _, v := range d.Imports {
		m.Add(v)
	}
	return m
}

func (d *DepInfo) DepsMap() packages {
	m := make(packages, len(d.Deps))
	for _, v := range d.Deps {
		m.Add(v)
	}
	return m
}

type packages map[string]struct{}

func (p packages) Contains(packageName string) bool {
	_, ok := p[packageName]
	return ok
}

func (p packages) Add(packageName string) {
	p[packageName] = struct{}{}
}

func (p packages) ContainsDir(dirName string) (packageName string, contains bool) {
	dirName = strings.TrimSuffix(dirName, "/")
	if p.Contains(dirName) {
		return dirName, true
	}
	dirName += "/"
	for packageName := range p {
		if strings.HasPrefix(packageName, dirName) {
			return packageName, true
		}
	}
	return "", false
}

type DepGraph struct {
	imports      map[string]packages
	allDeps      map[string]packages
	mainPackages packages
	testPackages packages
}

func (g *DepGraph) Add(d DepInfo) {
	if g.imports == nil {
		g.imports = make(map[string]packages)
		g.imports["main"] = make(packages)
	}
	if g.mainPackages == nil {
		g.mainPackages = make(packages)
	}
	if g.testPackages == nil {
		g.testPackages = make(packages)
	}
	if g.allDeps == nil {
		g.allDeps = make(map[string]packages)
	}
	if strings.HasSuffix(d.ImportPath, "]") { // skip test package
		return
	}
	isTestPackage := strings.HasSuffix(d.ImportPath, ".test")
	if d.Name == "main" {
		if isTestPackage {
			g.testPackages.Add(d.ImportPath)
		} else {
			g.mainPackages.Add(d.ImportPath)
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
		if g.allDeps[v].Contains(packageName) || v == packageName {
			packages = append(packages, v)
		}
	}
	return
}

func (g *DepGraph) SearchTest(packageName string) (packages []string) {
	for v := range g.testPackages {
		if g.allDeps[v].Contains(packageName) {
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
		if v.Contains(packageName) {
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
		if g.mainPackages.Contains(p) || g.testPackages.Contains(p) {
			continue
		}
		found := false
		for m := range g.allDeps {
			if g.allDeps[m].Contains(p) {
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
	return g.mainPackages.Contains(packageName)
}

func (g *DepGraph) IsTestPackage(packageName string) bool {
	return g.testPackages.Contains(packageName)
}

func (g *DepGraph) SearchChain(packageName string) (chains [][]string) {
	for _, p := range g.SearchMain(packageName) {
		if p == packageName {
			chains = append(chains, []string{"main", p})
			continue
		}
		chain := []string{}
		chain, found := g.search(p, packageName, chain)
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

func (g *DepGraph) search(start, packageName string, current []string) (after []string, found bool) {
	if !g.allDeps[start].Contains(packageName) {
		return
	}
	if g.imports[start].Contains(packageName) {
		found = true
		after = append(current, packageName)
		return
	}
	for p := range g.imports[start] {
		if after, ok := g.search(p, packageName, current); ok {
			after = append(after, p)
			return after, true
		}
	}
	return
}

func (g *DepGraph) SearchGraph(start, toSearch string) (result map[string][]string) {
	if !g.allDeps[start].Contains(toSearch) {
		return
	}
	result = make(map[string][]string)
	checked := make(map[string]bool)
	l := list.New()
	l.PushBack(start)
	checked[start] = true
	for e := l.Front(); e != nil; e = e.Next() {
		fromPackage := e.Value.(string)
		for p := range g.imports[fromPackage] {
			if p == toSearch {
				result[fromPackage] = append(result[fromPackage], p)
				continue
			}
			if g.allDeps[p].Contains(toSearch) {
				if !checked[p] {
					checked[p] = true
					l.PushBack(p)
				}
				result[fromPackage] = append(result[fromPackage], p)
			}
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
