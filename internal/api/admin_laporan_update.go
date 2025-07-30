package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

func AdminUpdate(w http.ResponseWriter, r *http.Request) {
	adminID, role, err := object.ParseJWT(r.Header.Get("Authorization"))
	if err != nil || role != "admin" {
		response := object.NewResponse(http.StatusUnauthorized, "Unauthorized", nil)
		if err := response.WriteJson(w); err != nil {
			return
		}
		return
	}

	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			return
		}
		return
	}

	var req struct {
		Pengguna          map[string]any `json:"pengguna,omitempty"`
		Balita            map[string]any `json:"balita,omitempty"`
		Keluarga          map[string]any `json:"keluarga,omitempty"`
		LaporanMasyarakat map[string]any `json:"laporan_masyarakat,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := object.NewResponse(http.StatusBadRequest, "Invalid request body", nil)
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

	// Helper: dynamic update
	updateTable := func(table string, idField string, data map[string]any) error {
		id, ok := data[idField]
		if !ok {
			return nil // skip if no id
		}
		delete(data, idField)
		if len(data) == 0 {
			return nil // nothing to update
		}
		// Tambahkan updated_id dan updated_date
		data["updated_id"] = adminID
		data["updated_date"] = nil // gunakan NOW() di SQL

		var sets []string
		var args []any
		for k, v := range data {
			if k == "updated_date" {
				sets = append(sets, k+" = NOW()")
			} else {
				sets = append(sets, k+" = ?")
				args = append(args, v)
			}
		}
		args = append(args, id)
		query := "UPDATE " + table + " SET " + strings.Join(sets, ", ") + " WHERE " + idField + " = ?"
		_, err := db.Exec(query, args...)
		return err
	}

	// Pengguna
	if req.Pengguna != nil {
		if err := updateTable("pengguna", "id", req.Pengguna); err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to update pengguna", nil)
			if err := response.WriteJson(w); err != nil {
				return
			}
			return
		}
	}

	// Balita
	if req.Balita != nil {
		if err := updateTable("balita", "id", req.Balita); err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to update balita", nil)
			if err := response.WriteJson(w); err != nil {
				return
			}
			return
		}
	}

	// Keluarga
	if req.Keluarga != nil {
		if err := updateTable("keluarga", "id", req.Keluarga); err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to update keluarga", nil)
			if err := response.WriteJson(w); err != nil {
				return
			}
			return
		}
	}

	// LaporanMasyarakat
	if req.LaporanMasyarakat != nil {
		if err := updateTable("laporan_masyarakat", "id", req.LaporanMasyarakat); err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to update laporan masyarakat", nil)
			if err := response.WriteJson(w); err != nil {
				return
			}
			return
		}
	}

	response := object.NewResponse(http.StatusOK, "Update successful", nil)
	if err := response.WriteJson(w); err != nil {
		return
	}
}
