package community

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type balitaResponse struct {
	Id           string `json:"id"`
	IdKeluarga   string `json:"id_keluarga"`
	NomorKk      string `json:"nomor_kk"`
	NamaAyah     string `json:"nama_ayah"`
	NamaIbu      string `json:"nama_ibu"`
	Nama         string `json:"nama"`
	TanggalLahir string `json:"tanggal_lahir"`
	JenisKelamin string `json:"jenis_kelamin"`
	BeratLahir   string `json:"berat_lahir"`
	TinggiLahir  string `json:"tinggi_lahir"`
	Umur         string `json:"umur"` // calculated field in months
	Kelurahan    string `json:"kelurahan"`
	Kecamatan    string `json:"kecamatan"`
	CreatedDate  string `json:"created_date"`
	UpdatedDate  string `json:"updated_date,omitempty"`

	// Status untuk masyarakat
	JumlahLaporan              int    `json:"jumlah_laporan"`
	StatusLaporanAktif         string `json:"status_laporan_aktif,omitempty"`
	StatusGiziTerakhir         string `json:"status_gizi_terakhir,omitempty"`
	TanggalPemeriksaanTerakhir string `json:"tanggal_pemeriksaan_terakhir,omitempty"`
	CanEdit                    bool   `json:"can_edit"`
	CanReport                  bool   `json:"can_report"`
}

type getAllBalitaResponse struct {
	Data  []balitaResponse `json:"data"`
	Total int              `json:"total"`
}

type getBalitaByIdResponse struct {
	Data balitaResponse `json:"data"`
}

// # BalitaGet handles getting balita data for masyarakat
//
// @Summary Get balita data (Community)
// @Description Get balita data for community/masyarakat users (own data only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all balita from user's keluarga
// @Description - With id parameter: Returns specific balita data (if owned by user)
// @Description
// @Description Data includes balita information, laporan status, medical history summary,
// @Description and action permissions (edit/report capabilities).
// @Description Users can only access balita from keluarga they have created themselves.
// @Tags community
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Balita ID"
// @Success 200 {object} object.Response{data=getAllBalitaResponse} "Balita data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Masyarakat role required"
// @Failure 404 {object} object.Response{data=nil} "Balita not found or not owned by user"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/community/balita/get [get]
func BalitaGet(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is masyarakat
	if role != "masyarakat" {
		response := object.NewResponse(http.StatusForbidden, "Access denied. Masyarakat role required", nil)
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

	// Check if ID parameter is provided
	idParam := r.URL.Query().Get("id")
	if idParam != "" {
		// Get specific balita by ID (only if owned by user)
		balita, err := getBalitaByIdForUser(db, idParam, userId)
		if err != nil {
			if err == sql.ErrNoRows {
				response := object.NewResponse(http.StatusNotFound, "Balita not found or not owned by you", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get balita", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Balita retrieved successfully", getBalitaByIdResponse{
			Data: balita,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Get all balita for this user's keluarga
		balitaList, total, err := getAllBalitaForUser(db, userId)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get balita list", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "All balita retrieved successfully", getAllBalitaResponse{
			Data:  balitaList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Helper function to get balita by ID for specific user (ownership check)
func getBalitaByIdForUser(db *sql.DB, id string, userId string) (balitaResponse, error) {
	var balita balitaResponse
	var updatedDate sql.NullString

	query := `
        SELECT 
            b.id, b.id_keluarga, b.nama, b.tanggal_lahir, b.jenis_kelamin,
            b.berat_lahir, b.tinggi_lahir, b.created_date, b.updated_date,
            k.nomor_kk, k.nama_ayah, k.nama_ibu,
            kel.kelurahan, kec.kecamatan
        FROM balita b
        JOIN keluarga k ON b.id_keluarga = k.id
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        WHERE b.id = ? AND k.created_id = ? AND b.deleted_date IS NULL AND k.deleted_date IS NULL
    `

	err := db.QueryRow(query, id, userId).Scan(
		&balita.Id,
		&balita.IdKeluarga,
		&balita.Nama,
		&balita.TanggalLahir,
		&balita.JenisKelamin,
		&balita.BeratLahir,
		&balita.TinggiLahir,
		&balita.CreatedDate,
		&updatedDate,
		&balita.NomorKk,
		&balita.NamaAyah,
		&balita.NamaIbu,
		&balita.Kelurahan,
		&balita.Kecamatan,
	)

	if err != nil {
		return balita, err
	}

	// Calculate age in months
	balita.Umur = calculateAgeInMonths(balita.TanggalLahir)

	// Handle nullable updated_date
	if updatedDate.Valid {
		balita.UpdatedDate = updatedDate.String
	}

	// Get additional information for masyarakat
	err = getBalitaAdditionalInfo(db, &balita)
	if err != nil {
		// Log error but don't fail the request
		// Additional info is not critical
	}

	return balita, nil
}

// Helper function to get all balita for specific user's keluarga
func getAllBalitaForUser(db *sql.DB, userId string) ([]balitaResponse, int, error) {
	var balitaList []balitaResponse

	query := `
        SELECT 
            b.id, b.id_keluarga, b.nama, b.tanggal_lahir, b.jenis_kelamin,
            b.berat_lahir, b.tinggi_lahir, b.created_date, b.updated_date,
            k.nomor_kk, k.nama_ayah, k.nama_ibu,
            kel.kelurahan, kec.kecamatan
        FROM balita b
        JOIN keluarga k ON b.id_keluarga = k.id
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        WHERE k.created_id = ? AND b.deleted_date IS NULL AND k.deleted_date IS NULL
        ORDER BY b.created_date DESC
    `

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var balita balitaResponse
		var updatedDate sql.NullString

		err := rows.Scan(
			&balita.Id,
			&balita.IdKeluarga,
			&balita.Nama,
			&balita.TanggalLahir,
			&balita.JenisKelamin,
			&balita.BeratLahir,
			&balita.TinggiLahir,
			&balita.CreatedDate,
			&updatedDate,
			&balita.NomorKk,
			&balita.NamaAyah,
			&balita.NamaIbu,
			&balita.Kelurahan,
			&balita.Kecamatan,
		)

		if err != nil {
			return nil, 0, err
		}

		// Calculate age in months
		balita.Umur = calculateAgeInMonths(balita.TanggalLahir)

		// Handle nullable updated_date
		if updatedDate.Valid {
			balita.UpdatedDate = updatedDate.String
		}

		// Get additional information for masyarakat
		err = getBalitaAdditionalInfo(db, &balita)
		if err != nil {
			// Log error but don't fail the request
			// Set default values
			balita.JumlahLaporan = 0
			balita.CanEdit = true
			balita.CanReport = true
		}

		balitaList = append(balitaList, balita)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count for this user
	var total int
	countQuery := `SELECT COUNT(*) 
        FROM balita b 
        JOIN keluarga k ON b.id_keluarga = k.id 
        WHERE k.created_id = ? AND b.deleted_date IS NULL AND k.deleted_date IS NULL`
	err = db.QueryRow(countQuery, userId).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return balitaList, total, nil
}

// Helper function to get additional information for balita (laporan status, medical history, permissions, etc.)
func getBalitaAdditionalInfo(db *sql.DB, balita *balitaResponse) error {
	// Get laporan count and active status
	var laporanCount int
	var activeStatus sql.NullString
	laporanQuery := `
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
        WHERE lm.id_balita = ? 
        AND lm.deleted_date IS NULL
    `
	err := db.QueryRow(laporanQuery, balita.Id).Scan(&laporanCount, &activeStatus)
	if err != nil {
		laporanCount = 0
	}
	balita.JumlahLaporan = laporanCount

	// Set active status if any
	if activeStatus.Valid && activeStatus.String != "" {
		balita.StatusLaporanAktif = activeStatus.String
	}

	// Get latest medical examination status
	var latestGiziStatus sql.NullString
	var latestExamDate sql.NullString
	medicalQuery := `
        SELECT rp.status_gizi, rp.tanggal
        FROM riwayat_pemeriksaan rp
        WHERE rp.id_balita = ? AND rp.deleted_date IS NULL
        ORDER BY rp.tanggal DESC, rp.created_date DESC
        LIMIT 1
    `
	err = db.QueryRow(medicalQuery, balita.Id).Scan(&latestGiziStatus, &latestExamDate)
	if err == nil {
		if latestGiziStatus.Valid {
			balita.StatusGiziTerakhir = latestGiziStatus.String
		}
		if latestExamDate.Valid {
			balita.TanggalPemeriksaanTerakhir = latestExamDate.String
		}
	}

	// Determine if balita can be edited
	// Can't edit if there are active reports
	canEdit := true
	if activeStatus.Valid && activeStatus.String != "" {
		canEdit = false
	}
	balita.CanEdit = canEdit

	// Determine if balita can be reported
	// Can report if no active reports or if latest report status allows new report
	canReport := true
	if activeStatus.Valid && activeStatus.String != "" {
		// Check if there's any "Belum diproses" status
		if activeStatus.String == "Belum diproses" {
			canReport = false
		}
	}
	balita.CanReport = canReport

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
		return strconv.Itoa(totalMonths) + " bulan"
	} else {
		years := totalMonths / 12
		remainingMonths := totalMonths % 12
		if remainingMonths == 0 {
			return strconv.Itoa(years) + " tahun"
		}
		return strconv.Itoa(years) + " tahun " + strconv.Itoa(remainingMonths) + " bulan"
	}
}
