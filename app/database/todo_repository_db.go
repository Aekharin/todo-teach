package database

import (
	"context"

	"todo/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, createTodoRequest *models.CreateTodoRequest) error
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

	stmt := `INSERT INTO todo_list (todo_name,is_check)
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
