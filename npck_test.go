package npck

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"testing"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type NPCKSuite struct {
	Dir string
}

var _ = Suite(&NPCKSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *NPCKSuite) TestUnpack(c *C) {
	c.Assert(Unpack(".testdata/data.txt.gz", c.MkDir()), IsNil)
	c.Assert(Unpack(".testdata/data.txt.bz2", c.MkDir()), IsNil)
	c.Assert(Unpack(".testdata/data.txt.xz", c.MkDir()), IsNil)
	c.Assert(Unpack(".testdata/data.txt.zst", c.MkDir()), IsNil)

	c.Assert(Unpack(".testdata/data.tar", c.MkDir()), IsNil)
	c.Assert(Unpack(".testdata/data.tgz", c.MkDir()), IsNil)
	c.Assert(Unpack(".testdata/data.tbz2", c.MkDir()), IsNil)
	c.Assert(Unpack(".testdata/data.txz", c.MkDir()), IsNil)
	c.Assert(Unpack(".testdata/data.tzst", c.MkDir()), IsNil)
	c.Assert(Unpack(".testdata/data.zip", c.MkDir()), IsNil)
	c.Assert(Unpack(".testdata/data.tar.gz", c.MkDir()), IsNil)

	c.Assert(Unpack(".testdata/data.jpg", c.MkDir()), NotNil)
}
