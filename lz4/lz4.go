// Package lz4 provides methods for unpacking files with LZ4 compression
package lz4

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

	"github.com/pierrec/lz4/v4"

	"github.com/essentialkaos/npck/v3/utils"
)

// DEFAULT_MAX_READ_LIMIT is default the maximum read limit (1GB)
const DEFAULT_MAX_READ_LIMIT int64 = 1024 * 1024 * 1024

// ////////////////////////////////////////////////////////////////////////////////// //

// Options is reader options
type Options struct {
	// MaxReadLimit is the maximum number of bytes the decompressed data can
	// contain (default: 1GB)
	MaxReadLimit int64
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilReader   = utils.ErrNilReader
	ErrEmptyInput  = utils.ErrEmptyInput
	ErrEmptyOutput = utils.ErrEmptyOutput
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpacks file to given directory
func Unpack(file, dir string, options Options) error {
	switch {
	case file == "":
		return ErrEmptyInput
	case dir == "":
		return ErrEmptyOutput
	}

	output := strings.TrimSuffix(filepath.Base(file), ".lz4")
	output = strings.TrimSuffix(output, ".LZ4")

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

	return read(bufio.NewReader(fd), ofd, options)
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

	return read(r, fd, options)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// read decompresses data from given reader and writes it to given file
func read(r io.Reader, fd *os.File, options Options) error {
	cr := lz4.NewReader(r)
	bw := bufio.NewWriter(fd)

	err := utils.CopyLimited(bw, cr, limit(options))

	if err != nil {
		return err
	}

	err = utils.IsLimitExceeded(cr, limit(options))

	if err != nil {
		return err
	}

	return bw.Flush()
}

// ////////////////////////////////////////////////////////////////////////////////// //

// limit returns limit from given options
func limit(options Options) int64 {
	if options.MaxReadLimit == 0 {
		return DEFAULT_MAX_READ_LIMIT
	}

	return options.MaxReadLimit
}
