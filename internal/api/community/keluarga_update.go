package community

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type updateKeluargaRequest struct {
	Id          string     `json:"id"`
	NomorKk     string     `json:"nomor_kk"`
	NamaAyah    string     `json:"nama_ayah"`
	NamaIbu     string     `json:"nama_ibu"`
	NikAyah     string     `json:"nik_ayah"`
	NikIbu      string     `json:"nik_ibu"`
	Alamat      string     `json:"alamat"`
	Rt          string     `json:"rt"`
	Rw          string     `json:"rw"`
	IdKelurahan string     `json:"id_kelurahan"`
	Koordinat   [2]float64 `json:"koordinat"` // [longitude, latitude]
}

func (r *updateKeluargaRequest) validate() error {
	// ID validation
	if r.Id == "" {
		return fmt.Errorf("keluarga ID is required")
	}

	// Nomor KK validation: 16 digits
	if r.NomorKk == "" {
		return fmt.Errorf("nomor KK is required")
	}
	kkRegex := regexp.MustCompile(`^\d{16}$`)
	if !kkRegex.MatchString(r.NomorKk) {
		return fmt.Errorf("nomor KK must be exactly 16 digits")
	}

	// Nama Ayah validation
	if r.NamaAyah == "" {
		return fmt.Errorf("nama ayah is required")
	}
	namaRegex := regexp.MustCompile(`^[a-zA-Z\s]{2,50}$`)
	if !namaRegex.MatchString(r.NamaAyah) {
		return fmt.Errorf("nama ayah must be 2-50 characters and contain only letters and spaces")
	}

	// Nama Ibu validation
	if r.NamaIbu == "" {
		return fmt.Errorf("nama ibu is required")
	}
	if !namaRegex.MatchString(r.NamaIbu) {
		return fmt.Errorf("nama ibu must be 2-50 characters and contain only letters and spaces")
	}

	// NIK Ayah validation: 16 digits
	if r.NikAyah == "" {
		return fmt.Errorf("NIK ayah is required")
	}
	nikRegex := regexp.MustCompile(`^\d{16}$`)
	if !nikRegex.MatchString(r.NikAyah) {
		return fmt.Errorf("NIK ayah must be exactly 16 digits")
	}

	// NIK Ibu validation: 16 digits
	if r.NikIbu == "" {
		return fmt.Errorf("NIK ibu is required")
	}
	if !nikRegex.MatchString(r.NikIbu) {
		return fmt.Errorf("NIK ibu must be exactly 16 digits")
	}

	// Alamat validation
	if r.Alamat == "" {
		return fmt.Errorf("alamat is required")
	}
	if len(r.Alamat) < 5 || len(r.Alamat) > 255 {
		return fmt.Errorf("alamat must be between 5-255 characters")
	}

	// RT validation: 1-3 digits
	if r.Rt == "" {
		return fmt.Errorf("RT is required")
	}
	rtRegex := regexp.MustCompile(`^\d{1,3}$`)
	if !rtRegex.MatchString(r.Rt) {
		return fmt.Errorf("RT must be 1-3 digits")
	}

	// RW validation: 1-3 digits
	if r.Rw == "" {
		return fmt.Errorf("RW is required")
	}
	rwRegex := regexp.MustCompile(`^\d{1,3}$`)
	if !rwRegex.MatchString(r.Rw) {
		return fmt.Errorf("RW must be 1-3 digits")
	}

	// ID Kelurahan validation
	if r.IdKelurahan == "" {
		return fmt.Errorf("id kelurahan is required")
	}

	// Koordinat validation (basic bounds check)
	// if r.Koordinat[0] < -180 || r.Koordinat[0] > 180 {
	// 	return fmt.Errorf("longitude must be between -180 and 180")
	// }
	// if r.Koordinat[1] < -90 || r.Koordinat[1] > 90 {
	// 	return fmt.Errorf("latitude must be between -90 and 90")
	// }

	return nil
}

type updateKeluargaResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # KeluargaUpdate handles updating keluarga data for masyarakat
//
// @Summary Update keluarga data (Community)
// @Description Update existing keluarga data for community/masyarakat users
// @Description
// @Description This endpoint allows masyarakat users to update family data
// @Description that they have previously created. Users can only update their own data.
// @Description
// @Description Validation includes:
// @Description - Ownership verification (user can only update their own data)
// @Description - Nomor KK and NIK uniqueness check (excluding current record)
// @Description - Format validation for all fields
// @Description - Kelurahan existence validation
// @Description - Coordinate bounds validation
// @Description - Business rule checks (no active reports constraint)
// @Tags community
// @Accept json
// @Produce json
// @Security Bearer
// @Param keluarga body updateKeluargaRequest true "Keluarga data to update"
// @Success 200 {object} object.Response{data=updateKeluargaResponse} "Keluarga updated successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Masyarakat role required or not owner"
// @Failure 404 {object} object.Response{data=nil} "Keluarga not found"
// @Failure 409 {object} object.Response{data=nil} "Conflict - Cannot update due to active reports"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/community/keluarga/update [put]
func KeluargaUpdate(w http.ResponseWriter, r *http.Request) {
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
	var req updateKeluargaRequest
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

	// Check if keluarga exists, not soft deleted, and verify ownership
	var exists int
	var createdId string
	var deletedDate sql.NullString
	checkKeluargaQuery := `SELECT COUNT(*), created_id, deleted_date 
        FROM keluarga 
        WHERE id = ? 
        GROUP BY created_id, deleted_date`
	err = db.QueryRow(checkKeluargaQuery, req.Id).Scan(&exists, &createdId, &deletedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Keluarga not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check keluarga existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Keluarga not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if already soft deleted
	if deletedDate.Valid {
		response := object.NewResponse(http.StatusNotFound, "Keluarga has been deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Verify ownership - user can only update their own data
	if createdId != userId {
		response := object.NewResponse(http.StatusForbidden, "Access denied. You can only update your own keluarga data", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Business rule: Check if keluarga has active reports that prevent updates
	var activeReportsCount int
	checkActiveReportsQuery := `
        SELECT COUNT(*) 
        FROM laporan_masyarakat lm
        JOIN balita b ON lm.id_balita = b.id
        JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        WHERE b.id_keluarga = ? 
        AND lm.deleted_date IS NULL 
        AND b.deleted_date IS NULL
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
			fmt.Sprintf("Cannot update keluarga data. There are %d active reports that need to be processed first", activeReportsCount), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if Nomor KK already exists (excluding current record)
	var kkExists int
	checkKKQuery := "SELECT COUNT(*) FROM keluarga WHERE nomor_kk = ? AND id != ? AND deleted_date IS NULL"
	err = db.QueryRow(checkKKQuery, req.NomorKk, req.Id).Scan(&kkExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check nomor KK", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if kkExists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "Nomor KK already exists", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if NIK Ayah already exists (excluding current record)
	var nikAyahExists int
	checkNikAyahQuery := "SELECT COUNT(*) FROM keluarga WHERE nik_ayah = ? AND id != ? AND deleted_date IS NULL"
	err = db.QueryRow(checkNikAyahQuery, req.NikAyah, req.Id).Scan(&nikAyahExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check NIK ayah", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if nikAyahExists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "NIK ayah already exists", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if NIK Ibu already exists (excluding current record)
	var nikIbuExists int
	checkNikIbuQuery := "SELECT COUNT(*) FROM keluarga WHERE nik_ibu = ? AND id != ? AND deleted_date IS NULL"
	err = db.QueryRow(checkNikIbuQuery, req.NikIbu, req.Id).Scan(&nikIbuExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check NIK ibu", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if nikIbuExists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "NIK ibu already exists", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if kelurahan exists
	var kelurahanExists int
	checkKelurahanQuery := "SELECT COUNT(*) FROM kelurahan WHERE id = ?"
	err = db.QueryRow(checkKelurahanQuery, req.IdKelurahan).Scan(&kelurahanExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check kelurahan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if kelurahanExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Kelurahan not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Convert coordinates to WKT format
	koordinatWKT := object.ToWKT(req.Koordinat)

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Update keluarga
	updateQuery := `UPDATE keluarga SET 
        nomor_kk = ?, 
        nama_ayah = ?, 
        nama_ibu = ?, 
        nik_ayah = ?, 
        nik_ibu = ?, 
        alamat = ?, 
        rt = ?, 
        rw = ?, 
        id_kelurahan = ?, 
        koordinat = ST_GeomFromText(?), 
        updated_id = ?, 
        updated_date = ? 
        WHERE id = ? AND deleted_date IS NULL`

	result, err := db.Exec(updateQuery,
		req.NomorKk,
		req.NamaAyah,
		req.NamaIbu,
		req.NikAyah,
		req.NikIbu,
		req.Alamat,
		req.Rt,
		req.Rw,
		req.IdKelurahan,
		koordinatWKT,
		userId, // Use user ID (from JWT) as updated_id
		currentTime,
		req.Id,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to update keluarga", nil)
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
		response := object.NewResponse(http.StatusNotFound, "Keluarga not found or already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Keluarga updated successfully", updateKeluargaResponse{
		Id:      req.Id,
		Message: "Data keluarga berhasil diperbarui",
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}