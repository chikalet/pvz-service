package auth
import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)
var secret = []byte("supersecret") 
func GenerateToken(role string) (string, error) {
	claims := jwt.MapClaims{
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
		"iss":  "pvz-service",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
func ParseToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role, _ := claims["role"].(string)
		return role, nil
	}
	return "", err
}