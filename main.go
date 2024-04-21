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

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxResponse struct {
	Tax       float64    `json:"tax"`
	TaxLevels []TaxLevel `json:"taxlevel"`
}

type Err struct {
	Message string `json:"message"`
}

var (
	personalDeduction = 60000.0
	kReceiptMax       = 50000.0
	donationMax       = 100000.0
)

func calculateTax(totalIncome, wht float64, allowances []Allowance) (float64, []TaxLevel) {

	taxableIncome := totalIncome - personalDeduction

	kReceiptAmount := 0.0
	donationAmount := 0.0
	// check k-receipt and donation
	for _, allowance := range allowances {
		switch allowance.AllowanceType {
		case "k-receipt":
			if allowance.Amount < kReceiptMax {
				kReceiptAmount = allowance.Amount
			} else {
				kReceiptAmount = kReceiptMax
			}

		case "donation":
			if allowance.Amount < donationMax {
				donationAmount = allowance.Amount
			} else {
				donationAmount = donationMax
			}
		}
	}

	// deduct k-receipt and donation
	taxableIncome -= kReceiptAmount
	taxableIncome -= donationAmount

	taxLevels := []TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 0.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}

	tax := 0.0
	switch {
	case taxableIncome <= 150000:
		tax = 0.0
	case taxableIncome > 150000 && taxableIncome <= 500000:
		tax = (taxableIncome - 150000) * 0.1
		taxLevels[1].Tax = tax
	case taxableIncome > 500000 && taxableIncome <= 1000000:
		taxLevel_2 := (500000 - 150000) * 0.1
		taxLevels[1].Tax = taxLevel_2
		taxLevel_3 := (taxableIncome - 500000) * 0.15
		taxLevels[2].Tax = taxLevel_3
		tax = taxLevel_2 + taxLevel_3
	case taxableIncome > 1000000 && taxableIncome <= 2000000:
		taxLevel_2 := (500000 - 150000) * 0.1
		taxLevels[1].Tax = taxLevel_2
		taxLevel_3 := (taxableIncome - 500000) * 0.15
		taxLevels[2].Tax = taxLevel_3
		taxLevel_4 := (taxableIncome - 1000000) * 0.2
		taxLevels[3].Tax = taxLevel_4
		tax = taxLevel_2 + taxLevel_3 + taxLevel_4
	default:
		taxLevel_2 := (500000 - 150000) * 0.1
		taxLevels[1].Tax = taxLevel_2
		taxLevel_3 := (taxableIncome - 500000) * 0.15
		taxLevels[2].Tax = taxLevel_3
		taxLevel_4 := (taxableIncome - 1000000) * 0.2
		taxLevels[3].Tax = taxLevel_4
		taxLevel_5 := (taxableIncome - 2000000) * 0.35
		taxLevels[4].Tax = taxLevel_5
		tax = taxLevel_2 + taxLevel_3 + taxLevel_4 + taxLevel_5
	}

	tax -= wht

	return tax, taxLevels
}

func calculateTaxHandler(c echo.Context) error {
	var t TaxRequest
	err := c.Bind(&t)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	tax, taxLevels := calculateTax(t.TotalIncome, t.WHT, t.Allowances)
	res := TaxResponse{Tax: tax, TaxLevels: taxLevels}

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
