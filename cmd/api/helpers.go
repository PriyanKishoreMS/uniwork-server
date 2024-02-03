package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type envelope map[string]interface{}

// func (app *application) background(fn func()) {
// 	app.wg.Add(1)

// 	go func() {

// 		defer app.wg.Done()

// 		defer func() {
// 			if err := recover(); err != nil {
// 				log.Error(err)
// 			}
// 		}()

// 		fn()
// 	}()
// }

func convertToInt64(str string) (int64, error) {
	integer, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return integer, nil
}

func (app *application) readIntParam(c echo.Context, str string) (int64, error) {
	param := c.Param("id")
	id, err := convertToInt64(param)
	if err != nil || id < 1 {
		return 0, errors.New("invalid parameter")
	}

	return id, err
}

func (app *application) readJSON(c echo.Context, dst interface{}) error {
	maxBytes := 1_048_576
	c.Request().Body = http.MaxBytesReader(c.Response(), c.Request().Body, int64(maxBytes))

	dec := json.NewDecoder(c.Request().Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)
		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

func (app *application) readStringQuery(qs url.Values, key string, defaultValue string) string {

	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readIntQuery(qs url.Values, key string, defaultValue int) int {

	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	res, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return res
}

func updateField[T any](user *T, input *T) {
	if input != nil {
		*user = *input
	}
}
