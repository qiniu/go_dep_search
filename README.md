# go_dep_search

golang dependency search tool.

### Install

```
go get -u github.com/ma6174/go_dep_search
```

### Usage

```
go list -json ./... | go_dep_search package_names
```

eg: find which command(main package) use `net/http` and `encoding/json` package in go source code:

```bash
~/go/src(!go1.12.5!)$ go list -json ./... | go_dep_search -main net/http encoding/json
Name	->	ImportPath	->	dep_package
main	->	cmd/compile	->	encoding/json
main	->	cmd/cover	->	encoding/json
main	->	cmd/dist	->	encoding/json
main	->	cmd/go	->	net/http
main	->	cmd/go	->	encoding/json
main	->	cmd/link	->	encoding/json
main	->	cmd/pprof	->	net/http
main	->	cmd/pprof	->	encoding/json
main	->	cmd/test2json	->	encoding/json
main	->	cmd/trace	->	net/http
main	->	cmd/trace	->	encoding/json
main	->	cmd/vet	->	encoding/json
```
