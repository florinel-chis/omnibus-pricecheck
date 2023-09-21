package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
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

	// Find the latest date in the prices table
	var latestDateString string
	err = db.QueryRow(`SELECT MAX(DATE(date)) FROM prices`).Scan(&latestDateString)
	if err != nil {
		log.Fatalf("Error fetching latest date: %v", err)
	}

	latestDate, err := time.Parse("2006-01-02", latestDateString)
	if err != nil {
		log.Fatalf("Error parsing latest date: %v", err)
	}	
	// Get list_prices for the latest date and there is a promotional price setup
	rows, err := db.Query(`SELECT sku, list_price FROM prices WHERE final_price < list_price AND DATE(date) = ? and sku='SKU44241681' GROUP BY sku,list_price LIMIT 5`, latestDate)
	if err != nil {
		log.Fatalf("Error fetching latest prices: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var sku string
		var latestListPrice float64
		
		if err := rows.Scan(&sku, &latestListPrice); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		fmt.Printf("Checking SKU %s...\n", sku)
		// Check if the latest list_price is the minimum over the past 30 days
		var minPrice sql.NullFloat64
		err = db.QueryRow(`SELECT MIN(final_price) FROM prices WHERE sku = ? AND DATE(date) BETWEEN ? AND ?`,
			sku, latestDate.AddDate(0, 0, -30).Format("2006-01-02"), latestDate.AddDate(0, 0, -1).Format("2006-01-02")).Scan(&minPrice)
		if err != nil {
			log.Fatalf("Error fetching minimum price for SKU %s: %v", sku, err)
		}
		
		if !minPrice.Valid {
			fmt.Printf("No records for SKU %s over the past 30 days.\n", sku)
			continue
		}
		
		if latestListPrice > minPrice.Float64 {
			fmt.Printf("Incompatibility detected for SKU %s! Latest list_price (%.4f) is not the minimum over the past 30 days: %.4f.\n",
			sku, latestListPrice, minPrice.Float64)
		} else {
			//fmt.Printf("SKU %s is compliant with the Omnibus directive.\n", sku)
			//fmt.Printf("Latest list_price: %.4f\n", latestListPrice)
			//fmt.Printf("Minimum over the past 30 days: %.4f\n", minPrice.Float64)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v", err)
	}
}
