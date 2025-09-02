package community

import (
	"database/sql"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type laporanResponse struct {
	Id                    string `json:"id"`
	IdBalita              string `json:"id_balita"`
	NamaBalita            string `json:"nama_balita"`
	NamaAyah              string `json:"nama_ayah"`
	NamaIbu               string `json:"nama_ibu"`
	NomorKk               string `json:"nomor_kk"`
	Alamat                string `json:"alamat"`
	Kelurahan             string `json:"kelurahan"`
	Kecamatan             string `json:"kecamatan"`
	IdStatusLaporan       string `json:"id_status_laporan"`
	StatusLaporan         string `json:"status_laporan"`
	TanggalLaporan        string `json:"tanggal_laporan"`
	HubunganDenganBalita  string `json:"hubungan_dengan_balita"`
	NomorHpPelapor        string `json:"nomor_hp_pelapor"`
	NomorHpKeluargaBalita string `json:"nomor_hp_keluarga_balita"`
	CreatedDate           string `json:"created_date"`
	UpdatedDate           string `json:"updated_date,omitempty"`

	// Status untuk masyarakat
	CanEdit                 bool   `json:"can_edit"`
	StatusKeterangan        string `json:"status_keterangan,omitempty"`
	RiwayatPemeriksaan      int    `json:"riwayat_pemeriksaan"`
	IntervensiTerkait       int    `json:"intervensi_terkait"`
	TanggalTerakhirDiproses string `json:"tanggal_terakhir_diproses,omitempty"`
}

type getAllLaporanResponse struct {
	Data  []laporanResponse `json:"data"`
	Total int               `json:"total"`
}

type getLaporanByIdResponse struct {
	Data laporanResponse `json:"data"`
}

// # LaporanGet handles getting laporan data for masyarakat
//
// @Summary Get laporan data (Community)
// @Description Get laporan data for community/masyarakat users (own reports only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all laporan created by the user
// @Description - With id parameter: Returns specific laporan data (if owned by user)
// @Description
// @Description Data includes laporan information, balita details, keluarga info, status tracking,
// @Description related medical records count, and action permissions.
// @Description Users can only access laporan they have created themselves.
// @Tags community
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Laporan ID"
// @Success 200 {object} object.Response{data=getAllLaporanResponse} "Laporan data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Masyarakat role required"
// @Failure 404 {object} object.Response{data=nil} "Laporan not found or not owned by user"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/community/laporan/get [get]
func LaporanGet(w http.ResponseWriter, r *http.Request) {
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
		// Get specific laporan by ID (only if owned by user)
		laporan, err := getLaporanByIdForUser(db, idParam, masyarakatId)
		if err != nil {
			if err == sql.ErrNoRows {
				response := object.NewResponse(http.StatusNotFound, "Laporan not found or not owned by you", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get laporan", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Laporan retrieved successfully", getLaporanByIdResponse{
			Data: laporan,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Get all laporan for this user
		laporanList, total, err := getAllLaporanForUser(db, masyarakatId)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get laporan list", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "All laporan retrieved successfully", getAllLaporanResponse{
			Data:  laporanList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Helper function to get laporan by ID for specific user (ownership check)
func getLaporanByIdForUser(db *sql.DB, id string, masyarakatId string) (laporanResponse, error) {
	var laporan laporanResponse
	var updatedDate sql.NullString

	query := `
        SELECT 
            lm.id, lm.id_balita, lm.id_status_laporan,
            lm.tanggal_laporan, lm.hubungan_dengan_balita, lm.nomor_hp_pelapor, 
            lm.nomor_hp_keluarga_balita, lm.created_date, lm.updated_date,
            b.nama as nama_balita,
            k.nomor_kk, k.nama_ayah, k.nama_ibu, k.alamat,
            kel.kelurahan, kec.kecamatan,
            sl.status as status_laporan
        FROM laporan_masyarakat lm
        JOIN balita b ON lm.id_balita = b.id
        JOIN keluarga k ON b.id_keluarga = k.id
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        WHERE lm.id = ? AND lm.id_masyarakat = ? AND lm.deleted_date IS NULL 
        AND b.deleted_date IS NULL AND k.deleted_date IS NULL
    `

	err := db.QueryRow(query, id, masyarakatId).Scan(
		&laporan.Id,
		&laporan.IdBalita,
		&laporan.IdStatusLaporan,
		&laporan.TanggalLaporan,
		&laporan.HubunganDenganBalita,
		&laporan.NomorHpPelapor,
		&laporan.NomorHpKeluargaBalita,
		&laporan.CreatedDate,
		&updatedDate,
		&laporan.NamaBalita,
		&laporan.NomorKk,
		&laporan.NamaAyah,
		&laporan.NamaIbu,
		&laporan.Alamat,
		&laporan.Kelurahan,
		&laporan.Kecamatan,
		&laporan.StatusLaporan,
	)

	if err != nil {
		return laporan, err
	}

	// Handle nullable updated_date
	if updatedDate.Valid {
		laporan.UpdatedDate = updatedDate.String
	}

	// Get additional information for masyarakat
	err = getLaporanAdditionalInfo(db, &laporan)
	if err != nil {
		// Log error but don't fail the request
		// Additional info is not critical
	}

	return laporan, nil
}

// Helper function to get all laporan for specific user
func getAllLaporanForUser(db *sql.DB, masyarakatId string) ([]laporanResponse, int, error) {
	var laporanList []laporanResponse

	query := `
        SELECT 
            lm.id, lm.id_balita, lm.id_status_laporan,
            lm.tanggal_laporan, lm.hubungan_dengan_balita, lm.nomor_hp_pelapor, 
            lm.nomor_hp_keluarga_balita, lm.created_date, lm.updated_date,
            b.nama as nama_balita,
            k.nomor_kk, k.nama_ayah, k.nama_ibu, k.alamat,
            kel.kelurahan, kec.kecamatan,
            sl.status as status_laporan
        FROM laporan_masyarakat lm
        JOIN balita b ON lm.id_balita = b.id AND b.deleted_date IS NULL
        JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        WHERE lm.id_masyarakat = ? AND lm.deleted_date IS NULL
        ORDER BY lm.created_date DESC
    `

	rows, err := db.Query(query, masyarakatId)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var laporan laporanResponse
		var updatedDate sql.NullString

		err := rows.Scan(
			&laporan.Id,
			&laporan.IdBalita,
			&laporan.IdStatusLaporan,
			&laporan.TanggalLaporan,
			&laporan.HubunganDenganBalita,
			&laporan.NomorHpPelapor,
			&laporan.NomorHpKeluargaBalita,
			&laporan.CreatedDate,
			&updatedDate,
			&laporan.NamaBalita,
			&laporan.NomorKk,
			&laporan.NamaAyah,
			&laporan.NamaIbu,
			&laporan.Alamat,
			&laporan.Kelurahan,
			&laporan.Kecamatan,
			&laporan.StatusLaporan,
		)

		if err != nil {
			return nil, 0, err
		}

		// Handle nullable updated_date
		if updatedDate.Valid {
			laporan.UpdatedDate = updatedDate.String
		}

		// Get additional information for masyarakat
		err = getLaporanAdditionalInfo(db, &laporan)
		if err != nil {
			// Log error but don't fail the request
			// Set default values
			laporan.CanEdit = false
			laporan.RiwayatPemeriksaan = 0
			laporan.IntervensiTerkait = 0
		}

		laporanList = append(laporanList, laporan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count for this user
	var total int
	countQuery := "SELECT COUNT(*) FROM laporan_masyarakat WHERE id_masyarakat = ? AND deleted_date IS NULL"
	err = db.QueryRow(countQuery, masyarakatId).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return laporanList, total, nil
}

// Helper function to get additional information for laporan (status explanation, related records, permissions, etc.)
func getLaporanAdditionalInfo(db *sql.DB, laporan *laporanResponse) error {
	// Determine if laporan can be edited
	// Can't edit if status is not "Belum diproses"
	canEdit := laporan.StatusLaporan == "Belum diproses"
	laporan.CanEdit = canEdit

	// Get status explanation for user understanding
	statusKeterangan := getStatusExplanation(laporan.StatusLaporan)
	laporan.StatusKeterangan = statusKeterangan

	// Get count of related riwayat pemeriksaan
	var riwayatCount int
	riwayatQuery := `
        SELECT COUNT(*) 
        FROM riwayat_pemeriksaan rp
        WHERE rp.id_balita = ? AND rp.deleted_date IS NULL
    `
	err := db.QueryRow(riwayatQuery, laporan.IdBalita).Scan(&riwayatCount)
	if err != nil {
		riwayatCount = 0
	}
	laporan.RiwayatPemeriksaan = riwayatCount

	// Get count of related intervensi
	var intervensiCount int
	intervensiQuery := `
        SELECT COUNT(*) 
        FROM intervensi i
        WHERE i.id_balita = ? AND i.deleted_date IS NULL
    `
	err = db.QueryRow(intervensiQuery, laporan.IdBalita).Scan(&intervensiCount)
	if err != nil {
		intervensiCount = 0
	}
	laporan.IntervensiTerkait = intervensiCount

	// Get latest processing date (when status was last changed from "Belum diproses")
	var latestProcessDate sql.NullString
	processDateQuery := `
        SELECT lm.updated_date 
        FROM laporan_masyarakat lm
        JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        WHERE lm.id = ? AND sl.status != 'Belum diproses' AND lm.updated_date IS NOT NULL
        ORDER BY lm.updated_date DESC
        LIMIT 1
    `
	err = db.QueryRow(processDateQuery, laporan.Id).Scan(&latestProcessDate)
	if err == nil && latestProcessDate.Valid {
		laporan.TanggalTerakhirDiproses = latestProcessDate.String
	}

	return nil
}

// Helper function to provide user-friendly status explanations
func getStatusExplanation(status string) string {
	switch status {
	case "Belum diproses":
		return "Laporan Anda sedang menunggu untuk diproses oleh petugas kesehatan"
	case "Diproses dan data tidak sesuai":
		return "Laporan telah diproses, namun data yang dilaporkan tidak sesuai dengan kondisi aktual"
	case "Diproses dan data sesuai":
		return "Laporan telah diproses dan data yang dilaporkan sesuai dengan kondisi aktual"
	case "Belum ditindaklanjuti":
		return "Laporan telah diverifikasi dan menunggu tindak lanjut dari petugas kesehatan"
	case "Sudah ditindaklanjuti":
		return "Laporan telah ditindaklanjuti dengan pemeriksaan dan intervensi yang diperlukan"
	case "Sudah perbaikan gizi":
		return "Balita telah menunjukkan perbaikan status gizi setelah intervensi yang diberikan"
	default:
		return "Status laporan tidak dikenali"
	}
}
