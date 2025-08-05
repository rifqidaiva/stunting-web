package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type registerRequest struct {
	Email    string `json:"email"`
	Nama     string `json:"nama"`
	Password string `json:"password"`
	Alamat   string `json:"alamat"`
}

func (r *registerRequest) validate() error {
	// Email validation: not empty, valid format
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Nama validation: not empty, min 2 chars, only letters and spaces
	if r.Nama == "" {
		return fmt.Errorf("nama is required")
	}
	namaRegex := regexp.MustCompile(`^[a-zA-Z\s]{2,}$`)
	if !namaRegex.MatchString(r.Nama) {
		return fmt.Errorf("nama must be at least 2 characters and contain only letters and spaces")
	}

	// Password validation: not empty, min 8 chars, at least 1 number, 1 letter
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(r.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	passLetter := regexp.MustCompile(`[A-Za-z]`)
	passNumber := regexp.MustCompile(`[0-9]`)
	if !passLetter.MatchString(r.Password) || !passNumber.MatchString(r.Password) {
		return fmt.Errorf("password must contain at least one letter and one number")
	}

	// Alamat validation: not empty, min 5 chars
	if r.Alamat == "" {
		return fmt.Errorf("alamat is required")
	}
	if len(r.Alamat) < 5 {
		return fmt.Errorf("alamat must be at least 5 characters")
	}

	return nil
}

// # Register handles user registration requests
//
// @Summary User registration
// @Description Register a new user with email, nama, password, and alamat
// @Tags auth
// @Accept json
// @Produce json
// @Param register body registerRequest true "Register request"
// @Success 200 {object} object.Response{data=nil} "User registered successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/auth/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var req registerRequest
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
		response := object.NewResponse(http.StatusBadRequest, "Email sudah terdaftar", nil)
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

	// Insert into pengguna table
	queryPengguna := "INSERT INTO pengguna (email, password_hash, role) VALUES (?, ?, ?)"
	result, err := db.Exec(queryPengguna, req.Email, hashedPassword, "masyarakat")
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to register user", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get the inserted pengguna id
	penggunaId, err := result.LastInsertId()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to retrieve user ID", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Insert into masyarakat table
	queryMasyarakat := "INSERT INTO masyarakat (id_pengguna, nama, alamat) VALUES (?, ?, ?)"
	_, err = db.Exec(queryMasyarakat, penggunaId, req.Nama, req.Alamat)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to register masyarakat", nil)
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
