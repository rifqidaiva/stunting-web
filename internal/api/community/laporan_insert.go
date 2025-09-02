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

type insertLaporanRequest struct {
	IdBalita              string `json:"id_balita"`
	TanggalLaporan        string `json:"tanggal_laporan"` // Format: YYYY-MM-DD
	HubunganDenganBalita  string `json:"hubungan_dengan_balita"`
	NomorHpPelapor        string `json:"nomor_hp_pelapor"`
	NomorHpKeluargaBalita string `json:"nomor_hp_keluarga_balita"`
}

func (r *insertLaporanRequest) validate() error {
	// ID Balita validation (wajib)
	if r.IdBalita == "" {
		return fmt.Errorf("id balita is required")
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

	// Check if date is not older than 1 year (business rule for community reports)
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	if laporanDate.Before(oneYearAgo) {
		return fmt.Errorf("tanggal laporan cannot be older than 1 year")
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

type insertLaporanResponse struct {
	Id string `json:"id"`
}

// # LaporanInsert handles inserting new laporan for masyarakat
//
// @Summary Insert new laporan (Community)
// @Description Insert new laporan for community/masyarakat users
// @Description
// @Description This endpoint allows masyarakat users to report balita for stunting assessment.
// @Description The laporan will be automatically set with "Belum diproses" status and linked
// @Description to the reporting masyarakat user.
// @Description
// @Description Validation includes:
// @Description - Balita ownership verification (user can only report balita from their own keluarga)
// @Description - Duplicate prevention (same balita and date)
// @Description - Business rule checks (no pending reports for same balita)
// @Description - Date validation (not future, not older than 1 year)
// @Description - Contact information validation
// @Tags community
// @Accept json
// @Produce json
// @Security Bearer
// @Param laporan body insertLaporanRequest true "Laporan data"
// @Success 200 {object} object.Response{data=insertLaporanResponse} "Laporan inserted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Masyarakat role required or not owner of balita"
// @Failure 409 {object} object.Response{data=nil} "Conflict - Pending report exists for this balita"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/community/laporan/insert [post]
func LaporanInsert(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is masyarakat
	if role != "masyarakat" {
		response := object.NewResponse(http.StatusForbidden, "Access denied. Masyarakat role required", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Parse request body
	var req insertLaporanRequest
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
	var keluargaCreatedId string
	var keluargaId string
	checkBalitaQuery := `
        SELECT COUNT(*), k.created_id, b.id_keluarga
        FROM balita b
        JOIN keluarga k ON b.id_keluarga = k.id
        WHERE b.id = ? AND b.deleted_date IS NULL AND k.deleted_date IS NULL
        GROUP BY k.created_id, b.id_keluarga
    `
	err = db.QueryRow(checkBalitaQuery, req.IdBalita).Scan(&balitaExists, &keluargaCreatedId, &keluargaId)
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

	// Verify ownership - user can only report balita from their own keluarga
	if keluargaCreatedId != userId {
		response := object.NewResponse(http.StatusForbidden, "Access denied. You can only report balita from your own keluarga", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Business rule: Check if balita has pending reports (status "Belum diproses")
	var pendingReportsCount int
	checkPendingQuery := `
        SELECT COUNT(*) 
        FROM laporan_masyarakat lm
        JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        WHERE lm.id_balita = ? 
        AND lm.deleted_date IS NULL 
        AND sl.status = 'Belum diproses'
    `
	err = db.QueryRow(checkPendingQuery, req.IdBalita).Scan(&pendingReportsCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check pending reports", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if pendingReportsCount > 0 {
		response := object.NewResponse(http.StatusConflict,
			"Cannot create new report. There is already a pending report for this balita that needs to be processed first", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for duplicate laporan (same balita, tanggal laporan, dan masyarakat yang sama)
	var duplicateExists int
	checkDuplicateQuery := `SELECT COUNT(*) FROM laporan_masyarakat 
        WHERE id_balita = ? AND tanggal_laporan = ? AND id_masyarakat = ? AND deleted_date IS NULL`
	err = db.QueryRow(checkDuplicateQuery, req.IdBalita, req.TanggalLaporan, masyarakatId).Scan(&duplicateExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check duplicate laporan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if duplicateExists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "You have already reported this balita on the same date", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get default status laporan "Belum diproses"
	var statusLaporanId string
	getStatusQuery := "SELECT id FROM status_laporan WHERE status = 'Belum diproses'"
	err = db.QueryRow(getStatusQuery).Scan(&statusLaporanId)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get default status laporan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Business logic: Check masyarakat report limit per month (optional)
	var monthlyReportsCount int
	currentMonth := time.Now().Format("2006-01")
	checkMonthlyQuery := `SELECT COUNT(*) FROM laporan_masyarakat 
        WHERE id_masyarakat = ? AND DATE_FORMAT(tanggal_laporan, '%Y-%m') = ? AND deleted_date IS NULL`
	err = db.QueryRow(checkMonthlyQuery, masyarakatId, currentMonth).Scan(&monthlyReportsCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check monthly reports", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Optional: Limit maximum reports per month (example: 10 reports)
	const maxReportsPerMonth = 10
	if monthlyReportsCount >= maxReportsPerMonth {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Monthly report limit reached. You have already submitted %d reports this month (max: %d)",
				monthlyReportsCount, maxReportsPerMonth), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Insert laporan with masyarakat as reporter
	insertQuery := `INSERT INTO laporan_masyarakat 
        (id_masyarakat, id_balita, id_status_laporan, tanggal_laporan, hubungan_dengan_balita, 
        nomor_hp_pelapor, nomor_hp_keluarga_balita, created_id, created_date) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(insertQuery,
		masyarakatId,
		req.IdBalita,
		statusLaporanId,
		req.TanggalLaporan,
		req.HubunganDenganBalita,
		req.NomorHpPelapor,
		req.NomorHpKeluargaBalita,
		userId, // Use user ID (from JWT) as created_id
		currentTime,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert laporan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get the inserted ID
	insertedId, err := result.LastInsertId()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to retrieve inserted ID", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Laporan inserted successfully", insertLaporanResponse{
		Id: strconv.FormatInt(insertedId, 10),
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
