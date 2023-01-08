// Package zip provides methods for unpacking zip files
package zip

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
	"strings"

	"github.com/klauspost/compress/zip"

	securejoin "github.com/cyphar/filepath-securejoin"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpacks file to given directory
func Unpack(file, dir string) error {
	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		return err
	}

	defer fd.Close()

	return Read(fd, dir)
}

// Read reads compressed data using given reader and unpacks it to
// the given directory
func Read(r io.ReaderAt, dir string) error {
	switch {
	case r == nil:
		return fmt.Errorf("Reader can not be nil")
	case dir == "":
		return fmt.Errorf("Path to output directory can not be empty")
	}

	zr, err := zip.NewReader(r, 4096)

	if err != nil {
		return err
	}

	for _, file := range zr.File {
		header := file.FileHeader

		if strings.Contains(header.Name, "..") {
			return fmt.Errorf("Path \"%s\" contains directory traversal element and cannot be used", header.Name)
		}

		info := header.FileInfo()
		path, err := securejoin.SecureJoin(dir, header.Name)

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

		fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())

		if err != nil {
			return err
		}

		bw := bufio.NewWriter(fd)
		_, err = io.Copy(bw, zfd)

		bw.Flush()
		fd.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
