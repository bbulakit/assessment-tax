package tax

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func TaxCalculationsHandler(c echo.Context) error {
	itd := IncomeTaxDetail{}
	err := c.Bind(&itd)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tcr := taxCalculate(itd)

	return c.JSON(http.StatusOK, tcr)
}

func TaxUploadCsvHandler(c echo.Context) error {
	file, err := c.FormFile("taxFile")
	if err != nil {
		fmt.Println(err)
		return err
	}

	src, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer src.Close()

	reader := csv.NewReader(src)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		return err
	}

	for i, record := range records {
		if i == 0 {
			continue // Skip header row
		}

		if err := validateCsvData(record); err != nil {
			fmt.Println(err)
			return err
		}

		totalIncome, _ := strconv.ParseFloat(record[0], 64)
		wht, _ := strconv.ParseFloat(record[1], 64)
		donation, _ := strconv.ParseFloat(record[2], 64)

		fmt.Printf("totalIncome: %.2f, wht: %.2f, donation: %.2f\n", totalIncome, wht, donation)
	}

	return c.String(http.StatusOK, "CSV file uploaded successfully")
}

func validateCsvData(record []string) error {
	for _, field := range record {
		if strings.TrimSpace(field) == "" {
			return fmt.Errorf("all values must be non-empty")
		}
	}

	for _, field := range record {
		if _, err := strconv.ParseFloat(field, 64); err != nil {
			return fmt.Errorf("invalid format: %s", err)
		}
	}

	return nil
}
