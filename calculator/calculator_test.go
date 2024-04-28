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

		assert.Equal(t, expectedTax, tax, "Tax should be 29000.0")
		assert.Equal(t, expectedTaxLevels, taxLevels, "Tax should be show in range 150,001-500,000")
		assert.Equal(t, nil, err, "Should not be error")
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

		assert.Equal(t, expectedTax, tax, "Tax should be 19000.0")
		assert.Equal(t, expectedTaxLevels, taxLevels, "Tax should be show in range 150,001-500,000")
		assert.Equal(t, nil, err, "Should not be error")
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

		assert.Equal(t, expectedTax, tax, "Tax should be 14000.0")
		assert.Equal(t, expectedTaxLevels, taxLevels, "Tax should be show in range 150,001-500,000")
		assert.Equal(t, nil, err, "Should not be error")
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

		assert.Equal(t, expectedTax, tax, "Tax should be 4000.0")
		assert.Equal(t, expectedTaxLevels, taxLevels, "Tax should be show in range 150,001-500,000")
		assert.Equal(t, nil, err, "Should not be error")
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

		assert.Equal(t, expectedTax, tax, "Tax should be 0.0")
		assert.Equal(t, expectedTaxLevels, taxLevels, "Tax should be show in range 150,001-500,000")
		assert.Equal(t, nil, err, "Should not be error")
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

		assert.Equal(t, expectedTax, tax, "Tax should be 0.0")
		assert.Equal(t, expectedTaxLevels, taxLevels, "Tax should be show in range 150,001-500,000")
		assert.Equal(t, nil, err, "Should not be error")
	})

}
