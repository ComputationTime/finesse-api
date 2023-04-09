package database

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	SecretKey = []byte("secret")
)

func GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["username"] = username
	// change this line to modify expiration
	claims["exp"] = time.Now().Add(48 * time.Hour).Unix()
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	} else {
		return "", err
	}
}

func VerifyToken(tokenStr string, username string) (bool, error) {
	token_username, err := ParseToken(tokenStr)
	if err != nil {
		return false, err
	}
	var is_same_user = token_username == username
	return is_same_user, nil
}