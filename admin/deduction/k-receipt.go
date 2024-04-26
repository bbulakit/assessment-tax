package admin

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
