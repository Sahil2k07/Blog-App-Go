package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Sahil2k07/Blog-App-Go/src/config"
	"github.com/Sahil2k07/Blog-App-Go/src/database"
	"github.com/Sahil2k07/Blog-App-Go/src/dto"
	"github.com/Sahil2k07/Blog-App-Go/src/middlewares"
	"github.com/Sahil2k07/Blog-App-Go/src/utils"
	"github.com/go-playground/validator/v10"
)

func parseQueryParam(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}

	if value, err := strconv.Atoi(param); err == nil {
		return value
	}

	return defaultValue
}

// Single Blog
func GetBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WrongMethod(w)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		utils.InvalidInput(w, "Blog id not found")
		return
	}

	var blog database.Blog
	var tagsJSON, createdAtBytes, updatedAtBytes []byte

	err := config.DB.QueryRow(
		`
			SELECT id, profileId, title, content, tags, published, createdAt, updatedAt
			FROM Blog
			WHERE id = ? AND published = true	
		`, id,
	).Scan(&blog.ID, &blog.ProfileID, &blog.Title, &blog.Content, &tagsJSON, &blog.Published, &createdAtBytes, &updatedAtBytes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.InvalidInput(w, "No Blog with the id")
			return
		}

		utils.InternalServerError(w, "Somethng went wrong while loading Blog Post")
		return
	}

	if err := json.Unmarshal(tagsJSON, &blog.Tags); err != nil {
		utils.InternalServerError(w, "Error parsing Blog Details")
		return
	}

	// Parse createdAt Time
	createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
	if err != nil {
		utils.InternalServerError(w, "Error parsing Blog Details")
		return
	}
	blog.CreatedAt = createdAt

	// Parse updatedAt Time
	updatedAt, err := time.Parse("2006-01-02 15:04:05", string(updatedAtBytes))
	if err != nil {
		utils.InternalServerError(w, "Error parsing Blog Details")
		return
	}
	blog.UpdatedAt = updatedAt

	response := map[string]interface{}{
		"success": true,
		"message": "Got Blog Successfully",
		"data":    blog,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// All Blogs
func GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WrongMethod(w)
		return
	}

	const limit = 25
	defaultOffset := 0

	offset := parseQueryParam(r.URL.Query().Get("offset"), defaultOffset)

	rows, err := config.DB.Query(
		`
			SELECT id, profileId, title, content, tags, published, createdAt, updatedAt
		 	FROM Blog
			WHERE published = true
			ORDER BY updatedAt DESC
			LIMIT ? OFFSET ?
		`, limit, offset,
	)
	if err != nil {
		utils.InternalServerError(w, "Failed to get User's Blogs")
		return
	}
	defer rows.Close()

	var blogs []database.Blog

	for rows.Next() {
		var blog database.Blog
		var tagsJSON, createdAtBytes, updatedAtBytes []byte

		err := rows.Scan(&blog.ID, &blog.ProfileID, &blog.Title, &blog.Content, &tagsJSON, &blog.Published, &createdAtBytes, &updatedAtBytes)
		if err != nil {
			utils.InternalServerError(w, "Error parsing Blog Details")
			return
		}

		if err := json.Unmarshal(tagsJSON, &blog.Tags); err != nil {
			utils.InternalServerError(w, "Error parsing Blog Details")
			return
		}

		// Parse createdAt Time
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			utils.InternalServerError(w, "Error parsing Blog Details")
			return
		}
		blog.CreatedAt = createdAt

		// Parse updatedAt Time
		updatedAt, err := time.Parse("2006-01-02 15:04:05", string(updatedAtBytes))
		if err != nil {
			utils.InternalServerError(w, "Error parsing Blog Details")
			return
		}
		blog.UpdatedAt = updatedAt

		blogs = append(blogs, blog)
	}

	// Prepare the response
	response := map[string]interface{}{
		"success": true,
		"message": "Got All Blogs Successfully",
		"data":    blogs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Users Blogs
func GetUserBlogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WrongMethod(w)
		return
	}

	user, ok := r.Context().Value(middlewares.UserContext).(*middlewares.UserAuthDetails)
	if !ok || user == nil {
		utils.UnAuthorized(w, "User is not Authenticated")
		return
	}

	const limit = 25
	defaultOffset := 0

	offset := parseQueryParam(r.URL.Query().Get("offset"), defaultOffset)

	rows, err := config.DB.Query(
		`
			SELECT id, profileId, title, content, tags, published, createdAt, updatedAt
		 	FROM Blog
			WHERE profileId = ?
			LIMIT ? OFFSET ?
		`, user.ProfileId, limit, offset,
	)
	if err != nil {
		utils.InternalServerError(w, "Failed to get User's Blogs")
		return
	}
	defer rows.Close()

	var blogs []database.Blog

	for rows.Next() {
		var blog database.Blog
		var tagsJSON, createdAtBytes, updatedAtBytes []byte

		err := rows.Scan(&blog.ID, &blog.ProfileID, &blog.Title, &blog.Content, &tagsJSON, &blog.Published, &createdAtBytes, &updatedAtBytes)
		if err != nil {
			utils.InternalServerError(w, "Error parsing Blog Details")
			return
		}

		if err := json.Unmarshal(tagsJSON, &blog.Tags); err != nil {
			utils.InternalServerError(w, "Error parsing Blog Details")
			return
		}

		// Parse createdAt Time
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			utils.InternalServerError(w, "Error parsing Blog Details")
			return
		}
		blog.CreatedAt = createdAt

		// Parse updatedAt Time
		updatedAt, err := time.Parse("2006-01-02 15:04:05", string(updatedAtBytes))
		if err != nil {
			utils.InternalServerError(w, "Error parsing Blog Details")
			return
		}
		blog.UpdatedAt = updatedAt

		blogs = append(blogs, blog)
	}

	if len(blogs) == 0 {
		response := map[string]interface{}{
			"success": true,
			"message": "No blogs found from the user",
			"data":    []database.Blog{},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	// Prepare the response
	response := map[string]interface{}{
		"success": true,
		"message": "Got All User's Blogs Successfully",
		"data":    blogs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Delete Blog
func DeleteBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WrongMethod(w)
		return
	}

	user, ok := r.Context().Value(middlewares.UserContext).(*middlewares.UserAuthDetails)
	if !ok || user == nil {
		utils.UnAuthorized(w, "User is not Authenticated")
		return
	}

	id := r.PathValue("id")
	if id == "" {
		utils.InvalidInput(w, "Blog id not found")
		return
	}

	result, err := config.DB.Exec(
		`
			DELETE FROM Blog
			WHERE id = ? AND profileId = ?	
		`, id, user.ProfileId,
	)

	if err != nil {
		utils.InternalServerError(w, "Something went wrong while deleting the Blog")
		return
	}

	// Check if rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.InternalServerError(w, "Error occurred while checking deletion result")
		return
	}

	if rowsAffected == 0 {
		utils.InvalidInput(w, "No Blog with the given ID found")
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Blog deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

// Update Blog
func UpdateBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WrongMethod(w)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		utils.InvalidInput(w, "Blog id not found")
		return
	}

	user, ok := r.Context().Value(middlewares.UserContext).(*middlewares.UserAuthDetails)
	if !ok || user == nil {
		utils.UnAuthorized(w, "User is not Authenticated")
		return
	}

	var req dto.UpdateBlogRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.InvalidInput(w)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.InvalidInput(w)
		return
	}

	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		utils.InternalServerError(w, "Failed to convert tags to JSON")
		return
	}

	result, err := config.DB.Exec(
		`
			UPDATE Blog
			SET title = ?, content = ?, tags = ?
			WHERE id = ? AND profileId = ?
		`, req.Title, req.Content, tagsJSON, id, user.ProfileId,
	)
	if err != nil {
		utils.InternalServerError(w, "Someting went Wrong while Updating Blog")
		return
	}

	// Check if rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.InternalServerError(w, "Error occurred while checking deletion result")
		return
	}

	if rowsAffected == 0 {
		utils.InvalidInput(w, "No Blog with the given ID found")
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Blog updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}

// Create Blog
func CreateBlog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WrongMethod(w)
		return
	}

	user, ok := r.Context().Value(middlewares.UserContext).(*middlewares.UserAuthDetails)
	if !ok || user == nil {
		utils.UnAuthorized(w, "User is not Authenticated")
		return
	}

	var req dto.CreateBlogRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.InvalidInput(w)
		return
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		utils.InvalidInput(w)
		return
	}

	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		utils.InternalServerError(w, "Failed to convert tags to JSON")
		return
	}

	_, e := config.DB.Exec(
		`
			INSERT INTO Blog (id, profileId, title, content, tags)
			VALUES (UUID(), ?, ?, ?, ?)

		`, user.ProfileId, req.Title, req.Content, tagsJSON,
	)
	if e != nil {
		utils.InternalServerError(w, "Error while creating Blog")
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Blog created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)

}
