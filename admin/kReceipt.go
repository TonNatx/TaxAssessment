package admin

import (
	"github.com/TonRat/assessment-tax/calculator"
	"github.com/labstack/echo/v4"
	"net/http"
)

type KReceiptRequest struct {
	Amount float64 `json:"amount"`
}

type KReceiptResepond struct {
	KReceipt float64 `json:"kReceipt"`
}

func KReceiptHandler(c echo.Context) error {
	var req KReceiptRequest
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if req.Amount < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "kReceipt amount must be greater than or equal to 0"})
	}

	calculator.InitialKReceipt = req.Amount

	res := KReceiptResepond{KReceipt: req.Amount}

	return c.JSON(http.StatusOK, res)
}
