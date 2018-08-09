package aes

import (
	"crypto/aes"
	"crypto/cipher"
)

func CBCEncrypt(plainText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	plainText = PKCS7Padding(plainText, block.BlockSize())
	blockModel := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])
	cipherText := make([]byte, len(plainText))
	blockModel.CryptBlocks(cipherText, plainText)
	return cipherText, nil
}

func CBCDecrypt(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	plainText := make([]byte, len(cipherText))
	blockModel.CryptBlocks(plainText, cipherText)
	plainText = PKCS7UnPadding(plainText, block.BlockSize())
	return plainText, nil
}
