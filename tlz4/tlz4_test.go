package tlz4

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"testing"

	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/hash"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type TLZ4Suite struct {
	Dir string
}

var _ = Suite(&TLZ4Suite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *TLZ4Suite) SetUpSuite(c *C) {
	s.Dir = c.MkDir()
}

func (s *TLZ4Suite) TestUnpack(c *C) {
	err := Unpack("../.testdata/data.tlz4", s.Dir)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data"), Equals, os.FileMode(0700))

	c.Assert(fsutil.IsExist(s.Dir+"/data/payload.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data/payload.txt"), Equals, os.FileMode(0644))

	c.Assert(hash.FileHash(s.Dir+"/data/payload.txt"), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *TLZ4Suite) TestErrors(c *C) {
	err := Unpack("../.testdata/unknown.tlz4", s.Dir)
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.tlz4", "/unknown")
	c.Assert(err, NotNil)

	err = Read(nil, "/unknown")
	c.Assert(err, NotNil)
}
