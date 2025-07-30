package api

import (
	"encoding/json"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

func AdminRiwayatPemeriksaanInsert(w http.ResponseWriter, r *http.Request) {
	adminID, role, err := object.ParseJWT(r.Header.Get("Authorization"))
	if err != nil || role != "admin" {
		response := object.NewResponse(http.StatusUnauthorized, "Unauthorized", nil)
		if err := response.WriteJson(w); err != nil {
			return
		}
	}

	if r.Method != http.MethodPost {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			return
		}
		return
	}

	var rp object.RiwayatPemeriksaan
	if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
		response := object.NewResponse(http.StatusBadRequest, "Invalid request body", nil)
		if err := response.WriteJson(w); err != nil {
			return
		}
		return
	}

	// Validasi minimal field wajib
	if err := rp.ValidateFields("IdBalita", "Tanggal", "BeratBadan", "TinggiBadan", "StatusGizi"); err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			return
		}
		return
	}

	db, err := object.ConnectDb()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Database connection error", nil)
		if err := response.WriteJson(w); err != nil {
			return
		}
		return
	}
	defer db.Close()

	query := `INSERT INTO riwayat_pemeriksaan 
        (id_balita, tanggal, berat_badan, tinggi_badan, status_gizi, keterangan, created_id, created_date) 
        VALUES (?, ?, ?, ?, ?, ?, ?, NOW())`
	_, err = db.Exec(query,
		rp.IdBalita,
		rp.Tanggal,
		rp.BeratBadan,
		rp.TinggiBadan,
		rp.StatusGizi,
		rp.Keterangan,
		adminID,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert riwayat pemeriksaan", nil)
		if err := response.WriteJson(w); err != nil {
			return
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Riwayat pemeriksaan berhasil ditambahkan", nil)
	if err := response.WriteJson(w); err != nil {
		return
	}
}
