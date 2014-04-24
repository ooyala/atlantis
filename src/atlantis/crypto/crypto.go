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
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

var block cipher.Block

func Init() error {
	var err error
	block, err = aes.NewCipher(AES_KEY)
	return err
}

func pad(data []byte) []byte {
	var size int
	if extra := len(data) % block.BlockSize(); extra > 0 {
		size = len(data) - extra + block.BlockSize()
	} else {
		size = len(data)
	}
	padded := make([]byte, size, size)
	for i, b := range data {
		padded[i] = b
	}
	return padded
}

func unpad(data []byte) []byte {
	unpadded := []byte{}
	for _, b := range data {
		if b == 0 { // assume that there are no 0 bytes inbetween the data
			return unpadded
		}
		unpadded = append(unpadded, b)
	}
	return unpadded
}

func chomp(data []byte) []byte {
	var size int
	if extra := len(data) % block.BlockSize(); extra > 0 {
		size = len(data) - extra
	} else {
		size = len(data)
	}
	chomped := make([]byte, size, size)
	for i, _ := range chomped {
		chomped[i] = data[i]
	}
	return chomped
}

func clone(data []byte) []byte {
	theClone := make([]byte, len(data))
	for i, b := range data {
		theClone[i] = b
	}
	return theClone
}

func Encrypt(decrypted []byte) []byte {
	encrypter := cipher.NewCBCEncrypter(block, AES_IV) // encrypter
	toEncrypt := pad(decrypted)                        // clone and pad decrypted so we can mutate it
	encrypted := toEncrypt                             // maintain pointer to backing array because we're encrypting in place
	for len(toEncrypt) > 0 {                           // encrypt block by block
		encryptedPart := toEncrypt[0:block.BlockSize()]     // fetch the block
		toEncrypt = toEncrypt[block.BlockSize():]           // cut the block off toEncrypt
		encrypter.CryptBlocks(encryptedPart, encryptedPart) // encrypt the block
	}
	// B64 encode the shits
	b64ed := make([]byte, base64.StdEncoding.EncodedLen(len(encrypted)))
	base64.StdEncoding.Encode(b64ed, encrypted)
	return b64ed
}

func Decrypt(b64ed []byte) []byte {
	decrypter := cipher.NewCBCDecrypter(block, AES_IV)
	// B64 decode the shits
	toDecrypt := make([]byte, base64.StdEncoding.DecodedLen(len(b64ed)))
	base64.StdEncoding.Decode(toDecrypt, b64ed)
	toDecrypt = chomp(toDecrypt) // chomp extra bytes because base64 is stupid
	decrypted := toDecrypt
	for len(toDecrypt) > 0 { // decrypt block by block
		decryptedPart := toDecrypt[0:block.BlockSize()]
		toDecrypt = toDecrypt[block.BlockSize():]           // cut the block off toDecrypt
		decrypter.CryptBlocks(decryptedPart, decryptedPart) // decrypt the block
	}
	return unpad(decrypted) // unpad and we're done
}

func EncryptString(decrypted string) string {
	return string(Encrypt([]byte(decrypted)))
}

func DecryptString(b64ed string) string {
	return string(Decrypt([]byte(b64ed)))
}
