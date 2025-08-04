package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
	"golang.org/x/crypto/bcrypt"
)

// Login handles user login and returns a JWT token.
// @Summary Login user and return JWT token
// @Description Login user with email and password, returns JWT token on success
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body object.LoginRequest true "User login details"
// @Success 200 {object} object.Response{data=object.TokenResponse} "Login successful"
// @Failure 400 {object} object.Response "Invalid request body or validation error"
// @Failure 401 {object} object.Response "Unauthorized - Invalid email or password"
// @Failure 500 {object} object.Response "Internal server error"
// @Router /auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var user object.Pengguna
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, "Invalid request body", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = user.ValidateFields("Email", "Password")
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	db, err := object.ConnectDb()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Database connection error", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	// Check user credentials
	var storedUser object.Pengguna
	query := "SELECT id, email, nama, password_hash, role, alamat FROM pengguna WHERE email = ?"
	err = db.QueryRow(query, user.Email).Scan(&storedUser.Id, &storedUser.Email, &storedUser.Nama, &storedUser.PasswordHash, &storedUser.Role, &storedUser.Alamat)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusUnauthorized, "Invalid email or password", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Database query error", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(user.Password))
	if err != nil {
		response := object.NewResponse(http.StatusUnauthorized, "Invalid email or password", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Generate a JWT token
	token, err := object.GenerateJWT(storedUser.Id, storedUser.Role)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to generate token", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Login successful", map[string]any{
		"token": token,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
