package tax

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func taxCalculate(itd IncomeTaxDetail) TaxCalculationResult {
	tcr := TaxCalculationResult{}
	taxCal := taxLevelDeduction(itd.TotalIncome)

	personalDeduction := GetDeduction("personal")
	if personalDeduction <= 0 {
		personalDeduction = 60_000.0
	}

	taxCal -= personalDeduction

	for _, deduction := range itd.Allowances {
		var actualDeduction float64
		if deduction.AllowanceType == "donation" {

			maxDonationDeduction := GetDeduction("donation")
			if maxDonationDeduction <= 0 {
				maxDonationDeduction = 100_000.0
			}

			actualDeduction = deduction.Amount
			if actualDeduction > maxDonationDeduction {
				actualDeduction = maxDonationDeduction
			}
		}

		if deduction.AllowanceType == "k-receipt" {

			maxKReceiptDeduction := GetDeduction("donation")
			if maxKReceiptDeduction <= 0 {
				maxKReceiptDeduction = 50_000.0
			}

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
	if tcr.TotalTax < 0 {
		tcr.TaxRefund = tcr.TotalTax * -1.0
		//tcr.TotalTax = 0
	}
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

	if totalTax < 0 {
		totalTax *= -1
	}

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

type Deduction struct {
	Name  string
	Value float64
}

func GetDeduction(name string) float64 {
	client := &http.Client{}
	apiPort := os.Getenv("PORT")
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s/admin/deductions/%s", apiPort, name), nil)

	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}

	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	req.SetBasicAuth(adminUsername, adminPassword)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}

	var deduction Deduction
	err = json.Unmarshal(body, &deduction)
	if err != nil {
		fmt.Println(err.Error())
		return 0.0
	}

	return deduction.Value
}
