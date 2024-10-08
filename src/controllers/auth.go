package controllers

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Sahil2k07/Blog-App-Go/src/config"
	"github.com/Sahil2k07/Blog-App-Go/src/dto"
	"github.com/Sahil2k07/Blog-App-Go/src/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Signup or Register the User
func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WrongMethod(w)
		return
	}

	var req dto.SignupRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.InvalidInput(w)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.InvalidInput(w)
		return
	}

	if config.CheckEmailInBloom(req.Email) {
		// Check in Bloom filter first
		utils.InvalidInput(w, "Email already Registered")
		return
	} else {
		var existingEmail string
		err := config.DB.QueryRow(`SELECT email FROM User WHERE email = ?`, req.Email).Scan(&existingEmail)
		if err == nil {
			utils.InvalidInput(w, "Email already Registered")
			return
		}
	}

	otp := func() string {
		otp := make([]byte, 3)

		_, err := rand.Read(otp)
		if err != nil {
			return "123456"
		}

		return fmt.Sprintf("%06d", int(otp[0])%10*100000+int(otp[1])%10000+int(otp[2])%1000)
	}()

	_, err = config.DB.Exec(
		`
			INSERT INTO Otp (id, email, otp)
			VALUES (UUID(), ?, ?)
		`, req.Email, otp,
	)
	if err != nil {
		utils.InternalServerError(w, "Problem generating OTP")
		return
	}

	go func(email, otp string) {
		if err := config.Mailer(email, otp); err != nil {
			fmt.Printf("Failed to send OTP to %s: %v\n", email, err)
		}
	}(req.Email, otp)

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.InternalServerError(w, "Something went Wrong in our Server")
		return
	}

	userId := uuid.New().String()

	_, err = config.DB.Exec(
		`
			INSERT INTO User (id, email, password)
			VALUES (?, ?, ?)
		`, userId, req.Email, hashedPassword)
	if err != nil {
		utils.InternalServerError(w, "Failed to create user")
		return
	}

	imageURL := fmt.Sprintf("https://api.dicebear.com/5.x/initials/svg?seed=%s %s", req.FirstName, req.LastName)

	_, err = config.DB.Exec(
		`
			INSERT INTO Profile(id, userId, firstName, lastName, image)
			VALUES (UUID(), ?, ?, ?, ?)
		`, userId, req.FirstName, req.LastName, imageURL,
	)
	if err != nil {
		utils.InternalServerError(w, "Failed to create user profile")
		return
	}

	config.AddEmailToBloom(req.Email)

	response := map[string]interface{}{
		"success": true,
		"message": "User Created Successfully. Please verify with the OTP to continue",
		"userId":  userId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

// Verify the User
func VerifyUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WrongMethod(w)
		return
	}

	var req dto.VerifyUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.InvalidInput(w)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.InvalidInput(w)
		return
	}

	var verified bool
	var dbOtp string

	err = config.DB.QueryRow(
		`
		SELECT verified FROM User
		WHERE email = ?
	`, req.Email,
	).Scan(&verified)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User does not exist, prompt to sign up
			utils.InvalidInput(w, "Please Signup First")
			return
		}
		utils.InternalServerError(w, "Problem in Identifying User")
		return
	}

	if verified {
		utils.InvalidInput(w, "User already verified. Please login.")
		return
	}

	err = config.DB.QueryRow(
		`
			SELECT otp FROM Otp
			WHERE email = ?
		`, req.Email,
	).Scan(&dbOtp)
	if err != nil {

	}

	if dbOtp != req.Otp {
		utils.InvalidInput(w, "Wrong OTP")
		return
	}

	_, err = config.DB.Exec(
		`
			UPDATE User
			SET verified = ?
			WHERE email = ?
		`, true, req.Email,
	)
	if err != nil {
		utils.InternalServerError(w, "Problems while verifying User")
		return
	}

	config.DB.Exec(
		`
			DELETE FROM Otp
			WHERE email = ?	
		`, req.Email,
	)

	response := map[string]interface{}{
		"success": true,
		"message": "User Verified Successfully, Can Login now",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)

}

// Re-Send OTP for Verification
func ReSendOtp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WrongMethod(w)
		return
	}

	var req dto.ResendOtpRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.InvalidInput(w)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.InvalidInput(w)
		return
	}

	// Check if the user exists and is verified
	var verified bool
	var email string

	err := config.DB.QueryRow(
		`
		SELECT verified, email FROM User
		WHERE email = ?
	`, req.Email,
	).Scan(&verified, &email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User does not exist, prompt to sign up
			utils.InvalidInput(w, "Please Signup First")
			return
		}
		utils.InternalServerError(w, "Problem in Identifying User")
		return
	}

	if verified {
		utils.InvalidInput(w, "User already verified. Please login.")
		return
	}

	newOtp := func() string {
		otp := make([]byte, 3)

		_, err := rand.Read(otp)
		if err != nil {
			return "123456"
		}

		return fmt.Sprintf("%06d", int(otp[0])%10*100000+int(otp[1])%10000+int(otp[2])%1000)
	}()

	// Update the OTP in the database
	_, err = config.DB.Exec(
		`
			UPDATE Otp SET otp = ? WHERE email = ?
		`, newOtp, req.Email,
	)
	if err != nil {
		utils.InternalServerError(w, "Failed to update OTP")
		return
	}

	go func(email, otp string) {
		if err := config.Mailer(email, otp); err != nil {
			fmt.Printf("Failed to send OTP to %s: %v\n", email, err)
		}
	}(req.Email, newOtp)

	response := map[string]interface{}{
		"success": true,
		"message": "OTP re-sent successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
