package object

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret_key")

// GenerateJWT generates a JWT token for the given user ID and role.
func GenerateJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token valid for 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ParseJWT parses the JWT token and returns the user ID and role.
func ParseJWT(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", "", errors.New("user_id not found in claims")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", "", errors.New("role not found in claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", "", errors.New("exp not found in claims")
	}
	if int64(exp) < time.Now().Unix() {
		return "", "", errors.New("token expired")
	}

	return userID, role, nil
}
