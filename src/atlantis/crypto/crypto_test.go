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

package crypto

import (
	"launchpad.net/gocheck"
	"testing"
)

func TestCrypto(t *testing.T) { gocheck.TestingT(t) }

type CryptoSuite struct{}

var _ = gocheck.Suite(&CryptoSuite{})

func (s *CryptoSuite) TestBytes(c *gocheck.C) {
	c.Check(Init(), gocheck.IsNil)
	test := "a" // simple test
	c.Check(string(Decrypt(Encrypt([]byte(test)))), gocheck.Equals, test)
	test = "abcdabcdabcdabcd" // 16 byte test (block size)
	c.Check(string(Decrypt(Encrypt([]byte(test)))), gocheck.Equals, test)
	test = "abcdabcdabcdabcda" // 17 byte test (block size + 1)
	c.Check(string(Decrypt(Encrypt([]byte(test)))), gocheck.Equals, test)
	test = "`1234567890-=qwertyuiop[]\\asdfghjkl;'zxcvbnm,./" // random character test
	c.Check(string(Decrypt(Encrypt([]byte(test)))), gocheck.Equals, test)
	test = "~!@#$%^&*()_+QWERTYUIOP{}|ASDFGHJKL:\"ZXCVBNM<>?" // more random characters
	c.Check(string(Decrypt(Encrypt([]byte(test)))), gocheck.Equals, test)
	// test long string
	test = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	c.Check(string(Decrypt(Encrypt([]byte(test)))), gocheck.Equals, test)
}

func (s *CryptoSuite) TestString(c *gocheck.C) {
	c.Check(Init(), gocheck.IsNil)
	test := "a" // simple test
	c.Check(DecryptString(EncryptString(test)), gocheck.Equals, test)
	test = "abcdabcdabcdabcd" // 16 byte test (block size)
	c.Check(DecryptString(EncryptString(test)), gocheck.Equals, test)
	test = "abcdabcdabcdabcda" // 17 byte test (block size + 1)
	c.Check(DecryptString(EncryptString(test)), gocheck.Equals, test)
	test = "`1234567890-=qwertyuiop[]\\asdfghjkl;'zxcvbnm,./" // random character test
	c.Check(DecryptString(EncryptString(test)), gocheck.Equals, test)
	test = "~!@#$%^&*()_+QWERTYUIOP{}|ASDFGHJKL:\"ZXCVBNM<>?" // more random characters
	c.Check(DecryptString(EncryptString(test)), gocheck.Equals, test)
	// test long string
	test = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	c.Check(DecryptString(EncryptString(test)), gocheck.Equals, test)
}
