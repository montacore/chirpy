package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	issued := time.Now().UTC()
	expires := issued.Add(expiresIn)
	
	claims := jwt.RegisteredClaims{
		IssuedAt: jwt.NewNumericDate(issued),
		ExpiresAt: jwt.NewNumericDate(expires), 
		Issuer: "chirpy-access",
		Subject: userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	keyForToken := func(_ *jwt.Token) (any, error){
		return []byte(tokenSecret), nil
	}
	token, err := jwt.ParseWithClaims(tokenString, claims, keyForToken)
	if err != nil {
		return uuid.Nil, err
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.Parse(subject)
}
