package tax

/*
	Tax Rate are fixed
							TAX Rate
	0 ~ 150,000 			= 0
	150,001 ~ 500,000		= 10 %
	500,001 ~ 1,000,000		= 15 %
	1,000,001 ~ 2,000,000	= 20 %
	> 2,000,000				= 35 %

	Allowances (value > 0)
	1.	'personalDeduction' as ค่าลดหย่อนส่วนตัว
	Default: 60,000 (static)
	Max mod.: 100,000
	Min.: > 10,000

	2.	'donation' as เงินบริจาค
	Max deduction: 100,000

	3.	'k-receipt' as ช้อปลดภาษี
	Default max: 50,000
	Max mod.: 100,000
	Min.: > 0

	ภาษีหัก ณ ที่จ่าย (เสมือนว่าจ่ายไปแล้วล่วงหน้า)
	wht: withholding tax (value > 0 && value < income)
	income - deduction - wht (negative value = taxRefund)

*/

type IncomeTaxDetail struct {
	TotalIncome    float64     `json:"totalIncome"`
	WithHoldingTax float64     `json:"wht"`
	Allowances     []Allowance `json:"allowances"`
	TaxRefund      float64     `json:"taxRefund"`
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

/* End of Allowance Type */

type Err struct {
	Message string `json:"message"`
}
