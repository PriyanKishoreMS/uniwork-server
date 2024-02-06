package main

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/priyankishorems/uniwork-server/internal/data"
)

func (app *application) refreshTokenHandler(c echo.Context) error {

	c.Response().Writer.Header().Add("Vary", "Authorization")

	authorizationHeader := c.Request().Header.Get("Authorization")
	if authorizationHeader == "" {
		app.UserUnAuthorizedResponse(c)
		return ErrUserUnauthorized
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		app.UserUnAuthorizedResponse(c)
		return ErrUserUnauthorized
	}

	token := headerParts[1]

	claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
	if err != nil {
		app.UserUnAuthorizedResponse(c)
		return ErrUserUnauthorized
	}

	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	accessToken, err := data.GenerateAccessToken(id, []byte(app.config.jwt.secret), app.config.jwt.issuer)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}
	return c.JSON(200, envelope{"access_token": string(accessToken)})
}
