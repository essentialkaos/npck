<p align="center"><a href="#readme"><img src="https://gh.kaos.st/npck.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/g/npck"><img src="https://gh.kaos.st/godoc.svg" alt="PkgGoDev" /></a>
  <a href="https://kaos.sh/w/npck/ci"><img src="https://kaos.sh/w/npck/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/npck/codeql"><img src="https://kaos.sh/w/npck/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="https://kaos.sh/c/npck"><img src="https://kaos.sh/c/npck.svg" alt="Coverage Status" /></a>
  <a href="https://kaos.sh/b/npck"><img src="https://kaos.sh/b/fc322f23-4913-4edd-8f0f-33a3ce029add.svg" alt="Codebeat badge" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#supported-formats">Supported formats</a> • <a href="#installation">Installation</a> • <a href="#usage-example">Usage example</a> • <a href="#ci-status">CI Status</a> • <a href="#license">License</a></p>

<br/>

`npck` is a Go package for unpacking various types of archives.

### Supported formats

* [tar](https://en.wikipedia.org/wiki/Tar_(computing)) (`.tar`)
* [Gzip](https://www.gnu.org/software/gzip/) (`.gz`, `.tgz`, `.tar.gz`)
* [bzip2](http://sourceware.org/bzip2/) (`.bz2`, `.tbz2`, `.tar.bz2`)
* [xz](https://tukaani.org/xz/) (`.xz`, `.txz`, `.tar.xz`)
* [Zstandart](https://facebook.github.io/zstd/) (`.zst`, `.tzst`, `.tar.zst`)
* [LZ4](https://lz4.github.io/lz4/) (`.lz4`, `.tlz4`, `.tar.lz4`)
* [ZIP](https://en.wikipedia.org/wiki/ZIP_(file_format)) (`.zip`)

### Installation

Make sure you have a working Go 1.18+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```bash
go get -u github.com/essentialkaos/npck
```

### Usage example

```go
package main

import (
  "fmt"
  "github.com/essentialkaos/npck"
)

func main() {
  err := npck.Unpack("file.tar.gz", "/home/john")

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
