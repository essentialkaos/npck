package zst

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"crypto/sha256"
	"os"
	"strings"
	"testing"

	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/hashutil"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type ZSTSuite struct {
	Dir string
}

var _ = Suite(&ZSTSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *ZSTSuite) SetUpSuite(c *C) {
	s.Dir = c.MkDir()
}

func (s *ZSTSuite) TestUnpack(c *C) {
	err := Unpack("../.testdata/data.txt.zst", s.Dir)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data.txt"), Equals, os.FileMode(0640))

	c.Assert(hashutil.File(s.Dir+"/data.txt", sha256.New()).String(), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *ZSTSuite) TestErrors(c *C) {
	err := Unpack("", "/_unknown")
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.txt.zst", "")
	c.Assert(err, NotNil)

	err = Unpack("/_unknown", s.Dir)
	c.Assert(err, NotNil)

	err = Unpack("/_unknown", "/root")
	c.Assert(err, NotNil)

	err = Read(nil, "/_unknown")
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), "")
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), "/_unknown")
	c.Assert(err, NotNil)
}
