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
