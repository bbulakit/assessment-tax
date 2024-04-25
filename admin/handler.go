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
