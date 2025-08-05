package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
	"golang.org/x/crypto/bcrypt"
)

type registerAdminRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *registerAdminRequest) Validate() error {
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var req registerAdminRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, "Invalid request body", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Validate required fields
	err = req.Validate()
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

	// Check if email already exists
	var exists int
	checkQuery := "SELECT COUNT(*) FROM pengguna WHERE email = ?"
	err = db.QueryRow(checkQuery, req.Email).Scan(&exists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check email", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "Email already registered", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to hash password", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	query := "INSERT INTO pengguna (email, password_hash, role) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to prepare statement", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(req.Email, hashedPassword, "admin")
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to register admin", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	adminID, err := result.LastInsertId()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get last insert ID", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Admin registered successfully", map[string]any{
		"admin_id": adminID,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
