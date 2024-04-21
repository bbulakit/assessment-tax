package admin

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
