// Package utils provides auxiliary methods for working with archives
package utils

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// ErrLimitExceeded is returned when data exceeds maximum read limit
var ErrLimitExceeded = errors.New("data exceeds maximum read limit")

var (
	ErrNilReader   = errors.New("reader is nil")
	ErrEmptyInput  = errors.New("path to input file is empty")
	ErrEmptyOutput = errors.New("path to output file is empty")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// LocalPath converts archive entry name to clean relative path which can be used
// with os.Root. Leading slashes are removed. Method returns error if name contains
// directory traversal elements.
func LocalPath(name string) (string, error) {
	for _, e := range strings.Split(filepath.ToSlash(name), "/") {
		if e == ".." {
			return "", fmt.Errorf(
				"path %q contains directory traversal element and cannot be used", name,
			)
		}
	}

	name = strings.TrimLeft(path.Clean("/"+filepath.ToSlash(name)), "/")

	if name == "" {
		return ".", nil
	}

	return filepath.FromSlash(name), nil
}

// DirMode returns permission and special bits of given mode which can be used
// for directory creation using os.Root
func DirMode(mode os.FileMode) os.FileMode {
	return mode & (os.ModePerm | os.ModeSetuid | os.ModeSetgid | os.ModeSticky)
}

// CreateFile creates (or truncates) file with given name inside given directory.
// Symbolic links are resolved only within the directory, so the file can't be
// created outside of it.
func CreateFile(dir, name string, perm os.FileMode) (*os.File, error) {
	root, err := os.OpenRoot(dir)

	if err != nil {
		return nil, err
	}

	defer root.Close()

	return root.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
}

// CopyLimited copies up to limit bytes of data from reader to writer.
// Use IsLimitExceeded to check if there is more data than allowed by limit.
func CopyLimited(w io.Writer, r io.Reader, limit int64) error {
	_, err := io.Copy(w, io.LimitReader(r, limit))

	if err != nil {
		return err
	}

	return nil
}

// IsLimitExceeded checks if reader has more data than given limit
func IsLimitExceeded(r io.Reader, limit int64) error {
	buf := make([]byte, 1)

	_, err := io.ReadFull(io.LimitReader(r, 1), buf)

	if err == io.EOF {
		return nil
	} else if err != nil {
		return err
	}

	return ErrLimitExceeded
}

// IsOverLimit checks if given size can be read with given limit and
// already processed number of bytes
func IsOverLimit(size, limit, total int64) bool {
	return limit > 0 && total+size > limit
}

// Join joins all elements of path, makes lexical processing, and evaluating all symlinks.
// Method returns error if final destination is not a child path of root.
//
// Deprecated: Join is vulnerable to escapes through symbolic links in intermediate
// path elements. Use os.Root instead.
func Join(root string, elem ...string) (string, error) {
	result, err := filepath.EvalSymlinks(root)

	if err != nil {
		result = root
	} else {
		root = result
	}

	for _, e := range elem {
		result = filepath.Clean(result + "/" + e)

		if isLink(result) {
			result, err = filepath.EvalSymlinks(result)

			if err != nil {
				return "", fmt.Errorf("can't eval symlinks: %w", err)
			}
		}
	}

	if isExternalPath(root, result) {
		return "", fmt.Errorf("final destination (%s) is outside root (%s)", result, root)
	}

	return result, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// isExternalPath returns true if given path is outside of root
func isExternalPath(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err != nil || strings.HasPrefix(rel, "..")
}

// isLink returns true if given object is symlink
func isLink(path string) bool {
	fi, err := os.Lstat(path)
	return err == nil && fi.Mode()&os.ModeSymlink != 0
}
