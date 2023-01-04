package tzst

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
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
	if r == nil {
		return fmt.Errorf("Reader can not be nil")
	}

	gr, err := zstd.NewReader(r)

	if err != nil {
		return err
	}

	return tar.UnpackReader(gr, dir)
}
