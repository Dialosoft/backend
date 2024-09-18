package jsonWebToken

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWT(secretKey string, id uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "dialosoft-api",
		Subject:   id.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateRefreshToken(secretKey string, id uuid.UUID) (string, uuid.UUID, error) {
	tokenID := uuid.New()
	claims := jwt.RegisteredClaims{
		Issuer:    "dialosoft-api",
		Subject:   id.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 720)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        tokenID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", uuid.UUID{}, err
	}

	return refreshToken, tokenID, nil
}

func ValidateJWT(tokenString, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token has expired")
			}
		}
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
