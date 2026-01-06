package helpers

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashToken hashes a token using SHA-256
// This is used to securely store refresh tokens in the database
// Returns the hex-encoded hash string
func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// VerifyTokenHash verifies that a token matches a given hash
// Returns true if the token's hash matches the provided hash
func VerifyTokenHash(token, hash string) bool {
	tokenHash := HashToken(token)
	return tokenHash == hash
}
