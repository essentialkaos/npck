package utils

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"testing"

	. "github.com/essentialkaos/check"
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Test(t *testing.T) { TestingT(t) }

// ////////////////////////////////////////////////////////////////////////////////// //

type UtilsSuite struct{}

var _ = Suite(&UtilsSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

func (s *UtilsSuite) TestJoin(c *C) {
	p, err := Join("/test", "myapp")
	c.Assert(err, IsNil)
	c.Assert(p, Equals, "/test/myapp")

	p, err = Join("/test", "myapp/config/../global.cfg")
	c.Assert(err, IsNil)
	c.Assert(p, Equals, "/test/myapp/global.cfg")

	p, err = Join("/unknown", "myapp/config/../global.cfg")
	c.Assert(err, IsNil)
	c.Assert(p, Equals, "/unknown/myapp/global.cfg")

	tmpDir := c.MkDir()
	os.Mkdir(tmpDir+"/test", 0755)
	os.Symlink(tmpDir+"/test", tmpDir+"/testlink")
	testDir := tmpDir + "/testlink"

	os.Symlink(testDir+"/test.log", testDir+"/test1.link")
	os.WriteFile(testDir+"/test.log", []byte("\n"), 0644)
	os.Symlink(testDir+"/test.log", testDir+"/test1.link")
	os.Symlink("/etc", testDir+"/test2.link")
	os.Symlink(testDir+"/test3.link", testDir+"/test3.link")

	p, err = Join(testDir, "mytest/../test1.link")
	c.Assert(err, IsNil)
	c.Assert(p, Matches, "*/test/test.log")

	p, err = Join(testDir, "mytest/../test2.link")
	c.Assert(err, NotNil)

	p, err = Join(testDir, "mytest/../test3.link")
	c.Assert(err, NotNil)
}
