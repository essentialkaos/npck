// Package zip provides methods for unpacking files with ZIP compression
package zip

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/klauspost/compress/zip"

	"github.com/essentialkaos/npck/utils"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// MaxReadLimit is the maximum read limit for decompression bomb
// protection (default: 1GB)
var MaxReadLimit int64 = 1024 * 1024 * 1024

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilReader   = errors.New("reader is nil")
	ErrEmptyOutput = errors.New("path to output directory is empty")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpack unpacks archive file to given directory
func Unpack(file, dir string) error {
	fi, err := os.Stat(file)

	if err != nil {
		return err
	}

	fd, err := os.Open(file)

	if err != nil {
		return err
	}

	defer fd.Close()

	return Read(fd, fi.Size(), dir)
}

// Read reads compressed data using given reader and unpacks it to
// the given directory
func Read(r io.ReaderAt, size int64, dir string) error {
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

		zfd, err := file.Open()

		if err != nil {
			return err
		}

		defer zfd.Close()

		fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())

		if err != nil {
			return err
		}

		defer fd.Close()

		bw := bufio.NewWriter(fd)
		_, err = io.Copy(bw, io.LimitReader(zfd, MaxReadLimit))

		if err != nil {
			return err
		}

		err = bw.Flush()

		if err != nil {
			return err
		}

		fd.Close()
	}

	return nil
}
