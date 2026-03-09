// Package gz provides methods for unpacking files with GZip compression
package gz

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
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/gzip"

	"github.com/essentialkaos/npck/utils"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// MaxReadLimit is the maximum read limit for decompression bomb
// protection (default: 1GB)
var MaxReadLimit int64 = 1024 * 1024 * 1024

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilReader   = errors.New("reader is nil")
	ErrEmptyInput  = errors.New("path to input file is empty")
	ErrEmptyOutput = errors.New("path to output file is empty")
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

	output := strings.TrimSuffix(filepath.Base(file), ".gz")
	output = strings.TrimSuffix(output, ".GZ")

	path, err := utils.Join(dir, output)

	if err != nil {
		return err
	}

	fd, err := os.Open(file)

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

	defer fd.Close()

	cr, err := gzip.NewReader(r)

	if err != nil {
		return err
	}

	bw := bufio.NewWriter(fd)
	_, err = io.Copy(bw, io.LimitReader(cr, MaxReadLimit))

	if err != nil {
		return err
	}

	return bw.Flush()
}
