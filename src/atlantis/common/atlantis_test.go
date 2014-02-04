/* Copyright 2014 Ooyala, Inc. All rights reserved.
 *
 * This file is licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
 * except in compliance with the License. You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is
 * distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

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
