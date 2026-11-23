package auth

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AdminClaims struct {
	jwt.RegisteredClaims
}

func GenerateJWT(adminID int) (string, error) {
	adminClaims := AdminClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(adminID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "online-shop",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, adminClaims)
	secret := os.Getenv("JWT_SECRET")

	secretBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode JWT secret: %w", err)
	}

	if len(secretBytes) != 32 {
		return "", fmt.Errorf("invalid JWT secret length: expected 32 bytes, got %d", len(secretBytes))
	}

	signedToken, err := token.SignedString(secretBytes)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signedToken, nil
}

func VerifyJWT(tokenString string) (*AdminClaims, error) {
	secret := os.Getenv("JWT_SECRET")

	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("failed to decode jwt secret string: %w", err)
	}

	if len(decodedSecret) != 32 {
		return nil, fmt.Errorf("invalid JWT secret length: expected 32 bytes, got %d", len(decodedSecret))
	}

	claims := &AdminClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return decodedSecret, nil
		},

		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("online-shop"),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}
