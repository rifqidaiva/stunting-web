package admin

import (
	"database/sql"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// ==================== STATUS LAPORAN ====================

type statusLaporanResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type getAllStatusLaporanResponse struct {
	Data  []statusLaporanResponse `json:"data"`
	Total int                     `json:"total"`
}

// # StatusLaporanGet handles getting all status laporan for master data
//
// @Summary Get status laporan master data
// @Description Get all status laporan for dropdown/reference (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object.Response{data=getAllStatusLaporanResponse} "Status laporan data retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/master-status-laporan [get]
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

// ==================== MASYARAKAT ====================

type masyarakatResponse struct {
	Id     string `json:"id"`
	Nama   string `json:"nama"`
	Alamat string `json:"alamat"`
	Email  string `json:"email"`
}

type getAllMasyarakatResponse struct {
	Data  []masyarakatResponse `json:"data"`
	Total int                  `json:"total"`
}

// # MasyarakatGet handles getting all masyarakat for master data
//
// @Summary Get masyarakat master data
// @Description Get all masyarakat for dropdown/reference (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object.Response{data=getAllMasyarakatResponse} "Masyarakat data retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/master-masyarakat [get]
func MasyarakatGet(w http.ResponseWriter, r *http.Request) {
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

	// Get all masyarakat
	masyarakatList, total, err := getAllMasyarakat(db)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get masyarakat list", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Masyarakat master data retrieved successfully", getAllMasyarakatResponse{
		Data:  masyarakatList,
		Total: total,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ==================== KECAMATAN ====================

type kecamatanResponse struct {
	Id        string `json:"id"`
	Kecamatan string `json:"kecamatan"`
}

type getAllKecamatanResponse struct {
	Data  []kecamatanResponse `json:"data"`
	Total int                 `json:"total"`
}

// # KecamatanGet handles getting all kecamatan for master data
//
// @Summary Get kecamatan master data
// @Description Get all kecamatan for dropdown/reference (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object.Response{data=getAllKecamatanResponse} "Kecamatan data retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/master-kecamatan [get]
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

// ==================== KELURAHAN ====================

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

// # KelurahanGet handles getting all kelurahan for master data
//
// @Summary Get kelurahan master data
// @Description Get kelurahan data with optional kecamatan filter (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id_kecamatan query string false "Filter by Kecamatan ID"
// @Success 200 {object} object.Response{data=getAllKelurahanResponse} "Kelurahan data retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/master-kelurahan [get]
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

// ==================== SKPD MASTER ====================

type skpdMasterResponse struct {
	Id    string `json:"id"`
	Skpd  string `json:"skpd"`
	Jenis string `json:"jenis"`
}

type getAllSkpdMasterResponse struct {
	Data  []skpdMasterResponse `json:"data"`
	Total int                  `json:"total"`
}

// # SkpdMasterGet handles getting all SKPD for master data
//
// @Summary Get SKPD master data
// @Description Get SKPD data with optional jenis filter (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param jenis query string false "Filter by SKPD jenis (puskesmas, kelurahan, skpd)"
// @Success 200 {object} object.Response{data=getAllSkpdMasterResponse} "SKPD data retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/master-skpd [get]
func SkpdMasterGet(w http.ResponseWriter, r *http.Request) {
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

	// Check for jenis filter
	jenisParam := r.URL.Query().Get("jenis")

	// Get SKPD data
	skpdList, total, err := getAllSkpdMaster(db, jenisParam)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get SKPD list", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	message := "SKPD master data retrieved successfully"
	if jenisParam != "" {
		message = "SKPD data by jenis retrieved successfully"
	}

	response := object.NewResponse(http.StatusOK, message, getAllSkpdMasterResponse{
		Data:  skpdList,
		Total: total,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ==================== GEOJSON ENDPOINTS ====================

// # KecamatanGeoJSONGet handles getting kecamatan area as GeoJSON
//
// @Summary Get kecamatan area GeoJSON
// @Description Get kecamatan boundary areas as GeoJSON MultiPolygon (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Specific Kecamatan ID"
// @Success 200 {object} object.Response{data=object.GeoJSONFeatureCollection} "Kecamatan GeoJSON retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/geojson-kecamatan [get]
func KecamatanGeoJSONGet(w http.ResponseWriter, r *http.Request) {
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

	// Check for specific kecamatan filter
	idParam := r.URL.Query().Get("id")

	// Get kecamatan GeoJSON
	geoJSONCollection, err := getKecamatanGeoJSON(db, idParam)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get kecamatan GeoJSON", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	message := "Kecamatan GeoJSON retrieved successfully"
	if idParam != "" {
		message = "Kecamatan GeoJSON by ID retrieved successfully"
	}

	response := object.NewResponse(http.StatusOK, message, geoJSONCollection)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// # KelurahanGeoJSONGet handles getting kelurahan area as GeoJSON
//
// @Summary Get kelurahan area GeoJSON
// @Description Get kelurahan boundary areas as GeoJSON MultiPolygon (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Specific Kelurahan ID"
// @Param id_kecamatan query string false "Filter by Kecamatan ID"
// @Success 200 {object} object.Response{data=object.GeoJSONFeatureCollection} "Kelurahan GeoJSON retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/geojson-kelurahan [get]
func KelurahanGeoJSONGet(w http.ResponseWriter, r *http.Request) {
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

	// Check for filters
	idParam := r.URL.Query().Get("id")
	idKecamatanParam := r.URL.Query().Get("id_kecamatan")

	// Get kelurahan GeoJSON
	geoJSONCollection, err := getKelurahanGeoJSON(db, idParam, idKecamatanParam)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get kelurahan GeoJSON", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	message := "Kelurahan GeoJSON retrieved successfully"
	if idParam != "" {
		message = "Kelurahan GeoJSON by ID retrieved successfully"
	} else if idKecamatanParam != "" {
		message = "Kelurahan GeoJSON by kecamatan retrieved successfully"
	}

	response := object.NewResponse(http.StatusOK, message, geoJSONCollection)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// # BalitaPointsGeoJSONGet handles getting balita points as GeoJSON
//
// @Summary Get balita points GeoJSON
// @Description Get balita locations as GeoJSON points with status laporan (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param status_laporan query string false "Filter by status laporan"
// @Param id_kecamatan query string false "Filter by kecamatan"
// @Param id_kelurahan query string false "Filter by kelurahan"
// @Success 200 {object} object.Response{data=object.GeoJSONFeatureCollection} "Balita points GeoJSON retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/geojson-balita-points [get]
func BalitaPointsGeoJSONGet(w http.ResponseWriter, r *http.Request) {
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

	// Check for filters
	statusLaporanParam := r.URL.Query().Get("status_laporan")
	idKecamatanParam := r.URL.Query().Get("id_kecamatan")
	idKelurahanParam := r.URL.Query().Get("id_kelurahan")

	// Get balita points GeoJSON
	geoJSONCollection, err := getBalitaPointsGeoJSON(db, statusLaporanParam, idKecamatanParam, idKelurahanParam)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get balita points GeoJSON", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	message := "Balita points GeoJSON retrieved successfully"

	response := object.NewResponse(http.StatusOK, message, geoJSONCollection)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ==================== HELPER FUNCTIONS ====================

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

// Helper function to get all masyarakat
func getAllMasyarakat(db *sql.DB) ([]masyarakatResponse, int, error) {
	var masyarakatList []masyarakatResponse

	query := `
        SELECT m.id, m.nama, m.alamat, p.email
        FROM masyarakat m
        LEFT JOIN pengguna p ON m.id_pengguna = p.id
        ORDER BY m.nama ASC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var masyarakat masyarakatResponse
		err := rows.Scan(&masyarakat.Id, &masyarakat.Nama, &masyarakat.Alamat, &masyarakat.Email)
		if err != nil {
			return nil, 0, err
		}
		masyarakatList = append(masyarakatList, masyarakat)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return masyarakatList, len(masyarakatList), nil
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

// Helper function to get all SKPD with optional jenis filter
func getAllSkpdMaster(db *sql.DB, jenis string) ([]skpdMasterResponse, int, error) {
	var skpdList []skpdMasterResponse
	var query string
	var args []any

	if jenis != "" {
		query = `
            SELECT id, skpd, jenis
            FROM skpd
            WHERE jenis = ? AND deleted_date IS NULL
            ORDER BY skpd ASC
        `
		args = append(args, jenis)
	} else {
		query = `
            SELECT id, skpd, jenis
            FROM skpd
            WHERE deleted_date IS NULL
            ORDER BY jenis ASC, skpd ASC
        `
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var skpd skpdMasterResponse
		err := rows.Scan(&skpd.Id, &skpd.Skpd, &skpd.Jenis)
		if err != nil {
			return nil, 0, err
		}
		skpdList = append(skpdList, skpd)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return skpdList, len(skpdList), nil
}

// Helper function to get kecamatan GeoJSON
func getKecamatanGeoJSON(db *sql.DB, idKecamatan string) (object.GeoJSONFeatureCollection, error) {
	var features []object.GeoJSONFeature
	var query string
	var args []any

	if idKecamatan != "" {
		query = `
            SELECT id, kecamatan, ST_AsText(area) as area_wkt
            FROM kecamatan
            WHERE id = ? AND area IS NOT NULL
        `
		args = append(args, idKecamatan)
	} else {
		query = `
            SELECT id, kecamatan, ST_AsText(area) as area_wkt
            FROM kecamatan
            WHERE area IS NOT NULL
            ORDER BY kecamatan ASC
        `
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return object.GeoJSONFeatureCollection{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, kecamatan, areaWKT string
		err := rows.Scan(&id, &kecamatan, &areaWKT)
		if err != nil {
			return object.GeoJSONFeatureCollection{}, err
		}

		// Create properties
		properties := map[string]any{
			"id":        id,
			"kecamatan": kecamatan,
			"type":      "kecamatan",
		}

		// Create GeoJSON feature
		feature, err := object.CreateGeoJSONFeature(areaWKT, properties)
		if err != nil {
			continue // Skip if unable to parse geometry
		}

		features = append(features, feature)
	}

	if err = rows.Err(); err != nil {
		return object.GeoJSONFeatureCollection{}, err
	}

	return object.CreateGeoJSONFeatureCollection(features), nil
}

// Helper function to get kelurahan GeoJSON
func getKelurahanGeoJSON(db *sql.DB, idKelurahan, idKecamatan string) (object.GeoJSONFeatureCollection, error) {
	var features []object.GeoJSONFeature
	var query string
	var args []any

	if idKelurahan != "" {
		query = `
            SELECT kel.id, kel.kelurahan, kec.kecamatan, ST_AsText(kel.area) as area_wkt
            FROM kelurahan kel
            LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
            WHERE kel.id = ? AND kel.area IS NOT NULL
        `
		args = append(args, idKelurahan)
	} else if idKecamatan != "" {
		query = `
            SELECT kel.id, kel.kelurahan, kec.kecamatan, ST_AsText(kel.area) as area_wkt
            FROM kelurahan kel
            LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
            WHERE kel.id_kecamatan = ? AND kel.area IS NOT NULL
            ORDER BY kel.kelurahan ASC
        `
		args = append(args, idKecamatan)
	} else {
		query = `
            SELECT kel.id, kel.kelurahan, kec.kecamatan, ST_AsText(kel.area) as area_wkt
            FROM kelurahan kel
            LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
            WHERE kel.area IS NOT NULL
            ORDER BY kec.kecamatan ASC, kel.kelurahan ASC
        `
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return object.GeoJSONFeatureCollection{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, kelurahan, kecamatan, areaWKT string
		err := rows.Scan(&id, &kelurahan, &kecamatan, &areaWKT)
		if err != nil {
			return object.GeoJSONFeatureCollection{}, err
		}

		// Create properties
		properties := map[string]any{
			"id":        id,
			"kelurahan": kelurahan,
			"kecamatan": kecamatan,
			"type":      "kelurahan",
		}

		// Create GeoJSON feature
		feature, err := object.CreateGeoJSONFeature(areaWKT, properties)
		if err != nil {
			continue // Skip if unable to parse geometry
		}

		features = append(features, feature)
	}

	if err = rows.Err(); err != nil {
		return object.GeoJSONFeatureCollection{}, err
	}

	return object.CreateGeoJSONFeatureCollection(features), nil
}

// Helper function to get balita points GeoJSON with status laporan
func getBalitaPointsGeoJSON(db *sql.DB, statusLaporan, idKecamatan, idKelurahan string) (object.GeoJSONFeatureCollection, error) {
	var features []object.GeoJSONFeature
	var query string
	var args []any

	query = `
        SELECT DISTINCT
            b.id, b.nama, b.jenis_kelamin,
            TIMESTAMPDIFF(MONTH, b.tanggal_lahir, CURDATE()) as umur_bulan,
            k.nomor_kk, k.nama_ayah, k.nama_ibu,
            kel.kelurahan, kec.kecamatan,
            ST_AsText(k.koordinat) as koordinat_wkt,
            COALESCE(sl.status, 'Tidak ada laporan') as status_laporan,
            COALESCE(lm.tanggal_laporan, '') as tanggal_laporan,
            CASE 
                WHEN lm.id_masyarakat IS NOT NULL THEN 'masyarakat'
                WHEN lm.id_masyarakat IS NULL AND lm.id IS NOT NULL THEN 'admin'
                ELSE 'tidak ada'
            END as jenis_laporan,
            COALESCE(rp_latest.status_gizi, 'Belum diperiksa') as status_gizi_terakhir,
            COALESCE(rp_latest.tanggal, '') as tanggal_pemeriksaan_terakhir
        FROM balita b
        LEFT JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        LEFT JOIN laporan_masyarakat lm ON b.id = lm.id_balita AND lm.deleted_date IS NULL
        LEFT JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        LEFT JOIN (
            SELECT rp.id_balita, rp.status_gizi, rp.tanggal,
                   ROW_NUMBER() OVER (PARTITION BY rp.id_balita ORDER BY rp.tanggal DESC) as rn
            FROM riwayat_pemeriksaan rp
            WHERE rp.deleted_date IS NULL
        ) rp_latest ON b.id = rp_latest.id_balita AND rp_latest.rn = 1
        WHERE b.deleted_date IS NULL AND k.koordinat IS NOT NULL
    `

	// Add filters
	conditions := []string{}
	if statusLaporan != "" {
		if statusLaporan == "Tidak ada laporan" {
			conditions = append(conditions, "lm.id IS NULL")
		} else {
			conditions = append(conditions, "sl.status = ?")
			args = append(args, statusLaporan)
		}
	}

	if idKecamatan != "" {
		conditions = append(conditions, "kec.id = ?")
		args = append(args, idKecamatan)
	}

	if idKelurahan != "" {
		conditions = append(conditions, "kel.id = ?")
		args = append(args, idKelurahan)
	}

	if len(conditions) > 0 {
		query += " AND " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " ORDER BY kec.kecamatan ASC, kel.kelurahan ASC, b.nama ASC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return object.GeoJSONFeatureCollection{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, nama, jenisKelamin, umurBulan string
		var nomorKk, namaAyah, namaIbu, kelurahan, kecamatan string
		var koordinatWKT, statusLaporanDB, tanggalLaporan, jenisLaporan string
		var statusGiziTerakhir, tanggalPemeriksaanTerakhir string

		err := rows.Scan(
			&id, &nama, &jenisKelamin, &umurBulan,
			&nomorKk, &namaAyah, &namaIbu,
			&kelurahan, &kecamatan,
			&koordinatWKT,
			&statusLaporanDB, &tanggalLaporan, &jenisLaporan,
			&statusGiziTerakhir, &tanggalPemeriksaanTerakhir,
		)
		if err != nil {
			return object.GeoJSONFeatureCollection{}, err
		}

		// Format umur
		umurFormatted := formatUmurBalita(umurBulan)

		// Determine color based on status gizi and laporan
		color := getBalitaPointColor(statusGiziTerakhir, statusLaporanDB)

		// Create properties
		properties := map[string]any{
			"id":                           id,
			"nama":                         nama,
			"jenis_kelamin":                jenisKelamin,
			"umur":                         umurFormatted,
			"nomor_kk":                     nomorKk,
			"nama_ayah":                    namaAyah,
			"nama_ibu":                     namaIbu,
			"kelurahan":                    kelurahan,
			"kecamatan":                    kecamatan,
			"status_laporan":               statusLaporanDB,
			"tanggal_laporan":              tanggalLaporan,
			"jenis_laporan":                jenisLaporan,
			"status_gizi_terakhir":         statusGiziTerakhir,
			"tanggal_pemeriksaan_terakhir": tanggalPemeriksaanTerakhir,
			"color":                        color,
			"type":                         "balita",
		}

		// Create GeoJSON feature
		feature, err := object.CreateGeoJSONFeature(koordinatWKT, properties)
		if err != nil {
			continue // Skip if unable to parse geometry
		}

		features = append(features, feature)
	}

	if err = rows.Err(); err != nil {
		return object.GeoJSONFeatureCollection{}, err
	}

	return object.CreateGeoJSONFeatureCollection(features), nil
}

// Helper function to determine balita point color based on status
func getBalitaPointColor(statusGizi, statusLaporan string) string {
	// Priority: Status gizi first, then status laporan
	switch statusGizi {
	case "gizi buruk":
		return "#FF0000" // Red - Critical
	case "stunting":
		return "#FF6600" // Orange - Warning
	case "normal":
		return "#00AA00" // Green - Good
	default:
		// Based on laporan status if no medical examination
		switch statusLaporan {
		case "Belum diproses":
			return "#FFFF00" // Yellow - Pending
		case "Diproses dan data tidak sesuai":
			return "#808080" // Gray - Invalid
		case "Diproses dan data sesuai":
			return "#0066FF" // Blue - Verified
		case "Belum ditindaklanjuti":
			return "#FF9900" // Orange - Needs follow-up
		case "Sudah ditindaklanjuti":
			return "#00CCCC" // Cyan - In progress
		case "Sudah perbaikan gizi":
			return "#00FF00" // Bright green - Recovered
		case "Tidak ada laporan":
			return "#CCCCCC" // Light gray - No report
		default:
			return "#999999" // Default gray
		}
	}
}
