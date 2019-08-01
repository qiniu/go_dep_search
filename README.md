# go_dep_search

golang dependency search tool.

### Install

```
go get -u github.com/ma6174/go_dep_search
```

### Usage

```
go list -json -deps -test all | go_dep_search package_names

Args:
  -chain
    	show dep chained
  -main
    	only show main package
  -unused
    	list unused packages
```

eg: find which command(main package) use `net/http` or `encoding/json` package in go source code:

```
root@b7e158d83ff2:/go# go list -json all | go_dep_search -main net/http encoding/json
main -> cmd/go -> net/http
main -> cmd/pprof -> net/http
main -> cmd/trace -> net/http
main -> cmd/pprof -> encoding/json
main -> cmd/dist -> encoding/json
main -> cmd/go -> encoding/json
main -> cmd/test2json -> encoding/json
main -> cmd/cover -> encoding/json
main -> cmd/trace -> encoding/json
main -> cmd/compile -> encoding/json
main -> cmd/link -> encoding/json
main -> cmd/vet -> encoding/json
```

eg: show chained package deps

```
root@b7e158d83ff2:/go# go list -json all | go_dep_search -chain net/http encoding/json
main -> cmd/go -> cmd/go/internal/bug -> cmd/go/internal/envcmd -> cmd/go/internal/modload -> cmd/go/internal/modfetch -> cmd/go/internal/web2 -> net/http
main -> cmd/pprof -> net/http
main -> cmd/trace -> net/http
main -> cmd/compile -> cmd/compile/internal/amd64 -> cmd/compile/internal/gc -> encoding/json
main -> cmd/test2json -> cmd/internal/test2json -> encoding/json
main -> cmd/cover -> encoding/json
main -> cmd/dist -> encoding/json
main -> cmd/go -> cmd/go/internal/envcmd -> encoding/json
main -> cmd/link -> cmd/link/internal/arm64 -> cmd/link/internal/ld -> encoding/json
main -> cmd/pprof -> cmd/vendor/github.com/google/pprof/driver -> cmd/vendor/github.com/google/pprof/internal/driver -> encoding/json
main -> cmd/trace -> encoding/json
main -> cmd/vet -> cmd/vendor/golang.org/x/tools/go/analysis/unitchecker -> encoding/json
```

eg: show unsed packages

```
root@b7e158d83ff2:/go# go list -json all | go_dep_search -unused
archive/tar
cmd/compile/internal/test
cmd/go/internal/txtar
cmd/go/internal/webtest
cmd/vendor/github.com/google/pprof/internal/proftest
cmd/vendor/golang.org/x/sys/windows
cmd/vendor/golang.org/x/sys/windows/registry
cmd/vendor/golang.org/x/tools/go/analysis/passes/pkgfact
compress/bzip2
container/ring
database/sql
encoding/ascii85
encoding/base32
encoding/csv
expvar
hash/crc64
hash/fnv
image/gif
image/jpeg
image/png
index/suffixarray
internal/syscall/windows
internal/syscall/windows/registry
internal/syscall/windows/sysdll
internal/testenv
internal/x/net/internal/nettest
internal/x/net/nettest
internal/x/text/secure
internal/x/text/unicode
log/syslog
math/cmplx
net/http/cookiejar
net/http/fcgi
net/http/httptest
net/http/httputil
net/internal/socktest
net/mail
net/rpc/jsonrpc
net/smtp
os/signal/internal/pty
plugin
runtime/pprof/internal/profile
runtime/race
testing/internal/testdeps
testing/iotest
testing/quick
```
