// Package lz4 provides methods for unpacking files with LZ4 compression
package lz4

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

	"github.com/pierrec/lz4/v4"

	"github.com/essentialkaos/npck/utils"
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

	output := strings.TrimSuffix(filepath.Base(file), ".lz4")
	output = strings.TrimSuffix(output, ".LZ4")

	path, err := utils.Join(dir, output)

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

	cr := lz4.NewReader(r)

	if err != nil {
		return err
	}

	bw := bufio.NewWriter(fd)
	_, err = io.Copy(bw, cr)

	bw.Flush()
	fd.Close()

	return err
}
