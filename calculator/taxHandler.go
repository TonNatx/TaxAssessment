package calculator

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxResponse struct {
	Tax       float64    `json:"tax"`
	TaxLevels []TaxLevel `json:"taxlevel"`
}
type TaxRefundRespond struct {
	TaxRefund float64    `json:"taxRefund"`
	TaxLevels []TaxLevel `json:"taxlevel"`
}

type Err struct {
	Message string `json:"message"`
}

func CalculateTaxHandler(c echo.Context) error {
	var t TaxRequest
	err := c.Bind(&t)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	tax, taxLevels := CalculateTax(t.TotalIncome, t.WHT, t.Allowances)
	if tax < 0 {
		res := TaxRefundRespond{TaxRefund: -tax, TaxLevels: taxLevels}
		return c.JSON(http.StatusOK, res)
	} else {
		res := TaxResponse{Tax: tax, TaxLevels: taxLevels}
		return c.JSON(http.StatusOK, res)
	}
}
