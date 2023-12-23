package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"cloud.google.com/go/firestore"
	"github.com/labstack/echo/v4"
)

type App struct {
	e  *echo.Echo
	db *firestore.Client
}

func NewApp(e *echo.Echo, firestore *firestore.Client) *App {
	return &App{
		e:  e,
		db: firestore,
	}
}

// Bind the handlers
func (app *App) Init() {
	app.e.GET("/", welcome)
	app.e.POST("/create", app.createShortURL)
	app.e.GET("/:short-url", app.redirectURL)
}

func (app *App) Start(port string, bgCtx context.Context) {
	ctx, cancel := context.WithCancel(bgCtx)

	// Graceful shutdown
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		cancel()
		if err := app.e.Shutdown(ctx); err != nil {
			app.e.Logger.Fatal(err)
		}
	}()

	if err := app.e.Start(":" + port); err != http.ErrServerClosed {
		app.e.Logger.Fatal(err)
	}

	fmt.Println("Server shutdown gracefully")
}
