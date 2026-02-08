package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fezcode/atlas.compass/internal/crypto"
	"github.com/fezcode/atlas.compass/pkg/model"
)

const (
	DirName  = ".atlas"
	FileName = "compass.enc"
)

// GetVaultPath returns the full path to the encrypted vault file.
func GetVaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, DirName, FileName), nil
}

// EnsureDir ensures the config directory exists.
func EnsureDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dirPath := filepath.Join(home, DirName)
	return os.MkdirAll(dirPath, 0700)
}

// Load reads and decrypts the vault.
func Load(password string) (*model.Vault, error) {
	path, err := GetVaultPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// Return empty vault if file doesn't exist
		return &model.Vault{Entries: []model.Entry{}}, nil
	}
	if err != nil {
		return nil, err
	}

	plaintext, err := crypto.Decrypt(data, password)
	if err != nil {
		return nil, err
	}

	var vault model.Vault
	if err := json.Unmarshal(plaintext, &vault); err != nil {
		return nil, fmt.Errorf("corrupted vault data: %w", err)
	}

	return &vault, nil
}

// Save encrypts and writes the vault to disk.
func Save(vault *model.Vault, password string) error {
	if err := EnsureDir(); err != nil {
		return err
	}

	path, err := GetVaultPath()
	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(vault)
	if err != nil {
		return err
	}

	encryptedData, err := crypto.Encrypt(jsonBytes, password)
	if err != nil {
		return err
	}

	return os.WriteFile(path, encryptedData, 0600)
}

// Exists checks if the vault file exists.
func Exists() bool {
	path, err := GetVaultPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return !os.IsNotExist(err)
}
