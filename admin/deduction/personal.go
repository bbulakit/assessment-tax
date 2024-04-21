package admin

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
