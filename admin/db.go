package admin

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DBHandler struct {
	DB *sql.DB
}

func NewDBHandler(dataSourceName string) (*DBHandler, error) {
	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Fatal("Connect to database error", err)
		return nil, err
	}

	createTb := `
		CREATE TABLE IF NOT EXISTS deductions ( id SERIAL PRIMARY KEY, name TEXT UNIQUE, value DOUBLE PRECISION);
	`

	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
		return nil, err
	}

	return &DBHandler{db}, nil
}

func (h *DBHandler) SeedInitialData() error {
	seedData := `
        INSERT INTO deductions (name, value) VALUES
        ('personalDeduction', 60000.0),
        ('donation', 100000.0),
        ('kReceipt', 50000.0)
        ON CONFLICT (name) DO UPDATE
        SET value = EXCLUDED.value
    `
	if _, err := h.DB.Exec(seedData); err != nil {
		log.Println("Error seeding initial data", err)
		return err
	}
	return nil
}
