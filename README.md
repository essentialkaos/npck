<p align="center"><a href="#readme"><img src="https://gh.kaos.st/npck.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/g/npck"><img src="https://gh.kaos.st/godoc.svg" alt="PkgGoDev" /></a>
  <a href="https://kaos.sh/r/npck"><img src="https://kaos.sh/r/npck.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/w/npck/ci"><img src="https://kaos.sh/w/npck/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/npck/codeql"><img src="https://kaos.sh/w/npck/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="https://kaos.sh/c/npck"><img src="https://kaos.sh/c/npck.svg" alt="Coverage Status" /></a>
  <a href="https://kaos.sh/b/npck"><img src="https://kaos.sh/b/fc322f23-4913-4edd-8f0f-33a3ce029add.svg" alt="Codebeat badge" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage-example">Usage example</a> • <a href="#build-status">CI Status</a> • <a href="#license">License</a></p>

<br/>

`npck` is a Go package for unpacking various types of archives.

### Installation

Make sure you have a working Go 1.18+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get -d github.com/essentialkaos/npck
```

For update to the latest stable release, do:

```
go get -d -u github.com/essentialkaos/npck
```

### Usage example

```go
package main

import (
  "fmt"

  "github.com/essentialkaos/npck/tgz"
)

func main() {
  err := tgz.Unpack("file.tar.gz", "/home/john")

  if err != nil {
    panic("Can't unpack file: %v", err)
  }

  fmt.Printf("File %s successfully unpacked!\n")
}
```

### CI Status

| Branch | Status |
|--------|--------|
| `master` | [![CI](https://kaos.sh/w/npck/ci.svg?branch=master)](https://kaos.sh/w/npck/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/npck/ci.svg?branch=develop)](https://kaos.sh/w/npck/ci?query=branch:develop) |

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
