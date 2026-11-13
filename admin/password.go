package admin

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	// Argon2id configuration
	const (
		memory      = 19 * 1024 // 19 MiB, expressed in KiB
		iterations  = 2
		parallelism = 1
		saltLength  = 16
		keyLength   = 32
	)

	// Generate a cryptographically secure random salt.
	salt := make([]byte, saltLength)

	_, err := io.ReadFull(rand.Reader, salt)

	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Generate the password hash.
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		keyLength,
	)

	// Encode salt and hash for database storage.
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	// Store the algorithm and its parameters alongside the hash.
	encodedPassword := fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		memory,
		iterations,
		parallelism,
		encodedSalt,
		encodedHash,
	)

	return encodedPassword, nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")

	if len(parts) != 6 {
		return false, errors.New("invalid encoded password format")
	}

	algorithm := parts[1]
	version := parts[2]
	parameters := parts[3]
	salt := parts[4]
	hash := parts[5]

	if algorithm != "argon2id" {
		return false, errors.New("wrong hashing algorithm is being used")
	}

	if version != "v=19" {
		return false, errors.New("wrong version is being used")
	}

	parameterParts := strings.Split(parameters, ",")
	if len(parameterParts) != 3 {
		return false, errors.New("invalid number of hashing parameters")
	}

	memoryParts := strings.Split(parameterParts[0], "=")
	if len(memoryParts) != 2 || memoryParts[0] != "m" {
		return false, errors.New("invalid memory parameter format")
	}

	iterationParts := strings.Split(parameterParts[1], "=")
	if len(iterationParts) != 2 || iterationParts[0] != "t" {
		return false, errors.New("invalid iteration parameter format")
	}

	parallelParts := strings.Split(parameterParts[2], "=")
	if len(parallelParts) != 2 || parallelParts[0] != "p" {
		return false, errors.New("invalid parallelism parameter format")
	}

	memoryValue, err := strconv.Atoi(memoryParts[1])
	if err != nil {
		return false, fmt.Errorf("invalid memory parameter: %w", err)
	}

	iterationValue, err := strconv.Atoi(iterationParts[1])
	if err != nil {
		return false, fmt.Errorf("invalid iteration parameter: %w", err)
	}

	parallelValue, err := strconv.Atoi(parallelParts[1])
	if err != nil {
		return false, fmt.Errorf("invalid parallelism parameter: %w", err)
	}

	if memoryValue != 19*1024 || iterationValue != 2 || parallelValue != 1 {
		return false, errors.New("unsupported hashing parameters")
	}

	decodedSalt, err := base64.RawStdEncoding.DecodeString(salt)
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	if len(decodedSalt) != 16 {
		return false, errors.New("wrong salt length")
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(hash)
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	if len(decodedHash) != 32 {
		return false, errors.New("wrong hash length")
	}

	calculatedHash := argon2.IDKey(
		[]byte(password),
		decodedSalt,
		uint32(iterationValue),
		uint32(memoryValue),
		uint8(parallelValue),
		uint32(len(decodedHash)),
	)

	if subtle.ConstantTimeCompare(calculatedHash, decodedHash) == 1 {
		return true, nil
	}

	return false, nil
}
