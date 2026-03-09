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
	"errors"
	"io"
	"os"

	"github.com/klauspost/compress/zstd"

	"github.com/essentialkaos/npck/tar"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var ErrNilReader = errors.New("reader is nil")

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpack unpacks archive file to given directory
func Unpack(file, dir string) error {
	fd, err := os.Open(file)

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
