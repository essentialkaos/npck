package tar

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"archive/tar"
	"os"
	"strings"
	"testing"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/hash"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type TarSuite struct {
	Dir string
}

var _ = Suite(&TarSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *TarSuite) SetUpSuite(c *C) {
	s.Dir = c.MkDir()
}

func (s *TarSuite) TestUnpack(c *C) {
	err := Unpack("../.testdata/data.tar", s.Dir)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data"), Equals, os.FileMode(0700))

	c.Assert(fsutil.IsExist(s.Dir+"/data/payload.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data/payload.txt"), Equals, os.FileMode(0644))

	c.Assert(hash.FileHash(s.Dir+"/data/payload.txt"), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *TarSuite) TestErrors(c *C) {
	err := Unpack("../.testdata/unknown.tar", s.Dir)
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
	err = updateAttrs(&tar.Header{Linkname: "/__unknown"}, "/_unknown")
	c.Assert(err, NotNil)

	UpdateTimes, UpdateOwner = false, true
	err = updateAttrs(&tar.Header{Linkname: "/__unknown"}, "/_unknown")
	c.Assert(err, NotNil)

	UpdateTimes, UpdateOwner = true, false

	err = createSymlink(&tar.Header{Linkname: "/root/test"}, "/unknown", "/_unknown")
	c.Assert(err, NotNil)

	err = createHardlink(&tar.Header{Linkname: "/root/test"}, "/unknown", "/_unknown")
	c.Assert(err, NotNil)

	c.Assert(isExternalLink("../unknown", "/root"), Equals, true)
}
