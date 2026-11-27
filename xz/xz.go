// Package xz provides methods for unpacking files with XZ compression
package xz

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
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"

	"github.com/essentialkaos/npck/v2/utils"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilReader   = utils.ErrNilReader
	ErrEmptyInput  = utils.ErrEmptyInput
	ErrEmptyOutput = utils.ErrEmptyOutput
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpack unpacks archive file to given directory
func Unpack(file, dir string) error {
	switch {
	case file == "":
		return ErrEmptyInput
	case dir == "":
		return ErrEmptyOutput
	}

	output := strings.TrimSuffix(filepath.Base(file), ".xz")
	output = strings.TrimSuffix(output, ".XZ")

	fd, err := os.Open(file)

	if err != nil {
		return err
	}

	defer fd.Close()

	ofd, err := utils.CreateFile(dir, output, 0640)

	if err != nil {
		return err
	}

	defer ofd.Close()

	return read(bufio.NewReader(fd), ofd)
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

	return read(r, fd)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// read decompresses data from given reader and writes it to given file
func read(r io.Reader, fd *os.File) error {
	cr, err := xz.NewReader(r)

	if err != nil {
		return err
	}

	bw := bufio.NewWriter(fd)
	_, err = io.Copy(bw, cr)

	if err != nil {
		return err
	}

	return bw.Flush()
}
