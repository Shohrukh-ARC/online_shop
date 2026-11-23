package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func TestGenerateAndVerifyJWT(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	jwt, err := GenerateJWT(1)
	if err != nil {
		t.Fatalf("failed to generate JWT %v", err)
	}

	claims, err := VerifyJWT(jwt)
	if err != nil {
		t.Fatalf("failed to verify JWT %v", err)
	}

	if claims.Subject != "1" {
		t.Errorf("Expected subject 1, got %s", claims.Subject)
	}

	if claims.Issuer != "online-shop" {
		t.Errorf("Expected issuer online-shop, got %s", claims.Issuer)
	}
}

func TestInvalidJWT(t *testing.T) {
	_, err := VerifyJWT("invalid.token")

	if err == nil {
		t.Errorf("Expected %v, got %v", jwt.ErrTokenMalformed, err)
	}
}

func TestExpiredJWT(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	adminClaims := AdminClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			Issuer:    "online-shop",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, adminClaims)
	secret := os.Getenv("JWT_SECRET")

	secretBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		t.Fatalf("failed to decode JWT secret: %v", err)
	}

	if len(secretBytes) != 32 {
		t.Fatalf("invalid JWT secret length: expected 32 bytes, got %d", len(secretBytes))
	}

	signedToken, err := token.SignedString(secretBytes)
	if err != nil {
		t.Fatalf("failed to sign JWT: %v", err)
	}

	_, err = VerifyJWT(signedToken)
	if err == nil {
		t.Fatal("Expected expired JWT verification to fail")
	}

	if !errors.Is(err, jwt.ErrTokenExpired) {
		t.Fatalf("Expected %v, got %v", jwt.ErrTokenExpired, err)
	}
}

func TestJWTWrongSecret(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Failed to load .env file: %v", err)
	}

	adminClaims := AdminClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "1",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "online-shop",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, adminClaims)

	wrongSecret := make([]byte, 32)
	_, err = rand.Read(wrongSecret)
	if err != nil {
		t.Fatalf("failed to generate random secret: %v", err)
	}

	signedToken, err := token.SignedString(wrongSecret)
	if err != nil {
		t.Fatalf("failed to sign JWT: %v", err)
	}

	_, err = VerifyJWT(signedToken)
	if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		t.Fatalf("Expected %v, got %v", jwt.ErrTokenSignatureInvalid, err)
	}
}
