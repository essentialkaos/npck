<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/g/npck"><img src=".github/images/godoc.svg"/></a>
  <a href="https://kaos.sh/y/npck"><img src="https://kaos.sh/y/fdbafdcb2caa4516afbd5feabebce511.svg" alt="Codacy" /></a>
  <a href="https://kaos.sh/c/npck"><img src="https://kaos.sh/c/npck.svg" alt="Coverage Status" /></a>
  <a href="https://kaos.sh/w/npck/ci"><img src="https://kaos.sh/w/npck/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/npck/codeql"><img src="https://kaos.sh/w/npck/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#supported-formats">Supported formats</a> • <a href="#usage-example">Usage example</a> • <a href="#ci-status">CI Status</a> • <a href="#license">License</a></p>

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

### Usage example

```go
package main

import (
  "fmt"
  "github.com/essentialkaos/npck"
)

func main() {
  file := "file.tar.gz"
  err := npck.Unpack(file, "/home/john")

  if err != nil {
    fmt.Printf("Error: Can't unpack %s: %v\n", file, err)
    return
  }

  fmt.Printf("File %s successfully unpacked!\n", file)
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
