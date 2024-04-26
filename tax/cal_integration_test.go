package tax

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type TestHTTPRequest struct {
	method   string
	target   string
	username string
	password string
	body     io.Reader
}

func setup() *echo.Echo {
	e := echo.New()
	e.POST("/tax/calculations", TaxCalculationsHandler)
	e.POST("/tax/calculations/upload-csv", TaxUploadCsvHandler)
	return e
}

func testHTTPRequest(e *echo.Echo, r TestHTTPRequest) (int, []byte) {
	req := httptest.NewRequest(r.method, r.target, r.body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	req.SetBasicAuth(r.username, r.password)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec.Code, rec.Body.Bytes()
}

func TestTaxCalulate(t *testing.T) {
	e := setup()
	body := bytes.NewBufferString(`{
		"totalIncome": 500000.0,
		"wht": 0.0,
		"allowances": [
		  {
			"allowanceType": "donation",
			"amount": 0.0
		  }
		]
	  }`)

	req := TestHTTPRequest{
		method:   http.MethodPost,
		target:   "/tax/calculations",
		username: "AdminTax", //os.Getenv("ADMIN_USERNAME")
		password: "admin!",   //os.Getenv("ADMIN_PASSWORD")
		body:     body,
	}
	code, responseBody := testHTTPRequest(e, req)

	assert.Equal(t, http.StatusOK, code)

	var response struct {
		TotalTax float64 `json:"tax"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}
	assert.Equal(t, 29000.0, response.TotalTax)
}

func TestTaxCalulateWithWht(t *testing.T) {
	e := setup()
	body := bytes.NewBufferString(`{
		"totalIncome": 500000.0,
		"wht": 25000.0,
		"allowances": [
		  {
			"allowanceType": "donation",
			"amount": 0.0
		  }
		]
	  }`)
	req := TestHTTPRequest{
		method:   http.MethodPost,
		target:   "/tax/calculations",
		username: "AdminTax", //os.Getenv("ADMIN_USERNAME")
		password: "admin!",   //os.Getenv("ADMIN_PASSWORD")
		body:     body,
	}
	code, responseBody := testHTTPRequest(e, req)

	assert.Equal(t, http.StatusOK, code)

	var response struct {
		TotalTax float64 `json:"tax"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}
	assert.Equal(t, 4_000.0, response.TotalTax)
}

func TestTaxCalulateWithDonation(t *testing.T) {
	e := setup()
	body := bytes.NewBufferString(`{
		"totalIncome": 500000.0,
		"wht": 0.0,
		"allowances": [
		  {
			"allowanceType": "donation",
			"amount": 200000.0
		  }
		]
	  }`)
	req := TestHTTPRequest{
		method:   http.MethodPost,
		target:   "/tax/calculations",
		username: "AdminTax", //os.Getenv("ADMIN_USERNAME")
		password: "admin!",   //os.Getenv("ADMIN_PASSWORD")
		body:     body,
	}
	code, responseBody := testHTTPRequest(e, req)

	assert.Equal(t, http.StatusOK, code)

	var response struct {
		TotalTax float64 `json:"tax"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}
	assert.Equal(t, 19_000.0, response.TotalTax)
}

func TestTaxCalulateAndGetTaxLevelDetail(t *testing.T) {
	e := setup() // Setup the Echo instance
	body := bytes.NewBufferString(`{
        "totalIncome": 500000.0,
        "wht": 0.0,
        "allowances": [
            {
                "allowanceType": "donation",
                "amount": 200000.0
            }
        ]
    }`)

	req := TestHTTPRequest{
		method:   http.MethodPost,
		target:   "/tax/calculations",
		username: "AdminTax", // You could use os.Getenv("ADMIN_USERNAME") here
		password: "admin!",   // And os.Getenv("ADMIN_PASSWORD") for production
		body:     body,
	}

	code, responseBody := testHTTPRequest(e, req)

	assert.Equal(t, http.StatusOK, code)

	var response TaxCalculationResult
	if err := json.Unmarshal(responseBody, &response); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}

	expectedTaxLevels := []TaxLevel{
		{Level: "0-150,000", Tax: 0.0},
		{Level: "150,001-500,000", Tax: 19_000.0},
		{Level: "500,001-1,000,000", Tax: 0.0},
		{Level: "1,000,001-2,000,000", Tax: 0.0},
		{Level: "2,000,001 ขึ้นไป", Tax: 0.0},
	}

	assert.Equal(t, 19000.0, response.TotalTax)
	assert.Equal(t, expectedTaxLevels, response.TaxLevels, "Mismatch in tax levels")
}

func TestTaxUploadCsvHandler(t *testing.T) {
	csvData := "totalIncome,wht,donation\n" +
		"500000,0,0\n" +
		"600000,40000,20000\n" +
		"750000,50000,15000"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("taxFile", "taxes.csv")
	if err != nil {
		t.Fatal("Failed to create form file:", err)
	}
	_, err = io.Copy(part, strings.NewReader(csvData))
	if err != nil {
		t.Fatal("Failed to write CSV data to form file:", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rec := httptest.NewRecorder()

	TaxUploadCsvHandler(echo.New().NewContext(req, rec))

	assert.Equal(t, http.StatusOK, rec.Code)

	var result TaxCsvResult
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}

	expected := TaxCsvResult{
		Taxes: []TaxCsvResultDetail{
			{TotalIncome: 500000, Tax: 29000},
			{TotalIncome: 600000, TaxRefund: 37000},
			{TotalIncome: 750000, TaxRefund: 23750},
		},
	}
	assert.Equal(t, expected, result)
}
