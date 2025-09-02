package community

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type insertKeluargaRequest struct {
	NomorKk     string     `json:"nomor_kk"`
	NamaAyah    string     `json:"nama_ayah"`
	NamaIbu     string     `json:"nama_ibu"`
	NikAyah     string     `json:"nik_ayah"`
	NikIbu      string     `json:"nik_ibu"`
	Alamat      string     `json:"alamat"`
	Rt          string     `json:"rt"`
	Rw          string     `json:"rw"`
	IdKelurahan string     `json:"id_kelurahan"`
	Koordinat   [2]float64 `json:"koordinat"` // [longitude, latitude]
}

func (r *insertKeluargaRequest) validate() error {
	// Nomor KK validation: 16 digits
	if r.NomorKk == "" {
		return fmt.Errorf("nomor KK is required")
	}
	kkRegex := regexp.MustCompile(`^\d{16}$`)
	if !kkRegex.MatchString(r.NomorKk) {
		return fmt.Errorf("nomor KK must be exactly 16 digits")
	}

	// Nama Ayah validation
	if r.NamaAyah == "" {
		return fmt.Errorf("nama ayah is required")
	}
	namaRegex := regexp.MustCompile(`^[a-zA-Z\s]{2,50}$`)
	if !namaRegex.MatchString(r.NamaAyah) {
		return fmt.Errorf("nama ayah must be 2-50 characters and contain only letters and spaces")
	}

	// Nama Ibu validation
	if r.NamaIbu == "" {
		return fmt.Errorf("nama ibu is required")
	}
	if !namaRegex.MatchString(r.NamaIbu) {
		return fmt.Errorf("nama ibu must be 2-50 characters and contain only letters and spaces")
	}

	// NIK Ayah validation: 16 digits
	if r.NikAyah == "" {
		return fmt.Errorf("NIK ayah is required")
	}
	nikRegex := regexp.MustCompile(`^\d{16}$`)
	if !nikRegex.MatchString(r.NikAyah) {
		return fmt.Errorf("NIK ayah must be exactly 16 digits")
	}

	// NIK Ibu validation: 16 digits
	if r.NikIbu == "" {
		return fmt.Errorf("NIK ibu is required")
	}
	if !nikRegex.MatchString(r.NikIbu) {
		return fmt.Errorf("NIK ibu must be exactly 16 digits")
	}

	// Alamat validation
	if r.Alamat == "" {
		return fmt.Errorf("alamat is required")
	}
	if len(r.Alamat) < 5 || len(r.Alamat) > 255 {
		return fmt.Errorf("alamat must be between 5-255 characters")
	}

	// RT validation: 1-3 digits
	if r.Rt == "" {
		return fmt.Errorf("RT is required")
	}
	rtRegex := regexp.MustCompile(`^\d{1,3}$`)
	if !rtRegex.MatchString(r.Rt) {
		return fmt.Errorf("RT must be 1-3 digits")
	}

	// RW validation: 1-3 digits
	if r.Rw == "" {
		return fmt.Errorf("RW is required")
	}
	rwRegex := regexp.MustCompile(`^\d{1,3}$`)
	if !rwRegex.MatchString(r.Rw) {
		return fmt.Errorf("RW must be 1-3 digits")
	}

	// ID Kelurahan validation
	if r.IdKelurahan == "" {
		return fmt.Errorf("id kelurahan is required")
	}

	// Koordinat validation (basic bounds check)
	// if r.Koordinat[0] < -180 || r.Koordinat[0] > 180 {
	// 	return fmt.Errorf("longitude must be between -180 and 180")
	// }
	// if r.Koordinat[1] < -90 || r.Koordinat[1] > 90 {
	// 	return fmt.Errorf("latitude must be between -90 and 90")
	// }

	return nil
}

type insertKeluargaResponse struct {
	Id string `json:"id"`
}

// # KeluargaInsert handles inserting new keluarga data for masyarakat
//
// @Summary Insert new keluarga (Community)
// @Description Insert new keluarga data for community/masyarakat users
// @Description
// @Description This endpoint allows masyarakat users to register family data
// @Description when reporting balita. The data will be linked to the reporting user.
// @Description
// @Description Validation includes:
// @Description - Nomor KK and NIK uniqueness check
// @Description - Format validation for all fields
// @Description - Kelurahan existence validation
// @Description - Coordinate bounds validation
// @Tags community
// @Accept json
// @Produce json
// @Security Bearer
// @Param keluarga body insertKeluargaRequest true "Keluarga data"
// @Success 200 {object} object.Response{data=insertKeluargaResponse} "Keluarga inserted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Masyarakat role required"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/community/keluarga/insert [post]
func KeluargaInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	// Parse request body
	var req insertKeluargaRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, "Invalid request body", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Validate request
	err = req.validate()
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
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

	// Check if Nomor KK already exists (not soft deleted)
	var exists int
	checkKKQuery := "SELECT COUNT(*) FROM keluarga WHERE nomor_kk = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkKKQuery, req.NomorKk).Scan(&exists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check nomor KK", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "Nomor KK already exists", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if NIK Ayah already exists (not soft deleted)
	checkNikAyahQuery := "SELECT COUNT(*) FROM keluarga WHERE nik_ayah = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkNikAyahQuery, req.NikAyah).Scan(&exists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check NIK ayah", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "NIK ayah already exists", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if NIK Ibu already exists (not soft deleted)
	checkNikIbuQuery := "SELECT COUNT(*) FROM keluarga WHERE nik_ibu = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkNikIbuQuery, req.NikIbu).Scan(&exists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check NIK ibu", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists > 0 {
		response := object.NewResponse(http.StatusBadRequest, "NIK ibu already exists", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if kelurahan exists
	var kelurahanExists int
	checkKelurahanQuery := "SELECT COUNT(*) FROM kelurahan WHERE id = ?"
	err = db.QueryRow(checkKelurahanQuery, req.IdKelurahan).Scan(&kelurahanExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check kelurahan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if kelurahanExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Kelurahan not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Convert coordinates to WKT format
	koordinatWKT := object.ToWKT(req.Koordinat)

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Insert keluarga with masyarakat as creator
	insertQuery := `INSERT INTO keluarga 
        (nomor_kk, nama_ayah, nama_ibu, nik_ayah, nik_ibu, alamat, rt, rw, id_kelurahan, koordinat, created_id, created_date) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ST_GeomFromText(?), ?, ?)`

	result, err := db.Exec(insertQuery,
		req.NomorKk,
		req.NamaAyah,
		req.NamaIbu,
		req.NikAyah,
		req.NikIbu,
		req.Alamat,
		req.Rt,
		req.Rw,
		req.IdKelurahan,
		koordinatWKT,
		userId, // Use user ID (from JWT) as created_id
		currentTime,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert keluarga", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get the inserted ID
	insertedId, err := result.LastInsertId()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to retrieve inserted ID", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Keluarga inserted successfully", insertKeluargaResponse{
		Id: strconv.FormatInt(insertedId, 10),
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
