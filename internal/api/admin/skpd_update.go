package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type updateSkpdRequest struct {
	Id    string `json:"id"`
	Skpd  string `json:"skpd"`
	Jenis string `json:"jenis"` // "puskesmas", "kelurahan", "skpd"
}

func (r *updateSkpdRequest) validate() error {
	// ID validation
	if r.Id == "" {
		return fmt.Errorf("SKPD ID is required")
	}

	// SKPD name validation: 2-100 characters
	if r.Skpd == "" {
		return fmt.Errorf("nama SKPD is required")
	}
	if len(r.Skpd) < 2 || len(r.Skpd) > 100 {
		return fmt.Errorf("nama SKPD must be between 2-100 characters")
	}
	// Validasi format nama SKPD (huruf, angka, spasi, dan tanda baca umum)
	skpdRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\.\-\,\(\)\/]{2,100}$`)
	if !skpdRegex.MatchString(r.Skpd) {
		return fmt.Errorf("nama SKPD contains invalid characters")
	}

	// Jenis validation: must be one of the allowed values
	if r.Jenis == "" {
		return fmt.Errorf("jenis SKPD is required")
	}
	allowedJenis := []string{"puskesmas", "kelurahan", "skpd"}
	isValidJenis := slices.Contains(allowedJenis, r.Jenis)
	if !isValidJenis {
		return fmt.Errorf("jenis SKPD must be one of: %v", allowedJenis)
	}

	return nil
}

type updateSkpdResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # UpdateSkpd handles updating SKPD data
//
// @Summary Update SKPD data
// @Description Update existing SKPD data (Admin only)
// @Description
// @Description Updates SKPD record with new data including:
// @Description - skpd (nama SKPD), jenis (puskesmas/kelurahan/skpd)
// @Description - Validates uniqueness of SKPD name within the same jenis (excluding current record)
// @Description - Checks for related petugas kesehatan before allowing jenis change
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param skpd body updateSkpdRequest true "Updated SKPD data"
// @Success 200 {object} object.Response{data=updateSkpdResponse} "SKPD updated successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "SKPD not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/skpd/update [put]
func AdminSkpdUpdate(w http.ResponseWriter, r *http.Request) {
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
	var req updateSkpdRequest
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

	// Check if SKPD exists and not soft deleted
	var exists int
	checkExistQuery := "SELECT COUNT(*) FROM skpd WHERE id = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkExistQuery, req.Id).Scan(&exists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check SKPD existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "SKPD not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get current SKPD data to check for changes
	var currentSkpd, currentJenis string
	getCurrentQuery := "SELECT skpd, jenis FROM skpd WHERE id = ?"
	err = db.QueryRow(getCurrentQuery, req.Id).Scan(&currentSkpd, &currentJenis)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get current SKPD data", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if SKPD name already exists within the same jenis (excluding current record and not soft deleted)
	var duplicateExists int
	checkDuplicateQuery := "SELECT COUNT(*) FROM skpd WHERE skpd = ? AND jenis = ? AND id != ? AND deleted_date IS NULL"
	err = db.QueryRow(checkDuplicateQuery, req.Skpd, req.Jenis, req.Id).Scan(&duplicateExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check SKPD duplication", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if duplicateExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("SKPD '%s' with jenis '%s' already exists", req.Skpd, req.Jenis), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if SKPD has related petugas kesehatan records (warn user if has staff)
	var petugasCount int
	checkPetugasQuery := "SELECT COUNT(*) FROM petugas_kesehatan WHERE id_skpd = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkPetugasQuery, req.Id).Scan(&petugasCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related petugas kesehatan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prevent jenis change if SKPD has related petugas kesehatan and jenis is being changed
	if petugasCount > 0 && currentJenis != req.Jenis {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot change SKPD jenis. There are %d active petugas kesehatan records related to this SKPD", petugasCount), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for similar SKPD names to warn about potential conflicts
	var similarCount int
	similarQuery := "SELECT COUNT(*) FROM skpd WHERE LOWER(skpd) LIKE LOWER(?) AND id != ? AND deleted_date IS NULL"
	similarPattern := "%" + req.Skpd + "%"
	err = db.QueryRow(similarQuery, similarPattern, req.Id).Scan(&similarCount)
	if err != nil {
		// Not critical, continue with update
		similarCount = 0
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Update SKPD
	updateQuery := `UPDATE skpd SET 
        skpd = ?, jenis = ?, updated_id = ?, updated_date = ?
        WHERE id = ? AND deleted_date IS NULL`

	result, err := db.Exec(updateQuery,
		req.Skpd,
		req.Jenis,
		userId,
		currentTime,
		req.Id,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to update SKPD", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check update result", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if rowsAffected == 0 {
		response := object.NewResponse(http.StatusNotFound, "SKPD not found, already deleted or no changes made", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with additional information
	message := "Data SKPD berhasil diperbarui"
	if petugasCount > 0 {
		message += fmt.Sprintf(" (Note: This SKPD has %d related petugas kesehatan)", petugasCount)
	}
	if similarCount > 0 {
		message += fmt.Sprintf(" (Warning: Found %d similar SKPD names)", similarCount)
	}

	response := object.NewResponse(http.StatusOK, "SKPD updated successfully", updateSkpdResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
