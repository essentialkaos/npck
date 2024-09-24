package lz4

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

	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/hash"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type LZ4Suite struct {
	Dir string
}

var _ = Suite(&LZ4Suite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *LZ4Suite) SetUpSuite(c *C) {
	s.Dir = c.MkDir()
}

func (s *LZ4Suite) TestUnpack(c *C) {
	err := Unpack("../.testdata/data.txt.lz4", s.Dir)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data.txt"), Equals, os.FileMode(0640))

	c.Assert(hash.FileHash(s.Dir+"/data.txt"), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *LZ4Suite) TestErrors(c *C) {
	err := Unpack("", "/_unknown")
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.txt.lz4", "")
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
