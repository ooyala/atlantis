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
	"crypto/rand"
	"io"
)

var randomChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// NOTE[jigish]: yes, i know this has modulo bias. i don't care. we don't need a truly random string, just
// one that won't collide often.
func CreateRandomID(size int) string {
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
