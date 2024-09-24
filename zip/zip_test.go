package zip

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

type ZipSuite struct {
	Dir string
}

var _ = Suite(&ZipSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *ZipSuite) SetUpSuite(c *C) {
	s.Dir = c.MkDir()
}

func (s *ZipSuite) TestUnpack(c *C) {
	err := Unpack("../.testdata/data.zip", s.Dir)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data"), Equals, os.FileMode(0700))

	c.Assert(fsutil.IsExist(s.Dir+"/data/payload.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data/payload.txt"), Equals, os.FileMode(0644))

	c.Assert(hash.FileHash(s.Dir+"/data/payload.txt"), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *ZipSuite) TestErrors(c *C) {
	err := Unpack("../.testdata/unknown.zip", s.Dir)
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.zip", "/unknown")
	c.Assert(err, NotNil)

	err = Read(nil, "/unknown")
	c.Assert(err, NotNil)

	err = Read(nil, "/unknown")
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), "")
	c.Assert(err, NotNil)

	err = Read(strings.NewReader(""), s.Dir)
	c.Assert(err, NotNil)
}
