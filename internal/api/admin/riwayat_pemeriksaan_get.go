package admin

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type riwayatPemeriksaanResponse struct {
	Id                  string `json:"id"`
	IdBalita            string `json:"id_balita"`
	NamaBalita          string `json:"nama_balita"`
	UmurBalita          string `json:"umur_balita"`
	JenisKelamin        string `json:"jenis_kelamin"`
	NamaAyah            string `json:"nama_ayah"`
	NamaIbu             string `json:"nama_ibu"`
	NomorKk             string `json:"nomor_kk"`
	IdIntervensi        string `json:"id_intervensi"`
	JenisIntervensi     string `json:"jenis_intervensi"`
	TanggalIntervensi   string `json:"tanggal_intervensi"`
	IdLaporanMasyarakat string `json:"id_laporan_masyarakat"` // <- Field baru
	StatusLaporan       string `json:"status_laporan"`        // <- Field baru
	TanggalLaporan      string `json:"tanggal_laporan"`       // <- Field baru
	JenisLaporan        string `json:"jenis_laporan"`         // <- Field baru (masyarakat/admin)
	Tanggal             string `json:"tanggal"`
	BeratBadan          string `json:"berat_badan"`
	TinggiBadan         string `json:"tinggi_badan"`
	StatusGizi          string `json:"status_gizi"`
	Keterangan          string `json:"keterangan"`
	Kelurahan           string `json:"kelurahan"`
	Kecamatan           string `json:"kecamatan"`
	CreatedDate         string `json:"created_date"`
	UpdatedDate         string `json:"updated_date,omitempty"`
	CreatedBy           string `json:"created_by,omitempty"`
	UpdatedBy           string `json:"updated_by,omitempty"`
}

type getAllRiwayatPemeriksaanResponse struct {
	Data  []riwayatPemeriksaanResponse `json:"data"`
	Total int                          `json:"total"`
}

type getRiwayatPemeriksaanByIdResponse struct {
	Data riwayatPemeriksaanResponse `json:"data"`
}

// # RiwayatPemeriksaanGet handles getting riwayat pemeriksaan data
//
// @Summary Get riwayat pemeriksaan data
// @Description Get riwayat pemeriksaan data based on query parameter (Admin only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without any parameter: Returns all riwayat pemeriksaan with total count
// @Description - With id parameter: Returns specific riwayat pemeriksaan data
// @Description - With id_balita parameter: Returns all riwayat pemeriksaan for specific balita
// @Description - With id_laporan_masyarakat parameter: Returns all riwayat pemeriksaan for specific laporan
// @Description - With id_intervensi parameter: Returns all riwayat pemeriksaan for specific intervensi
// @Description
// @Description Riwayat pemeriksaan data includes: balita info, intervensi info, laporan info, examination details, location info
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Riwayat Pemeriksaan ID"
// @Param id_balita query string false "Balita ID"
// @Param id_laporan_masyarakat query string false "Laporan Masyarakat ID"
// @Param id_intervensi query string false "Intervensi ID"
// @Success 200 {object} object.Response{data=getAllRiwayatPemeriksaanResponse} "Riwayat pemeriksaan data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Riwayat pemeriksaan not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/riwayat-pemeriksaan/get [get]
func RiwayatPemeriksaanGet(w http.ResponseWriter, r *http.Request) {
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

	// Check query parameters
	idParam := r.URL.Query().Get("id")
	idBalitaParam := r.URL.Query().Get("id_balita")
	idLaporanParam := r.URL.Query().Get("id_laporan_masyarakat")
	idIntervensiParam := r.URL.Query().Get("id_intervensi")

	if idParam != "" {
		// Get specific riwayat pemeriksaan by ID
		riwayat, err := getRiwayatPemeriksaanById(db, idParam)
		if err != nil {
			if err == sql.ErrNoRows {
				response := object.NewResponse(http.StatusNotFound, "Riwayat pemeriksaan not found", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get riwayat pemeriksaan", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Riwayat pemeriksaan retrieved successfully", getRiwayatPemeriksaanByIdResponse{
			Data: riwayat,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if idBalitaParam != "" {
		// Get all riwayat pemeriksaan for specific balita
		riwayatList, total, err := getRiwayatPemeriksaanByBalita(db, idBalitaParam)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get riwayat pemeriksaan by balita", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Riwayat pemeriksaan by balita retrieved successfully", getAllRiwayatPemeriksaanResponse{
			Data:  riwayatList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if idLaporanParam != "" {
		// Get all riwayat pemeriksaan for specific laporan masyarakat
		riwayatList, total, err := getRiwayatPemeriksaanByLaporan(db, idLaporanParam)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get riwayat pemeriksaan by laporan", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Riwayat pemeriksaan by laporan retrieved successfully", getAllRiwayatPemeriksaanResponse{
			Data:  riwayatList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if idIntervensiParam != "" {
		// Get all riwayat pemeriksaan for specific intervensi
		riwayatList, total, err := getRiwayatPemeriksaanByIntervensi(db, idIntervensiParam)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get riwayat pemeriksaan by intervensi", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Riwayat pemeriksaan by intervensi retrieved successfully", getAllRiwayatPemeriksaanResponse{
			Data:  riwayatList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Get all riwayat pemeriksaan
		riwayatList, total, err := getAllRiwayatPemeriksaan(db)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get riwayat pemeriksaan list", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "All riwayat pemeriksaan retrieved successfully", getAllRiwayatPemeriksaanResponse{
			Data:  riwayatList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Helper function to get riwayat pemeriksaan by ID
func getRiwayatPemeriksaanById(db *sql.DB, id string) (riwayatPemeriksaanResponse, error) {
	var riwayat riwayatPemeriksaanResponse
	var updatedDate sql.NullString
	var createdBy, updatedBy sql.NullString
	var idMasyarakatLaporan sql.NullString

	query := `
        SELECT 
            rp.id, rp.id_balita, rp.id_intervensi, rp.id_laporan_masyarakat, rp.tanggal,
            rp.berat_badan, rp.tinggi_badan, rp.status_gizi, rp.keterangan,
            rp.created_date, rp.updated_date,
            b.nama as nama_balita, b.jenis_kelamin,
            TIMESTAMPDIFF(MONTH, b.tanggal_lahir, CURDATE()) as umur_bulan,
            k.nomor_kk, k.nama_ayah, k.nama_ibu,
            kel.kelurahan, kec.kecamatan,
            i.jenis as jenis_intervensi, i.tanggal as tanggal_intervensi,
            lm.id_masyarakat, sl.status as status_laporan, lm.tanggal_laporan,
            pc.email as created_by,
            pu.email as updated_by
        FROM riwayat_pemeriksaan rp
        LEFT JOIN balita b ON rp.id_balita = b.id AND b.deleted_date IS NULL
        LEFT JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        LEFT JOIN intervensi i ON rp.id_intervensi = i.id AND i.deleted_date IS NULL
        LEFT JOIN laporan_masyarakat lm ON rp.id_laporan_masyarakat = lm.id AND lm.deleted_date IS NULL
        LEFT JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        LEFT JOIN pengguna pc ON rp.created_id = pc.id
        LEFT JOIN pengguna pu ON rp.updated_id = pu.id
        WHERE rp.id = ? AND rp.deleted_date IS NULL
    `

	err := db.QueryRow(query, id).Scan(
		&riwayat.Id,
		&riwayat.IdBalita,
		&riwayat.IdIntervensi,
		&riwayat.IdLaporanMasyarakat,
		&riwayat.Tanggal,
		&riwayat.BeratBadan,
		&riwayat.TinggiBadan,
		&riwayat.StatusGizi,
		&riwayat.Keterangan,
		&riwayat.CreatedDate,
		&updatedDate,
		&riwayat.NamaBalita,
		&riwayat.JenisKelamin,
		&riwayat.UmurBalita,
		&riwayat.NomorKk,
		&riwayat.NamaAyah,
		&riwayat.NamaIbu,
		&riwayat.Kelurahan,
		&riwayat.Kecamatan,
		&riwayat.JenisIntervensi,
		&riwayat.TanggalIntervensi,
		&idMasyarakatLaporan,
		&riwayat.StatusLaporan,
		&riwayat.TanggalLaporan,
		&createdBy,
		&updatedBy,
	)

	if err != nil {
		return riwayat, err
	}

	// Format umur balita
	riwayat.UmurBalita = formatUmurBalita(riwayat.UmurBalita)

	// Determine jenis laporan
	if idMasyarakatLaporan.Valid {
		riwayat.JenisLaporan = "masyarakat"
	} else {
		riwayat.JenisLaporan = "admin"
	}

	// Handle nullable fields
	if updatedDate.Valid {
		riwayat.UpdatedDate = updatedDate.String
	}
	if createdBy.Valid {
		riwayat.CreatedBy = createdBy.String
	}
	if updatedBy.Valid {
		riwayat.UpdatedBy = updatedBy.String
	}

	return riwayat, nil
}

// Helper function to get all riwayat pemeriksaan
func getAllRiwayatPemeriksaan(db *sql.DB) ([]riwayatPemeriksaanResponse, int, error) {
	var riwayatList []riwayatPemeriksaanResponse

	query := `
        SELECT 
            rp.id, rp.id_balita, rp.id_intervensi, rp.id_laporan_masyarakat, rp.tanggal,
            rp.berat_badan, rp.tinggi_badan, rp.status_gizi, rp.keterangan,
            rp.created_date, rp.updated_date,
            b.nama as nama_balita, b.jenis_kelamin,
            TIMESTAMPDIFF(MONTH, b.tanggal_lahir, CURDATE()) as umur_bulan,
            k.nomor_kk, k.nama_ayah, k.nama_ibu,
            kel.kelurahan, kec.kecamatan,
            i.jenis as jenis_intervensi, i.tanggal as tanggal_intervensi,
            lm.id_masyarakat, sl.status as status_laporan, lm.tanggal_laporan,
            pc.email as created_by,
            pu.email as updated_by
        FROM riwayat_pemeriksaan rp
        LEFT JOIN balita b ON rp.id_balita = b.id AND b.deleted_date IS NULL
        LEFT JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        LEFT JOIN intervensi i ON rp.id_intervensi = i.id AND i.deleted_date IS NULL
        LEFT JOIN laporan_masyarakat lm ON rp.id_laporan_masyarakat = lm.id AND lm.deleted_date IS NULL
        LEFT JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        LEFT JOIN pengguna pc ON rp.created_id = pc.id
        LEFT JOIN pengguna pu ON rp.updated_id = pu.id
        WHERE rp.deleted_date IS NULL
        ORDER BY rp.tanggal DESC, rp.created_date DESC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var riwayat riwayatPemeriksaanResponse
		var updatedDate sql.NullString
		var createdBy, updatedBy sql.NullString
		var idMasyarakatLaporan sql.NullString

		err := rows.Scan(
			&riwayat.Id,
			&riwayat.IdBalita,
			&riwayat.IdIntervensi,
			&riwayat.IdLaporanMasyarakat,
			&riwayat.Tanggal,
			&riwayat.BeratBadan,
			&riwayat.TinggiBadan,
			&riwayat.StatusGizi,
			&riwayat.Keterangan,
			&riwayat.CreatedDate,
			&updatedDate,
			&riwayat.NamaBalita,
			&riwayat.JenisKelamin,
			&riwayat.UmurBalita,
			&riwayat.NomorKk,
			&riwayat.NamaAyah,
			&riwayat.NamaIbu,
			&riwayat.Kelurahan,
			&riwayat.Kecamatan,
			&riwayat.JenisIntervensi,
			&riwayat.TanggalIntervensi,
			&idMasyarakatLaporan,
			&riwayat.StatusLaporan,
			&riwayat.TanggalLaporan,
			&createdBy,
			&updatedBy,
		)

		if err != nil {
			return nil, 0, err
		}

		// Format umur balita
		riwayat.UmurBalita = formatUmurBalita(riwayat.UmurBalita)

		// Determine jenis laporan
		if idMasyarakatLaporan.Valid {
			riwayat.JenisLaporan = "masyarakat"
		} else {
			riwayat.JenisLaporan = "admin"
		}

		// Handle nullable fields
		if updatedDate.Valid {
			riwayat.UpdatedDate = updatedDate.String
		}
		if createdBy.Valid {
			riwayat.CreatedBy = createdBy.String
		}
		if updatedBy.Valid {
			riwayat.UpdatedBy = updatedBy.String
		}

		riwayatList = append(riwayatList, riwayat)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM riwayat_pemeriksaan WHERE deleted_date IS NULL"
	err = db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return riwayatList, total, nil
}

// Helper function to format umur balita
func formatUmurBalita(umurBulan string) string {
	if umurBulan == "" {
		return "0 bulan"
	}

	// Parse umur in months
	months := 0
	fmt.Sscanf(umurBulan, "%d", &months)

	if months < 0 {
		return "0 bulan"
	}

	if months < 12 {
		return fmt.Sprintf("%d bulan", months)
	} else {
		years := months / 12
		remainingMonths := months % 12
		if remainingMonths == 0 {
			return fmt.Sprintf("%d tahun", years)
		}
		return fmt.Sprintf("%d tahun %d bulan", years, remainingMonths)
	}
}

// Helper function to get riwayat pemeriksaan by balita ID
func getRiwayatPemeriksaanByBalita(db *sql.DB, idBalita string) ([]riwayatPemeriksaanResponse, int, error) {
	var riwayatList []riwayatPemeriksaanResponse

	query := `
        SELECT 
            rp.id, rp.id_balita, rp.id_intervensi, rp.id_laporan_masyarakat, rp.tanggal,
            rp.berat_badan, rp.tinggi_badan, rp.status_gizi, rp.keterangan,
            rp.created_date, rp.updated_date,
            b.nama as nama_balita, b.jenis_kelamin,
            TIMESTAMPDIFF(MONTH, b.tanggal_lahir, CURDATE()) as umur_bulan,
            k.nomor_kk, k.nama_ayah, k.nama_ibu,
            kel.kelurahan, kec.kecamatan,
            i.jenis as jenis_intervensi, i.tanggal as tanggal_intervensi,
            lm.id_masyarakat, sl.status as status_laporan, lm.tanggal_laporan,
            pc.email as created_by,
            pu.email as updated_by
        FROM riwayat_pemeriksaan rp
        LEFT JOIN balita b ON rp.id_balita = b.id AND b.deleted_date IS NULL
        LEFT JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        LEFT JOIN intervensi i ON rp.id_intervensi = i.id AND i.deleted_date IS NULL
        LEFT JOIN laporan_masyarakat lm ON rp.id_laporan_masyarakat = lm.id AND lm.deleted_date IS NULL
        LEFT JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        LEFT JOIN pengguna pc ON rp.created_id = pc.id
        LEFT JOIN pengguna pu ON rp.updated_id = pu.id
        WHERE rp.id_balita = ? AND rp.deleted_date IS NULL
        ORDER BY rp.tanggal DESC, rp.created_date DESC
    `

	rows, err := db.Query(query, idBalita)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var riwayat riwayatPemeriksaanResponse
		var updatedDate sql.NullString
		var createdBy, updatedBy sql.NullString
		var idMasyarakatLaporan sql.NullString

		err := rows.Scan(
			&riwayat.Id,
			&riwayat.IdBalita,
			&riwayat.IdIntervensi,
			&riwayat.IdLaporanMasyarakat,
			&riwayat.Tanggal,
			&riwayat.BeratBadan,
			&riwayat.TinggiBadan,
			&riwayat.StatusGizi,
			&riwayat.Keterangan,
			&riwayat.CreatedDate,
			&updatedDate,
			&riwayat.NamaBalita,
			&riwayat.JenisKelamin,
			&riwayat.UmurBalita,
			&riwayat.NomorKk,
			&riwayat.NamaAyah,
			&riwayat.NamaIbu,
			&riwayat.Kelurahan,
			&riwayat.Kecamatan,
			&riwayat.JenisIntervensi,
			&riwayat.TanggalIntervensi,
			&idMasyarakatLaporan,
			&riwayat.StatusLaporan,
			&riwayat.TanggalLaporan,
			&createdBy,
			&updatedBy,
		)

		if err != nil {
			return nil, 0, err
		}

		// Format umur balita
		riwayat.UmurBalita = formatUmurBalita(riwayat.UmurBalita)

		// Determine jenis laporan
		if idMasyarakatLaporan.Valid {
			riwayat.JenisLaporan = "masyarakat"
		} else {
			riwayat.JenisLaporan = "admin"
		}

		// Handle nullable fields
		if updatedDate.Valid {
			riwayat.UpdatedDate = updatedDate.String
		}
		if createdBy.Valid {
			riwayat.CreatedBy = createdBy.String
		}
		if updatedBy.Valid {
			riwayat.UpdatedBy = updatedBy.String
		}

		riwayatList = append(riwayatList, riwayat)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return riwayatList, len(riwayatList), nil
}

// Helper function to get riwayat pemeriksaan by laporan masyarakat ID
func getRiwayatPemeriksaanByLaporan(db *sql.DB, idLaporan string) ([]riwayatPemeriksaanResponse, int, error) {
	var riwayatList []riwayatPemeriksaanResponse

	query := `
        SELECT 
            rp.id, rp.id_balita, rp.id_intervensi, rp.id_laporan_masyarakat, rp.tanggal,
            rp.berat_badan, rp.tinggi_badan, rp.status_gizi, rp.keterangan,
            rp.created_date, rp.updated_date,
            b.nama as nama_balita, b.jenis_kelamin,
            TIMESTAMPDIFF(MONTH, b.tanggal_lahir, CURDATE()) as umur_bulan,
            k.nomor_kk, k.nama_ayah, k.nama_ibu,
            kel.kelurahan, kec.kecamatan,
            i.jenis as jenis_intervensi, i.tanggal as tanggal_intervensi,
            lm.id_masyarakat, sl.status as status_laporan, lm.tanggal_laporan,
            pc.email as created_by,
            pu.email as updated_by
        FROM riwayat_pemeriksaan rp
        LEFT JOIN balita b ON rp.id_balita = b.id AND b.deleted_date IS NULL
        LEFT JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        LEFT JOIN intervensi i ON rp.id_intervensi = i.id AND i.deleted_date IS NULL
        LEFT JOIN laporan_masyarakat lm ON rp.id_laporan_masyarakat = lm.id AND lm.deleted_date IS NULL
        LEFT JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        LEFT JOIN pengguna pc ON rp.created_id = pc.id
        LEFT JOIN pengguna pu ON rp.updated_id = pu.id
        WHERE rp.id_laporan_masyarakat = ? AND rp.deleted_date IS NULL
        ORDER BY rp.tanggal DESC, rp.created_date DESC
    `

	rows, err := db.Query(query, idLaporan)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var riwayat riwayatPemeriksaanResponse
		var updatedDate sql.NullString
		var createdBy, updatedBy sql.NullString
		var idMasyarakatLaporan sql.NullString

		err := rows.Scan(
			&riwayat.Id,
			&riwayat.IdBalita,
			&riwayat.IdIntervensi,
			&riwayat.IdLaporanMasyarakat,
			&riwayat.Tanggal,
			&riwayat.BeratBadan,
			&riwayat.TinggiBadan,
			&riwayat.StatusGizi,
			&riwayat.Keterangan,
			&riwayat.CreatedDate,
			&updatedDate,
			&riwayat.NamaBalita,
			&riwayat.JenisKelamin,
			&riwayat.UmurBalita,
			&riwayat.NomorKk,
			&riwayat.NamaAyah,
			&riwayat.NamaIbu,
			&riwayat.Kelurahan,
			&riwayat.Kecamatan,
			&riwayat.JenisIntervensi,
			&riwayat.TanggalIntervensi,
			&idMasyarakatLaporan,
			&riwayat.StatusLaporan,
			&riwayat.TanggalLaporan,
			&createdBy,
			&updatedBy,
		)

		if err != nil {
			return nil, 0, err
		}

		// Format umur balita
		riwayat.UmurBalita = formatUmurBalita(riwayat.UmurBalita)

		// Determine jenis laporan
		if idMasyarakatLaporan.Valid {
			riwayat.JenisLaporan = "masyarakat"
		} else {
			riwayat.JenisLaporan = "admin"
		}

		// Handle nullable fields
		if updatedDate.Valid {
			riwayat.UpdatedDate = updatedDate.String
		}
		if createdBy.Valid {
			riwayat.CreatedBy = createdBy.String
		}
		if updatedBy.Valid {
			riwayat.UpdatedBy = updatedBy.String
		}

		riwayatList = append(riwayatList, riwayat)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return riwayatList, len(riwayatList), nil
}

// Helper function to get riwayat pemeriksaan by intervensi ID
func getRiwayatPemeriksaanByIntervensi(db *sql.DB, idIntervensi string) ([]riwayatPemeriksaanResponse, int, error) {
	var riwayatList []riwayatPemeriksaanResponse

	query := `
        SELECT 
            rp.id, rp.id_balita, rp.id_intervensi, rp.id_laporan_masyarakat, rp.tanggal,
            rp.berat_badan, rp.tinggi_badan, rp.status_gizi, rp.keterangan,
            rp.created_date, rp.updated_date,
            b.nama as nama_balita, b.jenis_kelamin,
            TIMESTAMPDIFF(MONTH, b.tanggal_lahir, CURDATE()) as umur_bulan,
            k.nomor_kk, k.nama_ayah, k.nama_ibu,
            kel.kelurahan, kec.kecamatan,
            i.jenis as jenis_intervensi, i.tanggal as tanggal_intervensi,
            lm.id_masyarakat, sl.status as status_laporan, lm.tanggal_laporan,
            pc.email as created_by,
            pu.email as updated_by
        FROM riwayat_pemeriksaan rp
        LEFT JOIN balita b ON rp.id_balita = b.id AND b.deleted_date IS NULL
        LEFT JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        LEFT JOIN intervensi i ON rp.id_intervensi = i.id AND i.deleted_date IS NULL
        LEFT JOIN laporan_masyarakat lm ON rp.id_laporan_masyarakat = lm.id AND lm.deleted_date IS NULL
        LEFT JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        LEFT JOIN pengguna pc ON rp.created_id = pc.id
        LEFT JOIN pengguna pu ON rp.updated_id = pu.id
        WHERE rp.id_intervensi = ? AND rp.deleted_date IS NULL
        ORDER BY rp.tanggal DESC, rp.created_date DESC
    `

	rows, err := db.Query(query, idIntervensi)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var riwayat riwayatPemeriksaanResponse
		var updatedDate sql.NullString
		var createdBy, updatedBy sql.NullString
		var idMasyarakatLaporan sql.NullString

		err := rows.Scan(
			&riwayat.Id,
			&riwayat.IdBalita,
			&riwayat.IdIntervensi,
			&riwayat.IdLaporanMasyarakat,
			&riwayat.Tanggal,
			&riwayat.BeratBadan,
			&riwayat.TinggiBadan,
			&riwayat.StatusGizi,
			&riwayat.Keterangan,
			&riwayat.CreatedDate,
			&updatedDate,
			&riwayat.NamaBalita,
			&riwayat.JenisKelamin,
			&riwayat.UmurBalita,
			&riwayat.NomorKk,
			&riwayat.NamaAyah,
			&riwayat.NamaIbu,
			&riwayat.Kelurahan,
			&riwayat.Kecamatan,
			&riwayat.JenisIntervensi,
			&riwayat.TanggalIntervensi,
			&idMasyarakatLaporan,
			&riwayat.StatusLaporan,
			&riwayat.TanggalLaporan,
			&createdBy,
			&updatedBy,
		)

		if err != nil {
			return nil, 0, err
		}

		// Format umur balita
		riwayat.UmurBalita = formatUmurBalita(riwayat.UmurBalita)

		// Determine jenis laporan
		if idMasyarakatLaporan.Valid {
			riwayat.JenisLaporan = "masyarakat"
		} else {
			riwayat.JenisLaporan = "admin"
		}

		// Handle nullable fields
		if updatedDate.Valid {
			riwayat.UpdatedDate = updatedDate.String
		}
		if createdBy.Valid {
			riwayat.CreatedBy = createdBy.String
		}
		if updatedBy.Valid {
			riwayat.UpdatedBy = updatedBy.String
		}

		riwayatList = append(riwayatList, riwayat)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return riwayatList, len(riwayatList), nil
}
