package utils

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2026 ESSENTIAL KAOS                          //
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

	_, err = Join(testDir, "mytest/../test2.link")
	c.Assert(err, NotNil)

	_, err = Join(testDir, "mytest/../test3.link")
	c.Assert(err, NotNil)
}

func (s *UtilsSuite) TestLocalPath(c *C) {
	for name, expected := range map[string]string{
		"":              ".",
		"./":            ".",
		"/":             ".",
		"a/b/":          "a/b",
		"/abs/file.txt": "abs/file.txt",
		"./a//b/./c":    "a/b/c",
		"v1..2.txt":     "v1..2.txt",
	} {
		p, err := LocalPath(name)
		c.Assert(err, IsNil, Commentf("name: %q", name))
		c.Assert(p, Equals, expected, Commentf("name: %q", name))
	}

	for _, name := range []string{"..", "../a", "a/../../b", "/../a", "a/.."} {
		_, err := LocalPath(name)
		c.Assert(err, NotNil, Commentf("name: %q", name))
	}
}

func (s *UtilsSuite) TestCreateFile(c *C) {
	base := c.MkDir()
	dir := base + "/out"

	c.Assert(os.Mkdir(dir, 0750), IsNil)
	c.Assert(os.WriteFile(base+"/victim.txt", []byte("original"), 0640), IsNil)
	c.Assert(os.Symlink(base+"/victim.txt", dir+"/link"), IsNil)

	fd, err := CreateFile(dir, "test.txt", 0640)
	c.Assert(err, IsNil)
	c.Assert(fd.Close(), IsNil)

	_, err = CreateFile(dir, "link", 0640)
	c.Assert(err, NotNil)

	victim, err := os.ReadFile(base + "/victim.txt")
	c.Assert(err, IsNil)
	c.Assert(string(victim), Equals, "original")

	_, err = CreateFile(base+"/unknown", "test.txt", 0640)
	c.Assert(err, NotNil)
}

func (s *UtilsSuite) TestDirMode(c *C) {
	c.Assert(DirMode(os.ModeDir|0750), Equals, os.FileMode(0750))
	c.Assert(DirMode(os.ModeDir|os.ModeSticky|0777), Equals, os.ModeSticky|0777)
}
