package data

import (
	"strconv"
	"time"

	"github.com/pascaldekloe/jwt"
)

func GenerateAuthTokens(id int64, secret string, issuer string) ([]byte, []byte, error) {
	byteSecret := []byte(secret)
	accessToken, err := GenerateAccessToken(id, byteSecret, issuer)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := GenerateRefreshToken(id, byteSecret, issuer)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func GenerateAccessToken(id int64, secret []byte, issuer string) ([]byte, error) {
	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(id, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(time.Hour * 24))
	claims.Issuer = issuer
	claims.Set = map[string]interface{}{
		"type": "access",
	}

	accessToken, err := claims.HMACSign(jwt.HS256, secret)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func GenerateRefreshToken(id int64, secret []byte, issuer string) ([]byte, error) {
	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(id, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add((time.Hour * 24) * 90))
	claims.Issuer = issuer
	claims.Set = map[string]interface{}{
		"type": "refresh",
	}

	refreshToken, err := claims.HMACSign(jwt.HS256, secret)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}
