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
