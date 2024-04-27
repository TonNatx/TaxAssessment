package calculator

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

var (
	initialPersonalDeduction = 60000.0
	initialKReceipt          = 50000.0
	donationMax              = 100000.0
)

func CalculateTax(totalIncome, wht float64, allowances []Allowance) (float64, []TaxLevel) {

	taxableIncome := totalIncome - initialPersonalDeduction

	kReceiptAmount := 0.0
	donationAmount := 0.0
	// check k-receipt and donation
	for _, allowance := range allowances {
		switch allowance.AllowanceType {
		case "k-receipt":
			if allowance.Amount < initialKReceipt {
				kReceiptAmount = allowance.Amount
			} else {
				kReceiptAmount = initialKReceipt
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
