package tax

type IncomeTaxDetail struct {
	TotalIncome    float64     `json:"totalIncome"`
	WithHoldingTax float64     `json:"wht"`
	Allowances     []Allowance `json:"allowances"`
	TaxRefund      float64     `json:"taxRefund,omitempty"`
}

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type TaxCalculationResult struct {
	TotalTax  float64    `json:"tax"`
	TaxLevels []TaxLevel `json:"taxLevel"`
}

type TaxCsv struct {
	TotalIncome float64 `form:"totalIncome"`
	Wht         float64 `form:"wht"`
	Donation    float64 `form:"donation"`
}

type TaxCsvResult struct {
	Taxes []TaxCsvResultDetail `json:"taxes"`
}

type TaxCsvResultDetail struct {
	TotalIncome float64 `json:"totalIncome"`
	Tax         float64 `json:"tax,omitempty"`
	TaxRefund   float64 `json:"taxRefund,omitempty"`
}

type Err struct {
	Message string `json:"message"`
}
