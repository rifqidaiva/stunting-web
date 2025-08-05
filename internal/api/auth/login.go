package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *loginRequest) validate() error {
	if l.Email == "" {
		return fmt.Errorf("email is required")
	}
	if l.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type loginResponse struct {
	Token string `json:"token"`
}

// # Login handles user login requests
//
// @Summary User login
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param login body loginRequest true "Login request"
// @Success 200 {object} object.Response{data=loginResponse} "Login successful"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, "Invalid request body", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = req.validate()
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
	query := "SELECT id, email, password_hash, role FROM pengguna WHERE email = ?"
	err = db.QueryRow(query, req.Email).Scan(&storedUser.Id, &storedUser.Email, &storedUser.PasswordHash, &storedUser.Role)
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
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.PasswordHash), []byte(req.Password))
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

	response := object.NewResponse(http.StatusOK, "Login successful", loginResponse{
		Token: token,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
