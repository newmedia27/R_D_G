package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func generateFingerprint(userAgent, ip string) string {
	subnet := isSubnet(ip)
	raw := userAgent + subnet
	hash := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hash[:])
}

func isSubnet(ip string) string {
	arr := strings.Split(ip, ".")
	if len(arr) < 3 {
		return ip
	}
	return strings.Join(arr[:3], ".")
}
