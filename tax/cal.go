package tax

func TaxCalculate(tax Tax) float64 {
	return 0
}

func TaxRate(totalIncome float64) float64 {
	if totalIncome <= 150000 {
		return 0.00
	} else if totalIncome <= 500000 {
		return 0.10
	} else if totalIncome <= 1000000 {
		return 0.15
	} else if totalIncome <= 2000000 {
		return 0.20
	}
	return 0.35
}
