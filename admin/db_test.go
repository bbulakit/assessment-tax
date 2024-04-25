package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var dbHandler *DBHandler

func setup(t *testing.T) func() {
	t.Parallel()
	dataSource := "host=localhost port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable"
	var err error
	dbHandler, err = NewDBHandler(dataSource)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	if err := dbHandler.SeedInitialData(); err != nil {
		t.Fatalf("Failed to seed initial data: %v", err)
	}

	teardown := func() {
		if _, err := dbHandler.DB.Exec("DROP TABLE IF EXISTS deductions;"); err != nil {
			t.Logf("Failed to drop table `deductions`: %v", err)
		}
		if err := dbHandler.DB.Close(); err != nil {
			t.Logf("Failed to close database connection: %v", err)
		}
	}

	return teardown
}

func TestGetDeductionsHandler(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/deductions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, dbHandler.GetDeductionsHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var deductions []Deduction
		if err := json.Unmarshal(rec.Body.Bytes(), &deductions); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		expected := []Deduction{
			{Name: "personalDeduction", Value: 60000.0},
			{Name: "donation", Value: 100000.0},
			{Name: "kreceipt", Value: 50000.0},
		}

		assert.Equal(t, expected, deductions)
	}
}

func TestGetDeductionHandler(t *testing.T) {
	teardown := setup(t)
	defer teardown()

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/deduction/personalDeduction", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("personalDeduction")

	if assert.NoError(t, dbHandler.GetDeductionHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var deduction Deduction
		if err := json.Unmarshal(rec.Body.Bytes(), &deduction); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		expected := Deduction{Name: "personalDeduction", Value: 60000.0}
		assert.Equal(t, expected, deduction)
	}
}
