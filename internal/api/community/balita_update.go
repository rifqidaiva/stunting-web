package community

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type updateBalitaRequest struct {
	Id           string `json:"id"`
	IdKeluarga   string `json:"id_keluarga"`
	Nama         string `json:"nama"`
	TanggalLahir string `json:"tanggal_lahir"` // Format: YYYY-MM-DD
	JenisKelamin string `json:"jenis_kelamin"` // "L" or "P"
	BeratLahir   string `json:"berat_lahir"`   // in grams
	TinggiLahir  string `json:"tinggi_lahir"`  // in cm
}

func (r *updateBalitaRequest) validate() error {
	// ID validation
	if r.Id == "" {
		return fmt.Errorf("balita ID is required")
	}

	// ID Keluarga validation
	if r.IdKeluarga == "" {
		return fmt.Errorf("id keluarga is required")
	}

	// Nama validation: 2-50 characters, only letters and spaces
	if r.Nama == "" {
		return fmt.Errorf("nama is required")
	}
	namaRegex := regexp.MustCompile(`^[a-zA-Z\s]{2,50}$`)
	if !namaRegex.MatchString(r.Nama) {
		return fmt.Errorf("nama must be 2-50 characters and contain only letters and spaces")
	}

	// Tanggal Lahir validation: YYYY-MM-DD format
	if r.TanggalLahir == "" {
		return fmt.Errorf("tanggal lahir is required")
	}
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(r.TanggalLahir) {
		return fmt.Errorf("tanggal lahir must be in YYYY-MM-DD format")
	}

	// Parse and validate the date
	_, err := time.Parse("2006-01-02", r.TanggalLahir)
	if err != nil {
		return fmt.Errorf("invalid tanggal lahir format")
	}

	// Check if date is not in the future
	birthDate, _ := time.Parse("2006-01-02", r.TanggalLahir)
	if birthDate.After(time.Now()) {
		return fmt.Errorf("tanggal lahir cannot be in the future")
	}

	// Check if child is not older than 5 years (balita criteria)
	fiveYearsAgo := time.Now().AddDate(-5, 0, 0)
	if birthDate.Before(fiveYearsAgo) {
		return fmt.Errorf("child must be under 5 years old (balita criteria)")
	}

	// Jenis Kelamin validation: L or P
	if r.JenisKelamin == "" {
		return fmt.Errorf("jenis kelamin is required")
	}
	if r.JenisKelamin != "L" && r.JenisKelamin != "P" {
		return fmt.Errorf("jenis kelamin must be 'L' (Laki-laki) or 'P' (Perempuan)")
	}

	// Berat Lahir validation: numeric, reasonable range (500-6000 grams)
	if r.BeratLahir == "" {
		return fmt.Errorf("berat lahir is required")
	}
	beratLahir, err := strconv.Atoi(r.BeratLahir)
	if err != nil {
		return fmt.Errorf("berat lahir must be a valid number (in grams)")
	}
	if beratLahir < 500 || beratLahir > 6000 {
		return fmt.Errorf("berat lahir must be between 500-6000 grams")
	}

	// Tinggi Lahir validation: numeric, reasonable range (25-65 cm)
	if r.TinggiLahir == "" {
		return fmt.Errorf("tinggi lahir is required")
	}
	tinggiLahir, err := strconv.Atoi(r.TinggiLahir)
	if err != nil {
		return fmt.Errorf("tinggi lahir must be a valid number (in cm)")
	}
	if tinggiLahir < 25 || tinggiLahir > 65 {
		return fmt.Errorf("tinggi lahir must be between 25-65 cm")
	}

	return nil
}

type updateBalitaResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # BalitaUpdate handles updating balita data for masyarakat
//
// @Summary Update balita data (Community)
// @Description Update existing balita data for community/masyarakat users
// @Description
// @Description This endpoint allows masyarakat users to update balita data
// @Description that they have previously created. Users can only update balita
// @Description from keluarga they own and only if there are no active reports.
// @Description
// @Description Validation includes:
// @Description - Balita ownership verification (through keluarga ownership)
// @Description - Age criteria check (must be under 5 years old for balita classification)
// @Description - Birth data validation (weight, height within reasonable ranges)
// @Description - Duplicate prevention (same name and birth date in same keluarga)
// @Description - Business rule checks (no active reports constraint)
// @Description - Format validation for all fields
// @Tags community
// @Accept json
// @Produce json
// @Security Bearer
// @Param balita body updateBalitaRequest true "Updated balita data"
// @Success 200 {object} object.Response{data=updateBalitaResponse} "Balita updated successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Masyarakat role required or not owner"
// @Failure 404 {object} object.Response{data=nil} "Balita not found"
// @Failure 409 {object} object.Response{data=nil} "Conflict - Cannot update due to active reports"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/community/balita/update [put]
func BalitaUpdate(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is masyarakat
	if role != "masyarakat" {
		response := object.NewResponse(http.StatusForbidden, "Access denied. Masyarakat role required", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Parse request body
	var req updateBalitaRequest
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

	// Verify user exists and get masyarakat ID
	var masyarakatId string
	checkUserQuery := "SELECT m.id FROM masyarakat m JOIN pengguna p ON m.id_pengguna = p.id WHERE p.id = ?"
	err = db.QueryRow(checkUserQuery, userId).Scan(&masyarakatId)
	if err != nil {
		response := object.NewResponse(http.StatusUnauthorized, "Masyarakat profile not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if balita exists, not soft deleted, and verify ownership through keluarga
	var balitaExists int
	var currentKeluargaId string
	var keluargaCreatedId string
	var deletedDate sql.NullString
	checkBalitaQuery := `
        SELECT COUNT(*), b.id_keluarga, k.created_id, b.deleted_date
        FROM balita b
        JOIN keluarga k ON b.id_keluarga = k.id
        WHERE b.id = ? AND k.deleted_date IS NULL
        GROUP BY b.id_keluarga, k.created_id, b.deleted_date
    `
	err = db.QueryRow(checkBalitaQuery, req.Id).Scan(&balitaExists, &currentKeluargaId, &keluargaCreatedId, &deletedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Balita not found or keluarga has been deleted", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check balita existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if balitaExists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Balita not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if balita is already soft deleted
	if deletedDate.Valid {
		response := object.NewResponse(http.StatusNotFound, "Balita has been deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Verify ownership - user can only update balita from their own keluarga
	if keluargaCreatedId != userId {
		response := object.NewResponse(http.StatusForbidden, "Access denied. You can only update balita from your own keluarga", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Business rule: Check if balita has active reports that prevent updates
	var activeReportsCount int
	checkActiveReportsQuery := `
        SELECT COUNT(*) 
        FROM laporan_masyarakat lm
        JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        WHERE lm.id_balita = ? 
        AND lm.deleted_date IS NULL 
        AND sl.status IN ('Belum diproses', 'Diproses dan data sesuai', 'Belum ditindaklanjuti')
    `
	err = db.QueryRow(checkActiveReportsQuery, req.Id).Scan(&activeReportsCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check active reports", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if activeReportsCount > 0 {
		response := object.NewResponse(http.StatusConflict,
			fmt.Sprintf("Cannot update balita data. There are %d active reports that need to be processed first", activeReportsCount), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Verify that the new keluarga exists, not soft deleted, and is owned by the user
	if req.IdKeluarga != currentKeluargaId {
		var newKeluargaExists int
		var newKeluargaCreatedId string
		checkNewKeluargaQuery := `SELECT COUNT(*), created_id 
            FROM keluarga 
            WHERE id = ? AND deleted_date IS NULL 
            GROUP BY created_id`
		err = db.QueryRow(checkNewKeluargaQuery, req.IdKeluarga).Scan(&newKeluargaExists, &newKeluargaCreatedId)
		if err != nil {
			response := object.NewResponse(http.StatusBadRequest, "New keluarga not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if newKeluargaExists == 0 {
			response := object.NewResponse(http.StatusBadRequest, "New keluarga not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Verify ownership of new keluarga
		if newKeluargaCreatedId != userId {
			response := object.NewResponse(http.StatusForbidden, "Access denied. You can only move balita to your own keluarga", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	// Check for duplicate balita (same name, birth date, and keluarga, excluding current record)
	var duplicateExists int
	checkDuplicateQuery := `SELECT COUNT(*) FROM balita 
        WHERE id_keluarga = ? AND nama = ? AND tanggal_lahir = ? AND id != ? AND deleted_date IS NULL`
	err = db.QueryRow(checkDuplicateQuery, req.IdKeluarga, req.Nama, req.TanggalLahir, req.Id).Scan(&duplicateExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check duplicate balita", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if duplicateExists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "Balita with same name and birth date already exists in this keluarga", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Business logic: Check if new keluarga already has maximum balita (if moving to different keluarga)
	if req.IdKeluarga != currentKeluargaId {
		var currentBalitaCount int
		checkBalitaCountQuery := "SELECT COUNT(*) FROM balita WHERE id_keluarga = ? AND deleted_date IS NULL"
		err = db.QueryRow(checkBalitaCountQuery, req.IdKeluarga).Scan(&currentBalitaCount)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to check balita count", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Optional: Limit maximum balita per keluarga (example: 10 balita)
		const maxBalitaPerKeluarga = 10
		if currentBalitaCount >= maxBalitaPerKeluarga {
			response := object.NewResponse(http.StatusBadRequest,
				fmt.Sprintf("Maximum balita limit reached. Target keluarga already has %d balita (max: %d)",
					currentBalitaCount, maxBalitaPerKeluarga), nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Update balita
	updateQuery := `UPDATE balita SET 
        id_keluarga = ?, 
        nama = ?, 
        tanggal_lahir = ?, 
        jenis_kelamin = ?,
        berat_lahir = ?, 
        tinggi_lahir = ?, 
        updated_id = ?, 
        updated_date = ?
        WHERE id = ? AND deleted_date IS NULL`

	result, err := db.Exec(updateQuery,
		req.IdKeluarga,
		req.Nama,
		req.TanggalLahir,
		req.JenisKelamin,
		req.BeratLahir,
		req.TinggiLahir,
		userId, // Use user ID (from JWT) as updated_id
		currentTime,
		req.Id,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to update balita", nil)
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
		response := object.NewResponse(http.StatusNotFound, "Balita not found or already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Balita updated successfully", updateBalitaResponse{
		Id:      req.Id,
		Message: "Data balita berhasil diperbarui",
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
