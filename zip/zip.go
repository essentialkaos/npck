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

	"github.com/klauspost/compress/zip"

	"github.com/essentialkaos/npck/v3/utils"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// DEFAULT_MAX_READ_LIMIT is default the maximum read limit (1GB)
const DEFAULT_MAX_READ_LIMIT int64 = 1024 * 1024 * 1024

// DEFAULT_MAX_ARCHIVE_SIZE is default the maximum total size of all extracted
// files (100GB)
const DEFAULT_MAX_ARCHIVE_SIZE int64 = 100 * 1024 * 1024 * 1024

// ////////////////////////////////////////////////////////////////////////////////// //

// Options is reader options
type Options struct {
	// MaxReadLimit is the maximum read limit for decompression bomb
	// protection (default: 1GB)
	MaxReadLimit int64

	// MaxTotalSize is the maximum total size of all extracted files
	// (default: 100GB)
	MaxTotalSize int64
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
		return fmt.Errorf("invalid data size (%d <= 0)", size)
	}

	zr, err := zip.NewReader(r, size)

	if err != nil {
		return err
	}

	root, err := os.OpenRoot(dir)

	if err != nil {
		return err
	}

	defer root.Close()

	if options.MaxTotalSize == 0 {
		options.MaxTotalSize = DEFAULT_MAX_ARCHIVE_SIZE
	}

	total := int64(0)

	for _, file := range zr.File {
		header := file.FileHeader
		path, err := utils.LocalPath(header.Name)
		info := header.FileInfo()

		if err != nil {
			return err
		}

		size := int64(header.UncompressedSize64)

		if size < 0 {
			return fmt.Errorf("entry %q declares invalid size (%d)", header.Name, size)
		}

		if utils.IsOverLimit(size, options.MaxTotalSize, total) {
			return fmt.Errorf("archive total size is bigger than %d bytes", options.MaxTotalSize)
		}

		total += size

		if info.IsDir() {
			err = root.MkdirAll(path, utils.DirMode(info.Mode()))

			if err != nil {
				return err
			}

			continue
		}

		err = extractFile(root, file, path, info.Mode(), limit(options))

		if err != nil {
			return err
		}
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// extractFile extracts file from archive to given path inside root
func extractFile(root *os.Root, file *zip.File, path string, mode os.FileMode, limit int64) error {
	zfd, err := file.Open()

	if err != nil {
		return err
	}

	defer zfd.Close()

	fd, err := root.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)

	if err != nil {
		return err
	}

	defer fd.Close()

	bw := bufio.NewWriter(fd)
	err = utils.CopyLimited(bw, zfd, limit)

	if err != nil {
		return err
	}

	err = utils.IsLimitExceeded(zfd, limit)

	if err != nil {
		return err
	}

	return bw.Flush()
}

// ////////////////////////////////////////////////////////////////////////////////// //

// limit returns per-file limit from given options
func limit(options Options) int64 {
	if options.MaxReadLimit == 0 {
		return DEFAULT_MAX_READ_LIMIT
	}

	return options.MaxReadLimit
}
