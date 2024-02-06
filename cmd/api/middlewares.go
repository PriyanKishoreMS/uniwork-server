package main

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
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

	type client struct {
		limiter  *rate.Limiter
		lastseen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// background routine to remove old entries from the map
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastseen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if app.config.limiter.enabled {
				ip, _, err := net.SplitHostPort(c.Request().RemoteAddr)
				if err != nil {
					app.InternalServerError(c, err)
					return err
				}

				mu.Lock()

				_, found := clients[ip]
				if !found {
					clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
				}

				clients[ip].lastseen = time.Now()

				if !clients[ip].limiter.Allow() {
					mu.Unlock()
					app.RateLimitExceededResponse(c)
					return errors.New("rate limit exceeded")
				}

				mu.Unlock()
			}

			return next(c)
		}
	}
}
