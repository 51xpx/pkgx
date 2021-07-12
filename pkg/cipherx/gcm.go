package cipherx

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// Gcm .
type Gcm struct {
	Key   string
	Nonce string
}

// NewKey new random secret key
func NewKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

// NewNonce new random nonce
func NewNonce() (string, error) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	return hex.EncodeToString(nonce), nil
}

// NewGcm new cipher with key
func NewGcm(key string) (*Gcm, error) {
	nonce, err := NewNonce()
	if err != nil {
		return nil, err
	}
	return &Gcm{
		Key:   key,
		Nonce: nonce,
	}, nil
}

// Encrypt aes encrypt
func (g *Gcm) Encrypt(content string) (string, error) {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	key, _ := hex.DecodeString(g.Key)
	plaintext := []byte(content)
	nonce, _ := hex.DecodeString(g.Nonce)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", nil
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", nil
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt aes decrypt
func (g *Gcm) Decrypt(content string) (string, error) {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal/Open calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	// When decoded the key should be 16 bytes (AES-128) or 32 (AES-256).
	key, _ := hex.DecodeString(g.Key)
	ciphertext, _ := hex.DecodeString(content)
	nonce, _ := hex.DecodeString(g.Nonce)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
