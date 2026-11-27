package gz

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bufio"
	"crypto/sha256"
	"github.com/essentialkaos/npck/v3/utils"
	"os"
	"strings"
	"testing"

	"github.com/essentialkaos/ek/v14/fsutil"
	"github.com/essentialkaos/ek/v14/hashutil"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type GZSuite struct {
	Dir string
}

var _ = Suite(&GZSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *GZSuite) SetUpSuite(c *C) {
	s.Dir = c.MkDir()
}

func (s *GZSuite) TestUnpack(c *C) {
	err := Unpack("../.testdata/data.txt.gz", s.Dir, Options{})

	c.Assert(err, IsNil)

	c.Assert(fsutil.IsExist(s.Dir+"/data.txt"), Equals, true)
	c.Assert(fsutil.GetMode(s.Dir+"/data.txt"), Equals, os.FileMode(0640))

	c.Assert(hashutil.File(s.Dir+"/data.txt", sha256.New()).String(), Equals, "918c03a211adc19a466c9db22efa575efb6c488fd41c70e57b1ec0920f1a1d8c")
}

func (s *GZSuite) TestErrors(c *C) {
	err := Unpack("", "/_unknown", Options{})
	c.Assert(err, NotNil)

	err = Unpack("../.testdata/data.txt.gz", "", Options{})
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

	err = Read(strings.NewReader(""), s.Dir+"/test", Options{})
	c.Assert(err, NotNil)
}

func (s *GZSuite) TestSymlinkEscape(c *C) {
	base := c.MkDir()
	dir := base + "/out"

	c.Assert(os.Mkdir(dir, 0750), IsNil)
	c.Assert(os.WriteFile(base+"/victim.txt", []byte("original"), 0640), IsNil)
	c.Assert(os.Symlink(base+"/victim.txt", dir+"/data.txt"), IsNil)

	err := Unpack("../.testdata/data.txt.gz", dir, Options{})
	c.Assert(err, NotNil)

	victim, err := os.ReadFile(base + "/victim.txt")
	c.Assert(err, IsNil)
	c.Assert(string(victim), Equals, "original")
}

func (s *GZSuite) TestReadLimit(c *C) {
	fd, err := os.Open("../.testdata/data.txt.gz")
	c.Assert(err, IsNil)

	err = Read(bufio.NewReader(fd), c.MkDir()+"/data.txt", Options{MaxReadLimit: 3})
	c.Assert(err, Equals, utils.ErrLimitExceeded)
}
