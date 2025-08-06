package api

import (
	"database/sql"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type laporanMasyarakatResponse struct {
	Id                    string `json:"id"`
	IdMasyarakat          string `json:"id_masyarakat,omitempty"`
	NamaPelapor           string `json:"nama_pelapor,omitempty"`
	EmailPelapor          string `json:"email_pelapor,omitempty"`
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
	JenisLaporan          string `json:"jenis_laporan"` // "masyarakat" atau "admin"
	CreatedDate           string `json:"created_date"`
	UpdatedDate           string `json:"updated_date,omitempty"`
}

type getAllLaporanMasyarakatResponse struct {
	Data  []laporanMasyarakatResponse `json:"data"`
	Total int                         `json:"total"`
}

type getLaporanMasyarakatByIdResponse struct {
	Data laporanMasyarakatResponse `json:"data"`
}

// # AdminLaporanMasyarakatGet handles getting laporan masyarakat data
//
// @Summary Get laporan masyarakat data
// @Description Get laporan masyarakat data based on query parameter (Admin only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all laporan masyarakat with total count
// @Description - With id parameter: Returns specific laporan masyarakat data
// @Description
// @Description Laporan masyarakat data includes: pelapor info, balita info, keluarga info, status laporan, contact details
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Laporan Masyarakat ID"
// @Success 200 {object} object.Response{data=getAllLaporanMasyarakatResponse} "Laporan masyarakat data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Laporan masyarakat not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/laporan-masyarakat/get [get]
func AdminLaporanMasyarakatGet(w http.ResponseWriter, r *http.Request) {
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

	_, role, err := object.ParseJWT(token)
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

	// Check if ID parameter is provided
	idParam := r.URL.Query().Get("id")
	if idParam != "" {
		// Get specific laporan masyarakat by ID
		laporan, err := getLaporanMasyarakatById(db, idParam)
		if err != nil {
			if err == sql.ErrNoRows {
				response := object.NewResponse(http.StatusNotFound, "Laporan masyarakat not found", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get laporan masyarakat", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Laporan masyarakat retrieved successfully", getLaporanMasyarakatByIdResponse{
			Data: laporan,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Get all laporan masyarakat
		laporanList, total, err := getAllLaporanMasyarakat(db)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get laporan masyarakat list", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "All laporan masyarakat retrieved successfully", getAllLaporanMasyarakatResponse{
			Data:  laporanList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Helper function to get laporan masyarakat by ID
func getLaporanMasyarakatById(db *sql.DB, id string) (laporanMasyarakatResponse, error) {
	var laporan laporanMasyarakatResponse
	var idMasyarakat, namaPelapor, emailPelapor sql.NullString
	var updatedDate sql.NullString

	query := `
        SELECT 
            lm.id, lm.id_masyarakat, lm.id_balita, lm.id_status_laporan,
            lm.tanggal_laporan, lm.hubungan_dengan_balita, lm.nomor_hp_pelapor, 
            lm.nomor_hp_keluarga_balita, lm.created_date, lm.updated_date,
            b.nama as nama_balita,
            k.nomor_kk, k.nama_ayah, k.nama_ibu, k.alamat,
            kel.kelurahan, kec.kecamatan,
            sl.status as status_laporan,
            m.nama as nama_pelapor,
            p.email as email_pelapor
        FROM laporan_masyarakat lm
        LEFT JOIN balita b ON lm.id_balita = b.id
        LEFT JOIN keluarga k ON b.id_keluarga = k.id
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        LEFT JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        LEFT JOIN masyarakat m ON lm.id_masyarakat = m.id
        LEFT JOIN pengguna p ON m.id_pengguna = p.id
        WHERE lm.id = ? AND lm.deleted_date IS NULL
    `

	err := db.QueryRow(query, id).Scan(
		&laporan.Id,
		&idMasyarakat,
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
		&namaPelapor,
		&emailPelapor,
	)

	if err != nil {
		return laporan, err
	}

	// Handle nullable fields
	if idMasyarakat.Valid {
		laporan.IdMasyarakat = idMasyarakat.String
		laporan.JenisLaporan = "masyarakat"
		if namaPelapor.Valid {
			laporan.NamaPelapor = namaPelapor.String
		}
		if emailPelapor.Valid {
			laporan.EmailPelapor = emailPelapor.String
		}
	} else {
		laporan.JenisLaporan = "admin"
		laporan.NamaPelapor = "Admin"
	}

	if updatedDate.Valid {
		laporan.UpdatedDate = updatedDate.String
	}

	return laporan, nil
}

// Helper function to get all laporan masyarakat
func getAllLaporanMasyarakat(db *sql.DB) ([]laporanMasyarakatResponse, int, error) {
	var laporanList []laporanMasyarakatResponse

	query := `
        SELECT 
            lm.id, lm.id_masyarakat, lm.id_balita, lm.id_status_laporan,
            lm.tanggal_laporan, lm.hubungan_dengan_balita, lm.nomor_hp_pelapor, 
            lm.nomor_hp_keluarga_balita, lm.created_date, lm.updated_date,
            b.nama as nama_balita,
            k.nomor_kk, k.nama_ayah, k.nama_ibu, k.alamat,
            kel.kelurahan, kec.kecamatan,
            sl.status as status_laporan,
            m.nama as nama_pelapor,
            p.email as email_pelapor
        FROM laporan_masyarakat lm
        LEFT JOIN balita b ON lm.id_balita = b.id AND b.deleted_date IS NULL
        LEFT JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        LEFT JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        LEFT JOIN masyarakat m ON lm.id_masyarakat = m.id
        LEFT JOIN pengguna p ON m.id_pengguna = p.id
        WHERE lm.deleted_date IS NULL
        ORDER BY lm.created_date DESC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var laporan laporanMasyarakatResponse
		var idMasyarakat, namaPelapor, emailPelapor sql.NullString
		var updatedDate sql.NullString

		err := rows.Scan(
			&laporan.Id,
			&idMasyarakat,
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
			&namaPelapor,
			&emailPelapor,
		)

		if err != nil {
			return nil, 0, err
		}

		// Handle nullable fields
		if idMasyarakat.Valid {
			laporan.IdMasyarakat = idMasyarakat.String
			laporan.JenisLaporan = "masyarakat"
			if namaPelapor.Valid {
				laporan.NamaPelapor = namaPelapor.String
			}
			if emailPelapor.Valid {
				laporan.EmailPelapor = emailPelapor.String
			}
		} else {
			laporan.JenisLaporan = "admin"
			laporan.NamaPelapor = "Admin"
		}

		if updatedDate.Valid {
			laporan.UpdatedDate = updatedDate.String
		}

		laporanList = append(laporanList, laporan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM laporan_masyarakat WHERE deleted_date IS NULL"
	err = db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return laporanList, total, nil
}