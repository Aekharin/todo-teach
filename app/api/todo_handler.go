package api

import (
	"net/http"
	"todo/models"

	"todo/app/database"
	"todo/constants"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TodoHandler struct {
	todoRepository database.TodoRepository
}

func NewTodoHandler(todoRepository database.TodoRepository) *TodoHandler {
	return &TodoHandler{
		todoRepository: todoRepository,
	}
}

func (h *TodoHandler) CreateTodo(ctx *fiber.Ctx) error {
	request := new(models.CreateTodoRequest)
	if err := ctx.BodyParser(&request); err != nil {
		return models.Response(constants.StatusCodeBadRequest, constants.BadRequestMessage, err.Error()).SendResponse(ctx, http.StatusBadRequest)
	}
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return models.Response(constants.StatusCodeBadRequest, nil, constants.BadRequestMessage).SendResponse(ctx, http.StatusBadRequest)
	}

	if err := h.todoRepository.CreateTodo(ctx.Context(), request); err != nil {
		return models.Response(constants.StatusCodeSystemError, nil, constants.StatusCodeSystemErrorMessage).SendResponse(ctx, http.StatusInternalServerError)
	}
	return models.ResponseSuccess(constants.StatusCodeSuccess, constants.SuccessMessage, nil).SendResponseSuccess(ctx, http.StatusOK)
}
