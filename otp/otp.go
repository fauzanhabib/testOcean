package otp

import (
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
	"oceanTest/model"
	"time"
)
const maxAttempts = 3
const cooldownDuration = 60 * time.Minute

var expectedOTP string
var attempts int
var cooldownTime time.Time

func GenerateOTPHandler(w http.ResponseWriter, r *http.Request,db *gorm.DB) {

	if attempts == maxAttempts && time.Now().Before(cooldownTime) {
		//set 60 minutes
		durationInMinutes := 60
		duration := time.Duration(durationInMinutes) * time.Minute
		minutes := int(duration.Minutes())
		http.Error(w, fmt.Sprintf("Maximum attemp generate OTP, please try again after %v%v.", minutes, " minutes"), http.StatusForbidden)
		return
	}

	attempts++
	//jika sudah maximal attemp set cooldowntime
	if attempts == maxAttempts {
		cooldownTime = time.Now().Add(cooldownDuration)
	}

	// set OTP dan save ke database
	randomOtp := generateRandomOTP(4)
	currentTimestamp := time.Now()
	otp := &model.Otp{
		Code: randomOtp,
		Attempts: 3,
		LastRetry: currentTimestamp,
	}

	result := db.Create(otp)
	if result.Error != nil {
		log.Printf("Error creating a new record: %v", result.Error)
		return
	}
	fmt.Fprintf(w, "result Generate OTP is : %04d\n", otp.Code)
}



func generateRandomOTP(digits int) int {
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999
	return rand.Intn(max-min+1) + min
}

func ValidateOTP(w http.ResponseWriter, r *http.Request,db *gorm.DB) {
	code := mux.Vars(r)["code"]

	// Cek apakah kode OTP valid
	var existingOTP model.Otp
	result := db.Where("code = ?", code).First(&existingOTP)
	if result.Error != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "Invalid OTP")
		return
	}

	// Cek batas percobaan attemp
	if existingOTP.Attempts == 0  {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintf(w, "You have reached the maximum attemp limit, please request again")
		return
	}


	// Validasi sukses
	existingOTP.Attempts = existingOTP.Attempts - 1
	existingOTP.LastRetry = time.Now()
	db.Save(&existingOTP)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Valid OTP",existingOTP.Code)
}
