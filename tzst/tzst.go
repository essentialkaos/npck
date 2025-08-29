// Package tzst provides methods for unpacking tar.zst files
package tzst

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/klauspost/compress/zstd"

	"github.com/essentialkaos/npck/tar"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var ErrNilReader = fmt.Errorf("Reader can not be nil")

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpacks file to given directory
func Unpack(file, dir string) error {
	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		return err
	}

	defer fd.Close()

	return Read(bufio.NewReader(fd), dir)
}

// Read reads compressed data using given reader and unpacks it to
// the given directory
func Read(r io.Reader, dir string) error {
	if r == nil {
		return ErrNilReader
	}

	cr, err := zstd.NewReader(r)

	if err != nil {
		return err
	}

	return tar.Read(cr, dir)
}
