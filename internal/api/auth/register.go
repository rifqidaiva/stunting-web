package auth

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// Register handles user registration.
// @Summary Register a new user
// @Description Register a new user with email, name, password, and address
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body object.RegisterRequest true "User registration details"
// @Success 200 {object} object.Response "Registration successful"
// @Failure 400 {object} object.Response "Invalid request body or validation error"
// @Failure 500 {object} object.Response "Internal server error"
// @Router /auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
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

	err = user.ValidateFields("Email", "Nama", "Password", "Alamat")
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

	// Insert user into the database
	query := "INSERT INTO pengguna (email, nama, password_hash, role, alamat) VALUES (?, ?, ?, ?, ?)"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to hash password", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	_, err = db.Exec(query, user.Email, user.Nama, hashedPassword, "masyarakat", user.Alamat)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to register user", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "User registered successfully", nil)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
