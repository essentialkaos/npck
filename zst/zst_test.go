package zst

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"crypto/sha256"
	"os"
	"strings"
	"testing"

	"github.com/essentialkaos/ek/v14/fsutil"
	"github.com/essentialkaos/ek/v14/hashutil"

	"github.com/essentialkaos/npck/v3/utils"

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
	err := Unpack("../.testdata/data.txt.zst", s.Dir, Options{})

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data.txt"), Equals, os.FileMode(0640))

	c.Assert(hashutil.File(s.Dir+"/data.txt", sha256.New()).String(), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *ZSTSuite) TestErrors(c *C) {
	err := Unpack("", "/_unknown", Options{})
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.txt.zst", "", Options{})
	c.Assert(err, NotNil)

	err = Unpack("/_unknown", s.Dir, Options{})
	c.Assert(err, NotNil)

	err = Unpack("/_unknown", "/root", Options{})
	c.Assert(err, NotNil)

	err = Read(nil, "/_unknown", Options{})
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), "", Options{})
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), "/_unknown", Options{})
	c.Assert(err, NotNil)
}

func (s *ZSTSuite) TestReadLimit(c *C) {
	err := Unpack("../.testdata/data.txt.zst", c.MkDir(), Options{MaxReadLimit: 3})
	c.Assert(err, Equals, utils.ErrLimitExceeded)
}
