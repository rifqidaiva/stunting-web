package healthworker

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type assignedIntervensiResponse struct {
	Id                 string `json:"id"`
	IdIntervensi       string `json:"id_intervensi"`
	IdBalita           string `json:"id_balita"`
	NamaBalita         string `json:"nama_balita"`
	TanggalLahirBalita string `json:"tanggal_lahir_balita"`
	JenisKelaminBalita string `json:"jenis_kelamin_balita"`
	UmurBalita         string `json:"umur_balita"`

	// Keluarga Info
	NomorKk   string `json:"nomor_kk"`
	NamaAyah  string `json:"nama_ayah"`
	NamaIbu   string `json:"nama_ibu"`
	Alamat    string `json:"alamat"`
	Kelurahan string `json:"kelurahan"`
	Kecamatan string `json:"kecamatan"`

	// Intervensi Info
	JenisIntervensi     string `json:"jenis_intervensi"`
	TanggalIntervensi   string `json:"tanggal_intervensi"`
	DeskripsiIntervensi string `json:"deskripsi_intervensi"`
	HasilIntervensi     string `json:"hasil_intervensi,omitempty"`
	StatusIntervensi    string `json:"status_intervensi"`

	// Latest Medical Info
	StatusGiziTerakhir         string `json:"status_gizi_terakhir,omitempty"`
	TanggalPemeriksaanTerakhir string `json:"tanggal_pemeriksaan_terakhir,omitempty"`
	BeratBadanTerakhir         string `json:"berat_badan_terakhir,omitempty"`
	TinggiBadanTerakhir        string `json:"tinggi_badan_terakhir,omitempty"`

	// Related Reports
	JumlahLaporanTerkait int    `json:"jumlah_laporan_terkait"`
	StatusLaporanAktif   string `json:"status_laporan_aktif,omitempty"`

	// Assignment Info
	TanggalPenugasan string `json:"tanggal_penugasan"`

	// Actions
	CanAddMedicalRecord bool `json:"can_add_medical_record"`
	CanUpdateStatus     bool `json:"can_update_status"`
}

type getAllAssignedIntervensiResponse struct {
	Data  []assignedIntervensiResponse `json:"data"`
	Total int                          `json:"total"`
}

type getAssignedIntervensiByIdResponse struct {
	Data assignedIntervensiResponse `json:"data"`
}

// # AssignmentGet handles getting assigned interventions for health workers
//
// @Summary Get assigned interventions (Health Worker)
// @Description Get all interventions assigned to the authenticated health worker
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all assigned interventions
// @Description - With id parameter: Returns specific intervention details (if assigned to user)
// @Description
// @Description Data includes intervention information, balita details, family info, medical history,
// @Description related reports, and action permissions based on intervention status.
// @Description Health workers can only access interventions assigned to them.
// @Tags health-worker
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Intervention ID"
// @Param status query string false "Filter by intervention status (pending, in_progress, completed)"
// @Success 200 {object} object.Response{data=getAllAssignedIntervensiResponse} "Assigned interventions retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Health worker role required"
// @Failure 404 {object} object.Response{data=nil} "Intervention not found or not assigned to user"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/health-worker/assignment/get [get]
func AssignmentGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	// Check if user is petugas kesehatan
	if role != "petugas kesehatan" {
		response := object.NewResponse(http.StatusForbidden, "Access denied. Health worker role required", nil)
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

	// Verify user exists and get petugas kesehatan ID
	var petugasKesehatanId string
	var skpdId string
	checkUserQuery := `
        SELECT pk.id, pk.id_skpd 
        FROM petugas_kesehatan pk 
        JOIN pengguna p ON pk.id_pengguna = p.id 
        WHERE p.id = ? AND pk.deleted_date IS NULL
    `
	err = db.QueryRow(checkUserQuery, userId).Scan(&petugasKesehatanId, &skpdId)
	if err != nil {
		response := object.NewResponse(http.StatusUnauthorized, "Health worker profile not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if ID parameter is provided
	idParam := r.URL.Query().Get("id")
	statusParam := r.URL.Query().Get("status")

	if idParam != "" {
		// Get specific assigned intervention by ID
		intervention, err := getAssignedIntervensiByIdForUser(db, idParam, petugasKesehatanId)
		if err != nil {
			if err == sql.ErrNoRows {
				response := object.NewResponse(http.StatusNotFound, "Intervention not found or not assigned to you", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get intervention", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Assigned intervention retrieved successfully", getAssignedIntervensiByIdResponse{
			Data: intervention,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Get all assigned interventions for this user
		interventionList, total, err := getAllAssignedIntervensiForUser(db, petugasKesehatanId, statusParam)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get assigned interventions", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		message := "All assigned interventions retrieved successfully"
		if statusParam != "" {
			message = "Filtered assigned interventions retrieved successfully"
		}

		response := object.NewResponse(http.StatusOK, message, getAllAssignedIntervensiResponse{
			Data:  interventionList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Helper function to get assigned intervention by ID for specific health worker
func getAssignedIntervensiByIdForUser(db *sql.DB, interventionId string, petugasKesehatanId string) (assignedIntervensiResponse, error) {
	var intervention assignedIntervensiResponse

	query := `
        SELECT 
            ip.id as assignment_id,
            i.id as id_intervensi,
            i.id_balita,
            i.jenis as jenis_intervensi,
            i.tanggal as tanggal_intervensi,
            i.deskripsi as deskripsi_intervensi,
            i.hasil as hasil_intervensi,
            i.created_date as tanggal_penugasan,
            b.nama as nama_balita,
            b.tanggal_lahir as tanggal_lahir_balita,
            b.jenis_kelamin as jenis_kelamin_balita,
            k.nomor_kk,
            k.nama_ayah,
            k.nama_ibu,
            k.alamat,
            kel.kelurahan,
            kec.kecamatan
        FROM intervensi_petugas ip
        JOIN intervensi i ON ip.id_intervensi = i.id
        JOIN balita b ON i.id_balita = b.id
        JOIN keluarga k ON b.id_keluarga = k.id
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        WHERE ip.id_petugas_kesehatan = ? 
        AND i.id = ? 
        AND i.deleted_date IS NULL 
        AND b.deleted_date IS NULL 
        AND k.deleted_date IS NULL
    `

	err := db.QueryRow(query, petugasKesehatanId, interventionId).Scan(
		&intervention.Id,
		&intervention.IdIntervensi,
		&intervention.IdBalita,
		&intervention.JenisIntervensi,
		&intervention.TanggalIntervensi,
		&intervention.DeskripsiIntervensi,
		&intervention.HasilIntervensi,
		&intervention.TanggalPenugasan,
		&intervention.NamaBalita,
		&intervention.TanggalLahirBalita,
		&intervention.JenisKelaminBalita,
		&intervention.NomorKk,
		&intervention.NamaAyah,
		&intervention.NamaIbu,
		&intervention.Alamat,
		&intervention.Kelurahan,
		&intervention.Kecamatan,
	)

	if err != nil {
		return intervention, err
	}

	// Calculate age
	intervention.UmurBalita = calculateAgeInMonths(intervention.TanggalLahirBalita)

	// Get additional information
	err = getIntervensiAdditionalInfo(db, &intervention)
	if err != nil {
		// Log error but don't fail the request
		// Additional info is not critical for basic functionality
	}

	return intervention, nil
}

// Helper function to get all assigned interventions for specific health worker
func getAllAssignedIntervensiForUser(db *sql.DB, petugasKesehatanId string, statusFilter string) ([]assignedIntervensiResponse, int, error) {
	var interventionList []assignedIntervensiResponse

	// Build query based on status filter
	query := `
        SELECT 
            ip.id as assignment_id,
            i.id as id_intervensi,
            i.id_balita,
            i.jenis as jenis_intervensi,
            i.tanggal as tanggal_intervensi,
            i.deskripsi as deskripsi_intervensi,
            i.hasil as hasil_intervensi,
            i.created_date as tanggal_penugasan,
            b.nama as nama_balita,
            b.tanggal_lahir as tanggal_lahir_balita,
            b.jenis_kelamin as jenis_kelamin_balita,
            k.nomor_kk,
            k.nama_ayah,
            k.nama_ibu,
            k.alamat,
            kel.kelurahan,
            kec.kecamatan
        FROM intervensi_petugas ip
        JOIN intervensi i ON ip.id_intervensi = i.id
        JOIN balita b ON i.id_balita = b.id AND b.deleted_date IS NULL
        JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        WHERE ip.id_petugas_kesehatan = ? 
        AND i.deleted_date IS NULL
    `

	args := []any{petugasKesehatanId}

	// Add status filter if provided
	if statusFilter != "" {
		switch statusFilter {
		case "pending":
			query += " AND (i.hasil IS NULL OR i.hasil = '')"
		case "in_progress":
			query += " AND i.hasil IS NOT NULL AND i.hasil != '' AND i.hasil NOT LIKE '%selesai%' AND i.hasil NOT LIKE '%completed%'"
		case "completed":
			query += " AND (i.hasil LIKE '%selesai%' OR i.hasil LIKE '%completed%')"
		}
	}

	query += " ORDER BY i.tanggal DESC, i.created_date DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var intervention assignedIntervensiResponse

		err := rows.Scan(
			&intervention.Id,
			&intervention.IdIntervensi,
			&intervention.IdBalita,
			&intervention.JenisIntervensi,
			&intervention.TanggalIntervensi,
			&intervention.DeskripsiIntervensi,
			&intervention.HasilIntervensi,
			&intervention.TanggalPenugasan,
			&intervention.NamaBalita,
			&intervention.TanggalLahirBalita,
			&intervention.JenisKelaminBalita,
			&intervention.NomorKk,
			&intervention.NamaAyah,
			&intervention.NamaIbu,
			&intervention.Alamat,
			&intervention.Kelurahan,
			&intervention.Kecamatan,
		)

		if err != nil {
			return nil, 0, err
		}

		// Calculate age
		intervention.UmurBalita = calculateAgeInMonths(intervention.TanggalLahirBalita)

		// Get additional information
		err = getIntervensiAdditionalInfo(db, &intervention)
		if err != nil {
			// Log error but don't fail the request
			// Set default values for failed additional info
			intervention.CanAddMedicalRecord = true
			intervention.CanUpdateStatus = true
			intervention.JumlahLaporanTerkait = 0
		}

		interventionList = append(interventionList, intervention)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countQuery := `
        SELECT COUNT(*) 
        FROM intervensi_petugas ip
        JOIN intervensi i ON ip.id_intervensi = i.id
        JOIN balita b ON i.id_balita = b.id AND b.deleted_date IS NULL
        JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        WHERE ip.id_petugas_kesehatan = ? 
        AND i.deleted_date IS NULL
    `

	countArgs := []any{petugasKesehatanId}

	// Add status filter to count query if provided
	if statusFilter != "" {
		switch statusFilter {
		case "pending":
			countQuery += " AND (i.hasil IS NULL OR i.hasil = '')"
		case "in_progress":
			countQuery += " AND i.hasil IS NOT NULL AND i.hasil != '' AND i.hasil NOT LIKE '%selesai%' AND i.hasil NOT LIKE '%completed%'"
		case "completed":
			countQuery += " AND (i.hasil LIKE '%selesai%' OR i.hasil LIKE '%completed%')"
		}
	}

	err = db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return interventionList, total, nil
}

// Helper function to get additional information for intervention
func getIntervensiAdditionalInfo(db *sql.DB, intervention *assignedIntervensiResponse) error {
	// Determine intervention status based on hasil field
	if intervention.HasilIntervensi == "" {
		intervention.StatusIntervensi = "pending"
	} else if intervention.HasilIntervensi != "" &&
		!(intervention.HasilIntervensi == "selesai" || intervention.HasilIntervensi == "completed") {
		intervention.StatusIntervensi = "in_progress"
	} else {
		intervention.StatusIntervensi = "completed"
	}

	// Get latest medical examination data
	var latestGiziStatus sql.NullString
	var latestExamDate sql.NullString
	var latestWeight sql.NullString
	var latestHeight sql.NullString

	medicalQuery := `
        SELECT rp.status_gizi, rp.tanggal, rp.berat_badan, rp.tinggi_badan
        FROM riwayat_pemeriksaan rp
        WHERE rp.id_balita = ? AND rp.deleted_date IS NULL
        ORDER BY rp.tanggal DESC, rp.created_date DESC
        LIMIT 1
    `
	err := db.QueryRow(medicalQuery, intervention.IdBalita).Scan(
		&latestGiziStatus, &latestExamDate, &latestWeight, &latestHeight)
	if err == nil {
		if latestGiziStatus.Valid {
			intervention.StatusGiziTerakhir = latestGiziStatus.String
		}
		if latestExamDate.Valid {
			intervention.TanggalPemeriksaanTerakhir = latestExamDate.String
		}
		if latestWeight.Valid {
			intervention.BeratBadanTerakhir = latestWeight.String
		}
		if latestHeight.Valid {
			intervention.TinggiBadanTerakhir = latestHeight.String
		}
	}

	// Get related reports count and active status
	var reportCount int
	var activeStatus sql.NullString
	reportQuery := `
        SELECT COUNT(*), 
            GROUP_CONCAT(
                CASE 
                    WHEN sl.status IN ('Belum diproses', 'Diproses dan data sesuai', 'Belum ditindaklanjuti') 
                    THEN sl.status 
                    ELSE NULL 
                END
            ) as active_statuses
        FROM laporan_masyarakat lm
        JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        WHERE lm.id_balita = ? AND lm.deleted_date IS NULL
    `
	err = db.QueryRow(reportQuery, intervention.IdBalita).Scan(&reportCount, &activeStatus)
	if err != nil {
		reportCount = 0
	}
	intervention.JumlahLaporanTerkait = reportCount

	if activeStatus.Valid && activeStatus.String != "" {
		intervention.StatusLaporanAktif = activeStatus.String
	}

	// Determine permissions
	// Can add medical record if intervention is in progress or completed
	intervention.CanAddMedicalRecord = intervention.StatusIntervensi != "pending"

	// Can update status if intervention is not completed
	intervention.CanUpdateStatus = intervention.StatusIntervensi != "completed"

	return nil
}

// Helper function to calculate age in months
func calculateAgeInMonths(birthDateStr string) string {
	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return "0 bulan"
	}

	now := time.Now()
	years := now.Year() - birthDate.Year()
	months := int(now.Month()) - int(birthDate.Month())

	// Adjust if the current day is before the birth day
	if now.Day() < birthDate.Day() {
		months--
	}

	// Convert to total months
	totalMonths := years*12 + months

	if totalMonths < 0 {
		return "0 bulan"
	}

	if totalMonths < 12 {
		return fmt.Sprintf("%d bulan", totalMonths)
	} else {
		years := totalMonths / 12
		remainingMonths := totalMonths % 12
		if remainingMonths == 0 {
			return fmt.Sprintf("%d tahun", years)
		}
		return fmt.Sprintf("%d tahun %d bulan", years, remainingMonths)
	}
}
