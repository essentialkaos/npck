package tar

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"archive/tar"
	"crypto/sha256"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/hashutil"

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
	err := Unpack("../.testdata/data.tar", dir)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(dir+"/data"), Equals, true)
	c.Assert(fsutil.GetMode(dir+"/data"), Equals, os.FileMode(0700))

	c.Assert(fsutil.IsExist(dir+"/data/payload.txt"), Equals, true)
	c.Assert(fsutil.GetMode(dir+"/data/payload.txt"), Equals, os.FileMode(0644))

	c.Assert(hashutil.File(dir+"/data/payload.txt", sha256.New()).String(), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *TarSuite) TestCPIOUnpack(c *C) {
	dir := c.MkDir()
	err := Unpack("../.testdata/data-cpio.tar", dir)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(dir+"/data"), Equals, true)
	c.Assert(fsutil.GetMode(dir+"/data"), Equals, os.FileMode(0755))

	c.Assert(fsutil.IsExist(dir+"/data/payload.txt"), Equals, true)
	c.Assert(fsutil.GetMode(dir+"/data/payload.txt"), Equals, os.FileMode(0644))

	c.Assert(hashutil.File(dir+"/data/payload.txt", sha256.New()).String(), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *TarSuite) TestErrors(c *C) {
	dir := c.MkDir()

	err := Unpack("../.testdata/unknown.tar", dir)
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.tar", "/unknown")
	c.Assert(err, NotNil)

	err = Read(nil, "/unknown")
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), "")
	c.Assert(err, NotNil)

	err = createDir(&tar.Header{}, "/_unknown")
	c.Assert(err, NotNil)

	err = createFile(&tar.Header{}, nil, "/_unknown")
	c.Assert(err, NotNil)

	err = createSymlink(&tar.Header{Linkname: "/__unknown"}, "", "/_unknown")
	c.Assert(err, NotNil)

	err = createHardlink(&tar.Header{Linkname: "/__unknown"}, "", "/_unknown")
	c.Assert(err, NotNil)

	UpdateTimes, UpdateOwner = true, false
	err = updateAttrs(&tar.Header{
		Linkname:   "/__unknown",
		AccessTime: time.Now(),
		ModTime:    time.Now(),
	}, "/_unknown")
	c.Assert(err, NotNil)

	UpdateTimes, UpdateOwner = false, true
	err = updateAttrs(&tar.Header{
		Linkname:   "/__unknown",
		AccessTime: time.Now(),
		ModTime:    time.Now(),
	}, "/_unknown")
	c.Assert(err, NotNil)

	UpdateTimes, UpdateOwner = true, false

	err = createSymlink(&tar.Header{Linkname: "/root/test"}, "/unknown", "/_unknown")
	c.Assert(err, NotNil)

	err = createHardlink(&tar.Header{Linkname: "/root/test"}, "/unknown", "/_unknown")
	c.Assert(err, NotNil)

	c.Assert(isExternalLink("../../unknown", "/root/test", "/root"), Equals, true)
}
