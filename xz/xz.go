// Package xz provides methods for unpacking xz files
package xz

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
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"
)

// Unpacks file to given directory
func Unpack(file, dir string) error {
	switch {
	case file == "":
		return fmt.Errorf("Path to input file can not be empty")
	case dir == "":
		return fmt.Errorf("Path to output file can not be empty")
	}

	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		return err
	}

	defer fd.Close()

	return Read(
		bufio.NewReader(fd),
		filepath.Join(
			filepath.Clean(dir),
			strings.TrimSuffix(filepath.Base(file), ".xz"),
		),
	)
}

// Read reads compressed data using given reader and unpacks it to
// the given directory
func Read(r io.Reader, output string) error {
	switch {
	case r == nil:
		return fmt.Errorf("Reader can not be nil")
	case output == "":
		return fmt.Errorf("Path to output file can not be empty")
	}

	fd, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)

	if err != nil {
		return err
	}

	xr, err := xz.NewReader(r)

	if err != nil {
		return err
	}

	bw := bufio.NewWriter(fd)
	_, err = io.Copy(bw, xr)

	bw.Flush()
	fd.Close()

	return err
}
