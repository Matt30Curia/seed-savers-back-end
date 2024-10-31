package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJson(r *http.Request, payload any) error{

	if r.Body == nil{
		return fmt.Errorf("missing request body")
	}
	
	return json.NewDecoder(r.Body).Decode(payload)
}


func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")
	
	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}


func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}



func DecodePayload[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
    var payload T

    // Parse the JSON payload from the request
    if err := ParseJson(r, &payload); err != nil {
        return nil, fmt.Errorf("failed to parse JSON payload: %w", err)
    }

    // Validate the decoded payload using the validator package
    if err := validatePayload(&payload); err != nil{
		return nil, err
	}

    return &payload, nil
}

func validatePayload(payload any) error{
	if err := Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return fmt.Errorf("invalid payload %v, ", errors)
	}
	return nil
}


func GetAESDecrypted(encrypted string) ([]byte, error) {
	key := "my32digitkey12345678901234567890"
	iv := "my16digitIvKey12"

	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("block size cant be zero")
	}

	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return ciphertext, nil
}


// GetAESEncrypted encrypts given text in AES 256 CBC
func GetAESEncrypted(plaintext string) (string, error) {
	key := "my32digitkey12345678901234567890"
	iv := "my16digitIvKey12"

	var plainTextBlock []byte
	length := len(plaintext)

	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, plaintext)
	block, err := aes.NewCipher([]byte(key))

	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

// PKCS5UnPadding  pads a certain blob of data with necessary data to be used in AES block cipher
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}
