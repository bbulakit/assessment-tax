package admin

import (
	"database/sql"

	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *DBHandler) GetDeductionsHandler(c echo.Context) error {
	stmt, err := h.DB.Prepare("SELECT name, value FROM deductions;")
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	deductions := []Deduction{}
	for rows.Next() {
		var d Deduction
		if err := rows.Scan(&d.Name, &d.Value); err != nil {
			return err
		}
		deductions = append(deductions, d)
	}

	return c.JSON(http.StatusOK, deductions)
}

func (h *DBHandler) GetDeductionHandler(c echo.Context) error {
	name := c.Param("name")
	name = convertDeductionPathToName(name)
	stmt, err := h.DB.Prepare("SELECT name, value FROM deductions WHERE name = $1")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	row := stmt.QueryRow(name)
	u := Deduction{}
	err = row.Scan(&u.Name, &u.Value)
	switch err {
	case sql.ErrNoRows:
		return c.JSON(http.StatusNotFound, err.Error())
	case nil:
		return c.JSON(http.StatusOK, u)
	default:
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
}

func (h *DBHandler) PostDeductionHandler(c echo.Context) error {
	name := c.Param("name")
	name = convertDeductionPathToName(name)
	type TempDeduction struct {
		Amount float64 `json:"amount"`
	}

	var td TempDeduction
	if err := c.Bind(&td); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	deduction := Deduction{
		Name:  name,
		Value: td.Amount,
	}

	stmt, err := h.DB.Prepare("UPDATE deductions SET value = $2 WHERE name = $1 RETURNING name, value")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer stmt.Close()

	err = stmt.QueryRow(name, deduction.Value).Scan(&deduction.Name, &deduction.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Deduction not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := map[string]float64{deduction.Name: deduction.Value}

	return c.JSON(http.StatusOK, response)
}

func convertDeductionPathToName(path string) string {
	if path == "personal" {
		return "personalDeduction"
	}

	if path == "donation" {
		return "donation"
	}

	if path == "k-receipt" {
		return "kReceipt"
	}

	return path
}
