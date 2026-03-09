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
	"os"
	"path/filepath"
	"strings"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var (
	ErrNilReader   = errors.New("reader is nil")
	ErrEmptyInput  = errors.New("path to input file is empty")
	ErrEmptyOutput = errors.New("path to output file is empty")
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Join joins all elements of path, makes lexical processing, and evaluating all symlinks.
// Method returns error if final destination is not a child path of root.
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

	if IsExternalPath(root, result) {
		return "", fmt.Errorf("final destination (%s) is outside root (%s)", result, root)
	}

	return result, nil
}

// IsExternalPath returns true if given path is outside of root
func IsExternalPath(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err != nil || strings.HasPrefix(rel, "..")
}

// ////////////////////////////////////////////////////////////////////////////////// //

func isLink(path string) bool {
	fi, err := os.Lstat(path)
	return err == nil && fi.Mode()&os.ModeSymlink != 0
}
