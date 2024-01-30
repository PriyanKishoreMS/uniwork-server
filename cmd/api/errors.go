package main

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
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
	message := "the request is not found"
	app.resposeError(c, http.StatusNotFound, message)
}

func (app *application) EditConflictResponse(c echo.Context) {
	message := "unable to update the record due to edit conflict, please try again"
	app.resposeError(c, http.StatusConflict, message)
}

func (app *application) ValidationError(c echo.Context, err error) {
	validationError := make(map[string]interface{})
	validErrs := err.(validator.ValidationErrors)
	for _, e := range validErrs {
		var errMsg string

		switch e.Tag() {
		case "required":
			errMsg = "is required"
		case "email":
			errMsg = fmt.Sprint(e.Field(), " must be a type of email")
		case "gte":
			errMsg = "value must be greater than 0"
		case "lte":
			errMsg = "value must be lesser than the given value"

		default:
			errMsg = fmt.Sprintf("Validation error on %s: %s", e.Field(), e.Tag())
		}

		validationError[e.Field()] = errMsg
	}
	app.resposeError(c, http.StatusUnprocessableEntity, validationError)
}
