// Package npck provides methods for unpacking various types of archives
package npck

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/essentialkaos/npck/bz2"
	"github.com/essentialkaos/npck/gz"
	"github.com/essentialkaos/npck/lz4"
	"github.com/essentialkaos/npck/tar"
	"github.com/essentialkaos/npck/tbz2"
	"github.com/essentialkaos/npck/tgz"
	"github.com/essentialkaos/npck/tlz4"
	"github.com/essentialkaos/npck/txz"
	"github.com/essentialkaos/npck/tzst"
	"github.com/essentialkaos/npck/xz"
	"github.com/essentialkaos/npck/zip"
	"github.com/essentialkaos/npck/zst"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var ErrUnsupportedFormat = errors.New("unknown or unsupported archive type")

// ////////////////////////////////////////////////////////////////////////////////// //

// Unpack unpacks archive file to given directory
func Unpack(file, dir string) error {
	fileLower := strings.ToLower(file)
	ext := filepath.Ext(fileLower)

	if strings.HasSuffix(fileLower, ".tar"+ext) {
		ext = ".tar" + ext
	}

	switch ext {
	case ".tgz", ".tar.gz":
		return tgz.Unpack(file, dir, tar.DefaultOptions)
	case ".tbz2", ".tar.bz2":
		return tbz2.Unpack(file, dir, tar.DefaultOptions)
	case ".txz", ".tar.xz":
		return txz.Unpack(file, dir, tar.DefaultOptions)
	case ".tzst", ".tar.zst":
		return tzst.Unpack(file, dir, tar.DefaultOptions)
	case ".tlz4", ".tar.lz4":
		return tlz4.Unpack(file, dir, tar.DefaultOptions)
	case ".zip":
		return zip.Unpack(file, dir, zip.Options{})
	case ".tar":
		return tar.Unpack(file, dir, tar.DefaultOptions)
	case ".gz":
		return gz.Unpack(file, dir, gz.Options{})
	case ".bz2":
		return bz2.Unpack(file, dir, bz2.Options{})
	case ".xz":
		return xz.Unpack(file, dir)
	case ".zst":
		return zst.Unpack(file, dir)
	case ".lz4":
		return lz4.Unpack(file, dir)
	}

	return ErrUnsupportedFormat
}
