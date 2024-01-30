package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/priyankishorems/uniwork-server/internal/data"
	"github.com/priyankishorems/uniwork-server/internal/helpers"
)

type envelope map[string]interface{}

func (app *application) createCollegeHandler(c echo.Context) error {
	clg := new(data.College)

	if err := c.Bind(&clg); err != nil {
		app.InternalServerError(c, err)
		return err
	}

	err := app.validate.Struct(clg)
	if err != nil {
		app.ValidationError(c, err)
		return err
	}

	res, err := app.models.Colleges.Create(clg)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": fmt.Sprint("row created successfully with id: ", res),
	})
}

func (app *application) getCollegeHandler(c echo.Context) error {
	id, err := helpers.ParamToInt64(c, "id")
	if err != nil {
		app.NotFoundResponse(c)
		return err
	}

	res, err := app.models.Colleges.Get(id)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"data": res})

}

func (app *application) updateCollegeHandler(c echo.Context) error {
	id, err := helpers.ParamToInt64(c, "id")
	if err != nil {
		app.NotFoundResponse(c)
		return err
	}

	college, err := app.models.Colleges.Get(id)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	log.Info("college info", college)

	var input struct {
		Name   *string `json:"name"`
		Domain *string `json:"domain"`
	}

	if err := c.Bind(&input); err != nil {
		app.InternalServerError(c, err)
		return err
	}

	if input.Name != nil {
		college.Name = *input.Name
	}

	if input.Domain != nil {
		college.Domain = *input.Domain
	}

	res, err := app.models.Colleges.Update(college)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.EditConflictResponse(c)
		default:
			app.InternalServerError(c, err)
		}
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprint(res, " row updated successfully"),
	})
}

func (app *application) deleteCollegeHandler(c echo.Context) error {
	id, err := helpers.ParamToInt64(c, "id")
	if err != nil {
		app.NotFoundResponse(c)
		return err
	}

	res, err := app.models.Colleges.Delete(id)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprint(res, " row deleted successfully"),
	})
}
