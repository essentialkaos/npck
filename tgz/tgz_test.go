package tgz

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"crypto/sha256"
	"os"
	"testing"

	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/hashutil"

	"github.com/essentialkaos/npck/tar"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type TGZSuite struct {
	Dir string
}

var _ = Suite(&TGZSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *TGZSuite) SetUpSuite(c *C) {
	s.Dir = c.MkDir()
}

func (s *TGZSuite) TestUnpack(c *C) {
	err := Unpack("../.testdata/data.tgz", s.Dir, tar.DefaultOptions)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data"), Equals, os.FileMode(0700))

	c.Assert(fsutil.IsExist(s.Dir+"/data/payload.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data/payload.txt"), Equals, os.FileMode(0644))

	c.Assert(hashutil.File(s.Dir+"/data/payload.txt", sha256.New()).String(), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *TGZSuite) TestErrors(c *C) {
	err := Unpack("../.testdata/unknown.tgz", s.Dir, tar.DefaultOptions)
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.tgz", "/unknown", tar.DefaultOptions)
	c.Assert(err, NotNil)

	err = Read(nil, "/unknown", tar.DefaultOptions)
	c.Assert(err, NotNil)
}
