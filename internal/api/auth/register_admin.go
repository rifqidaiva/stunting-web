package auth

import (
	"encoding/json"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
	"golang.org/x/crypto/bcrypt"
)

// RegisterAdmin handles admin registration.
func RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var admin object.Pengguna
	err := json.NewDecoder(r.Body).Decode(&admin)
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, "Invalid request body", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	admin.Role = "admin" // Set role to admin

	err = admin.ValidateFields("Email", "Nama", "Password", "Alamat")
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

	query := "INSERT INTO pengguna (email, nama, password_hash, role, alamat) VALUES (?, ?, ?, ?, ?)"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to hash password", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to prepare statement", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(admin.Email, admin.Nama, hashedPassword, admin.Role, admin.Alamat)
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
