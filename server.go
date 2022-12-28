package main

import (
	"fmt"
	"os"

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

	e.POST("/expenses", expense.CreateExpenseHadler)
	e.Start(os.Getenv("PORT"))

	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))
	// log.Fatal(e.Start(":2565"))
}
