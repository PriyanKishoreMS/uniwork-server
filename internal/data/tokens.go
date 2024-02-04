package data

import (
	"strconv"
	"time"

	"github.com/pascaldekloe/jwt"
)

func GenerateAuthTokens(id int64, secret string, issuer string) ([]byte, []byte, error) {
	byteSecret := []byte(secret)
	accessToken, err := generateAccessToken(id, byteSecret, issuer)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := generateRefreshToken(id, byteSecret, issuer)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func generateAccessToken(id int64, secret []byte, issuer string) ([]byte, error) {
	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(id, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(time.Hour * 12))
	claims.Issuer = issuer
	claims.Set = map[string]interface{}{
		"type": "access",
	}

	accessToken, err := claims.HMACSign(jwt.HS256, secret)
	if err != nil {
		return nil, err
	}

	return accessToken, err
}

func generateRefreshToken(id int64, secret []byte, issuer string) ([]byte, error) {
	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(id, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add((time.Hour * 24) * 30))
	claims.Issuer = issuer
	claims.Set = map[string]interface{}{
		"type": "refresh",
	}

	refreshToken, err := claims.HMACSign(jwt.HS256, secret)
	if err != nil {
		return nil, err
	}

	return refreshToken, err
}
