package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bbulakit/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	apiPort := os.Getenv("PORT")
	adminUsername := "adminTax" //os.Getenv("ADMIN_USERNAME")
	adminPassword := "admin!"   //os.Getenv("ADMIN_PASSWORD")

	e.Use(middleware.Logger())

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == adminUsername && password == adminPassword {
			return true, nil
		}
		return false, nil
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	})

	e.POST("/tax/calculations", tax.TaxCalculationsHandler)

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
