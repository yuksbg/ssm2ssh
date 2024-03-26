package keys

import (
	"os"
	"path/filepath"
	"strings"
)

func GetDefaultSSHKey() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	sshDir := filepath.Join(homeDir, ".ssh")
	publicKeyPath := filepath.Join(sshDir, "id_rsa.pub")

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return "", err
	}

	publicKey := strings.TrimSpace(string(publicKeyBytes))
	return publicKey, nil
}
