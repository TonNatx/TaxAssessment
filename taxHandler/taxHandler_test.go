package taxHandler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/TonRat/assessment-tax/calculator"
	"github.com/labstack/echo/v4"
)

func TestCalculateTaxHandler(t *testing.T) {
	// สร้าง JSON request
	reqBody := map[string]interface{}{
		"totalIncome": 500000.0,
		"wht":         0.0,
		"allowances": []calculator.Allowance{
			{AllowanceType: "donation", Amount: 200000.0},
		},
	}
	reqJSON, _ := json.Marshal(reqBody)

	// เรียกฟังก์ชัน CalculateTaxHandler ด้วย JSON request
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec) // สร้าง context ขึ้นมา
	err := CalculateTaxHandler(c)
	assert.NoError(t, err)

	// ตรวจสอบ response
	assert.Equal(t, http.StatusOK, rec.Code)

	var res TaxResponse
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err)

	expectedRes := TaxResponse{
		Tax: 19000.0,
		TaxLevels: []calculator.TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 19000.0},
			{Level: "500,001-1,000,000", Tax: 0.0},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		},
	}
	assert.Equal(t, expectedRes, res)
}
