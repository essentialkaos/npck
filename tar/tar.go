// Package tar provides method for unpacking tar files
package tar

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// UpdateOwner is flag for restoring owner for files and directories
var UpdateOwner = false

// UpdateOwner is flag for restoring mtime and atime
var UpdateTimes = true

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpacks file to given directory
func Unpack(file, dir string) error {
	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		return err
	}

	defer fd.Close()

	return Read(bufio.NewReader(fd), dir)
}

// Read reads compressed data using given reader and unpacks it to
// the given directory
func Read(r io.Reader, dir string) error {
	switch {
	case r == nil:
		return fmt.Errorf("Reader can not be nil")
	case dir == "":
		return fmt.Errorf("Path to output directory can not be empty")
	}

	_, err := os.Stat(dir)

	if err != nil {
		return err
	}

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(dir, header.Name)

		if strings.Contains(path, "..") {
			return fmt.Errorf("Path \"%s\" contains directory traversal element and cannot be used", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeReg:
			err = createFile(header, tr, path)
		case tar.TypeDir:
			err = createDir(header, path)
		case tar.TypeLink:
			err = createHardlink(header, path)
		case tar.TypeSymlink:
			err = createSymlink(header, path)
		default:
			err = fmt.Errorf(
				"Object %s has unsupported type (%d)",
				header.Name, header.Typeflag,
			)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// createDir creates new directory
func createDir(h *tar.Header, path string) error {
	err := os.MkdirAll(path, h.FileInfo().Mode())

	if err != nil {
		return err
	}

	return updateAttrs(h, path)
}

// createFile creates new file
func createFile(h *tar.Header, r io.Reader, path string) error {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, h.FileInfo().Mode())

	if err != nil {
		return err
	}

	bw := bufio.NewWriter(fd)
	_, err = io.Copy(bw, r)

	bw.Flush()
	fd.Close()

	return updateAttrs(h, path)
}

// createSymlink creates symbolic link
func createSymlink(h *tar.Header, path string) error {
	return os.Symlink(h.Linkname, path)
}

// createHardlink creates hard link
func createHardlink(h *tar.Header, path string) error {
	return os.Link(h.Linkname, path)
}

// updateAttrs updates target attributes
func updateAttrs(h *tar.Header, path string) error {
	var err error

	if UpdateTimes {
		err = os.Chtimes(path, h.AccessTime, h.ModTime)

		if err != nil {
			return err
		}
	}

	if UpdateOwner {
		err = os.Chown(path, h.Uid, h.Gid)

		if err != nil {
			return err
		}
	}

	return nil
}
