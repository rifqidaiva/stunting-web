package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
	"golang.org/x/crypto/bcrypt"
)

type insertPetugasKesehatanRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nama     string `json:"nama"`
	IdSkpd   string `json:"id_skpd"`
}

func (r *insertPetugasKesehatanRequest) validate() error {
	// Email validation: not empty, valid format
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return fmt.Errorf("invalid email format")
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

	// Nama validation: not empty, min 2 chars, only letters and spaces
	if r.Nama == "" {
		return fmt.Errorf("nama is required")
	}
	if len(r.Nama) < 2 || len(r.Nama) > 100 {
		return fmt.Errorf("nama must be between 2-100 characters")
	}
	namaRegex := regexp.MustCompile(`^[a-zA-Z\s]{2,100}$`)
	if !namaRegex.MatchString(r.Nama) {
		return fmt.Errorf("nama must contain only letters and spaces")
	}

	// ID SKPD validation
	if r.IdSkpd == "" {
		return fmt.Errorf("id SKPD is required")
	}

	return nil
}

type insertPetugasKesehatanResponse struct {
	Id string `json:"id"`
}

// # PetugasKesehatanInsert handles inserting new petugas kesehatan data
//
// @Summary Insert new petugas kesehatan
// @Description Insert new petugas kesehatan data (Admin only)
// @Description
// @Description Creates a complete petugas kesehatan record including:
// @Description - Creates pengguna account with email and hashed password
// @Description - Creates petugas_kesehatan record linked to SKPD
// @Description - Sets role as "petugas kesehatan" in pengguna table
// @Description - Validates email uniqueness and SKPD existence
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param petugas body insertPetugasKesehatanRequest true "Petugas kesehatan data"
// @Success 200 {object} object.Response{data=insertPetugasKesehatanResponse} "Petugas kesehatan inserted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/petugas-kesehatan/insert [post]
func PetugasKesehatanInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Extract and validate JWT token
	authHeader := r.Header.Get("Authorization")
	token, err := object.GetJWTFromHeader(authHeader)
	if err != nil {
		response := object.NewResponse(http.StatusUnauthorized, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	userId, role, err := object.ParseJWT(token)
	if err != nil {
		response := object.NewResponse(http.StatusUnauthorized, "Invalid or expired token", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if user is admin
	if role != "admin" {
		response := object.NewResponse(http.StatusForbidden, "Access denied. Admin role required", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Parse request body
	var req insertPetugasKesehatanRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, "Invalid request body", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Validate request
	err = req.validate()
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Connect to database
	db, err := object.ConnectDb()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Database connection error", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	// Check if email already exists in pengguna table
	var emailExists int
	checkEmailQuery := "SELECT COUNT(*) FROM pengguna WHERE email = ?"
	err = db.QueryRow(checkEmailQuery, req.Email).Scan(&emailExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check email existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if emailExists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "Email already exists", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if SKPD exists and not soft deleted
	var skpdExists int
	var skpdName, jenisSkpd string
	checkSkpdQuery := "SELECT COUNT(*), skpd, jenis FROM skpd WHERE id = ? AND deleted_date IS NULL GROUP BY skpd, jenis"
	err = db.QueryRow(checkSkpdQuery, req.IdSkpd).Scan(&skpdExists, &skpdName, &jenisSkpd)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusBadRequest, "SKPD not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check SKPD existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if skpdExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "SKPD not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if nama already exists in the same SKPD (prevent duplicate names in same organization)
	var nameExists int
	checkNameQuery := "SELECT COUNT(*) FROM petugas_kesehatan WHERE nama = ? AND id_skpd = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkNameQuery, req.Nama, req.IdSkpd).Scan(&nameExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check name existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if nameExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Petugas kesehatan with name '%s' already exists in %s '%s'", req.Nama, jenisSkpd, skpdName), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to hash password", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to begin transaction", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer tx.Rollback()

	// Insert into pengguna table
	insertPenggunaQuery := "INSERT INTO pengguna (email, password_hash, role) VALUES (?, ?, ?)"
	resultPengguna, err := tx.Exec(insertPenggunaQuery, req.Email, hashedPassword, "petugas kesehatan")
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert pengguna", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get the inserted pengguna ID
	penggunaId, err := resultPengguna.LastInsertId()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to retrieve pengguna ID", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Insert into petugas_kesehatan table
	insertPetugasQuery := `INSERT INTO petugas_kesehatan 
        (id_pengguna, id_skpd, nama, created_id, created_date) 
        VALUES (?, ?, ?, ?, ?)`

	resultPetugas, err := tx.Exec(insertPetugasQuery,
		penggunaId,
		req.IdSkpd,
		req.Nama,
		userId,
		currentTime,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert petugas kesehatan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get the inserted petugas kesehatan ID
	petugasId, err := resultPetugas.LastInsertId()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to retrieve petugas kesehatan ID", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to commit transaction", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare success response with additional info
	message := fmt.Sprintf("Petugas kesehatan '%s' successfully created for %s '%s'", req.Nama, jenisSkpd, skpdName)

	response := object.NewResponse(http.StatusOK, message, insertPetugasKesehatanResponse{
		Id: strconv.FormatInt(petugasId, 10),
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
