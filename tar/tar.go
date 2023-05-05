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

	"github.com/essentialkaos/npck/utils"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// UpdateOwner is flag for restoring owner for files and directories
var UpdateOwner = false

// UpdateOwner is flag for restoring mtime and atime
var UpdateTimes = true

// AllowExternalLinks is flag for protection against links to files and directories
// outside target directory
var AllowExternalLinks = false

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilReader   = fmt.Errorf("Reader can not be nil")
	ErrEmptyOutput = fmt.Errorf("Path to output directory can not be empty")
)

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
		return ErrNilReader
	case dir == "":
		return ErrEmptyOutput
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

		if strings.Contains(header.Name, "..") {
			return fmt.Errorf("Path \"%s\" contains directory traversal element and cannot be used", header.Name)
		}

		path, err := utils.Join(dir, header.Name)

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			err = createFile(header, tr, path)
		case tar.TypeDir:
			err = createDir(header, path)
		case tar.TypeLink:
			err = createHardlink(header, dir, path)
		case tar.TypeSymlink:
			err = createSymlink(header, dir, path)
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
	dir := filepath.Dir(path)
	_, err := os.Stat(dir)

	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)

		if err != nil {
			return err
		}
	}

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
func createSymlink(h *tar.Header, dir, path string) error {
	if !AllowExternalLinks && isExternalLink(h.Linkname, dir) {
		return fmt.Errorf("Symbolic link %s points to target outside of target directory (%s)", h.Name, h.Linkname)
	}

	return os.Symlink(h.Linkname, path)
}

// createHardlink creates hard link
func createHardlink(h *tar.Header, dir, path string) error {
	if !AllowExternalLinks && isExternalLink(h.Linkname, dir) {
		return fmt.Errorf("Hard link %s points to target outside of target directory (%s)", h.Name, h.Linkname)
	}

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

// isExternalLink checks if link leads to object outside target directory
func isExternalLink(path, dir string) bool {
	if filepath.IsAbs(path) && !strings.HasPrefix(path, dir) {
		return true
	}

	realPath, err := utils.Join(dir, path)

	if err != nil {
		return true
	}

	return !strings.HasPrefix(realPath, dir)
}
