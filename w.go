package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"runtime"
	"strconv"
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

const (
	insertSQL = "INSERT INTO prices (sku, date, list_price, final_price) VALUES (?, ?, ?, ?)"
)

var db *sql.DB

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic in webhookHandler: %v", r)
            stack := make([]byte, 1024*8)
            length := runtime.Stack(stack, true)
            log.Printf("%s", stack[:length])
        }
    }()
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
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	} else {
		log.Println("Successfully connected to the database!")
	}

	sku := r.URL.Query().Get("sku")
	listPriceStr := r.URL.Query().Get("list_price")
	finalPriceStr := r.URL.Query().Get("final_price")
	listPriceVal, err := strconv.ParseFloat(listPriceStr, 64)
	if err != nil {
		http.Error(w, "Invalid list_price value", http.StatusBadRequest)
		return
	}	
	finalPriceVal, err := strconv.ParseFloat(finalPriceStr, 64)
	if err != nil {
		http.Error(w, "Invalid list_price value", http.StatusBadRequest)
		return
	}
	if sku == "" || listPriceStr == "" || finalPriceStr == "" {
		http.Error(w, "sku, list_price, and final_price are required parameters", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("Executing SQL: %s with parameters: SKU: %s, Time: %s, ListPrice: %s, FinalPrice: %s", 
    insertSQL, sku, time.Now().Format("2006-01-02 15:04:05"), listPriceVal, finalPriceVal)
	log.Println(query)


	result, err := db.Exec(insertSQL, sku, time.Now(), listPriceVal, finalPriceVal)
	if err != nil || result == nil {
		http.Error(w, "Failed to insert data into the database", http.StatusInternalServerError)
		if err != nil {
			log.Println("Database error:", err)
		}
		return
	}

	fmt.Fprintf(w, "Data inserted successfully")
}

func main() {


	http.HandleFunc("/webhook", webhookHandler)
	log.Println("Server started on :8082")
	http.ListenAndServe(":8082", nil)
}
