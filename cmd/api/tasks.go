package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/priyankishorems/uniwork-server/internal/data"
)

var (
	ErrInvalidQuery = errors.New("invalid Query")
)

func (app *application) addNewTaskHandler(c echo.Context) error {
	input := new(data.Task)
	user := app.contextGetUser(c)

	err := app.readFormData(c, input)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	input.UserID = user.ID
	input.CollegeID = user.CollegeID
	input.Expiry = time.Now().Add(time.Hour * 24)

	err = app.validate.Struct(input)
	if err != nil {
		app.ValidationError(c, err)
		return err
	}

	if !slices.Contains(data.TaskCategories, input.Category) {
		app.CustomErrorResponse(c, envelope{"invalid": "Invalid category"}, http.StatusBadRequest, ErrInvalidQuery)
		return err
	}

	imageURLs, err := app.HandleFiles(c, "images", user.ID, user.CollegeID)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	fileURLs, err := app.HandleFiles(c, "files", user.ID, user.CollegeID)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	input.Images = imageURLs
	input.Files = fileURLs
	// app.awsS3.UploadFile(c.Request().Context(), "uniwork", file.Filename, os.NewFile(0, file.Filename))

	err = app.models.Tasks.Create(input)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	c.JSON(http.StatusCreated, envelope{"created": fmt.Sprintf("Row created with ID: %d", input.ID)})

	return nil
}

func (app *application) getTaskHandler(c echo.Context) error {
	id, err := app.readIntParam(c, "id")
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	data, err := app.models.Tasks.Get(id)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"data": data})
}

func (app *application) deleteTaskHandler(c echo.Context) error {
	id, err := app.readIntParam(c, "id")
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	requestedUserID := app.contextGetUser(c).ID
	task, err := app.models.Tasks.Get(id)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	taskOwnerID := task.UserID

	if requestedUserID != taskOwnerID {
		app.CustomErrorResponse(c, envelope{"unauthorized": "You are not authorized to delete this task"}, http.StatusUnauthorized, ErrUserUnauthorized)
		return ErrUserUnauthorized
	}

	err = app.models.Tasks.Delete(id)
	if err != nil {
		app.InternalServerError(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"deleted": fmt.Sprintf("Row deleted with ID: %d", id)})
}

func (app *application) listAllTasksHandler(c echo.Context) error {
	var input struct {
		Category string
		data.Filters
	}

	college_id := app.contextGetUser(c).CollegeID

	qs := c.Request().URL.Query()
	input.Category = app.readStringQuery(qs, "category", "")
	input.Filters.Page = app.readIntQuery(qs, "page", 1)
	input.Filters.PageSize = app.readIntQuery(qs, "page_size", 10)
	input.Filters.Sort = app.readStringQuery(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name", "rating", "-rating", "price", "-price", "created_at", "-created_at", "expiry", "-expiry"}

	if input.Category != "" {
		categoryList := strings.Split(input.Category, ",")

		for _, category := range categoryList {
			category = strings.Trim(category, `'`)
			if !slices.Contains(data.TaskCategories, category) {
				app.BadRequest(c, ErrInvalidQuery)
				return ErrInvalidQuery
			}
		}
	}

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

	res, metadata, err := app.models.Tasks.GetAllTasksInCollege(input.Category, int64(college_id), input.Filters)
	if err != nil {
		app.BadRequest(c, err)
		return err
	}

	return c.JSON(http.StatusOK, envelope{"metadata": metadata, "data": res})
}
