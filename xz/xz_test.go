package xz

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
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

type XZSuite struct {
	Dir string
}

var _ = Suite(&XZSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *XZSuite) SetUpSuite(c *C) {
	s.Dir = c.MkDir()
}

func (s *XZSuite) TestUnpack(c *C) {
	err := Unpack("../.testdata/data.txt.xz", s.Dir)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data.txt"), Equals, os.FileMode(0640))

	c.Assert(hash.FileHash(s.Dir+"/data.txt"), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *XZSuite) TestErrors(c *C) {
	err := Unpack("", "/_unknown")
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.txt.xz", "")
	c.Assert(err, NotNil)

	err = Unpack("/_unknown", s.Dir)
	c.Assert(err, NotNil)

	err = Read(nil, "/_unknown")
	c.Assert(err, NotNil)

	err = Unpack("/_unknown", "/root")
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), "")
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), "/_unknown")
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), s.Dir+"/test")
	c.Assert(err, NotNil)
}
