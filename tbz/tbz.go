package tbz

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"compress/bzip2"
	"io"
	"os"

	"github.com/essentialkaos/npck/tar"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpacks file to given directory
func Unpack(file, dir string) error {
	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		return err
	}

	defer fd.Close()

	return UnpackReader(bufio.NewReader(fd), dir)
}

// UnpackReader reads packed data using given reader and unpacks it to
// the given directory
func UnpackReader(r io.Reader, dir string) error {
	return tar.UnpackReader(bzip2.NewReader(r), dir)
}
