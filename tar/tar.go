// Package tar provides method for unpacking tar files
package tar

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
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

// DEFAULT_MAX_READ_LIMIT is default the maximum read limit (1GB)
const DEFAULT_MAX_READ_LIMIT int64 = 1024 * 1024 * 1024

// DEFAULT_DIR_MODE is default mode for directories
const DEFAULT_DIR_MODE os.FileMode = 0750

// ////////////////////////////////////////////////////////////////////////////////// //

type Options struct {
	// MaxReadLimit is the maximum read limit for decompression bomb
	// protection (default: 1GB)
	MaxReadLimit int64

	// DirMode is mode for all created directories (default: 0750)
	DirMode os.FileMode

	// AllowExternalLinks is flag for protection against links to files and directories
	// outside target directory
	AllowExternalLinks bool

	// UpdateTimes is flag for restoring mtime and atime
	UpdateTimes bool

	// UpdateOwner is flag for restoring owner for files and directories
	UpdateOwner bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilReader   = utils.ErrNilReader
	ErrEmptyOutput = utils.ErrEmptyOutput
)

// DefaultOptions is default unpacking options
var DefaultOptions = Options{
	MaxReadLimit: DEFAULT_MAX_READ_LIMIT,
	DirMode:      0750,
	UpdateTimes:  true,
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpack unpacks archive file to given directory
func Unpack(file, dir string, options Options) error {
	fd, err := os.OpenFile(file, os.O_RDONLY, 0)

	if err != nil {
		return err
	}

	defer fd.Close()

	return Read(bufio.NewReader(fd), dir, options)
}

// Read reads compressed data using given reader and unpacks it to
// the given directory
func Read(r io.Reader, dir string, options Options) error {
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
			return fmt.Errorf("path %q contains directory traversal element and cannot be used", header.Name)
		}

		path, err := utils.Join(dir, header.Name)

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			err = createFile(header, tr, path, options)
		case tar.TypeDir:
			err = createDir(header, path, options)
		case tar.TypeLink:
			err = createHardlink(header, dir, path, options.AllowExternalLinks)
		case tar.TypeSymlink:
			err = createSymlink(header, dir, path, options.AllowExternalLinks)
		default:
			err = fmt.Errorf(
				"object %q has unsupported type (%d)",
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
func createDir(h *tar.Header, path string, options Options) error {
	err := os.MkdirAll(path, h.FileInfo().Mode())

	if err != nil {
		return err
	}

	return updateAttrs(h, path, options)
}

// createFile creates new file
func createFile(h *tar.Header, r io.Reader, path string, options Options) error {
	dir := filepath.Dir(path)
	_, err := os.Stat(dir)

	mode := options.DirMode

	if mode == 0 {
		mode = DEFAULT_DIR_MODE
	}

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		err = os.MkdirAll(dir, mode)

		if err != nil {
			return err
		}
	}

	fd, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, h.FileInfo().Mode())

	if err != nil {
		return err
	}

	defer fd.Close()

	limit := options.MaxReadLimit

	if limit == 0 {
		limit = DEFAULT_MAX_READ_LIMIT
	}

	bw := bufio.NewWriter(fd)
	_, err = io.Copy(bw, io.LimitReader(r, limit))

	if err != nil {
		return err
	}

	err = bw.Flush()

	if err != nil {
		return err
	}

	return updateAttrs(h, path, options)
}

// createSymlink creates symbolic link
func createSymlink(h *tar.Header, dir, path string, allowExternalLinks bool) error {
	if !allowExternalLinks && isExternalLink(h.Linkname, path, dir) {
		return fmt.Errorf("symbolic link %q points to object outside of target directory (%q)", h.Name, h.Linkname)
	}

	return os.Symlink(h.Linkname, path)
}

// createHardlink creates hard link
func createHardlink(h *tar.Header, dir, path string, allowExternalLinks bool) error {
	if !allowExternalLinks && isExternalLink(h.Linkname, path, dir) {
		return fmt.Errorf("hard link %q points to object outside of target directory (%q)", h.Name, h.Linkname)
	}

	linkTarget, err := utils.Join(dir, h.Linkname)

	if err != nil {
		return err
	}

	return os.Link(linkTarget, path)
}

// updateAttrs updates target attributes
func updateAttrs(h *tar.Header, path string, options Options) error {
	var err error

	if options.UpdateTimes {
		err = os.Chtimes(path, h.AccessTime, h.ModTime)

		if err != nil {
			return err
		}
	}

	if options.UpdateOwner {
		err = os.Chown(path, h.Uid, h.Gid)

		if err != nil {
			return err
		}
	}

	return nil
}

// isExternalLink checks if link leads to object outside target directory
func isExternalLink(linkPath, objPath, targetDir string) bool {
	if !filepath.IsAbs(linkPath) {
		linkPath = filepath.Clean(filepath.Join(filepath.Dir(objPath), linkPath))
	}

	return linkPath != targetDir &&
		!strings.HasPrefix(linkPath, targetDir+string(filepath.Separator))
}
