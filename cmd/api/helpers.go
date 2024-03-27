package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type envelope map[string]interface{}

var uploadDir string = "./public"

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
	param := c.Param(str)
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
func (app *application) readFormData(c echo.Context, dst interface{}) error {
	err := c.Request().ParseMultipartForm(10 << 20)
	if err != nil {
		return fmt.Errorf("failed to parse multipart form: %v", err)
	}

	dstValue := reflect.ValueOf(dst).Elem()
	dstType := dstValue.Type()

	for i := 0; i < dstValue.NumField(); i++ {
		field := dstType.Field(i)
		fieldValue := dstValue.Field(i)
		formValue := c.FormValue(strings.ToLower(field.Name))

		if fieldValue.CanSet() {
			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.SetString(formValue)
			case reflect.Int64:
				if formValue == "" {
					fieldValue.SetInt(0)
				} else {
					value, err := convertToInt64(formValue)
					if err != nil {
						return fmt.Errorf("invalid value for field %s: %v", field.Name, err)
					}
					fieldValue.SetInt(value)
				}
			}
		}
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

func (app *application) HandleFiles(c echo.Context, key string, userID int64, collegeID int64) ([]string, error) {
	files := c.Request().MultipartForm.File[key]
	if len(files) == 0 {
		return []string{}, nil
	}
	filePaths := []string{}
	uploadDir := uploadDir + "/" + key
	fmt.Println(uploadDir, "uploadDir")

	for _, fileHeader := range files {
		file, err := fileHeader.Open()

		if err != nil {
			return []string{}, err
		}
		defer file.Close()

		b := make([]byte, 4)
		rand.Read(b)
		suffix := hex.EncodeToString(b)
		filename := fmt.Sprintf("%d_%d_%s%s", userID, collegeID, suffix, filepath.Ext(fileHeader.Filename))

		dst, err := os.Create(filepath.Join(uploadDir, filename))
		if err != nil {
			return []string{}, err
		}
		defer dst.Close()
		filePaths = append(filePaths, uploadDir[1:]+"/"+filename)

		_, err = io.Copy(dst, file)
		if err != nil {
			return []string{}, err
		}
	}
	return filePaths, nil
}
