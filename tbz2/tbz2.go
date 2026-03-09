// Package tbz2 provides methods for unpacking tar.bz2 files
package tbz2

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"compress/bzip2"
	"io"
	"os"

	"github.com/essentialkaos/npck/tar"
	"github.com/essentialkaos/npck/utils"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var ErrNilReader = utils.ErrNilReader

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpack unpacks archive file to given directory
func Unpack(file, dir string, options tar.Options) error {
	fd, err := os.Open(file)

	if err != nil {
		return err
	}

	defer fd.Close()

	return Read(bufio.NewReader(fd), dir, options)
}

// Read reads compressed data using given reader and unpacks it to
// the given directory
func Read(r io.Reader, dir string, options tar.Options) error {
	if r == nil {
		return ErrNilReader
	}

	return tar.Read(bzip2.NewReader(r), dir, options)
}
