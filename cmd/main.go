package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/OrIX219/todo/pkg"
	"github.com/OrIX219/todo/pkg/config"
	"github.com/OrIX219/todo/pkg/handler"
	"github.com/OrIX219/todo/pkg/repository"
	"github.com/OrIX219/todo/pkg/service"
)

func main() {
	db, err := repository.NewPostgres(repository.Config{
		Host:     config.Config["POSTGRES_HOST"],
		Port:     config.Config["POSTGRES_PORT"],
		Username: config.Config["POSTGRES_USER"],
		Password: config.Config["POSTGRES_PASSWORD"],
		DBName:   config.Config["POSTGRES_DB"],
		SSLMode:  config.Config["POSTGRES_SSLMODE"],
	})
	if err != nil {
		log.Fatalf("config.Config to init DB: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	server := new(todo.Server)
	go func() {
		err := server.Run(config.Config["PORT"], handlers.InitRoutes())
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %s", err.Error())
		}
	}()
	log.Println("Server is up and running")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Println("Server is shutting down")

	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Error: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Printf("Error: %s", err.Error())
	}
}
