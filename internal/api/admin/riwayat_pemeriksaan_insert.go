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
)

type insertRiwayatPemeriksaanRequest struct {
	IdBalita            string `json:"id_balita"`
	IdIntervensi        string `json:"id_intervensi"`
	IdLaporanMasyarakat string `json:"id_laporan_masyarakat"` // <- Field baru wajib
	Tanggal             string `json:"tanggal"`               // Format: YYYY-MM-DD
	BeratBadan          string `json:"berat_badan"`           // in kg (decimal)
	TinggiBadan         string `json:"tinggi_badan"`          // in cm (decimal)
	StatusGizi          string `json:"status_gizi"`           // "normal", "stunting", "gizi buruk"
	Keterangan          string `json:"keterangan"`
}

func (r *insertRiwayatPemeriksaanRequest) validate() error {
	// ID Balita validation (wajib)
	if r.IdBalita == "" {
		return fmt.Errorf("id balita is required")
	}

	// ID Intervensi validation (wajib)
	if r.IdIntervensi == "" {
		return fmt.Errorf("id intervensi is required")
	}

	// ID Laporan Masyarakat validation (wajib) <- Tambahan validasi baru
	if r.IdLaporanMasyarakat == "" {
		return fmt.Errorf("id laporan masyarakat is required")
	}

	// Tanggal validation: YYYY-MM-DD format
	if r.Tanggal == "" {
		return fmt.Errorf("tanggal pemeriksaan is required")
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
	pemeriksaanDate, _ := time.Parse("2006-01-02", r.Tanggal)
	if pemeriksaanDate.After(time.Now()) {
		return fmt.Errorf("tanggal pemeriksaan cannot be in the future")
	}

	// Berat Badan validation: decimal, reasonable range (1-50 kg)
	if r.BeratBadan == "" {
		return fmt.Errorf("berat badan is required")
	}
	beratBadan, err := strconv.ParseFloat(r.BeratBadan, 64)
	if err != nil {
		return fmt.Errorf("berat badan must be a valid decimal number (in kg)")
	}
	if beratBadan < 1.0 || beratBadan > 50.0 {
		return fmt.Errorf("berat badan must be between 1.0-50.0 kg")
	}

	// Tinggi Badan validation: decimal, reasonable range (30-150 cm)
	if r.TinggiBadan == "" {
		return fmt.Errorf("tinggi badan is required")
	}
	tinggiBadan, err := strconv.ParseFloat(r.TinggiBadan, 64)
	if err != nil {
		return fmt.Errorf("tinggi badan must be a valid decimal number (in cm)")
	}
	if tinggiBadan < 30.0 || tinggiBadan > 150.0 {
		return fmt.Errorf("tinggi badan must be between 30.0-150.0 cm")
	}

	// Status Gizi validation: must be one of the allowed status
	if r.StatusGizi == "" {
		return fmt.Errorf("status gizi is required")
	}
	allowedStatus := []string{"normal", "stunting", "gizi buruk"}
	statusValid := false
	for _, status := range allowedStatus {
		if r.StatusGizi == status {
			statusValid = true
			break
		}
	}
	if !statusValid {
		return fmt.Errorf("status gizi must be one of: normal, stunting, gizi buruk")
	}

	// Keterangan validation
	if r.Keterangan == "" {
		return fmt.Errorf("keterangan is required")
	}
	if len(r.Keterangan) < 5 || len(r.Keterangan) > 500 {
		return fmt.Errorf("keterangan must be between 5-500 characters")
	}

	return nil
}

type insertRiwayatPemeriksaanResponse struct {
	Id string `json:"id"`
}

// # RiwayatPemeriksaanInsert handles inserting new riwayat pemeriksaan data
//
// @Summary Insert new riwayat pemeriksaan
// @Description Insert new riwayat pemeriksaan data (Admin only)
// @Description
// @Description Creates a new riwayat pemeriksaan record with:
// @Description - id_balita: balita being examined
// @Description - id_intervensi: related intervention program
// @Description - id_laporan_masyarakat: related masyarakat report
// @Description - tanggal: examination date (YYYY-MM-DD format)
// @Description - berat_badan: weight in kg (decimal)
// @Description - tinggi_badan: height in cm (decimal)
// @Description - status_gizi: nutritional status (normal, stunting, gizi buruk)
// @Description - keterangan: examination notes and recommendations
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param riwayat body insertRiwayatPemeriksaanRequest true "Riwayat pemeriksaan data"
// @Success 200 {object} object.Response{data=insertRiwayatPemeriksaanResponse} "Riwayat pemeriksaan inserted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/riwayat-pemeriksaan/insert [post]
func RiwayatPemeriksaanInsert(w http.ResponseWriter, r *http.Request) {
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
	var req insertRiwayatPemeriksaanRequest
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

	// Check if balita exists and not soft deleted
	var balitaExists int
	var namaBalita, umurBalita string
	checkBalitaQuery := `SELECT COUNT(*), b.nama, 
        TIMESTAMPDIFF(MONTH, b.tanggal_lahir, CURDATE()) as umur_bulan
        FROM balita b WHERE b.id = ? AND b.deleted_date IS NULL 
        GROUP BY b.nama`
	err = db.QueryRow(checkBalitaQuery, req.IdBalita).Scan(&balitaExists, &namaBalita, &umurBalita)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusBadRequest, "Balita not found", nil)
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
		response := object.NewResponse(http.StatusBadRequest, "Balita not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if intervensi exists and not soft deleted
	var intervensiExists int
	var jenisIntervensi, tanggalIntervensi string
	checkIntervensiQuery := `SELECT COUNT(*), jenis, tanggal 
        FROM intervensi WHERE id = ? AND deleted_date IS NULL 
        GROUP BY jenis, tanggal`
	err = db.QueryRow(checkIntervensiQuery, req.IdIntervensi).Scan(&intervensiExists, &jenisIntervensi, &tanggalIntervensi)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusBadRequest, "Intervensi not found", nil)
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
	if intervensiExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Intervensi not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if laporan masyarakat exists and not soft deleted <- Validasi baru
	var laporanExists int
	var statusLaporan, jenisLaporan, tanggalLaporan string
	var idMasyarakatLaporan sql.NullString
	checkLaporanQuery := `SELECT COUNT(*), lm.id_masyarakat, sl.status, lm.tanggal_laporan
        FROM laporan_masyarakat lm
        LEFT JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        WHERE lm.id = ? AND lm.deleted_date IS NULL
        GROUP BY lm.id_masyarakat, sl.status, lm.tanggal_laporan`
	err = db.QueryRow(checkLaporanQuery, req.IdLaporanMasyarakat).Scan(&laporanExists, &idMasyarakatLaporan, &statusLaporan, &tanggalLaporan)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusBadRequest, "Laporan masyarakat not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check laporan masyarakat existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if laporanExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Laporan masyarakat not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Determine jenis laporan
	if idMasyarakatLaporan.Valid {
		jenisLaporan = "masyarakat"
	} else {
		jenisLaporan = "admin"
	}

	// Validate that laporan is related to the same balita <- Validasi business logic
	var laporanBalitaId string
	checkLaporanBalitaQuery := "SELECT id_balita FROM laporan_masyarakat WHERE id = ?"
	err = db.QueryRow(checkLaporanBalitaQuery, req.IdLaporanMasyarakat).Scan(&laporanBalitaId)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check laporan balita relationship", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if laporanBalitaId != req.IdBalita {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Laporan masyarakat is not related to balita '%s'. Please select the correct laporan.", namaBalita), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Business logic: Check if laporan status allows medical examination
	allowedStatuses := []string{"Diproses dan data sesuai", "Belum ditindaklanjuti", "Sudah ditindaklanjuti"}
	statusAllowed := false
	for _, allowedStatus := range allowedStatuses {
		if statusLaporan == allowedStatus {
			statusAllowed = true
			break
		}
	}
	if !statusAllowed {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot create medical examination. Laporan status '%s' does not allow medical examination. Required status: %v",
				statusLaporan, allowedStatuses), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for duplicate riwayat pemeriksaan (same balita, intervensi, laporan, and tanggal)
	var duplicateExists int
	checkDuplicateQuery := `SELECT COUNT(*) FROM riwayat_pemeriksaan 
        WHERE id_balita = ? AND id_intervensi = ? AND id_laporan_masyarakat = ? AND tanggal = ? AND deleted_date IS NULL`
	err = db.QueryRow(checkDuplicateQuery, req.IdBalita, req.IdIntervensi, req.IdLaporanMasyarakat, req.Tanggal).Scan(&duplicateExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check duplicate riwayat pemeriksaan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if duplicateExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Riwayat pemeriksaan for balita '%s' on date '%s' in intervensi '%s' with this laporan already exists",
				namaBalita, req.Tanggal, jenisIntervensi), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Validate examination date constraints
	pemeriksaanDate, _ := time.Parse("2006-01-02", req.Tanggal)

	// Should not be before intervensi starts
	intervensiDate, err := time.Parse("2006-01-02", tanggalIntervensi)
	if err == nil && pemeriksaanDate.Before(intervensiDate) {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Tanggal pemeriksaan (%s) cannot be before intervensi date (%s)", req.Tanggal, tanggalIntervensi), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Should not be before laporan date
	laporanDate, err := time.Parse("2006-01-02", tanggalLaporan)
	if err == nil && pemeriksaanDate.Before(laporanDate) {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Tanggal pemeriksaan (%s) cannot be before laporan date (%s)", req.Tanggal, tanggalLaporan), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Insert riwayat pemeriksaan dengan id_laporan_masyarakat
	insertQuery := `INSERT INTO riwayat_pemeriksaan 
        (id_balita, id_intervensi, id_laporan_masyarakat, tanggal, berat_badan, tinggi_badan, status_gizi, keterangan, created_id, created_date) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(insertQuery,
		req.IdBalita,
		req.IdIntervensi,
		req.IdLaporanMasyarakat, // <- Parameter baru
		req.Tanggal,
		req.BeratBadan,
		req.TinggiBadan,
		req.StatusGizi,
		req.Keterangan,
		userId,
		currentTime,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert riwayat pemeriksaan", nil)
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

	// Prepare success response with additional context information
	message := fmt.Sprintf("Riwayat pemeriksaan balita '%s' berhasil ditambahkan untuk intervensi %s pada tanggal %s (Status: %s, Laporan: %s)",
		namaBalita, jenisIntervensi, req.Tanggal, req.StatusGizi, jenisLaporan)

	response := object.NewResponse(http.StatusOK, message, insertRiwayatPemeriksaanResponse{
		Id: strconv.FormatInt(insertedId, 10),
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
