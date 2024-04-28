package tax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaxRate(t *testing.T) {
	tests := []struct {
		name        string
		totalIncome float64
		want        float64
	}{
		{name: "Income=0 then TaxRate=0", totalIncome: 0, want: 0.00},
		{"Income<=150,000 then TaxRate=0", 150_000, 0.00},
		{"Income=150,001 then TaxRate=10%", 150_001, 0.10},
		{"Income<=500,000 then TaxRate=10%", 500_000, 0.10},
		{"Income=500,001 then TaxRate=15%", 500_001, 0.15},
		{"Income<=1,000,000 then TaxRate=15%", 1_000_000, 0.15},
		{"Income=1,000,001 then TaxRate=20%", 1_000_001, 0.20},
		{"Income<=2,000,000 then TaxRate=20%", 2_000_000, 0.20},
		{"Income>2,000,000 then TaxRate=35%", 2_000_001, 0.35},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Arrange

			//Act
			got := taxRate(test.totalIncome)

			//Assert
			if test.want != got {
				t.Errorf("TaxRate(%f) = %f; want %f", test.totalIncome, got, test.want)
			}
		})
	}

}

func TestInitialTaxLevelDetail(t *testing.T) {
	want := []TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{"150,001-500,000", 0.0},
		{"500,001-1,000,000", 0.0},
		{"1,000,001-2,000,000", 0.0},
		{"2,000,001 ขึ้นไป", 0.0},
	}

	got := initialTaxLevelDetail()

	assert.Equal(t, want, got)
}

func TestTaxLevelDeduction(t *testing.T) {
	tests := []struct {
		name           string
		totalIncome    float64
		expectedResult float64
	}{
		{name: "income < 150,000", totalIncome: 100_000, expectedResult: 0.00},
		{"income between 150,000 and 500,000", 300_000, 150_000},
		{"income between 500,000 and 1,000,000", 750_000, 250_000},
		{"income between 1,000,000 and 2,000,000", 1_500_000, 500_000},
		{"income > 2,000,000", 2_500_000, 500_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taxLevelDeduction(tt.totalIncome)
			assert.Equal(t, tt.expectedResult, got)
		})
	}
}

func TestValidateCsvData(t *testing.T) {
	tests := []struct {
		record []string
		errMsg string
	}{
		{record: []string{"500000", "0", "0"}, errMsg: ""},
		{record: []string{"600000", "40000", ""}, errMsg: "all values must be non-empty"},
		{record: []string{"700000", "50000", "abc"}, errMsg: "invalid format"},
	}

	for _, test := range tests {
		err := validateCsvData(test.record)
		if test.errMsg == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.errMsg)
		}
	}
}

func TestValidateTaxValues(t *testing.T) {
	tests := []struct {
		incomeTaxDetail IncomeTaxDetail
		errMsg          string
	}{
		{incomeTaxDetail: IncomeTaxDetail{TotalIncome: 500000, WithHoldingTax: 0, Allowances: []Allowance{{AllowanceType: "donation", Amount: 200000}}}, errMsg: ""},                                                                     // Valid income tax detail
		{incomeTaxDetail: IncomeTaxDetail{TotalIncome: -500000, WithHoldingTax: 0, Allowances: []Allowance{{AllowanceType: "donation", Amount: 200000}}}, errMsg: "total income (-500000.00) cannot be negative"},                        // Negative total income
		{incomeTaxDetail: IncomeTaxDetail{TotalIncome: 500000, WithHoldingTax: -10000, Allowances: []Allowance{{AllowanceType: "donation", Amount: 200000}}}, errMsg: "wht (-10000.00) cannot be negative"},                              // Negative withholding tax
		{incomeTaxDetail: IncomeTaxDetail{TotalIncome: 500000, WithHoldingTax: 600000, Allowances: []Allowance{{AllowanceType: "donation", Amount: 200000}}}, errMsg: "wht (600000.00) cannot be greater than total income (500000.00)"}, // WHT greater than total income
		{incomeTaxDetail: IncomeTaxDetail{TotalIncome: 500000, WithHoldingTax: 0, Allowances: []Allowance{{AllowanceType: "invalid", Amount: 200000}}}, errMsg: "invalid allowance type: invalid"},                                       // Invalid allowance type
		{incomeTaxDetail: IncomeTaxDetail{TotalIncome: 500000, WithHoldingTax: 0, Allowances: []Allowance{{AllowanceType: "donation", Amount: -200000}}}, errMsg: "allowance amount (-200000.00) cannot be negative"},                    // Negative allowance amount
	}

	for _, test := range tests {
		err := validateTaxValues(&test.incomeTaxDetail)
		if test.errMsg == "" {
			assert.NoError(t, err)
		} else {
			assert.EqualError(t, err, test.errMsg)
		}
	}
}
