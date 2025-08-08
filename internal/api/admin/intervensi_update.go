package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type updateIntervensiRequest struct {
	Id        string `json:"id"`
	Jenis     string `json:"jenis"`   // "gizi", "kesehatan", "sosial"
	Tanggal   string `json:"tanggal"` // Format: YYYY-MM-DD
	Deskripsi string `json:"deskripsi"`
	Hasil     string `json:"hasil"`
}

func (r *updateIntervensiRequest) validate() error {
	// ID validation
	if r.Id == "" {
		return fmt.Errorf("intervensi ID is required")
	}

	// Jenis validation: must be one of the allowed types
	if r.Jenis == "" {
		return fmt.Errorf("jenis intervensi is required")
	}
	allowedJenis := []string{"gizi", "kesehatan", "sosial"}
	jenisValid := slices.Contains(allowedJenis, r.Jenis)
	if !jenisValid {
		return fmt.Errorf("jenis must be one of: gizi, kesehatan, sosial")
	}

	// Tanggal validation: YYYY-MM-DD format
	if r.Tanggal == "" {
		return fmt.Errorf("tanggal intervensi is required")
	}
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(r.Tanggal) {
		return fmt.Errorf("tanggal must be in YYYY-MM-DD format")
	}

	// Parse and validate the date
	_, err := time.Parse("2006-01-02", r.Tanggal)
	if err != nil {
		return fmt.Errorf("invalid tanggal format")
	}

	// Check if date is not in the future
	intervensiDate, _ := time.Parse("2006-01-02", r.Tanggal)
	if intervensiDate.After(time.Now()) {
		return fmt.Errorf("tanggal intervensi cannot be in the future")
	}

	// Deskripsi validation
	if r.Deskripsi == "" {
		return fmt.Errorf("deskripsi is required")
	}
	if len(r.Deskripsi) < 10 || len(r.Deskripsi) > 1000 {
		return fmt.Errorf("deskripsi must be between 10-1000 characters")
	}

	// Hasil validation
	if r.Hasil == "" {
		return fmt.Errorf("hasil is required")
	}
	if len(r.Hasil) < 5 || len(r.Hasil) > 500 {
		return fmt.Errorf("hasil must be between 5-500 characters")
	}

	return nil
}

type updateIntervensiResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # UpdateIntervensi handles updating intervensi data
//
// @Summary Update intervensi data
// @Description Update existing intervensi data (Admin only)
// @Description
// @Description Updates intervensi record including:
// @Description - jenis: type of intervention (gizi, kesehatan, sosial)
// @Description - tanggal: intervention date (YYYY-MM-DD format)
// @Description - deskripsi: detailed description of the intervention
// @Description - hasil: results or outcomes of the intervention
// @Description - Validates intervention type and date constraints
// @Description - Checks for related records before allowing changes
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param intervensi body updateIntervensiRequest true "Updated intervensi data"
// @Success 200 {object} object.Response{data=updateIntervensiResponse} "Intervensi updated successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Intervensi not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/intervensi/update [put]
func AdminIntervensiUpdate(w http.ResponseWriter, r *http.Request) {
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
	var req updateIntervensiRequest
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

	// Check if intervensi exists and not soft deleted, also get current data
	var exists int
	var currentJenis, currentTanggal, currentDeskripsi, currentHasil string
	checkExistQuery := `SELECT COUNT(*), jenis, tanggal, deskripsi, hasil 
        FROM intervensi 
        WHERE id = ? AND deleted_date IS NULL 
        GROUP BY jenis, tanggal, deskripsi, hasil`
	err = db.QueryRow(checkExistQuery, req.Id).Scan(&exists, &currentJenis, &currentTanggal, &currentDeskripsi, &currentHasil)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Intervensi not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check intervensi existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Intervensi not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for duplicate intervensi (same jenis, tanggal, and similar deskripsi, excluding current record)
	var duplicateExists int
	checkDuplicateQuery := `SELECT COUNT(*) FROM intervensi 
        WHERE jenis = ? AND tanggal = ? AND deskripsi = ? AND id != ? AND deleted_date IS NULL`
	err = db.QueryRow(checkDuplicateQuery, req.Jenis, req.Tanggal, req.Deskripsi, req.Id).Scan(&duplicateExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check duplicate intervensi", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if duplicateExists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "Intervensi with same type, date, and description already exists", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if intervensi has related petugas assigned
	var petugasCount int
	checkPetugasQuery := "SELECT COUNT(*) FROM intervensi_petugas WHERE id_intervensi = ?"
	err = db.QueryRow(checkPetugasQuery, req.Id).Scan(&petugasCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related petugas", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if intervensi has related riwayat pemeriksaan records
	var riwayatCount int
	checkRiwayatQuery := "SELECT COUNT(*) FROM riwayat_pemeriksaan WHERE id_intervensi = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkRiwayatQuery, req.Id).Scan(&riwayatCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related riwayat pemeriksaan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Update intervensi
	updateQuery := `UPDATE intervensi SET 
        jenis = ?, tanggal = ?, deskripsi = ?, hasil = ?, updated_id = ?, updated_date = ?
        WHERE id = ? AND deleted_date IS NULL`

	result, err := db.Exec(updateQuery,
		req.Jenis,
		req.Tanggal,
		req.Deskripsi,
		req.Hasil,
		userId,
		currentTime,
		req.Id,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to update intervensi", nil)
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
		response := object.NewResponse(http.StatusNotFound, "Intervensi not found, already deleted or no changes made", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with additional information
	message := fmt.Sprintf("Intervensi %s berhasil diperbarui untuk tanggal %s", req.Jenis, req.Tanggal)

	// Add information about changes made
	var changes []string
	if currentJenis != req.Jenis {
		changes = append(changes, fmt.Sprintf("jenis: %s → %s", currentJenis, req.Jenis))
	}
	if currentTanggal != req.Tanggal {
		changes = append(changes, fmt.Sprintf("tanggal: %s → %s", currentTanggal, req.Tanggal))
	}
	if currentDeskripsi != req.Deskripsi {
		changes = append(changes, "deskripsi updated")
	}
	if currentHasil != req.Hasil {
		changes = append(changes, "hasil updated")
	}

	if len(changes) > 0 {
		message += " (Changes made)"
	}

	// Add warnings about related records
	if petugasCount > 0 {
		message += fmt.Sprintf(" (Note: %d petugas assigned to this intervensi)", petugasCount)
	}
	if riwayatCount > 0 {
		message += fmt.Sprintf(" (Note: %d riwayat pemeriksaan related to this intervensi)", riwayatCount)
	}

	response := object.NewResponse(http.StatusOK, "Intervensi updated successfully", updateIntervensiResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
