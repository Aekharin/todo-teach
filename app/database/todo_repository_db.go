package database

import (
	"context"

	"todo/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, createTodoRequest *models.CreateTodoRequest) error
	ReadTodo(ctx context.Context) (*[]models.ResponseReadtodo, error)
}

type TodoRepositoryDB struct {
	pool *pgxpool.Pool
}

func NewTodoRepositoryDB(pool *pgxpool.Pool) TodoRepository {
	return &TodoRepositoryDB{
		pool: pool,
	}
}

func (r *TodoRepositoryDB) CreateTodo(ctx context.Context, createTodoRequest *models.CreateTodoRequest) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		switch err {
		case nil:
			err = tx.Commit(ctx)
		default:
			_ = tx.Rollback(ctx)
		}
	}()

	stmt := `INSERT INTO todo_list (activity,status)
        VALUES(@todo_name, @is_check);`
	args := pgx.NamedArgs{
		"todo_name": createTodoRequest.TodoName,
		"is_check":  createTodoRequest.IsCheck,
	}

	_, err = tx.Exec(ctx, stmt, args)
	if err != nil {
		return err
	}

	return err
}

// ฟังก์ชชั่น read
func (r *TodoRepositoryDB) ReadTodo(ctx context.Context) (*[]models.ResponseReadtodo, error) {

	queryy := "SELECT tl.id, tl.status, tl.activity FROM todo_list tl;"

	rows, err := r.pool.Query(ctx, queryy)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ResponseReadtodolist []models.ResponseReadtodo
	for rows.Next() {
		var ResponseReadtodo models.ResponseReadtodo
		err := rows.Scan(
			&ResponseReadtodo.Id,
			&ResponseReadtodo.IsCheck,
			&ResponseReadtodo.TodoName,
			//สนใจตำแหน่ง
		)
		if err != nil {
			return nil, err
		}
		ResponseReadtodolist = append(ResponseReadtodolist, ResponseReadtodo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(ResponseReadtodolist) == 0 {
		return &[]models.ResponseReadtodo{}, nil
	}

	return &ResponseReadtodolist, nil
}
