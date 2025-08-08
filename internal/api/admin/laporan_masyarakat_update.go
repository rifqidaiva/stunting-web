package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type updateLaporanMasyarakatRequest struct {
	Id                    string `json:"id"`
	IdMasyarakat          string `json:"id_masyarakat"` // dapat null jika laporan dari admin
	IdBalita              string `json:"id_balita"`
	IdStatusLaporan       string `json:"id_status_laporan"`
	TanggalLaporan        string `json:"tanggal_laporan"` // Format: YYYY-MM-DD
	HubunganDenganBalita  string `json:"hubungan_dengan_balita"`
	NomorHpPelapor        string `json:"nomor_hp_pelapor"`
	NomorHpKeluargaBalita string `json:"nomor_hp_keluarga_balita"`
}

func (r *updateLaporanMasyarakatRequest) validate() error {
	// ID validation
	if r.Id == "" {
		return fmt.Errorf("laporan masyarakat ID is required")
	}

	// ID Balita validation (wajib)
	if r.IdBalita == "" {
		return fmt.Errorf("id balita is required")
	}

	// ID Status Laporan validation (wajib)
	if r.IdStatusLaporan == "" {
		return fmt.Errorf("id status laporan is required")
	}

	// Tanggal Laporan validation: YYYY-MM-DD format
	if r.TanggalLaporan == "" {
		return fmt.Errorf("tanggal laporan is required")
	}
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(r.TanggalLaporan) {
		return fmt.Errorf("tanggal laporan must be in YYYY-MM-DD format")
	}

	// Parse and validate the date
	_, err := time.Parse("2006-01-02", r.TanggalLaporan)
	if err != nil {
		return fmt.Errorf("invalid tanggal laporan format")
	}

	// Check if date is not in the future
	laporanDate, _ := time.Parse("2006-01-02", r.TanggalLaporan)
	if laporanDate.After(time.Now()) {
		return fmt.Errorf("tanggal laporan cannot be in the future")
	}

	// Hubungan Dengan Balita validation
	if r.HubunganDenganBalita == "" {
		return fmt.Errorf("hubungan dengan balita is required")
	}
	if len(r.HubunganDenganBalita) < 2 || len(r.HubunganDenganBalita) > 50 {
		return fmt.Errorf("hubungan dengan balita must be between 2-50 characters")
	}
	// Validasi format hubungan (contoh: orang tua, kerabat, tetangga, dll)
	hubunganRegex := regexp.MustCompile(`^[a-zA-Z\s]{2,50}$`)
	if !hubunganRegex.MatchString(r.HubunganDenganBalita) {
		return fmt.Errorf("hubungan dengan balita must contain only letters and spaces")
	}

	// Nomor HP Pelapor validation
	if r.NomorHpPelapor == "" {
		return fmt.Errorf("nomor HP pelapor is required")
	}
	// Validasi format nomor HP Indonesia (08xxxxxxxxx atau +628xxxxxxxxx)
	hpRegex := regexp.MustCompile(`^(\+628|08)\d{8,11}$`)
	if !hpRegex.MatchString(r.NomorHpPelapor) {
		return fmt.Errorf("nomor HP pelapor must be valid Indonesian phone number (08xxxxxxxxx or +628xxxxxxxxx)")
	}

	// Nomor HP Keluarga Balita validation
	if r.NomorHpKeluargaBalita == "" {
		return fmt.Errorf("nomor HP keluarga balita is required")
	}
	if !hpRegex.MatchString(r.NomorHpKeluargaBalita) {
		return fmt.Errorf("nomor HP keluarga balita must be valid Indonesian phone number (08xxxxxxxxx or +628xxxxxxxxx)")
	}

	return nil
}

type updateLaporanMasyarakatResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # UpdateLaporanMasyarakat handles updating laporan masyarakat data
//
// @Summary Update laporan masyarakat data
// @Description Update existing laporan masyarakat data (Admin only)
// @Description
// @Description Updates laporan masyarakat record with new data including:
// @Description - id_masyarakat (optional, null if admin report), id_balita, id_status_laporan
// @Description - tanggal_laporan, hubungan_dengan_balita, nomor_hp_pelapor, nomor_hp_keluarga_balita
// @Description - Validates balita existence, status laporan, masyarakat (if provided), and prevents duplicates
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param laporan body updateLaporanMasyarakatRequest true "Updated laporan masyarakat data"
// @Success 200 {object} object.Response{data=updateLaporanMasyarakatResponse} "Laporan masyarakat updated successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Laporan masyarakat not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/laporan-masyarakat/update [put]
func AdminLaporanMasyarakatUpdate(w http.ResponseWriter, r *http.Request) {
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
	var req updateLaporanMasyarakatRequest
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

	// Check if laporan masyarakat exists and not soft deleted
	var exists int
	checkExistQuery := "SELECT COUNT(*) FROM laporan_masyarakat WHERE id = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkExistQuery, req.Id).Scan(&exists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check laporan masyarakat existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Laporan masyarakat not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if balita exists and not soft deleted
	var balitaExists int
	checkBalitaQuery := "SELECT COUNT(*) FROM balita WHERE id = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkBalitaQuery, req.IdBalita).Scan(&balitaExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check balita existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if balitaExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Balita not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if status laporan exists
	var statusExists int
	checkStatusQuery := "SELECT COUNT(*) FROM status_laporan WHERE id = ?"
	err = db.QueryRow(checkStatusQuery, req.IdStatusLaporan).Scan(&statusExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check status laporan existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if statusExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Status laporan not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if masyarakat exists (if provided)
	if req.IdMasyarakat != "" {
		var masyarakatExists int
		checkMasyarakatQuery := "SELECT COUNT(*) FROM masyarakat WHERE id = ?"
		err = db.QueryRow(checkMasyarakatQuery, req.IdMasyarakat).Scan(&masyarakatExists)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to check masyarakat existence", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if masyarakatExists == 0 {
			response := object.NewResponse(http.StatusBadRequest, "Masyarakat not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	// Check for duplicate laporan (same balita, tanggal laporan yang sama, dan pelapor yang sama, excluding current record)
	var duplicateExists int
	var checkDuplicateQuery string
	if req.IdMasyarakat != "" {
		checkDuplicateQuery = `SELECT COUNT(*) FROM laporan_masyarakat 
            WHERE id_balita = ? AND tanggal_laporan = ? AND id_masyarakat = ? AND id != ? AND deleted_date IS NULL`
		err = db.QueryRow(checkDuplicateQuery, req.IdBalita, req.TanggalLaporan, req.IdMasyarakat, req.Id).Scan(&duplicateExists)
	} else {
		checkDuplicateQuery = `SELECT COUNT(*) FROM laporan_masyarakat 
            WHERE id_balita = ? AND tanggal_laporan = ? AND id_masyarakat IS NULL AND id != ? AND deleted_date IS NULL`
		err = db.QueryRow(checkDuplicateQuery, req.IdBalita, req.TanggalLaporan, req.Id).Scan(&duplicateExists)
	}

	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check duplicate laporan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if duplicateExists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "Laporan with same balita and date already exists from this reporter", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if laporan has related riwayat pemeriksaan records (warn user)
	var riwayatCount int
	checkRiwayatQuery := `SELECT COUNT(*) FROM riwayat_pemeriksaan 
        WHERE id_balita = ? AND deleted_date IS NULL`
	err = db.QueryRow(checkRiwayatQuery, req.IdBalita).Scan(&riwayatCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related riwayat pemeriksaan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Prepare update query based on whether id_masyarakat is provided
	var updateQuery string
	var result sql.Result

	if req.IdMasyarakat != "" {
		// Update laporan with masyarakat
		updateQuery = `UPDATE laporan_masyarakat SET 
            id_masyarakat = ?, id_balita = ?, id_status_laporan = ?, tanggal_laporan = ?,
            hubungan_dengan_balita = ?, nomor_hp_pelapor = ?, nomor_hp_keluarga_balita = ?,
            updated_id = ?, updated_date = ?
            WHERE id = ? AND deleted_date IS NULL`

		result, err = db.Exec(updateQuery,
			req.IdMasyarakat,
			req.IdBalita,
			req.IdStatusLaporan,
			req.TanggalLaporan,
			req.HubunganDenganBalita,
			req.NomorHpPelapor,
			req.NomorHpKeluargaBalita,
			userId,
			currentTime,
			req.Id,
		)
	} else {
		// Update laporan tanpa masyarakat (admin report)
		updateQuery = `UPDATE laporan_masyarakat SET 
            id_masyarakat = NULL, id_balita = ?, id_status_laporan = ?, tanggal_laporan = ?,
            hubungan_dengan_balita = ?, nomor_hp_pelapor = ?, nomor_hp_keluarga_balita = ?,
            updated_id = ?, updated_date = ?
            WHERE id = ? AND deleted_date IS NULL`

		result, err = db.Exec(updateQuery,
			req.IdBalita,
			req.IdStatusLaporan,
			req.TanggalLaporan,
			req.HubunganDenganBalita,
			req.NomorHpPelapor,
			req.NomorHpKeluargaBalita,
			userId,
			currentTime,
			req.Id,
		)
	}

	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to update laporan masyarakat", nil)
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
		response := object.NewResponse(http.StatusNotFound, "Laporan masyarakat not found, already deleted or no changes made", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with warnings if applicable
	message := "Data laporan masyarakat berhasil diperbarui"
	if riwayatCount > 0 {
		message += fmt.Sprintf(" (Note: This balita has %d related riwayat pemeriksaan)", riwayatCount)
	}

	response := object.NewResponse(http.StatusOK, "Laporan masyarakat updated successfully", updateLaporanMasyarakatResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
