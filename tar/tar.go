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
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/essentialkaos/npck/v2/utils"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// DEFAULT_MAX_READ_LIMIT is default the maximum read limit (1GB)
const DEFAULT_MAX_READ_LIMIT int64 = 1024 * 1024 * 1024

// DEFAULT_DIR_MODE is default mode for directories
const DEFAULT_DIR_MODE os.FileMode = 0750

// MAX_LINK_DEPTH is the maximum number of symbolic links followed while checking
// symbolic link target
const MAX_LINK_DEPTH = 255

// ////////////////////////////////////////////////////////////////////////////////// //

// Options is reader options
type Options struct {
	// MaxReadLimit is the maximum read limit for decompression bomb
	// protection (default: 1GB)
	MaxReadLimit int64

	// DirMode is mode for all created directories (default: 0750)
	DirMode os.FileMode

	// AllowExternalLinks is flag for allowing symbolic links to files and directories
	// outside target directory. Hard links are always confined to target directory.
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

	rootPath, err := realPath(dir)

	if err != nil {
		return err
	}

	root, err := os.OpenRoot(rootPath)

	if err != nil {
		return err
	}

	defer root.Close()

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path, err := utils.LocalPath(header.Name)

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			err = createFile(root, header, tr, path, options)
		case tar.TypeDir:
			err = createDir(root, header, path, options)
		case tar.TypeLink:
			err = createHardlink(root, rootPath, header, path, options)
		case tar.TypeSymlink:
			err = createSymlink(root, rootPath, header, path, options)
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
func createDir(root *os.Root, h *tar.Header, path string, options Options) error {
	err := root.MkdirAll(path, utils.DirMode(h.FileInfo().Mode()))

	if err != nil {
		return err
	}

	return updateAttrs(root, h, path, options)
}

// createFile creates new file
func createFile(root *os.Root, h *tar.Header, r io.Reader, path string, options Options) error {
	err := createParentDir(root, path, options)

	if err != nil {
		return err
	}

	fd, err := root.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, h.FileInfo().Mode())

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

	return updateAttrs(root, h, path, options)
}

// createSymlink creates symbolic link
func createSymlink(
	root *os.Root, rootPath string, h *tar.Header,
	path string, options Options,
) error {
	err := createParentDir(root, path, options)

	if err != nil {
		return err
	}

	if !options.AllowExternalLinks && isExternalLink(root, rootPath, path, h.Linkname) {
		return fmt.Errorf(
			"symbolic link %q points to object outside of target directory (%q)",
			h.Name, h.Linkname,
		)
	}

	return root.Symlink(h.Linkname, path)
}

// createHardlink creates hard link
func createHardlink(
	root *os.Root, rootPath string, h *tar.Header,
	path string, options Options,
) error {
	target, err := utils.LocalPath(h.Linkname)

	if err != nil {
		return err
	}

	err = createParentDir(root, path, options)

	if err != nil {
		return err
	}

	fi, err := root.Lstat(target)

	if err != nil {
		return err
	}

	// Hard link to symbolic link is a copy of this link at the new location, so
	// its target must be checked again relative to the new location
	if fi.Mode()&fs.ModeSymlink != 0 && !options.AllowExternalLinks {
		linkTarget, err := root.Readlink(target)

		if err != nil {
			return err
		}

		if isExternalLink(root, rootPath, path, linkTarget) {
			return fmt.Errorf(
				"hard link %q points to symbolic link leading outside of target directory (%q)",
				h.Name, linkTarget,
			)
		}
	}

	// Both paths are resolved inside the root, so hard link to object outside
	// of target directory can't be created
	return root.Link(target, path)
}

// createParentDir creates parent directory for given path
func createParentDir(root *os.Root, path string, options Options) error {
	mode := options.DirMode

	if mode == 0 {
		mode = DEFAULT_DIR_MODE
	}

	return root.MkdirAll(filepath.Dir(path), mode)
}

// updateAttrs updates target attributes
func updateAttrs(root *os.Root, h *tar.Header, path string, options Options) error {
	var err error

	if options.UpdateTimes {
		err = root.Chtimes(path, h.AccessTime, h.ModTime)

		if err != nil {
			return err
		}
	}

	if options.UpdateOwner {
		err = root.Chown(path, h.Uid, h.Gid)

		if err != nil {
			return err
		}
	}

	return nil
}

// realPath returns absolute path with all symbolic links resolved
func realPath(path string) (string, error) {
	path, err := filepath.Abs(path)

	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(path)
}

// isExternalLink checks if symbolic link with given target created at given path
// (relative to root) leads to object outside of root
func isExternalLink(root *os.Root, rootPath, path, target string) bool {
	if filepath.IsAbs(target) {
		rel, err := filepath.Rel(rootPath, filepath.Clean(target))

		if err != nil || !filepath.IsLocal(rel) {
			return true
		}

		return !isLocalTarget(root, rootPath, splitPath(rel))
	}

	return !isLocalTarget(
		root, rootPath,
		append(splitPath(filepath.Dir(path)), splitPath(target)...),
	)
}

// isLocalTarget walks over given path elements (relative to root) using the current
// state of the file system and checks that resolved path is inside of the root.
// All symbolic links found during the walk are followed. Path elements which don't
// exist yet are processed lexically, but going up ("..") from such element is
// forbidden, because the element may be created later as a symbolic link.
func isLocalTarget(root *os.Root, rootPath string, elems []string) bool {
	var stack []string

	isMissing := false
	links := 0

	for len(elems) > 0 {
		elem := elems[0]
		elems = elems[1:]

		if elem == ".." {
			if isMissing || len(stack) == 0 {
				return false
			}

			stack = stack[:len(stack)-1]
			continue
		}

		stack = append(stack, elem)

		if isMissing {
			continue
		}

		cur := filepath.Join(stack...)
		fi, err := root.Lstat(cur)

		switch {
		case errors.Is(err, fs.ErrNotExist):
			isMissing = true

		case err != nil:
			return false

		case fi.Mode()&fs.ModeSymlink != 0:
			links++

			if links > MAX_LINK_DEPTH {
				return false
			}

			target, err := root.Readlink(cur)

			if err != nil {
				return false
			}

			stack = stack[:len(stack)-1]

			if filepath.IsAbs(target) {
				rel, err := filepath.Rel(rootPath, filepath.Clean(target))

				if err != nil || !filepath.IsLocal(rel) {
					return false
				}

				stack, target = nil, rel
			}

			elems = append(splitPath(target), elems...)

		case !fi.IsDir() && len(elems) > 0:
			return false
		}
	}

	return true
}

// splitPath splits path to elements omitting empty and "." elements
func splitPath(path string) []string {
	var result []string

	for _, e := range strings.Split(path, string(filepath.Separator)) {
		if e != "" && e != "." {
			result = append(result, e)
		}
	}

	return result
}
