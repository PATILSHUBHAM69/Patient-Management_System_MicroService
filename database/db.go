package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func Init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUsername, dbPassword, dbName))
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MySQL database!")

	createPatientsTable()
}

func createPatientsTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS patients2 (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(80),
        dob DATE,
        gender VARCHAR(10),
        contact VARCHAR(100),
        medical_history TEXT,
        attainder VARCHAR(50),
        isinsurance VARCHAR(25),
        type VARCHAR(40),
        payer varchar(80),
        claim_no VARCHAR(50),
        claim_status VARCHAR(50),
        claim_amt INT,
        settled_amt INT,
        PR_amt INT AS (claim_amt - settled_amt)
    );`)

	if err != nil {
		panic(err)
	}

	fmt.Println("Patients Details Table Successfully Created")
}
