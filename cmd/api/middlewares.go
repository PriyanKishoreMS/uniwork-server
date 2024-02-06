package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pascaldekloe/jwt"
	"github.com/priyankishorems/uniwork-server/internal/data"
	"golang.org/x/time/rate"
)

var (
	ErrUserUnauthorized = errors.New("middleware error: user unauthorized")
)

func (app *application) authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
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

			if !claims.Valid(time.Now()) {
				app.CustomErrorResponse(c, envelope{"token expired": "Send refresh token"}, http.StatusUnauthorized, ErrUserUnauthorized)
				return ErrUserUnauthorized
			}

			if claims.Issuer != app.config.jwt.issuer {
				app.UserUnAuthorizedResponse(c)
				return ErrUserUnauthorized
			}

			userID, err := strconv.ParseInt(claims.Subject, 10, 64)
			if err != nil {
				app.InternalServerError(c, err)
				return err
			}

			user, err := app.models.Users.Get(userID)
			if err != nil {
				switch {
				case errors.Is(err, data.ErrRecordNotFound):
					app.UserUnAuthorizedResponse(c)
				default:
					app.InternalServerError(c, err)
				}
				return err
			}

			r := app.contextSetUser(c, user)

			return next(r)
		}
	}
}

func (app *application) rateLimit() echo.MiddlewareFunc {
	limiter := rate.NewLimiter(20, 5)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if !limiter.Allow() {
				app.RateLimitExceededResponse(c)
				return nil
			}
			return next(c)
		}
	}
}
