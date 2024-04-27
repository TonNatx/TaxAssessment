package main

import (
	"context"

	"fmt"

	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/TonRat/assessment-tax/admin"
	"github.com/TonRat/assessment-tax/calculator"
	"github.com/TonRat/assessment-tax/uploadCSV"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file couldn't be load")
	}

	e := echo.New()

	e.POST("/tax/calculations", calculator.CalculateTaxHandler)
	e.POST("/tax/calculations/upload-csv", uploadcsv.UploadCSVHandler)

	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {

		if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
			return true, nil
		}
		return false, nil
	}))
	g.POST("/deductions/personal", admin.PersonalDeductionHandler)
	g.POST("/deductions/k-receipt", admin.KReceiptHandler)
	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown
	fmt.Println("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
