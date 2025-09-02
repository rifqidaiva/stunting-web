package community

import (
	"database/sql"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type kecamatanResponse struct {
	Id        string `json:"id"`
	Kecamatan string `json:"kecamatan"`
}

type getAllKecamatanResponse struct {
	Data  []kecamatanResponse `json:"data"`
	Total int                 `json:"total"`
}

// # KecamatanGet handles getting all kecamatan for master data (Community)
//
// @Summary Get kecamatan master data (Community)
// @Description Get all kecamatan for dropdown/reference (Masyarakat only)
// @Tags community
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object.Response{data=getAllKecamatanResponse} "Kecamatan data retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Masyarakat role required"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/community/master-kecamatan [get]
func KecamatanGet(w http.ResponseWriter, r *http.Request) {
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

	// Get all kecamatan
	kecamatanList, total, err := getAllKecamatan(db)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get kecamatan list", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Kecamatan master data retrieved successfully", getAllKecamatanResponse{
		Data:  kecamatanList,
		Total: total,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type kelurahanResponse struct {
	Id          string `json:"id"`
	IdKecamatan string `json:"id_kecamatan"`
	Kelurahan   string `json:"kelurahan"`
	Kecamatan   string `json:"kecamatan"`
}

type getAllKelurahanResponse struct {
	Data  []kelurahanResponse `json:"data"`
	Total int                 `json:"total"`
}

// # KelurahanGet handles getting all kelurahan for master data (Community)
//
// @Summary Get kelurahan master data (Community)
// @Description Get kelurahan data with optional kecamatan filter (Masyarakat only)
// @Tags community
// @Accept json
// @Produce json
// @Security Bearer
// @Param id_kecamatan query string false "Filter by Kecamatan ID"
// @Success 200 {object} object.Response{data=getAllKelurahanResponse} "Kelurahan data retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Masyarakat role required"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/community/master-kelurahan [get]
func KelurahanGet(w http.ResponseWriter, r *http.Request) {
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

	// Check for kecamatan filter
	idKecamatanParam := r.URL.Query().Get("id_kecamatan")

	// Get kelurahan data
	kelurahanList, total, err := getAllKelurahan(db, idKecamatanParam)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get kelurahan list", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	message := "Kelurahan master data retrieved successfully"
	if idKecamatanParam != "" {
		message = "Kelurahan data by kecamatan retrieved successfully"
	}

	response := object.NewResponse(http.StatusOK, message, getAllKelurahanResponse{
		Data:  kelurahanList,
		Total: total,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type statusLaporanResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type getAllStatusLaporanResponse struct {
	Data  []statusLaporanResponse `json:"data"`
	Total int                     `json:"total"`
}

// # StatusLaporanGet handles getting all status laporan for master data (Community)
//
// @Summary Get status laporan master data (Community)
// @Description Get all status laporan for reference/display (Masyarakat only)
// @Description This is primarily for display purposes to show status meanings
// @Tags community
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object.Response{data=getAllStatusLaporanResponse} "Status laporan data retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Masyarakat role required"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/community/master-status-laporan [get]
func StatusLaporanGet(w http.ResponseWriter, r *http.Request) {
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

	// Get all status laporan
	statusList, total, err := getAllStatusLaporan(db)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get status laporan list", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Status laporan master data retrieved successfully", getAllStatusLaporanResponse{
		Data:  statusList,
		Total: total,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Helper function to get all kecamatan
func getAllKecamatan(db *sql.DB) ([]kecamatanResponse, int, error) {
	var kecamatanList []kecamatanResponse

	query := "SELECT id, kecamatan FROM kecamatan ORDER BY kecamatan ASC"

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var kecamatan kecamatanResponse
		err := rows.Scan(&kecamatan.Id, &kecamatan.Kecamatan)
		if err != nil {
			return nil, 0, err
		}
		kecamatanList = append(kecamatanList, kecamatan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return kecamatanList, len(kecamatanList), nil
}

// Helper function to get all kelurahan with optional kecamatan filter
func getAllKelurahan(db *sql.DB, idKecamatan string) ([]kelurahanResponse, int, error) {
	var kelurahanList []kelurahanResponse
	var query string
	var args []any

	if idKecamatan != "" {
		query = `
            SELECT kel.id, kel.id_kecamatan, kel.kelurahan, kec.kecamatan
            FROM kelurahan kel
            LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
            WHERE kel.id_kecamatan = ?
            ORDER BY kel.kelurahan ASC
        `
		args = append(args, idKecamatan)
	} else {
		query = `
            SELECT kel.id, kel.id_kecamatan, kel.kelurahan, kec.kecamatan
            FROM kelurahan kel
            LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
            ORDER BY kec.kecamatan ASC, kel.kelurahan ASC
        `
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var kelurahan kelurahanResponse
		err := rows.Scan(&kelurahan.Id, &kelurahan.IdKecamatan, &kelurahan.Kelurahan, &kelurahan.Kecamatan)
		if err != nil {
			return nil, 0, err
		}
		kelurahanList = append(kelurahanList, kelurahan)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return kelurahanList, len(kelurahanList), nil
}

// Helper function to get all status laporan
func getAllStatusLaporan(db *sql.DB) ([]statusLaporanResponse, int, error) {
	var statusList []statusLaporanResponse

	query := "SELECT id, status FROM status_laporan ORDER BY id ASC"

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var status statusLaporanResponse
		err := rows.Scan(&status.Id, &status.Status)
		if err != nil {
			return nil, 0, err
		}
		statusList = append(statusList, status)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return statusList, len(statusList), nil
}
