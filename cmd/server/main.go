package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/app/handler"
	"github.com/mytheresa/go-hiring-challenge/app/middleware"
	"github.com/mytheresa/go-hiring-challenge/app/repository"
	"github.com/mytheresa/go-hiring-challenge/app/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := godotenv.Load(".env"); err != nil {
		logger.Error("Error loading .env file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	prodRepo := repository.NewProductsRepository(db)
	productService := service.NewProductService(prodRepo)
	cat := handler.NewCatalogHandler(productService)

	catRepo := repository.NewCategoriesRepository(db)
	categoryService := service.NewCategoryService(catRepo)
	categories := handler.NewCategoriesHandler(categoryService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", cat.HandleGetCatalog)
	mux.HandleFunc("GET /catalog/{code}", cat.HandleGetProductDetails)
	mux.HandleFunc("GET /categories", categories.HandleGetCategories)
	mux.HandleFunc("POST /categories", categories.HandleCreateCategory)

	h := middleware.Chain(mux,
		middleware.Logging(logger),
		middleware.Recovery(logger),
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: h,
	}

	go func() {
		logger.Info("Starting server", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("Shutting down server...")
	srv.Shutdown(ctx)
	stop()
}
