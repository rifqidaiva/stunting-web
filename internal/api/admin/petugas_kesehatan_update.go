package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
	"golang.org/x/crypto/bcrypt"
)

type updatePetugasKesehatanRequest struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"` // Required - wajib diisi
	Nama     string `json:"nama"`
	IdSkpd   string `json:"id_skpd"`
}

func (r *updatePetugasKesehatanRequest) validate() error {
	// ID validation
	if r.Id == "" {
		return fmt.Errorf("petugas kesehatan ID is required")
	}

	// Email validation: not empty, valid format
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Password validation: required and must be valid
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

type updatePetugasKesehatanResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # PetugasKesehatanUpdate handles updating petugas kesehatan data
//
// @Summary Update petugas kesehatan data
// @Description Update existing petugas kesehatan data (Admin only)
// @Description
// @Description Updates petugas kesehatan record including:
// @Description - Updates pengguna account (email and password)
// @Description - Updates petugas_kesehatan record (nama and SKPD)
// @Description - Validates email uniqueness (excluding current record)
// @Description - Validates SKPD existence and name uniqueness within SKPD
// @Description - Checks for related intervensi before allowing SKPD change
// @Description - Password is required and will be updated if different from current
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param petugas body updatePetugasKesehatanRequest true "Updated petugas kesehatan data"
// @Success 200 {object} object.Response{data=updatePetugasKesehatanResponse} "Petugas kesehatan updated successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Petugas kesehatan not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/petugas-kesehatan/update [put]
func PetugasKesehatanUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
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
	var req updatePetugasKesehatanRequest
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

	// Check if petugas kesehatan exists and not soft deleted, also get current data
	var exists int
	var currentIdPengguna, currentIdSkpd, currentNama, currentPasswordHash string
	checkExistQuery := `SELECT COUNT(*), pk.id_pengguna, pk.id_skpd, pk.nama, p.password_hash 
        FROM petugas_kesehatan pk
        JOIN pengguna p ON pk.id_pengguna = p.id
        WHERE pk.id = ? AND pk.deleted_date IS NULL 
        GROUP BY pk.id_pengguna, pk.id_skpd, pk.nama, p.password_hash`
	err = db.QueryRow(checkExistQuery, req.Id).Scan(&exists, &currentIdPengguna, &currentIdSkpd, &currentNama, &currentPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check petugas kesehatan existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if email already exists (excluding current record)
	var emailExists int
	checkEmailQuery := "SELECT COUNT(*) FROM pengguna WHERE email = ? AND id != ?"
	err = db.QueryRow(checkEmailQuery, req.Email, currentIdPengguna).Scan(&emailExists)
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

	// Check if new SKPD exists and not soft deleted
	var newSkpdExists int
	var newSkpdName, newJenisSkpd string
	checkNewSkpdQuery := "SELECT COUNT(*), skpd, jenis FROM skpd WHERE id = ? AND deleted_date IS NULL GROUP BY skpd, jenis"
	err = db.QueryRow(checkNewSkpdQuery, req.IdSkpd).Scan(&newSkpdExists, &newSkpdName, &newJenisSkpd)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusBadRequest, "SKPD not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check new SKPD existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if newSkpdExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "SKPD not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if nama already exists in the new SKPD (excluding current record)
	var nameExists int
	checkNameQuery := "SELECT COUNT(*) FROM petugas_kesehatan WHERE nama = ? AND id_skpd = ? AND id != ? AND deleted_date IS NULL"
	err = db.QueryRow(checkNameQuery, req.Nama, req.IdSkpd, req.Id).Scan(&nameExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check name existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if nameExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Petugas kesehatan with name '%s' already exists in %s '%s'", req.Nama, newJenisSkpd, newSkpdName), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if petugas has related intervensi records when changing SKPD
	var intervensiCount int
	if currentIdSkpd != req.IdSkpd {
		checkIntervensiQuery := "SELECT COUNT(*) FROM intervensi_petugas WHERE id_petugas_kesehatan = ?"
		err = db.QueryRow(checkIntervensiQuery, req.Id).Scan(&intervensiCount)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to check related intervensi", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if intervensiCount > 0 {
			response := object.NewResponse(http.StatusBadRequest,
				fmt.Sprintf("Cannot change SKPD. Petugas kesehatan has %d related intervensi records", intervensiCount), nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	// Check if password is same as current password
	var passwordChanged bool
	err = bcrypt.CompareHashAndPassword([]byte(currentPasswordHash), []byte(req.Password))
	if err != nil {
		// Password is different, needs to be updated
		passwordChanged = true
	} else {
		// Password is the same
		passwordChanged = false
	}

	// Get current SKPD info for response message
	var currentSkpdName, currentJenisSkpd string
	getCurrentSkpdQuery := "SELECT skpd, jenis FROM skpd WHERE id = ?"
	err = db.QueryRow(getCurrentSkpdQuery, currentIdSkpd).Scan(&currentSkpdName, &currentJenisSkpd)
	if err != nil {
		// Not critical, continue with update
		currentSkpdName = "Unknown"
		currentJenisSkpd = "unknown"
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

	// Update pengguna table (email and password if changed)
	var updatePenggunaQuery string
	var updatePenggunaArgs []any

	if passwordChanged {
		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to hash password", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		updatePenggunaQuery = "UPDATE pengguna SET email = ?, password_hash = ? WHERE id = ?"
		updatePenggunaArgs = []any{req.Email, hashedPassword, currentIdPengguna}
	} else {
		// Only update email, keep current password
		updatePenggunaQuery = "UPDATE pengguna SET email = ? WHERE id = ?"
		updatePenggunaArgs = []any{req.Email, currentIdPengguna}
	}

	_, err = tx.Exec(updatePenggunaQuery, updatePenggunaArgs...)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to update pengguna", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Update petugas_kesehatan table
	updatePetugasQuery := `UPDATE petugas_kesehatan SET 
        id_skpd = ?, nama = ?, updated_id = ?, updated_date = ?
        WHERE id = ? AND deleted_date IS NULL`

	result, err := tx.Exec(updatePetugasQuery,
		req.IdSkpd,
		req.Nama,
		userId,
		currentTime,
		req.Id,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to update petugas kesehatan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if any rows were affected in petugas_kesehatan table
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check update result", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if there were actual changes to petugas_kesehatan data
	petugasDataChanged := (currentIdSkpd != req.IdSkpd) || (currentNama != req.Nama)

	// If no changes to petugas_kesehatan table but password changed, it's still a valid update
	if rowsAffected == 0 && !passwordChanged {
		response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found, already deleted or no changes made", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// If no changes to petugas_kesehatan but password changed, that's still a successful update
	if rowsAffected == 0 && passwordChanged && !petugasDataChanged {
		// This is fine - only password was updated
	} else if rowsAffected == 0 {
		// No changes at all
		response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found or already deleted", nil)
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

	// Prepare response message with additional information
	message := fmt.Sprintf("Data petugas kesehatan '%s' berhasil diperbarui", req.Nama)

	// Add SKPD change information
	if currentIdSkpd != req.IdSkpd {
		message += fmt.Sprintf(" (SKPD changed from %s '%s' to %s '%s')",
			currentJenisSkpd, currentSkpdName, newJenisSkpd, newSkpdName)
	} else {
		message += fmt.Sprintf(" (SKPD: %s '%s')", newJenisSkpd, newSkpdName)
	}

	// Add password change information
	if passwordChanged {
		message += " (Password updated)"
	} else {
		message += " (Password unchanged)"
	}

	response := object.NewResponse(http.StatusOK, "Petugas kesehatan updated successfully", updatePetugasKesehatanResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
