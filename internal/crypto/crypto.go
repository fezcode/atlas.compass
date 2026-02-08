package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/argon2"
)

const (
	SaltSize   = 16
	NonceSize  = 12
	KeySize    = 32 // AES-256
	ArgonTime  = 1
	ArgonMem   = 64 * 1024
	ArgonThreads = 4
)

// DeriveKey derives a 32-byte key from the password and salt using Argon2id.
func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, ArgonTime, ArgonMem, ArgonThreads, KeySize)
}

// Encrypt encrypts the plaintext using the password.
// It returns a byte slice containing Salt + Nonce + Ciphertext.
func Encrypt(plaintext []byte, password string) ([]byte, error) {
	// 1. Generate Salt
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	// 2. Derive Key
	key := DeriveKey(password, salt)

	// 3. Create Cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 4. Generate Nonce
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 5. Encrypt
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	// 6. Pack: Salt + Nonce + Ciphertext
	result := make([]byte, SaltSize+NonceSize+len(ciphertext))
	copy(result[0:SaltSize], salt)
	copy(result[SaltSize:SaltSize+NonceSize], nonce)
	copy(result[SaltSize+NonceSize:], ciphertext)

	return result, nil
}

// Decrypt decrypts the data using the password.
// It expects data to be Salt + Nonce + Ciphertext.
func Decrypt(data []byte, password string) ([]byte, error) {
	if len(data) < SaltSize+NonceSize {
		return nil, errors.New("invalid data length")
	}

	// 1. Extract Salt
	salt := data[:SaltSize]
	
	// 2. Extract Nonce
	nonce := data[SaltSize : SaltSize+NonceSize]

	// 3. Extract Ciphertext
	ciphertext := data[SaltSize+NonceSize:]

	// 4. Derive Key
	key := DeriveKey(password, salt)

	// 5. Create Cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 6. Decrypt
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed: invalid password or corrupted data")
	}

	return plaintext, nil
}
