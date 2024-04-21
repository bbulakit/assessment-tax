package tax

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Response struct {
	*http.Response
	err error
}

func TestTaxCalulate(t *testing.T) {
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
	tr := struct {
		TotalTax float64 `json:"tax"`
	}{}

	res := request(http.MethodPost, uri("tax/calculations"), body)
	err := res.Decode(&tr)
	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 29_000.0, tr.TotalTax)
}

func TestTaxCalulateWithWht(t *testing.T) {
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
	tr := struct {
		TotalTax float64 `json:"tax"`
	}{}

	res := request(http.MethodPost, uri("tax/calculations"), body)
	err := res.Decode(&tr)
	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 4_000.0, tr.TotalTax)
}

func uri(paths ...string) string {
	apiPort := os.Getenv("PORT")
	host := "http://localhost:" + apiPort
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	username := "adminTax" //os.Getenv("ADMIN_USERNAME")
	password := "admin!"   //os.Getenv("ADMIN_PASSWORD")

	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v) //ไม่ก็ json.Unmarshal
}
