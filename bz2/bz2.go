// Package bz2 provides methods for unpacking files with BZip2 compression
package bz2

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	securejoin "github.com/cyphar/filepath-securejoin"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilReader   = fmt.Errorf("Reader can not be nil")
	ErrEmptyInput  = fmt.Errorf("Path to input file can not be empty")
	ErrEmptyOutput = fmt.Errorf("Path to output file can not be empty")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpacks file to given directory
func Unpack(file, dir string) error {
	switch {
	case file == "":
		return ErrEmptyInput
	case dir == "":
		return ErrEmptyOutput
	}

	output := strings.TrimSuffix(filepath.Base(file), ".bz2")
	output = strings.TrimSuffix(output, ".BZ2")

	path, err := securejoin.SecureJoin(dir, output)

	if err != nil {
		return err
	}

	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		return err
	}

	defer fd.Close()

	return Read(bufio.NewReader(fd), path)
}

// Read reads compressed data using given reader and unpacks it to
// the given directory
func Read(r io.Reader, output string) error {
	switch {
	case r == nil:
		return ErrNilReader
	case output == "":
		return ErrEmptyOutput
	}

	fd, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)

	if err != nil {
		return err
	}

	bw := bufio.NewWriter(fd)
	_, err = io.Copy(bw, bzip2.NewReader(r))

	bw.Flush()
	fd.Close()

	return err
}
