// config/common.go
package config

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateRandomState(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	state := base64.URLEncoding.EncodeToString(bytes)
	return state, nil
}
