package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"todo/app/api"
	"todo/app/database"
	"todo/config"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {
	runtime.GOMAXPROCS(1)
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	cfg := config.InitialConfig()

	postgresClient, err := newPostgresClient(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresClient.Close()

	app := initFiber()

	todoHandler := api.NewTodoHandler(database.NewTodoRepositoryDB(postgresClient))

	app.Post("/create-todo", todoHandler.CreateTodo)
	healthCheck(app, postgresClient)

	log.Printf("Listening on port: %s", cfg.Server.Port)
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
			log.Fatal(err)
		}
	}()

	gracefulShutdown(app)
}

func healthCheck(app *fiber.App, pool *pgxpool.Pool) {
	app.Get("/health", func(c *fiber.Ctx) error {
		if err := pool.Ping(context.Background()); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fiber.ErrInternalServerError.Error())
		}
		return c.Status(fiber.StatusOK).SendString("Ready!!")
	})
}

func gracefulShutdown(f *fiber.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	if err := f.Shutdown(); err != nil {
		log.Fatal("shutdown server:", err)
	}
}

func initFiber() *fiber.App {
	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
			CaseSensitive:         true,
			StrictRouting:         true,
		},
	)
	return app
}

func newPostgresClient(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Database,
	)
	pgxpoolConfig, err := pgxpool.ParseConfig(psqlInfo)

	if err != nil {
		return nil, fmt.Errorf("postgres client parse config error: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxpoolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres client connect error: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres client Ping error: %v", err)
	}

	return pool, nil
}
