package handlers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ResponseResult struct {
	Data    map[string]interface{} `json:"data"`
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
}

type ResponseResults struct {
	ResponseResult
	Data []map[string]interface{} `json:"data"`
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
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message":    err.Error(),
		"success":    false,
		"invalid_id": c.Params("id"),
	})
}

func ErorrAndDataResponse(c *fiber.Ctx, err error, data interface{}) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": err.Error(),
		"success": false,
		"data":    data,
	})
}

func DataResponse(c *fiber.Ctx, data interface{}) error {
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func DeletedResponse(c *fiber.Ctx, message string, id interface{}) error {
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": message,
		"id":      id,
	})
}
