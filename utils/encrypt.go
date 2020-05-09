package utils

import (
	//"fmt"
	//"os"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

/*
func main() {
	originalText := "encrypt this golang"
	fmt.Println(originalText)

	key := "example key 1234"

	// encrypt value to base64
	cryptoText, _ := encrypt(originalText, key)
	fmt.Println(cryptoText)

	// encrypt base64 crypto to original value
	text, _ := decrypt(cryptoText, key)
	fmt.Println(text)
}
*/
// Takes two strings, cryptoText and keyString.
// cryptoText is the text to be decrypted and the keyString is the key to use for the decryption.
// The function will output the resulting plain text string with an error variable.
func decrypt(cryptoText string, keyString string) (plainTextString string, err error) {

	// Format the keyString so that it's 32 bytes.
	newKeyString, err := hashTo32Bytes(keyString)

	// Encode the cryptoText to base 64.
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher([]byte(newKeyString))

	if err != nil {
		panic(err)
	}

	if len(cipherText) < aes.BlockSize {
		panic("*** cipherText too short ***\n")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// Takes two string, plainText and keyString.
// plainText is the text that needs to be encrypted by keyString.
// The function will output the resulting crypto text and an error variable.
func encrypt(plainText string, keyString string) (cipherTextString string, err error) {

	// Format the keyString so that it's 32 bytes.
	newKeyString, err := hashTo32Bytes(keyString)

	if err != nil {
		return "", err
	}

	key := []byte(newKeyString)
	value := []byte(plainText)

	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(value))

	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], value)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// As we cannot use a variable length key, we must cut the users key
// up to or down to 32 bytes. To do this the function takes a hash
// of the key and cuts it down to 32 bytes.
func hashTo32Bytes(input string) (output string, err error) {

	if len(input) == 0 {
		return "", errors.New("No input supplied")
	}

	hasher := sha256.New()
	hasher.Write([]byte(input))

	stringToSHA256 := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Cut the length down to 32 bytes and return.
	return stringToSHA256[:32], nil
}
