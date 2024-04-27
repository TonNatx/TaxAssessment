package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/TonRat/assessment-tax/calculator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TaxRecord struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax"`
}

type TaxRecordRefund struct {
	TotalIncome float64 `json:"totalIncome"`
	TaxRefund   float64 `json:"taxRefund"`
}

type TaxResponseCSV struct {
	Taxes []interface{} `json:"taxes"`
}

type DeductionRequest struct {
	Amount float64 `json:"amount"`
}

type DeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type KReceiptRequest struct {
	Amount float64 `json:"amount"`
}

type KReceiptResepond struct {
	KReceipt float64 `json:"kReceipt"`
}

type Err struct {
	Message string `json:"message"`
}

var (
	initialPersonalDeduction = 60000.0
	initialKReceipt          = 50000.0
	donationMax              = 100000.0
)

func uploadCSVHandler(c echo.Context) error {
	// Get uploaded file
	file, err := c.FormFile("taxFile")
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	defer src.Close()
	// Create a CSV reader
	reader := csv.NewReader(src)
	// Read the header row
	header, err := reader.Read()
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "Failed to read CSV header"})
	}

	// Check if the header matches the expected format
	expectedHeader := []string{"totalIncome", "wht", "donation"}
	for i, col := range header {
		if col != expectedHeader[i] {
			return c.JSON(http.StatusBadRequest, Err{Message: "Invalid CSV header format"})
		}
	}

	// Read and process CSV records
	var taxes []interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
		}

		// Convert CSV data to float64
		totalIncome, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: "Invalid totalIncome format"})
		}
		wht, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: "Invalid WHT format"})
		}
		donation, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, Err{Message: "Invalid donation format"})
		}

		// Perform tax calculation
		allowances := []calculator.Allowance{{AllowanceType: "donation", Amount: donation}}
		tax, _ := calculator.CalculateTax(totalIncome, wht, allowances)
		if tax < 0 {
			refundRecord := TaxRecordRefund{TotalIncome: totalIncome, TaxRefund: -tax}
			taxes = append(taxes, refundRecord)
		} else {
			normalRecord := TaxRecord{TotalIncome: totalIncome, Tax: tax}
			taxes = append(taxes, normalRecord)
		}
	}

	res := TaxResponseCSV{Taxes: taxes}

	return c.JSON(http.StatusOK, res)
}

func personalDeductionHandler(c echo.Context) error {
	var req DeductionRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if req.Amount < 10000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "personalDeduction amount must be greater than or equal to 10,000"})
	}

	initialPersonalDeduction = req.Amount

	res := DeductionResponse{PersonalDeduction: initialPersonalDeduction}

	return c.JSON(http.StatusOK, res)
}

func kReceiptHandler(c echo.Context) error {
	var req KReceiptRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if req.Amount < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "kReceipt amount must be greater than or equal to 0"})
	}

	initialKReceipt = req.Amount

	res := KReceiptResepond{KReceipt: initialKReceipt}

	return c.JSON(http.StatusOK, res)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file couldn't be load")
	}

	e := echo.New()

	e.POST("/tax/calculations", calculator.CalculateTaxHandler)
	e.POST("/tax/calculations/upload-csv", uploadCSVHandler)

	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {

		if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
			return true, nil
		}
		return false, nil
	}))
	g.POST("/deductions/personal", personalDeductionHandler)
	g.POST("/deductions/k-receipt", kReceiptHandler)
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
