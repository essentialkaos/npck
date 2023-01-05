// Package npck provides methods for unpacking various types of archives
package npck

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

	"github.com/essentialkaos/npck/bz2"
	"github.com/essentialkaos/npck/gz"
	"github.com/essentialkaos/npck/tar"
	"github.com/essentialkaos/npck/tbz2"
	"github.com/essentialkaos/npck/tgz"
	"github.com/essentialkaos/npck/txz"
	"github.com/essentialkaos/npck/tzst"
	"github.com/essentialkaos/npck/xz"
	"github.com/essentialkaos/npck/zst"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpack unpacks given file
func Unpack(file, dir string) error {
	ext := filepath.Ext(file)

	if strings.Contains(file, ".tar.") {
		ext = ".tar" + ext
	}

	switch ext {
	case ".tgz", ".tar.gz":
		return tgz.Unpack(file, dir)
	case ".tbz2", ".tar.bz2":
		return tbz2.Unpack(file, dir)
	case ".txz", ".tar.xz":
		return txz.Unpack(file, dir)
	case ".tzst", ".tar.zst":
		return tzst.Unpack(file, dir)
	case ".tar":
		return tar.Unpack(file, dir)
	case ".gz":
		return gz.Unpack(file, dir)
	case ".bz2":
		return bz2.Unpack(file, dir)
	case ".xz":
		return xz.Unpack(file, dir)
	case ".zst":
		return zst.Unpack(file, dir)
	}

	return fmt.Errorf("Unknown or unsupported archive type")
}
