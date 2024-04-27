package admin

import (
	"github.com/TonRat/assessment-tax/calculator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type DeductionRequest struct {
	Amount float64 `json:"amount"`
}

type DeductionResponse struct {
	PersonalDeduction float64 `json:"personalDeduction"`
}

type Err struct {
	Message string `json:"message"`
}

func PersonalDeductionHandler(c echo.Context) error {
	var req DeductionRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if req.Amount < 10000 {
		return c.JSON(http.StatusBadRequest, Err{Message: "personalDeduction amount must be greater than or equal to 10,000"})
	}

	calculator.InitialPersonalDeduction = req.Amount

	res := DeductionResponse{PersonalDeduction: req.Amount}

	return c.JSON(http.StatusOK, res)
}
