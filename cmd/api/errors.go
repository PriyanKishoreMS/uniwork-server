package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func (app *application) logError(c echo.Context, err error) {
	log.Error(err)
}

func (app *application) resposeError(c echo.Context, status int, message interface{}) {
	err := c.JSON(status, envelope{"error": message})
	if err != nil {
		app.logError(c, err)
		c.Response().WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) InternalServerError(c echo.Context, err error) {
	app.logError(c, err)
	message := "server encountered an error and could not process your request"
	app.resposeError(c, http.StatusInternalServerError, message)
}

func (app *application) BadRequest(c echo.Context, err error) {
	app.logError(c, err)
	message := "you have given a bad or invalid request, please try again"
	app.resposeError(c, http.StatusBadRequest, message)
}

func (app *application) MethodNotFound(c echo.Context) {
	message := "the method is not allowed"
	app.resposeError(c, http.StatusMethodNotAllowed, message)
}

func (app *application) NotFoundResponse(c echo.Context) {
	message := "the request is no found"
	app.resposeError(c, http.StatusNotFound, message)
}

func (app *application) editConflictResponse(c echo.Context) {
	message := "unable to update the record due to edit conflict, please try again"
	app.resposeError(c, http.StatusConflict, message)
}
