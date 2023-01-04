package tgz

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"testing"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/hash"

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
	err := Unpack("../.testdata/data.tgz", s.Dir)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data"), Equals, os.FileMode(0700))

	c.Assert(fsutil.IsExist(s.Dir+"/data/payload.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data/payload.txt"), Equals, os.FileMode(0644))

	c.Assert(hash.FileHash(s.Dir+"/data/payload.txt"), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *TGZSuite) TestErrors(c *C) {
	err := Unpack("../.testdata/unknown.tgz", s.Dir)
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.tgz", "/unknown")
	c.Assert(err, NotNil)

	err = UnpackReader(nil, "/unknown")
	c.Assert(err, NotNil)
}
