package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bbulakit/assessment-tax/admin"

	"github.com/bbulakit/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	apiPort := os.Getenv("PORT")
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	databaseUrl := os.Getenv("DATABASE_URL")
	dbHandler, err := admin.NewDBHandler(databaseUrl)

	if err != nil {
		fmt.Println("Cannot create new database handlers")
		return
	}

	dbHandler.SeedInitialData()

	e.Use(middleware.Logger())

	adminGroup := e.Group("/admin", middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == adminUsername && password == adminPassword {
			return true, nil
		}
		return false, nil
	}))

	e.POST("/tax/calculations", tax.TaxCalculationsHandler)
	e.POST("/tax/calculations/upload-csv", tax.TaxUploadCsvHandler)
	adminGroup.GET("/deductions/:name", dbHandler.GetDeductionHandler)
	adminGroup.GET("/deductions/", dbHandler.GetDeductionsHandler)
	adminGroup.POST("/deductions/:name", dbHandler.PostDeductionHandler)

	go func() {
		if err := e.Start(":" + apiPort); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
