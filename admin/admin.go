package admin

/* Allowance Type */
type DeductionInterfacer interface {
	Deduction() float64
	Amount() float64
}

type Deduction struct {
	Name  string
	Value float64
}
