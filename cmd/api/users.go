package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/priyankishorems/uniwork-server/internal/data"
)

func (app *application) registerUserHandler(c echo.Context) error {
	user := new(data.User)

	if err := app.readJSON(c, &user); err != nil {
		app.BadRequest(c, err)
		return err
	}

	// TODO: Implement Google Siginin and return the user details

	// TODO: Implement regex for each college emails to verify

	err := app.validate.Struct(user)
	if err != nil {
		app.ValidationError(c, err)
		return err
	}

	err = app.models.Users.Register(user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				{
					app.CustomErrorResponse(c, envelope{"message": "duplicate entry"}, http.StatusConflict, err)
					return err
				}
			case "23503":
				{
					app.CustomErrorResponse(c, envelope{"message": "college not found"}, http.StatusNotFound, err)
					return err
				}
			default:
				{
					app.InternalServerError(c, err)
					return err
				}
			}
		}
	}

	accessToken, RefreshToken, err := data.GenerateAuthTokens(user.ID, app.config.jwt.secret, app.config.jwt.issuer)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	authTokens := envelope{
		"accessToken":  string(accessToken),
		"refreshToken": string(RefreshToken),
	}

	return c.JSON(http.StatusOK, envelope{"authTokens": authTokens, "data": user})
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

func (app *application) getRequestedUserHandler(c echo.Context) error {
	user := app.contextGetUser(c)
	if user == nil {
		app.NotFoundResponse(c)
		return errors.New("not found")
	}

	return c.JSON(http.StatusOK, envelope{"data": user})
}

func (app *application) updateUserHandler(c echo.Context) error {
	user := app.contextGetUser(c)

	var input struct {
		CollegeID *int64  `json:"college_id"`
		Name      *string `json:"name" validate:"required"`
		Email     *string `json:"email" validate:"required,email"`
		Mobile    *string `json:"mobile"`
		Avatar    *string `json:"avatar"`
		Dept      *string `json:"dept" validate:"required"`
	}

	err := app.readJSON(c, &input)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	updateField(&user.CollegeID, input.CollegeID)
	updateField(&user.Name, input.Name)
	updateField(&user.Email, input.Email)
	updateField(&user.Mobile, input.Mobile)
	updateField(&user.Avatar, input.Avatar)
	updateField(&user.Dept, input.Dept)

	err = app.validate.Struct(user)
	if err != nil {
		app.ValidationError(c, err)
		return err
	}

	err = app.models.Users.Update(user)
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
		"message": fmt.Sprintf("row updated successfully with id: %d", user.ID),
	})
}

func (app *application) deleteUserHandler(c echo.Context) error {
	user := app.contextGetUser(c)

	err := app.models.Users.Delete(user.ID)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": fmt.Sprintf("row deleted successfully with id: %d", user.ID),
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

	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name", "tasks_completed", "-tasks_completed", "rating", "-rating", "earned", "-earned"}

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
