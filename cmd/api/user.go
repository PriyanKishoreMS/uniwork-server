package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/priyankishorems/uniwork-server/internal/data"
)

func (app *application) registerUserHandler(c echo.Context) error {
	user := new(data.User)

	if err := app.readJSON(c, &user); err != nil {
		app.BadRequest(c, err)
		return err
	}

	err := app.validate.Struct(user)
	if err != nil {
		app.ValidationError(c, err)
		return err
	}

	res, err := app.models.Users.Register(user)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	data, err := app.models.Users.Get(res)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"data": data})
}

func (app *application) getUserHandler(c echo.Context) error {
	id, err := app.readIntParam(c, "id")
	if err != nil {
		app.NotFoundResponse(c)
		return err
	}

	res, err := app.models.Users.Get(id)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"data": res})

}

func (app *application) updateUserHandler(c echo.Context) error {
	id, err := app.readIntParam(c, "id")
	if err != nil {
		app.NotFoundResponse(c)
		return err
	}

	user, err := app.models.Users.Get(id)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	var input struct {
		CollegeID  *int64   `json:"college_id"`
		Name       *string  `json:"name" validate:"required"`
		Email      *string  `json:"email" validate:"required,email"`
		Mobile     *string  `json:"mobile"`
		ProfilePic **string `json:"profile_pic"`
		Dept       *string  `json:"dept" validate:"required"`
	}

	err = app.readJSON(c, &input)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	updateField(&user.CollegeID, input.CollegeID)
	updateField(&user.Name, input.Name)
	updateField(&user.Email, input.Email)
	updateField(&user.Mobile, input.Mobile)
	updateField(&user.ProfilePic, input.ProfilePic)
	updateField(&user.Dept, input.Dept)

	err = app.validate.Struct(user)
	if err != nil {
		app.ValidationError(c, err)
		return err
	}

	res, err := app.models.Users.Update(user)
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

func (app *application) deleteUserHandler(c echo.Context) error {
	id, err := app.readIntParam(c, "id")
	if err != nil {
		app.NotFoundResponse(c)
		return err
	}

	res, err := app.models.Users.Delete(id)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprint(res, " row deleted successfully"),
	})
}

func (app *application) listAllUsersInCollegeHandler(c echo.Context) error {
	var input struct {
		Name string
		data.Filters
	}

	college_id, err := app.readIntParam(c, "id")
	if err != nil {
		app.NotFoundResponse(c)
		return err
	}

	qs := c.Request().URL.Query()
	input.Name = app.readStringQuery(qs, "name", "")
	input.Filters.Page = app.readIntQuery(qs, "page", 1)
	input.Filters.PageSize = app.readIntQuery(qs, "page_size", 10)
	input.Filters.Sort = app.readStringQuery(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}

	err = app.validate.Struct(input)
	if err != nil {
		app.ValidationError(c, err)
		return err
	}

	if !slices.Contains(input.Filters.SortSafelist, input.Filters.Sort) {
		err := errors.New("unsafe query parameter")
		app.BadRequest(c, err)
		return err
	}

	res, metadata, err := app.models.Users.GetAllInCollege(input.Name, int64(college_id), input.Filters)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"metadata": metadata, "data": res})
}
