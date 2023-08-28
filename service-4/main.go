package main

import (
	"database/sql"
	"demo/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

func DeletePatient(c *gin.Context) {
	id := c.Param("id")
	patientID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	_, err = db.Exec("DELETE FROM patients2 WHERE id=?", patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patient deleted successfully"})
}

func main() {
	router := gin.Default()
	database.Init()
	connectDB()
	router.DELETE("/service4/patients/:id", DeletePatient)
	err := router.Run(":8083")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
