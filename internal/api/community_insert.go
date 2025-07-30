package api

import (
	"encoding/json"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// CommunityInsert handles the insertion of new reports by masyarakat (community members).
func CommunityInsert(w http.ResponseWriter, r *http.Request) {
	userID, role, err := object.ParseJWT(r.Header.Get("Authorization"))
	if err != nil || role != "masyarakat" {
		response := object.NewResponse(http.StatusUnauthorized, "Unauthorized", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodPost {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var report object.Report
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		response := object.NewResponse(http.StatusBadRequest, "Invalid request body", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = report.Balita.ValidateFields("Nama", "TanggalLahir", "JenisKelamin", "BeratLahir", "TinggiLahir")
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = report.Keluarga.ValidateFields("NomorKK", "NamaAyah", "NamaIbu", "NikAyah", "NikIbu", "Alamat", "Rt", "Rw", "IdKelurahan", "Koordinat")
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = report.LaporanMasyarakat.ValidateFields("NomorHpPelapor", "NomorHpKeluargaBalita", "HubunganDenganBalita")
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Insert the report into the database
	db, err := object.ConnectDb()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Database connection error", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	// 1. Insert ke tabel balita
	balitaQuery := `INSERT INTO balita (nama, tanggal_lahir, jenis_kelamin, berat_lahir, tinggi_lahir, created_id, created_date) VALUES (?, ?, ?, ?, ?, ?, NOW())`
	balitaStmt, err := db.Prepare(balitaQuery)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to prepare balita statement", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer balitaStmt.Close()

	resBalita, err := balitaStmt.Exec(
		report.Balita.Nama,
		report.Balita.TanggalLahir,
		report.Balita.JenisKelamin,
		report.Balita.BeratLahir,
		report.Balita.TinggiLahir,
		userID, // created_id
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert balita", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	idBalita, err := resBalita.LastInsertId()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get balita ID", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 2. Insert ke tabel keluarga
	keluargaQuery := `INSERT INTO keluarga (nomor_kk, nama_ayah, nama_ibu, nik_ayah, nik_ibu, alamat, rt, rw, id_kelurahan, koordinat, created_id, created_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`
	keluargaStmt, err := db.Prepare(keluargaQuery)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to prepare keluarga statement", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer keluargaStmt.Close()

	resKeluarga, err := keluargaStmt.Exec(
		report.Keluarga.NomorKk,
		report.Keluarga.NamaAyah,
		report.Keluarga.NamaIbu,
		report.Keluarga.NikAyah,
		report.Keluarga.NikIbu,
		report.Keluarga.Alamat,
		report.Keluarga.Rt,
		report.Keluarga.Rw,
		report.Keluarga.IdKelurahan,
		object.ToWKT(report.Keluarga.Koordinat),
		userID, // created_id
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert keluarga", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	idKeluarga, err := resKeluarga.LastInsertId()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get keluarga ID", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// 3. Insert ke tabel laporan
	laporanQuery := `INSERT INTO laporan_masyarakat (id_pengguna, id_balita, id_keluarga, nomor_hp_pelapor, nomor_hp_keluarga_balita, hubungan_dengan_balita, created_id, created_date) VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`
	laporanStmt, err := db.Prepare(laporanQuery)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to prepare laporan statement", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer laporanStmt.Close()

	_, err = laporanStmt.Exec(
		userID,
		idBalita,
		idKeluarga,
		report.LaporanMasyarakat.NomorHpPelapor,
		report.LaporanMasyarakat.NomorHpKeluargaBalita,
		report.LaporanMasyarakat.HubunganDenganBalita,
		userID, // created_id
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert laporan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Report inserted successfully", nil)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
