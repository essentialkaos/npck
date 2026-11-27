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

	"github.com/essentialkaos/npck/v3/bz2"
	"github.com/essentialkaos/npck/v3/gz"
	"github.com/essentialkaos/npck/v3/lz4"
	"github.com/essentialkaos/npck/v3/tar"
	"github.com/essentialkaos/npck/v3/tbz2"
	"github.com/essentialkaos/npck/v3/tgz"
	"github.com/essentialkaos/npck/v3/tlz4"
	"github.com/essentialkaos/npck/v3/txz"
	"github.com/essentialkaos/npck/v3/tzst"
	"github.com/essentialkaos/npck/v3/xz"
	"github.com/essentialkaos/npck/v3/zip"
	"github.com/essentialkaos/npck/v3/zst"
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
		return xz.Unpack(file, dir, xz.Options{})
	case ".zst":
		return zst.Unpack(file, dir, zst.Options{})
	case ".lz4":
		return lz4.Unpack(file, dir, lz4.Options{})
	}

	return ErrUnsupportedFormat
}
