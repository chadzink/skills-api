package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ResponseResult[T any] struct {
	Data    T      `json:"data"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ResponseResults[T any] struct {
	ResponseResult[T]
	Data []T `json:"data"`
}

type InvalidIdResult[T any] struct {
	ResponseResult[T]
	Success   bool        `json:"success" example:"false"`
	InvalidId interface{} `json:"invalid_id"`
}

type ErrorResult[T any] struct {
	ResponseResult[T]
	Success bool `json:"success" example:"false"`
}

type DeletResponse[T any] struct {
	ResponseResult[T]
	Id interface{} `json:"id"`
}

func GetValidId(c *fiber.Ctx) (uint, error) {
	value64, err := strconv.ParseUint(c.Params("id"), 10, 0)
	var value = uint(value64)

	if err != nil {
		return 0, err
	}

	if value <= 0 {
		return 0, errors.New("invalid id value of zero or less")
	}

	return value, nil
}

func HandleInvalidId(c *fiber.Ctx, err error) error {
	response := InvalidIdResult[interface{}]{
		ResponseResult: ResponseResult[interface{}]{
			Success: false,
			Message: err.Error(),
		},
		InvalidId: c.Params("id"),
	}

	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

func ErorrAndDataResponse(c *fiber.Ctx, err error, data interface{}) error {
	response := ErrorResult[interface{}]{
		ResponseResult: ResponseResult[interface{}]{
			Data:    data,
			Success: false,
			Message: err.Error(),
		},
	}

	return c.Status(fiber.StatusInternalServerError).JSON(response)
}

func DataResponse(c *fiber.Ctx, data interface{}) error {
	response := ResponseResult[interface{}]{
		Data:    data,
		Success: true,
		Message: "Operation was successful",
	}

	return c.Status(200).JSON(response)
}

func DeletedResponse(c *fiber.Ctx, message string, id interface{}) error {
	response := DeletResponse[interface{}]{
		ResponseResult: ResponseResult[interface{}]{
			Success: true,
			Message: message,
		},
		Id: id,
	}
	return c.Status(200).JSON(response)
}
