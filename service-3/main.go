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
	Name           string      `json:"name"`
	DOB            string      `json:"dob"`
	Gender         string      `json:"gender" validate:"oneof=Male Female"`
	Contact        string      `json:"contact" validate:"len=10"`
	MedicalHistory string      `json:"medical_history"`
	Attainder      string      `json:"attainder"`
	IsInsurance    string      `json:"isinsurance"`
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

func UpdatePatient(c *gin.Context) {
	patientIDStr := c.Query("id")
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil || patientID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}
	var patient Patient
	err = json.NewDecoder(c.Request.Body).Decode(&patient)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("UPDATE patients2 SET name=?, dob=?, gender=?, contact=?, medical_history=?,attainder=?,isinsurance=?,payer=?,type=?,claim_no=?,claim_status=?,claim_amt=?,settled_amt=? WHERE id=?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(patient.Name, patient.DOB, patient.Gender, patient.Contact, patient.MedicalHistory, patient.Attainder, patient.IsInsurance, patient.Payer, patient.Type, patient.ClaimNo, patient.ClaimStatus, patient.ClaimAmt, patient.SettledAmt, patientID)
	if err != nil {
		log.Printf("Error in updating value :%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Patient updated successfully"})
}

func GetPatientForUpdate(c *gin.Context) {
	patientIDStr := c.Query("id")
	patientID, err := strconv.Atoi(patientIDStr)
	if err != nil || patientID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	row := db.QueryRow("SELECT * FROM patients2 WHERE id=?", patientID)

	var patient Patient
	var claimAmt, settledAmt, prAmt int64

	err = row.Scan(
		&patient.ID, &patient.Name, &patient.DOB, &patient.Gender, &patient.Contact, &patient.MedicalHistory,
		&patient.Attainder, &patient.IsInsurance, &patient.Type, &patient.Payer, &patient.ClaimNo, &patient.ClaimStatus,
		&claimAmt, &settledAmt, &prAmt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	patient.ClaimAmt = json.Number(strconv.FormatInt(claimAmt, 10))
	patient.SettledAmt = json.Number(strconv.FormatInt(settledAmt, 10))
	patient.PRamt = json.Number(strconv.FormatInt(prAmt, 10))

	c.JSON(http.StatusOK, patient)
}

func main() {
	router := gin.Default()
	database.Init()
	connectDB()

	router.GET("/service3/patients", GetPatientForUpdate)
	router.PUT("/service3/patients", UpdatePatient)

	err := router.Run(":8082")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
