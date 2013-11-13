package common

import (
	"crypto/rand"
	"io"
)

var randomChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// NOTE[jigish]: yes, i know this has modulo bias. i don't care. we don't need a truly random string, just
// one that won't collide often.
func CreateRandomId(size int) string {
	randomBytes := make([]byte, size)
	randomCharsLen := byte(len(randomChars))
	// ignore error here because manas said so. randomBytes is static so if there was an error here we'd be
	// completely screwed anyways.
	io.ReadFull(rand.Reader, randomBytes)
	for i, b := range randomBytes {
		randomBytes[i] = randomChars[b%randomCharsLen]
	}
	return string(randomBytes)
}
