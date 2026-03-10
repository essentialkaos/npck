// Package txz provides methods for unpacking tar.xz files
package txz

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"io"
	"os"

	"github.com/ulikunitz/xz"

	"github.com/essentialkaos/npck/v2/tar"
	"github.com/essentialkaos/npck/v2/utils"
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

	cr, err := xz.NewReader(r)

	if err != nil {
		return err
	}

	return tar.Read(cr, dir, options)
}
