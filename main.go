package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	db2 "oceanTest/db"
	"oceanTest/model"
	"oceanTest/otp"
)


func main(){
	db, err := db2.ConnectDb()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Migrasi database untuk membuat tabel OTP
	db.AutoMigrate(&model.Otp{})

	r := mux.NewRouter()
	r.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		otp.GenerateOTPHandler(w, r, db) // Inject the database connection
	}).Methods("GET")
	r.HandleFunc("/request/{code}", func(w http.ResponseWriter, r *http.Request) {
		otp.ValidateOTP(w, r, db) // Inject the database connection
	}).Methods("POST")

	http.Handle("/", r)
	log.Println("Server is running on :8181")
	http.ListenAndServe(":8181", nil)
}

