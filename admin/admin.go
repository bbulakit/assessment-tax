package admin

/* Allowance Type */
type AllowanceInterfacer interface {
	Deduction() float64
	Amount() float64
}
