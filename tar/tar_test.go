package tar

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"archive/tar"

	"bytes"
	"crypto/sha256"
	"github.com/essentialkaos/npck/v3/utils"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/essentialkaos/ek/v14/fsutil"
	"github.com/essentialkaos/ek/v14/hashutil"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type TarSuite struct{}

var _ = Suite(&TarSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *TarSuite) TestUnpack(c *C) {
	dir := c.MkDir()
	err := Unpack("../.testdata/data.tar", dir, Options{UpdateTimes: true})

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(dir+"/data"), Equals, true)
	c.Assert(fsutil.GetMode(dir+"/data"), Equals, os.FileMode(0700))

	c.Assert(fsutil.IsExist(dir+"/data/payload.txt"), Equals, true)
	c.Assert(fsutil.GetMode(dir+"/data/payload.txt"), Equals, os.FileMode(0644))

	c.Assert(hashutil.File(dir+"/data/payload.txt", sha256.New()).String(), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *TarSuite) TestCPIOUnpack(c *C) {
	dir := c.MkDir()
	err := Unpack("../.testdata/data-cpio.tar", dir, Options{UpdateTimes: true})

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(dir+"/data"), Equals, true)
	c.Assert(fsutil.GetMode(dir+"/data"), Equals, os.FileMode(0750))

	c.Assert(fsutil.IsExist(dir+"/data/payload.txt"), Equals, true)
	c.Assert(fsutil.GetMode(dir+"/data/payload.txt"), Equals, os.FileMode(0644))

	c.Assert(hashutil.File(dir+"/data/payload.txt", sha256.New()).String(), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *TarSuite) TestErrors(c *C) {
	dir := c.MkDir()

	err := Unpack("../.testdata/unknown.tar", dir, Options{UpdateTimes: true})
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.tar", "/unknown", Options{UpdateTimes: true})
	c.Assert(err, NotNil)

	err = Read(nil, "/unknown", Options{UpdateTimes: true})
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), "", Options{UpdateTimes: true})
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), "/_unknown", Options{})
	c.Assert(err, NotNil)

	root, err := os.OpenRoot(dir)
	c.Assert(err, IsNil)
	root.Close()

	err = createDir(root, &tar.Header{}, "_unknown", Options{UpdateTimes: true})
	c.Assert(err, NotNil)

	err = createFile(root, &tar.Header{}, nil, "_unknown", Options{UpdateTimes: true})
	c.Assert(err, NotNil)

	err = createSymlink(root, dir, &tar.Header{Linkname: "test"}, "_unknown", Options{})
	c.Assert(err, NotNil)

	err = createHardlink(root, dir, &tar.Header{Linkname: "test"}, "_unknown", Options{})
	c.Assert(err, NotNil)

	err = createHardlink(root, dir, &tar.Header{Linkname: "../test"}, "_unknown", Options{})
	c.Assert(err, NotNil)

	root, err = os.OpenRoot(dir)
	c.Assert(err, IsNil)
	defer root.Close()

	err = updateAttrs(root, &tar.Header{
		AccessTime: time.Now(),
		ModTime:    time.Now(),
	}, "_unknown", Options{UpdateTimes: true, UpdateOwner: false})
	c.Assert(err, NotNil)

	err = updateAttrs(root, &tar.Header{
		AccessTime: time.Now(),
		ModTime:    time.Now(),
	}, "_unknown", Options{UpdateTimes: false, UpdateOwner: true})
	c.Assert(err, NotNil)

	err = createSymlink(root, dir, &tar.Header{Linkname: "/root/test"}, "_unknown", Options{})
	c.Assert(err, NotNil)

	err = createSymlink(root, dir, &tar.Header{Linkname: "test"}, "../_unknown", Options{})
	c.Assert(err, NotNil)

	err = createHardlink(root, dir, &tar.Header{Linkname: "/root/test"}, "_unknown", Options{})
	c.Assert(err, NotNil)

	_, err = realPath("/_unknown")
	c.Assert(err, NotNil)
}

func (s *TarSuite) TestSymlinkChainEscape(c *C) {
	base := c.MkDir()
	dir := base + "/a/b/c/out"

	c.Assert(os.MkdirAll(dir, 0750), IsNil)
	c.Assert(os.WriteFile(base+"/victim.txt", []byte("original"), 0640), IsNil)

	data := makeTar(c,
		&tar.Header{Name: "d", Typeflag: tar.TypeSymlink, Linkname: "."},
		&tar.Header{Name: "d/d/d/d/l", Typeflag: tar.TypeSymlink, Linkname: "../../../.."},
		&tar.Header{Name: "l/victim.txt", Typeflag: tar.TypeReg, Mode: 0644, Size: 5},
	)

	err := Read(bytes.NewReader(data), dir, Options{})
	c.Assert(err, NotNil)

	victim, err := os.ReadFile(base + "/victim.txt")
	c.Assert(err, IsNil)
	c.Assert(string(victim), Equals, "original")

	// Even if external links are allowed, writing through them must fail
	err = Read(bytes.NewReader(data), c.MkDir(), Options{AllowExternalLinks: true})
	c.Assert(err, NotNil)
}

func (s *TarSuite) TestLaterSymlinkEscape(c *C) {
	data := makeTar(c,
		&tar.Header{Name: "x", Typeflag: tar.TypeSymlink, Linkname: "y/.."},
		&tar.Header{Name: "y", Typeflag: tar.TypeSymlink, Linkname: "."},
	)

	err := Read(bytes.NewReader(data), c.MkDir(), Options{})
	c.Assert(err, ErrorMatches, `symbolic link "x" points to object outside .*`)
}

func (s *TarSuite) TestHardlinkEscape(c *C) {
	base := c.MkDir()
	dir := base + "/a/b/c/out"

	c.Assert(os.MkdirAll(dir, 0750), IsNil)
	c.Assert(os.WriteFile(base+"/secret", []byte("secret"), 0640), IsNil)

	data := makeTar(c,
		&tar.Header{Name: "l", Typeflag: tar.TypeSymlink, Linkname: base},
		&tar.Header{Name: "h", Typeflag: tar.TypeLink, Linkname: "l/secret"},
	)

	err := Read(bytes.NewReader(data), dir, Options{AllowExternalLinks: true})
	c.Assert(err, NotNil)
	c.Assert(fsutil.IsExist(dir+"/h"), Equals, false)
}

func (s *TarSuite) TestHardlinkToSymlinkEscape(c *C) {
	data := makeTar(c,
		&tar.Header{Name: "a/b/c/l", Typeflag: tar.TypeSymlink, Linkname: "../../../etc/passwd"},
		&tar.Header{Name: "evil", Typeflag: tar.TypeLink, Linkname: "a/b/c/l"},
	)

	dir := c.MkDir()
	err := Read(bytes.NewReader(data), dir, Options{})
	c.Assert(err, ErrorMatches, `hard link "evil" points to symbolic link leading outside .*`)
	c.Assert(fsutil.IsLink(dir+"/evil"), Equals, false)

	data = makeTar(c,
		&tar.Header{Name: "a/b/c/l", Typeflag: tar.TypeSymlink, Linkname: "../../../etc/passwd"},
		&tar.Header{Name: "a/b/c/h", Typeflag: tar.TypeLink, Linkname: "a/b/c/l"},
	)

	dir = c.MkDir()
	err = Read(bytes.NewReader(data), dir, Options{})
	c.Assert(err, IsNil)
	c.Assert(fsutil.IsLink(dir+"/a/b/c/h"), Equals, true)
}

func (s *TarSuite) TestPaths(c *C) {
	data := makeTar(c,
		&tar.Header{Name: "./", Typeflag: tar.TypeDir, Mode: 0750},
		&tar.Header{Name: "/abs/file.txt", Typeflag: tar.TypeReg, Mode: 0644, Size: 5},
		&tar.Header{Name: "v1..2.txt", Typeflag: tar.TypeReg, Mode: 0644, Size: 5},
		&tar.Header{Name: "sub/link", Typeflag: tar.TypeSymlink, Linkname: "../v1..2.txt"},
		&tar.Header{Name: "hard", Typeflag: tar.TypeLink, Linkname: "/abs/file.txt"},
	)

	base := c.MkDir()
	c.Assert(os.Mkdir(base+"/real", 0750), IsNil)
	c.Assert(os.Symlink(base+"/real", base+"/link"), IsNil)

	for _, dir := range []string{base + "/real/", base + "/link"} {
		err := Read(bytes.NewReader(data), dir, Options{})
		c.Assert(err, IsNil)

		c.Assert(fsutil.IsExist(base+"/real/abs/file.txt"), Equals, true)
		c.Assert(fsutil.IsExist(base+"/real/v1..2.txt"), Equals, true)
		c.Assert(fsutil.IsLink(base+"/real/sub/link"), Equals, true)
		c.Assert(fsutil.IsExist(base+"/real/hard"), Equals, true)

		os.RemoveAll(base + "/real/sub")
		os.RemoveAll(base + "/real/hard")
	}

	data = makeTar(c,
		&tar.Header{Name: "sub/abs", Typeflag: tar.TypeSymlink, Linkname: base},
	)

	err := Read(bytes.NewReader(data), base+"/real", Options{})
	c.Assert(err, NotNil)

	data = makeTar(c,
		&tar.Header{Name: "a/../../x", Typeflag: tar.TypeReg, Size: 5},
	)

	err = Read(bytes.NewReader(data), base+"/real", Options{})
	c.Assert(err, ErrorMatches, `path .* contains directory traversal element .*`)
}

func (s *TarSuite) TestLinkTargetCheck(c *C) {
	dir := c.MkDir()
	root, err := os.OpenRoot(dir)
	c.Assert(err, IsNil)
	defer root.Close()

	c.Assert(os.Mkdir(dir+"/d", 0750), IsNil)
	c.Assert(os.WriteFile(dir+"/f", nil, 0640), IsNil)
	c.Assert(os.Symlink(".", dir+"/dot"), IsNil)
	c.Assert(os.Symlink(dir+"/d", dir+"/absd"), IsNil)
	c.Assert(os.Symlink("/etc", dir+"/etc"), IsNil)
	c.Assert(os.Symlink("loop", dir+"/loop"), IsNil)

	c.Assert(isExternalLink(root, dir, "l", "d"), Equals, false)
	c.Assert(isExternalLink(root, dir, "d/l", ".."), Equals, false)
	c.Assert(isExternalLink(root, dir, "d/l", "../.."), Equals, true)
	c.Assert(isExternalLink(root, dir, "l", "dot/.."), Equals, true)
	c.Assert(isExternalLink(root, dir, "l", "absd/.."), Equals, false)
	c.Assert(isExternalLink(root, dir, "l", "absd/../.."), Equals, true)
	c.Assert(isExternalLink(root, dir, "l", "etc/passwd"), Equals, true)
	c.Assert(isExternalLink(root, dir, "l", "loop"), Equals, true)
	c.Assert(isExternalLink(root, dir, "l", "f/.."), Equals, true)
	c.Assert(isExternalLink(root, dir, "l", "missing/x"), Equals, false)
	c.Assert(isExternalLink(root, dir, "l", "missing/.."), Equals, true)
	c.Assert(isExternalLink(root, dir, "l", dir+"/d"), Equals, false)
	c.Assert(isExternalLink(root, dir, "l", "/etc"), Equals, true)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// makeTar creates tar archive with given headers, all regular files are filled
// with "PWNED" string
func makeTar(c *C, headers ...*tar.Header) []byte {
	var buf bytes.Buffer

	tw := tar.NewWriter(&buf)

	for _, h := range headers {
		c.Assert(tw.WriteHeader(h), IsNil)

		if h.Typeflag == tar.TypeReg {
			_, err := tw.Write([]byte("PWNED"))
			c.Assert(err, IsNil)
		}
	}

	c.Assert(tw.Close(), IsNil)

	return buf.Bytes()
}

func (s *TarSuite) TestReadLimit(c *C) {
	// 4 bytes of decompressed data (TypeRegA + payload)
	data := makeTar(c,
		&tar.Header{Name: "file.txt", Typeflag: tar.TypeReg, Mode: 0644, Size: 5},
	)

	err := Read(bytes.NewReader(data), c.MkDir(), Options{MaxReadLimit: 3})
	c.Assert(err, Equals, utils.ErrLimitExceeded)

	err = Read(bytes.NewReader(data), c.MkDir(), Options{MaxReadLimit: 5})
	c.Assert(err, IsNil)

	err = Read(bytes.NewReader(data), c.MkDir(), Options{MaxTotalSize: 3})
	c.Assert(err, ErrorMatches, "archive total size is bigger than .*")
}

func (s *TarSuite) TestMaxTotalSizeDefault(c *C) {
	// DefaultOptions must enforce the total cap out of the box
	c.Assert(DefaultOptions.MaxTotalSize, Equals, int64(100*1024*1024*1024))

	data := makeTar(c,
		&tar.Header{Name: "file.txt", Typeflag: tar.TypeReg, Mode: 0644, Size: 5},
		&tar.Header{Name: "file2.txt", Typeflag: tar.TypeReg, Mode: 0644, Size: 5},
	)

	err := Read(bytes.NewReader(data), c.MkDir(), DefaultOptions)
	c.Assert(err, IsNil)

	err = Unpack("../.testdata/data.tar", c.MkDir(), DefaultOptions)
	c.Assert(err, IsNil)
}
