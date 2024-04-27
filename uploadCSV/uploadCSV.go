package uploadcsv

import (
	"encoding/csv"
	"github.com/TonRat/assessment-tax/calculator"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
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

type Err struct {
	Message string `json:"message"`
}

func UploadCSVHandler(c echo.Context) error {
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
