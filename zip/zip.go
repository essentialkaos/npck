// Package zip provides methods for unpacking files with ZIP compression
package zip

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/klauspost/compress/zip"

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
	ErrEmptyOutput = utils.ErrEmptyOutput
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpack unpacks archive file to given directory
func Unpack(file, dir string, options Options) error {
	fi, err := os.Stat(file)

	if err != nil {
		return err
	}

	fd, err := os.Open(file)

	if err != nil {
		return err
	}

	defer fd.Close()

	return Read(fd, fi.Size(), dir, options)
}

// Read reads compressed data using given reader and unpacks it to
// the given directory
func Read(r io.ReaderAt, size int64, dir string, options Options) error {
	switch {
	case r == nil:
		return ErrNilReader
	case dir == "":
		return ErrEmptyOutput
	case size <= 0:
		return fmt.Errorf("invalid data size (%d < 0)", size)
	}

	zr, err := zip.NewReader(r, size)

	if err != nil {
		return err
	}

	limit := options.MaxReadLimit

	if limit == 0 {
		limit = DEFAULT_MAX_READ_LIMIT
	}

	for _, file := range zr.File {
		header := file.FileHeader

		if strings.Contains(header.Name, "..") {
			return fmt.Errorf("path %q contains directory traversal element and cannot be used", header.Name)
		}

		info := header.FileInfo()
		path, err := utils.Join(dir, header.Name)

		if err != nil {
			return err
		}

		if info.IsDir() {
			err = os.MkdirAll(path, info.Mode())

			if err != nil {
				return err
			}

			continue
		}

		err = extractFile(file, path, info.Mode(), limit)

		if err != nil {
			return err
		}
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func extractFile(file *zip.File, path string, mode os.FileMode, limit int64) error {
	zfd, err := file.Open()

	if err != nil {
		return err
	}

	defer zfd.Close()

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)

	if err != nil {
		return err
	}

	defer fd.Close()

	bw := bufio.NewWriter(fd)
	_, err = io.Copy(bw, io.LimitReader(zfd, limit))

	if err != nil {
		return err
	}

	return bw.Flush()
}
