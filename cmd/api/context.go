package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/priyankishorems/uniwork-server/internal/data"
)

type ContextKey string

const userContextKey = ContextKey("user")

func (app *application) contextSetUser(c echo.Context, user *data.User) echo.Context {
	ctx := context.WithValue(c.Request().Context(), userContextKey, user)
	c.SetRequest(c.Request().WithContext(ctx))
	return c
}

func (app *application) contextGetUser(c echo.Context) *data.User {
	user, ok := c.Request().Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
