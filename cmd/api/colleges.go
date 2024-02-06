package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/priyankishorems/uniwork-server/internal/data"
)

func (app *application) createCollegeHandler(c echo.Context) error {
	clg := new(data.College)

	if err := app.readJSON(c, &clg); err != nil {
		app.BadRequest(c, err)
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
	id, err := app.readIntParam(c, "id")
	if err != nil {
		app.NotFoundResponse(c)
		return err
	}

	res, err := app.models.Colleges.Get(id)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"data": res})

}

func (app *application) updateCollegeHandler(c echo.Context) error {
	id, err := app.readIntParam(c, "id")
	if err != nil {
		app.NotFoundResponse(c)
		return err
	}

	college, err := app.models.Colleges.Get(id)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	var input struct {
		Name   *string `json:"name"`
		Domain *string `json:"domain" validate:"email"`
	}

	err = app.readJSON(c, &input)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	updateField(&college.Name, input.Name)
	updateField(&college.Domain, input.Domain)

	err = app.validate.Struct(college)
	if err != nil {
		app.ValidationError(c, err)
		return err
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
	id, err := app.readIntParam(c, "id")
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

func (app *application) listAllCollegesHandler(c echo.Context) error {
	var input struct {
		Name string
		data.Filters
	}

	user := app.contextGetUser(c)
	log.Info("user: ", user)

	qs := c.Request().URL.Query()
	input.Name = app.readStringQuery(qs, "name", "")
	input.Filters.Page = app.readIntQuery(qs, "page", 1)
	input.Filters.PageSize = app.readIntQuery(qs, "page_size", 10)
	input.Filters.Sort = app.readStringQuery(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}

	err := app.validate.Struct(input)
	if err != nil {
		app.ValidationError(c, err)
		return err
	}

	if !slices.Contains(input.Filters.SortSafelist, input.Filters.Sort) {
		err := errors.New("unsafe query parameter")
		app.BadRequest(c, err)
		return err
	}

	res, metadata, err := app.models.Colleges.GetAll(input.Name, input.Filters)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"metadata": metadata, "data": res})
}
