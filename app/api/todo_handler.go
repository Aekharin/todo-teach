package api

import (
	"log"
	"net/http"
	"todo/models"

	"todo/app/database"
	"todo/constants"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TodoHandler struct {
	todoRepository database.TodoRepository
} //ตัวเชื่อมไว้เรียกใช้

func NewTodoHandler(todoRepository database.TodoRepository) *TodoHandler {
	return &TodoHandler{
		todoRepository: todoRepository,
	}
} //

func (h *TodoHandler) CreateTodo(ctx *fiber.Ctx) error {
	request := new(models.CreateTodoRequest)
	if err := ctx.BodyParser(&request); err != nil {
		return models.Response(constants.StatusCodeBadRequest, constants.BadRequestMessage, err.Error()).SendResponse(ctx, http.StatusBadRequest)
	}
	//สร้าง modle มารับตัวแปร opject .หาค่าได้
	//แปร json opject

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return models.Response(constants.StatusCodeBadRequest, nil, constants.BadRequestMessage).SendResponse(ctx, http.StatusBadRequest)
	}
	// เช็คค่าที่ต้องการว่ามีค่าไหม

	if err := h.todoRepository.CreateTodo(ctx.Context(), request); err != nil {
		return models.Response(constants.StatusCodeSystemError, nil, constants.StatusCodeSystemErrorMessage).SendResponse(ctx, http.StatusInternalServerError)
	}

	//เรียกใช้ ฟังก์ชั่น createtodo
	return models.ResponseSuccess(constants.StatusCodeSuccess, constants.SuccessMessage, nil).SendResponseSuccess(ctx, http.StatusOK)

}

func (h *TodoHandler) ReadTodo(ctx *fiber.Ctx) error {

	data, err := h.todoRepository.ReadTodo(ctx.Context())
	if err != nil {
		log.Print(err.Error())
		return models.Response(constants.StatusCodeBadRequest, nil, constants.BadRequestMessage).SendResponse(ctx, http.StatusBadRequest)
	}

	return models.ResponseSuccess(constants.StatusCodeSuccess, constants.SuccessMessage, data).SendResponseSuccess(ctx, http.StatusOK)

}
