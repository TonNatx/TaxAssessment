package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxResponse struct {
	Tax float64 `json:"tax"`
}

type Err struct {
	Message string `json:"message"`
}

var (
	personalDeduction = 60000.0
)

func calculateTax(totalIncome, wht float64, allowances []Allowance) float64 {

	taxableIncome := totalIncome - personalDeduction

	tax := 0.0
	switch {
	case taxableIncome <= 150000:
		tax = 0.0
	case taxableIncome > 150000 && taxableIncome <= 500000:
		tax = (taxableIncome - 150000) * 0.1
	case taxableIncome > 500000 && taxableIncome <= 1000000:
		tax = (500000-150000)*0.1 + (taxableIncome-500000)*0.15
	case taxableIncome > 1000000 && taxableIncome <= 2000000:
		tax = (500000-150000)*0.1 + (1000000-500000)*0.15 + (taxableIncome-1000000)*0.2
	default:
		tax = (500000-150000)*0.1 + (1000000-500000)*0.15 + (2000000-1000000)*0.2 + (taxableIncome-2000000)*0.35
	}

	return tax
}

func calculateTaxHandler(c echo.Context) error {
	var t TaxRequest
	err := c.Bind(&t)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	tax := calculateTax(t.TotalIncome, t.WHT, t.Allowances)
	res := TaxResponse{Tax: tax}

	return c.JSON(http.StatusOK, res)
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})
	e.POST("/tax/calculation", calculateTaxHandler)
	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
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
