package common

import (
	"launchpad.net/gocheck"
	"testing"
)

func TestAtlantis(t *testing.T) { gocheck.TestingT(t) }

type AtlantisSuite struct{}

var _ = gocheck.Suite(&AtlantisSuite{})

func (s *AtlantisSuite) TestDiffSlices(c *gocheck.C) {
	var s1, s2, onlyInS1, onlyInS2 []string
	s1 = nil
	s2 = nil
	onlyInS1, onlyInS2 = DiffSlices(s1, s2)
	c.Check(onlyInS1, gocheck.DeepEquals, []string{})
	c.Check(onlyInS2, gocheck.DeepEquals, []string{})
	s1 = []string{}
	onlyInS1, onlyInS2 = DiffSlices(s1, s2)
	c.Check(onlyInS1, gocheck.DeepEquals, []string{})
	c.Check(onlyInS2, gocheck.DeepEquals, []string{})
	s2 = []string{}
	onlyInS1, onlyInS2 = DiffSlices(s1, s2)
	c.Check(onlyInS1, gocheck.DeepEquals, []string{})
	c.Check(onlyInS2, gocheck.DeepEquals, []string{})
	s1 = []string{"a", "b"}
	s2 = []string{}
	onlyInS1, onlyInS2 = DiffSlices(s1, s2)
	c.Check(onlyInS1, gocheck.DeepEquals, []string{"a", "b"})
	c.Check(onlyInS2, gocheck.DeepEquals, []string{})
	s1 = []string{"a", "b"}
	s2 = []string{"b", "c"}
	onlyInS1, onlyInS2 = DiffSlices(s1, s2)
	c.Check(onlyInS1, gocheck.DeepEquals, []string{"a"})
	c.Check(onlyInS2, gocheck.DeepEquals, []string{"c"})
	s1 = []string{"a", "b", "b", "b"}
	s2 = []string{"b", "c"}
	onlyInS1, onlyInS2 = DiffSlices(s1, s2)
	c.Check(onlyInS1, gocheck.DeepEquals, []string{"a", "b", "b"})
	c.Check(onlyInS2, gocheck.DeepEquals, []string{"c"})
}
