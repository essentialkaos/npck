// Package bz2 provides methods for unpacking files with BZip2 compression
package bz2

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
	"path/filepath"
	"strings"

	"github.com/essentialkaos/npck/v2/utils"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// DEFAULT_MAX_READ_LIMIT is default the maximum read limit (1GB)
const DEFAULT_MAX_READ_LIMIT int64 = 1024 * 1024 * 1024

// ////////////////////////////////////////////////////////////////////////////////// //

// Options is reader options
type Options struct {
	// MaxReadLimit is the maximum read limit for decompression bomb
	// protection (default: 1GB)
	MaxReadLimit int64
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilReader   = utils.ErrNilReader
	ErrEmptyInput  = utils.ErrEmptyInput
	ErrEmptyOutput = utils.ErrEmptyOutput
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpack unpacks archive file to given directory
func Unpack(file, dir string, options Options) error {
	switch {
	case file == "":
		return ErrEmptyInput
	case dir == "":
		return ErrEmptyOutput
	}

	output := strings.TrimSuffix(filepath.Base(file), ".bz2")
	output = strings.TrimSuffix(output, ".BZ2")

	path, err := utils.Join(dir, output)

	if err != nil {
		return err
	}

	fd, err := os.Open(file)

	if err != nil {
		return err
	}

	defer fd.Close()

	return Read(bufio.NewReader(fd), path, options)
}

// Read reads compressed data using given reader and unpacks it to
// the given directory
func Read(r io.Reader, output string, options Options) error {
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

	limit := options.MaxReadLimit

	if limit == 0 {
		limit = DEFAULT_MAX_READ_LIMIT
	}

	bw := bufio.NewWriter(fd)
	_, err = io.Copy(bw, io.LimitReader(bzip2.NewReader(r), limit))

	if err != nil {
		return err
	}

	return bw.Flush()
}
