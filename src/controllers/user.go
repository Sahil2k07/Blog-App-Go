package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Sahil2k07/Blog-App-Go/src/config"
	"github.com/Sahil2k07/Blog-App-Go/src/database"
	"github.com/Sahil2k07/Blog-App-Go/src/dto"
	"github.com/Sahil2k07/Blog-App-Go/src/middlewares"
	"github.com/Sahil2k07/Blog-App-Go/src/utils"

	"github.com/go-playground/validator/v10"
)

// Update the Users Profile
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WrongMethod(w)
		return
	}

	user, ok := r.Context().Value(middlewares.UserContext).(*middlewares.UserAuthDetails)
	if !ok || user == nil {
		utils.UnAuthorized(w, "User is not authenticated")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.InvalidInput(w, "Failed to parse form data. Ensure the size limit is correct.")
		return
	}

	var req dto.UpdateProfileRequest
	req.FirstName = r.FormValue("firstName")
	req.LastName = r.FormValue("lastName")

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.InvalidInput(w, "Invalid input for profile update")
		return
	}

	var profile database.Profile

	err := config.DB.QueryRow(
		`
			SELECT firstName, lastName, image
			FROM Profile
			WHERE userId = ?
		`, user.Id,
	).Scan(&profile.FirstName, &profile.LastName, &profile.Image)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.InvalidInput(w, "Profile not found")
			return
		}
		utils.InternalServerError(w, "Error retrieving profile information")
		return
	}

	// Handle image file upload, if provided
	var imageUrl string

	imageFile, _, err := r.FormFile("image")
	if err == nil {
		imageUrl, err = config.Cloudinary(imageFile)
		if err != nil {
			utils.InternalServerError(w, "Failed to upload image to Cloudinary")
			return
		}
	} else if err != http.ErrMissingFile {
		utils.InvalidInput(w, "Failed to retrieve image")
		return
	} else {
		imageUrl = profile.Image
	}

	// Use existing values if no new input is provided
	if req.FirstName == "" {
		req.FirstName = profile.FirstName
	}
	if req.LastName == "" {
		req.LastName = profile.LastName
	}

	_, err = config.DB.Exec(
		`
			UPDATE Profile
			SET firstName = ?, lastName = ?, image = ?
			WHERE userId = ?
		`, req.FirstName, req.LastName, imageUrl, user.Id,
	)
	if err != nil {
		utils.InternalServerError(w, "Failed to update profile")
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
		"data": map[string]interface{}{
			"userId":    user.Id,
			"firstName": req.FirstName,
			"lastName":  req.LastName,
			"image":     imageUrl,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Get Users Profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WrongMethod(w)
		return
	}

	user, ok := r.Context().Value(middlewares.UserContext).(*middlewares.UserAuthDetails)
	if !ok || user == nil {
		utils.UnAuthorized(w, "User is not Authenticated")
		return
	}

	var profile database.Profile
	var createdAtBytes, updatedAtBytes []byte

	err := config.DB.QueryRow(
		`
			SELECT id, userId, firstName, lastName, image, createdAt, updatedAt
			FROM Profile
			WHERE userId = ?	
		`, user.Id,
	).Scan(&profile.ID, &profile.UserID, &profile.FirstName, &profile.LastName, &profile.Image, &createdAtBytes, &updatedAtBytes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.InvalidInput(w, "User profile data not found")
			return
		}

		utils.InternalServerError(w, "Error getting User's profile")
		return
	}

	// Parse createdAt Time
	createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
	if err != nil {
		utils.InternalServerError(w, "Error parsing Blog Details")
		return
	}
	profile.CreatedAt = createdAt

	// Parse updatedAt Time
	updatedAt, err := time.Parse("2006-01-02 15:04:05", string(updatedAtBytes))
	if err != nil {
		utils.InternalServerError(w, "Error parsing Blog Details")
		return
	}
	profile.UpdatedAt = updatedAt

	response := map[string]interface{}{
		"success": true,
		"message": "Got User's profile successfully",
		"data":    profile,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)

}

// Login the User
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WrongMethod(w)
		return
	}

	var req dto.LoginRequest

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

	var user struct {
		ID             string
		ProfileId      string
		HashedPassword string
		Verified       bool
	}

	err = config.DB.QueryRow(`
		SELECT u.id, u.password, u.verified, p.id AS profileId
		FROM User u
		JOIN Profile p ON u.id = p.userId
		WHERE u.email = ?
	`, req.Email).Scan(&user.ID, &user.HashedPassword, &user.Verified, &user.ProfileId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.InvalidInput(w, "Invalid email or password")
			return
		}
		utils.InternalServerError(w, "Error finding user")
		return
	}

	if !user.Verified {
		utils.InvalidInput(w, "Please verify your email")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.HashedPassword) {
		utils.InvalidInput(w, "Wrong password")
		return
	}

	token, err := utils.GenerateJWT(user.ID, req.Email, user.ProfileId, user.Verified)
	if err != nil {
		utils.InternalServerError(w, "Error generating access token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   3 * 24 * 60 * 60, // 3 days
	})

	response := map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"user": map[string]interface{}{
			"id":        user.ID,
			"email":     req.Email,
			"verified":  user.Verified,
			"profileId": user.ProfileId,
		},
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
