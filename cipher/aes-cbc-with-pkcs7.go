package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func pkcs7Padding(plainText []byte) []byte {
	padLen := aes.BlockSize - (len(plainText) % aes.BlockSize)
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(plainText, padText...)
}

func pkcs7Unpadding(plainText []byte) []byte {
	if textLen := len(plainText); textLen != 0 {
		padLen := int(plainText[textLen-1])
		if padLen >= textLen || padLen > aes.BlockSize {
			return []byte{}
		}
		return plainText[:textLen-padLen]
	}
	return []byte{}
}

func aesEncrypter(plainText, key []byte) []byte {
	plainText = pkcs7Padding(plainText)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)
	return cipherText
}

func aesDecrypter(cipherText, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if len(cipherText) < aes.BlockSize || len(cipherText)%aes.BlockSize != 0 {
		panic("invalid cipherText")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	plainText := pkcs7Unpadding(cipherText)
	return plainText
}

func main() {
	// 256bit key
	// key, _ := hex.DecodeString("6368616e6765207468697320706173736368616e676520746869732070617373")
	// 128bit key
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plainText := []byte("this is an example of aes-cbc-with-pkcs7")
	cipherText := aesEncrypter(plainText, key)
	fmt.Printf("%x\n", cipherText)
	plainText2 := aesDecrypter(cipherText, key)
	fmt.Printf("%s\n", plainText2)
}
