package tbz2

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"crypto/sha256"
	"os"
	"testing"

	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/hashutil"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type TBZ2Suite struct {
	Dir string
}

var _ = Suite(&TBZ2Suite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *TBZ2Suite) SetUpSuite(c *C) {
	s.Dir = c.MkDir()
}

func (s *TBZ2Suite) TestUnpack(c *C) {
	err := Unpack("../.testdata/data.tbz2", s.Dir)

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data"), Equals, os.FileMode(0700))

	c.Assert(fsutil.IsExist(s.Dir+"/data/payload.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data/payload.txt"), Equals, os.FileMode(0644))

	c.Assert(hashutil.File(s.Dir+"/data/payload.txt", sha256.New()), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *TBZ2Suite) TestErrors(c *C) {
	err := Unpack("../.testdata/unknown.tbz2", s.Dir)
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.tbz2", "/unknown")
	c.Assert(err, NotNil)

	err = Read(nil, "/unknown")
	c.Assert(err, NotNil)
}
