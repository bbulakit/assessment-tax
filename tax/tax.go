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
	Default: 50,000
	Max mod.: 100,000
	Min.: > 0

	ภาษีหัก ณ ที่จ่าย (เสมือนว่าจ่ายไปแล้วล่วงหน้า)
	wht: withholding tax (value > 0 && value < income)
	income - deduction - wht (negative value = taxRefund)

*/

type IncomeTaxDetails struct {
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

type TaxResponse struct {
	TotalTax  float64    `json:"tax"`
	TaxLevels []TaxLevel `json:"taxLevel"`
}

/* Allowance Type */
type AllowanceInterfacer interface {
	Deduction() float64
	Amount() float64
}

type Donation struct {
	amount float64
}

func (d Donation) Deduction() float64 {
	if d.amount >= 100_000 {
		return 100_000
	}

	if d.amount < 0 {
		return 0
	}

	return d.amount
}

func (d Donation) Amount() float64 {
	return d.amount
}

type PersonalDeduction struct {
	amount float64
}

func (p PersonalDeduction) Deduction() float64 {
	if p.amount >= 100_000 {
		return 100_000
	}

	if p.amount <= 10_000 {
		return 10_000
	}
	return p.amount
}

func (p PersonalDeduction) Amount() float64 {
	return p.amount
}

type KReceipt struct {
	amount float64
}

func (k KReceipt) Deduction() float64 {
	if k.amount >= 100_000 {
		return 100_000
	}

	if k.amount < 0 {
		return 0
	}
	return k.amount
}

func (k KReceipt) Amount() float64 {
	return k.amount
}

/* End of Allowance Type */

type Err struct {
	Message string `json:"message"`
}
