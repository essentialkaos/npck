// Package utils provides auxiliary methods for working with archives
package utils

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
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
				return "", fmt.Errorf("Can't eval symlinks: %w", err)
			}
		}
	}

	if !strings.HasPrefix(result, root) {
		return "", fmt.Errorf("Final destination (%s) is outside root (%s)", result, root)
	}

	return result, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func isLink(path string) bool {
	var buf = make([]byte, 1)
	_, err := syscall.Readlink(path, buf)

	return err == nil
}
