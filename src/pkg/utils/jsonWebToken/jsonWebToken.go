package jsonWebToken

import (
	"fmt"
	"time"

	"github.com/Dialosoft/src/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateAccessJWT generates a signed JWT access token using the given secret key.
// It includes claims such as the user ID, role ID, expiration time (5 minutes), and issued at time.
// Returns the signed JWT as a string or an error if the signing process fails.
func GenerateAccessJWT(secretKey string, id uuid.UUID, roleID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"iss": "dialosoft-api",
		"sub": id.String(),
		"rid": roleID.String(),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Minute * 5)).Unix(),
		"iat": jwt.NewNumericDate(time.Now()).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GenerateRefreshToken generates a refresh token and returns the token along with a TokenEntity.
// The refresh token has an expiration time of 720 hours (30 days) and is signed using the provided secret key.
// Returns the signed refresh token, a TokenEntity containing metadata, or an error if token creation fails.
func GenerateRefreshToken(secretKey string, userID uuid.UUID) (string, models.TokenEntity, error) {
	tokenID := uuid.New()
	claims := jwt.MapClaims{
		"iss": "dialosoft-api",
		"sub": userID.String(),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour * 720)),
		"iat": jwt.NewNumericDate(time.Now()),
		"jti": tokenID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", models.TokenEntity{}, err
	}

	tokenEntity := models.TokenEntity{
		Token:     refreshToken,
		ID:        tokenID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour * 720),
		CreatedAt: time.Now(),
	}

	return refreshToken, tokenEntity, nil
}

// ValidateJWT validates the given JWT token string using the provided secret key.
// It checks if the token's signing method is HMAC and verifies the token's expiration time.
// Returns the token claims if valid, or an error if the token is invalid or expired.
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
