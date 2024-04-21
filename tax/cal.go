package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func TaxCalculationsHandler(c echo.Context) error {
	u := Tax{}
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	// row := db.QueryRow("INSERT INTO users (name, age) values ($1, $2)  RETURNING id", u.Allowances, u.Allowances)
	// err = row.Scan(&u.Allowances)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	// }
	tr := TaxResponse{}
	tr.TotalTax = 29_000.0

	return c.JSON(http.StatusOK, tr)
}

func taxCalculate(tax Tax) float64 {
	cal := tax.TotalIncome
	for _, deduction := range tax.Allowances {
		cal -= deduction.Amount
	}

	cal -= tax.WithHoldingTax

	cal *= taxRate(tax.TotalIncome)

	return cal
}

func taxRate(totalIncome float64) float64 {
	if totalIncome <= 150000 {
		return 0.00
	} else if totalIncome <= 500000 {
		return 0.10
	} else if totalIncome <= 1000000 {
		return 0.15
	} else if totalIncome <= 2000000 {
		return 0.20
	}
	return 0.35
}
