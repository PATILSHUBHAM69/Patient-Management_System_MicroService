package main

import (
	"database/sql"
	"demo/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Patient struct {
	ID             int         `json:"id"`
	Name           string      `json:"name" validate:"required"`
	DOB            string      `json:"dob" validate:"required"`
	Gender         string      `json:"gender" validate:"required,oneof=Male Female"`
	Contact        string      `json:"contact" validate:"required,len=10"`
	MedicalHistory string      `json:"medical_history" validate:"required"`
	Attainder      string      `json:"attainder" validate:"required"`
	IsInsurance    string      `json:"isinsurance" validate:"required"`
	Payer          string      `json:"payer"`
	Type           string      `json:"type"`
	ClaimNo        string      `json:"claimno"`
	ClaimStatus    string      `json:"claimstatus"`
	ClaimAmt       json.Number `json:"claimamt"`
	SettledAmt     json.Number `json:"settledamt"`
	PRamt          json.Number `json:"pramt"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func connectDB() {
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
}

func CreatePatient(c *gin.Context) {
	var patient Patient
	err := c.ShouldBindJSON(&patient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	result, err := db.Exec("INSERT INTO patients2(name, dob, gender, contact, medical_history,attainder, isinsurance,type,payer,claim_no,claim_status,claim_amt,settled_amt) VALUES (?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?)", patient.Name, patient.DOB, patient.Gender, patient.Contact, patient.MedicalHistory, patient.Attainder, patient.IsInsurance, patient.Type, patient.Payer, patient.ClaimNo, patient.ClaimStatus, patient.ClaimAmt, patient.SettledAmt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	patient.ID = int(id)
	c.JSON(http.StatusCreated, gin.H{"New patient id : ": id})
}

func main() {
	router := gin.Default()
	database.Init()
	connectDB()

	router.POST("service1/patients", CreatePatient)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
