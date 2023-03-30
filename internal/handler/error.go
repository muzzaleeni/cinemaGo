package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func newInvalidJSONErr(err error) *fiber.Error {
	msg := fmt.Sprintf("invalid JSON input: %v", err)
	return fiber.NewError(http.StatusBadRequest, msg)
}

func newQueryParamErr(err error) *fiber.Error {
	prefix := "invalid input query parameter"
	code := http.StatusBadRequest

	var merr fiber.MultiError
	if errors.As(err, &merr) {
		for _, err := range merr {
			var convErr fiber.ConversionError
			if errors.As(err, &convErr) {
				msg := fmt.Sprintf("%s: %q expected to be %v", prefix, convErr.Key, convErr.Type)
				return fiber.NewError(code, msg)
			}
		}
	}

	msg := fmt.Sprintf("%s: %v", prefix, err)
	return fiber.NewError(code, msg)
}

func newNotFoundErr(entityName string, id int) *fiber.Error {
	msg := fmt.Sprintf("%s with id=%d not found", entityName, id)
	return fiber.NewError(http.StatusNotFound, msg)
}

func newRouteParamErr(msg string) *fiber.Error {
	prefix := "invalid route parameter"
	errMsg := fmt.Sprintf("%s: %s", prefix, msg)
	return fiber.NewError(http.StatusBadRequest, errMsg)
}
