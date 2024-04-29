package models

import (
	"github.com/gofiber/fiber/v2"
)

type response struct {
	Code  string      `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
type responseSuccess struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type responseSuccessPage struct {
	Code       string      `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
}
type errorDes struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func Response(code string, data interface{}, errorMessage string) *response {
	if errorMessage != "" {
		return &response{
			Code:  code,
			Error: errorMessage}
	}
	return &response{
		Code: code,
		Data: data,
	}
}

func ResponseSuccess(code string, message string, data interface{}) *responseSuccess {
	return &responseSuccess{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func ResponseSuccessPage(code string, message string, data interface{}, page interface{}) *responseSuccessPage {
	return &responseSuccessPage{
		Code:       code,
		Message:    message,
		Data:       data,
		Pagination: page,
	}
}
func ResponseError(code string, message string, errorMessage string) *errorDes {
	return &errorDes{
		Code:    code,
		Message: message,
		Error:   errorMessage,
	}
}
func (r *response) SendResponse(ctx *fiber.Ctx, httpStatus int) error {
	return ctx.Status(httpStatus).JSON(r)
}
func (r *responseSuccess) SendResponseSuccess(ctx *fiber.Ctx, httpStatus int) error {
	return ctx.Status(httpStatus).JSON(r)
}
func (r *responseSuccessPage) SendResponseSuccessPage(ctx *fiber.Ctx, httpStatus int) error {
	return ctx.Status(httpStatus).JSON(r)
}
func (r *errorDes) SendResponseError(ctx *fiber.Ctx, httpStatus int) error {
	return ctx.Status(httpStatus).JSON(r)
}

type ApplicationError struct {
	Code string
	Desc string
}

func (a ApplicationError) Error() string {
	return a.Desc
}
