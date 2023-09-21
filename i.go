package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"time"
	"log"
	"github.com/spf13/viper"
	_ "github.com/go-sql-driver/mysql"
)


type Config struct {
	Database struct {
		User     string
		Password string
		Name     string
		Host     string
		Port     string
	}
}
// Price represents a single row in the prices table.
type Price struct {
	Sku        string
	Date       time.Time
	ListPrice  float64
	FinalPrice float64
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".") // Look for the config in the current directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		log.Fatalf("Error: %v", err)
	}

	var configuration Config
	err = viper.Unmarshal(&configuration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	// Establish connection to the MySQL database
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configuration.Database.User, configuration.Database.Password, configuration.Database.Host, configuration.Database.Port, configuration.Database.Name)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Read from CSV
	file, err := os.Open("prices2.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	// Assuming the CSV has headers, we read one line to skip it
	_, err = r.Read()
	if err != nil {
		panic(err)
	}

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO prices (sku, date, list_price, final_price) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE list_price=values(list_price), final_price=values(final_price)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	// Iterate over the records and insert them into the MySQL table
	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		// Assuming the CSV columns are in the order: sku, list_price, final_price
		price := Price{
			Sku:        record[0],
			Date:       time.Now(),   // Using the current timestamp for simplicity
			ListPrice:  ParseFloat(record[1]),
			FinalPrice: ParseFloat(record[2]),
		}

		_, err = stmt.Exec(price.Sku, price.Date, price.ListPrice, price.FinalPrice)
		if err != nil {
			panic(err)
		}
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

// ParseFloat converts string to float64
func ParseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
