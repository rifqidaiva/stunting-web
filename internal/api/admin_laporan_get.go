package api

import (
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// AdminGet handles the retrieval of all reports by community members and admin users.
func AdminGet(w http.ResponseWriter, r *http.Request) {
	_, role, err := object.ParseJWT(r.Header.Get("Authorization"))
	if err != nil || role != "admin" {
		response := object.NewResponse(http.StatusUnauthorized, "Unauthorized", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if r.Method != http.MethodGet {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	db, err := object.ConnectDb()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Database connection error", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	var reports []object.Report

	query := `
    SELECT
        l.id_pengguna, l.id_balita, l.id_status_laporan, l.tanggal_laporan,
        l.hubungan_dengan_balita, l.nomor_hp_pelapor, l.nomor_hp_keluarga_balita,
        l.created_id, l.created_date, l.updated_id, l.updated_date, l.deleted_id, l.deleted_date,

        b.id, b.id_keluarga, b.nama, b.tanggal_lahir, b.jenis_kelamin, b.berat_lahir, b.tinggi_lahir,
        b.created_id, b.created_date, b.updated_id, b.updated_date, b.deleted_id, b.deleted_date,

        k.id, k.nomor_kk, k.nama_ayah, k.nama_ibu, k.nik_ayah, k.nik_ibu, k.alamat, k.rt, k.rw, k.id_kelurahan, k.koordinat,
        k.created_id, k.created_date, k.updated_id, k.updated_date, k.deleted_id, k.deleted_date

    FROM laporan_masyarakat l
    JOIN balita b ON l.id_balita = b.id
    JOIN keluarga k ON b.id_keluarga = k.id
    WHERE l.deleted_date IS NULL`

	rows, err := db.Query(query)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to query reports", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var report object.Report
		var koordinatWKT string

		err := rows.Scan(
			// LaporanMasyarakat
			&report.LaporanMasyarakat.IdPengguna,
			&report.LaporanMasyarakat.IdBalita,
			&report.LaporanMasyarakat.IdStatusLaporan,
			&report.LaporanMasyarakat.TanggalLaporan,
			&report.LaporanMasyarakat.HubunganDenganBalita,
			&report.LaporanMasyarakat.NomorHpPelapor,
			&report.LaporanMasyarakat.NomorHpKeluargaBalita,
			&report.LaporanMasyarakat.CreatedId,
			&report.LaporanMasyarakat.CreatedDate,
			&report.LaporanMasyarakat.UpdatedId,
			&report.LaporanMasyarakat.UpdatedDate,
			&report.LaporanMasyarakat.DeletedId,
			&report.LaporanMasyarakat.DeletedDate,

			// Balita
			&report.Balita.Id,
			&report.Balita.IdKeluarga,
			&report.Balita.Nama,
			&report.Balita.TanggalLahir,
			&report.Balita.JenisKelamin,
			&report.Balita.BeratLahir,
			&report.Balita.TinggiLahir,
			&report.Balita.CreatedId,
			&report.Balita.CreatedDate,
			&report.Balita.UpdatedId,
			&report.Balita.UpdatedDate,
			&report.Balita.DeletedId,
			&report.Balita.DeletedDate,

			// Keluarga
			&report.Keluarga.Id,
			&report.Keluarga.NomorKk,
			&report.Keluarga.NamaAyah,
			&report.Keluarga.NamaIbu,
			&report.Keluarga.NikAyah,
			&report.Keluarga.NikIbu,
			&report.Keluarga.Alamat,
			&report.Keluarga.Rt,
			&report.Keluarga.Rw,
			&report.Keluarga.IdKelurahan,
			&koordinatWKT,
			&report.Keluarga.CreatedId,
			&report.Keluarga.CreatedDate,
			&report.Keluarga.UpdatedId,
			&report.Keluarga.UpdatedDate,
			&report.Keluarga.DeletedId,
			&report.Keluarga.DeletedDate,
		)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to scan report", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		report.Keluarga.Koordinat = object.ParseWKT(koordinatWKT)
		reports = append(reports, report)
	}

	response := object.NewResponse(http.StatusOK, "Reports retrieved successfully", reports)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
