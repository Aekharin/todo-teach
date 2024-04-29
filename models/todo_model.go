package models

type CreateTodoRequest struct {
	TodoName string `json:"todo_name" validate:"required"`
	IsCheck  bool   `json:"is_check"`
}
