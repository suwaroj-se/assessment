package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/suwaroj-se/assessment/expense"
)

func main() {
	expense.InitDB()
	defer expense.CloseDB()
	expense := expense.Connection(expense.GetDB())

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/expenses", expense.CreateExpenseHandler)
	e.GET("/expenses", expense.GetAllExpenseHandler)
	e.GET("/expenses/:id", expense.GetExpenseHandlerByID)
	e.PUT("/expenses/:id", expense.PutExpenseHandlerByID)

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	go func() {
		if err := e.Start(os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// e.Start(os.Getenv("PORT"))
	// log.Fatal(e.Start(":2565"))
}
