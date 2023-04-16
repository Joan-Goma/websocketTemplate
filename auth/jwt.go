package auth

import (
	"time"

	engine "github.com/JoanGTSQ/api"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"websocketTemplate/models"
)

var jwtKey = []byte("kedg5b0ays1ekngsg18ruawcekgvcnz6")

type Context struct {
	User models.User `json:"user"`
}

type JWTClaim struct {
	Context Context `json:"context"`
	jwt.StandardClaims
}

func GenerateJWT(user models.User) (tokenString string, err error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	claims := &JWTClaim{
		Context: Context{
			User: user,
		},
		StandardClaims: jwt.StandardClaims{
			Id:        tokenID.String(),
			Issuer:    "neftsec",
			Subject:   user.UserName,
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	tokenString, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = engine.ERR_JWT_CLAIMS_INVALID
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = engine.ERR_JWT_TOKEN_EXPIRED
		return
	}
	return
}

func ReturnClaims(signedToken string) (claim *JWTClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return nil, engine.ERR_JWT_CLAIMS_INVALID
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = engine.ERR_JWT_CLAIMS_INVALID
		return nil, err
	}
	return claims, nil
}
