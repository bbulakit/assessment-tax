package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func TaxCalculationsHandler(c echo.Context) error {
	itd := IncomeTaxDetail{}
	err := c.Bind(&itd)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	tcr := taxCalculate(itd)

	return c.JSON(http.StatusOK, tcr)
}

func taxCalculate(itd IncomeTaxDetail) TaxCalculationResult {
	tcr := TaxCalculationResult{}
	taxCal := taxLevelDeduction(itd.TotalIncome)

	//TODO: Get default personal deduction
	personalDeduction := 60_000.0
	taxCal -= personalDeduction

	for _, deduction := range itd.Allowances {
		var actualDeduction float64
		if deduction.AllowanceType == "donation" {
			//TODO: Get max. donation deduction
			maxDonationDeduction := 100_000.0
			actualDeduction = deduction.Amount
			if actualDeduction > maxDonationDeduction {
				actualDeduction = maxDonationDeduction
			}
		}

		if deduction.AllowanceType == "k-receipt" {
			//TODO: Get max. k-receipt deduction
			maxKReceiptDeduction := 50_000.0 //Default @ 50_000
			actualDeduction = deduction.Amount
			if actualDeduction > maxKReceiptDeduction {
				actualDeduction = maxKReceiptDeduction
			}
		}
		taxCal -= actualDeduction
	}

	taxCal *= taxRate(itd.TotalIncome)
	taxCal -= itd.WithHoldingTax

	tcr.TotalTax = taxCal
	tcr.TaxLevels = taxLevelDetail(itd.TotalIncome, tcr.TotalTax)
	return tcr
}

func taxLevelDeduction(totalIncome float64) float64 {
	if totalIncome <= 150_000 {
		return 0.00
	} else if totalIncome <= 500_000 {
		return totalIncome - 150_000
	} else if totalIncome <= 1_000_000 {
		return totalIncome - 500_000
	} else if totalIncome <= 2_000_000 {
		return totalIncome - 1_000_000
	}
	return totalIncome - 2_000_000
}

func taxRate(totalIncome float64) float64 {
	if totalIncome <= 150_000 {
		return 0.00
	} else if totalIncome <= 500_000 {
		return 0.10
	} else if totalIncome <= 1_000_000 {
		return 0.15
	} else if totalIncome <= 2_000_000 {
		return 0.20
	}
	return 0.35
}

func taxLevelDetail(totalIncome float64, totalTax float64) []TaxLevel {
	taxLevels := initialTaxLevelDetail()

	if totalIncome <= 150_000 {
		taxLevels[0].Tax = totalTax
		return taxLevels
	}

	if totalIncome <= 500_000 {
		taxLevels[1].Tax = totalTax
		return taxLevels
	}

	if totalIncome <= 1_000_000 {
		taxLevels[2].Tax = totalTax
		return taxLevels
	}

	if totalIncome <= 2_000_000 {
		taxLevels[3].Tax = totalTax
		return taxLevels
	}

	taxLevels[4].Tax = totalTax
	return taxLevels
}

func initialTaxLevelDetail() []TaxLevel {
	return []TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{"150,001-500,000", 0.0},
		{"500,001-1,000,000", 0.0},
		{"1,000,001-2,000,000", 0.0},
		{"2,000,001 ขึ้นไป", 0.0},
	}
}
