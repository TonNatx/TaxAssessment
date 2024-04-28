package calculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateTaxNoWht(t *testing.T) {
	t.Run("NoAllowance", func(t *testing.T) {
		totalIncome := 500000.0
		wht := 0.0
		allowances := []Allowance{
			{AllowanceType: "donation", Amount: 0.0},
		}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 29000.0
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 29000.0},
			{Level: "500,001-1,000,000", Tax: 0.0},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Wrong tax level")
		assert.Nil(t, err, "Should not be error")
	})

	t.Run("OneAllowance", func(t *testing.T) {
		totalIncome := 500000.0
		wht := 0.0
		allowances := []Allowance{
			{AllowanceType: "donation", Amount: 200000.0},
		}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 19000.0
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 19000.0},
			{Level: "500,001-1,000,000", Tax: 0.0},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Wrong tax level")
		assert.Nil(t, err, "Should not be error")
	})

	t.Run("TwoAllowance", func(t *testing.T) {
		totalIncome := 500000.0
		wht := 0.0
		allowances := []Allowance{
			{AllowanceType: "donation", Amount: 200000.0},
			{AllowanceType: "k-receipt", Amount: 100000.0},
		}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 14000.0
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 14000.0},
			{Level: "500,001-1,000,000", Tax: 0.0},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Wrong tax level")
		assert.Nil(t, err, "Should not be error")
	})

}

func TestCalculateTaxWithWht(t *testing.T) {
	t.Run("NoAllowance", func(t *testing.T) {
		totalIncome := 500000.0
		wht := 25000.0
		allowances := []Allowance{
			{AllowanceType: "donation", Amount: 0.0},
		}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 4000.0
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 29000.0},
			{Level: "500,001-1,000,000", Tax: 0.0},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Wrong tax level")
		assert.Nil(t, err, "Should not be error")
	})

	t.Run("OneAllowance", func(t *testing.T) {
		totalIncome := 500000.0
		wht := 19000.0
		allowances := []Allowance{
			{AllowanceType: "donation", Amount: 200000.0},
		}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 0.0
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 19000.0},
			{Level: "500,001-1,000,000", Tax: 0.0},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Wrong tax level")
		assert.Nil(t, err, "Should not be error")
	})

	t.Run("TwoAllowance", func(t *testing.T) {
		totalIncome := 500000.0
		wht := 14000.0
		allowances := []Allowance{
			{AllowanceType: "donation", Amount: 200000.0},
			{AllowanceType: "k-receipt", Amount: 100000.0},
		}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 0.0
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 14000.0},
			{Level: "500,001-1,000,000", Tax: 0.0},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Wrong tax level")
		assert.Nil(t, err, "Should not be error")
	})

}
func TestCalculateTaxError(t *testing.T) {
	t.Run("WhtLessThanZero", func(t *testing.T) {
		totalIncome := 500000.0
		wht := -25000.0
		allowances := []Allowance{
			{AllowanceType: "donation", Amount: 0.0},
		}

		_, _, err := CalculateTax(totalIncome, wht, allowances)

		expectedErrorMsg := "wht must be between 0 and totalIncome"

		assert.Equal(t, expectedErrorMsg, err.Error(), "Should be error")
	})

	t.Run("WhtMoreThanTotalIncome", func(t *testing.T) {
		totalIncome := 500000.0
		wht := 500001.0
		allowances := []Allowance{
			{AllowanceType: "donation", Amount: 0.0},
		}

		_, _, err := CalculateTax(totalIncome, wht, allowances)

		expectedErrorMsg := "wht must be between 0 and totalIncome"

		assert.Equal(t, expectedErrorMsg, err.Error(), "Should be error")
	})

	t.Run("DonationLessThanZero", func(t *testing.T) {
		totalIncome := 500000.0
		wht := 0.0
		allowances := []Allowance{
			{AllowanceType: "donation", Amount: -100000.0},
		}

		_, _, err := CalculateTax(totalIncome, wht, allowances)

		expectedErrorMsg := "donation must be greater than 0"

		assert.Equal(t, expectedErrorMsg, err.Error(), "Should be error")
	})
	t.Run("KReceiptLessThanZero", func(t *testing.T) {
		totalIncome := 500000.0
		wht := 0.0
		allowances := []Allowance{
			{AllowanceType: "k-receipt", Amount: -100000.0},
		}

		_, _, err := CalculateTax(totalIncome, wht, allowances)

		expectedErrorMsg := "kReceiptAmount must be greater than 0"

		assert.Equal(t, expectedErrorMsg, err.Error(), "Should be error")
	})
}
func TestCalculateTotalIncome(t *testing.T) {
	t.Run("TotalIncomeIsZero", func(t *testing.T) {
		totalIncome := 0.0
		wht := 0.0
		allowances := []Allowance{
			{AllowanceType: "donation", Amount: 200000.0},
			{AllowanceType: "k-receipt", Amount: 100000.0},
		}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 0.0
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 0.0},
			{Level: "500,001-1,000,000", Tax: 0.0},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Tax level should be empty")
		assert.Nil(t, err, "Should not be error")
	})
}
func TestCalculateTaxLevel(t *testing.T) {
	//Assume that PersonalDeduction is 60,000 constance
	t.Run("TotalIncomeIs150000", func(t *testing.T) {
		totalIncome := 150000.0
		wht := 0.0
		allowances := []Allowance{}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 0.0
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 0.0},
			{Level: "500,001-1,000,000", Tax: 0.0},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Tax level should be empty")
		assert.Nil(t, err, "Should not be error")
	})
	t.Run("TotalIncomeIs150001", func(t *testing.T) {
		totalIncome := 210001.0
		wht := 0.0
		allowances := []Allowance{}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 0.1
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 0.1},
			{Level: "500,001-1,000,000", Tax: 0.0},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Wrong tax level")
		assert.Nil(t, err, "Should not be error")
	})
	t.Run("TotalIncomeIs500001", func(t *testing.T) {
		totalIncome := 560001.0
		wht := 0.0
		allowances := []Allowance{}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 35000.15
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 35000.0},
			{Level: "500,001-1,000,000", Tax: 0.15},
			{Level: "1,000,001-2,000,000", Tax: 0.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Wrong tax level")
		assert.Nil(t, err, "Should not be error")
	})

	t.Run("TotalIncomeIs1,000,001", func(t *testing.T) {
		totalIncome := 1060001.0
		wht := 0.0
		allowances := []Allowance{}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 110000.2
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 35000.0},
			{Level: "500,001-1,000,000", Tax: 75000.0},
			{Level: "1,000,001-2,000,000", Tax: 0.2},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Wrong tax level")
		assert.Nil(t, err, "Should not be error")
	})

	t.Run("TotalIncomeIs2,000,001", func(t *testing.T) {
		totalIncome := 2060001.0
		wht := 0.0
		allowances := []Allowance{}

		tax, taxLevels, err := CalculateTax(totalIncome, wht, allowances)

		expectedTax := 310000.35
		expectedTaxLevels := []TaxLevel{
			{Level: "0-150,000", Tax: 0.0},
			{Level: "150,001-500,000", Tax: 35000.0},
			{Level: "500,001-1,000,000", Tax: 75000.0},
			{Level: "1,000,001-2,000,000", Tax: 200000.0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0.35},
		}

		assert.Equal(t, expectedTax, tax, "Tax should be %.1f", expectedTax)
		assert.Equal(t, expectedTaxLevels, taxLevels, "Wrong tax level")
		assert.Nil(t, err, "Should not be error")
	})
}
