package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type insertIntervensiRequest struct {
	IdBalita  string `json:"id_balita"`   // <- Field baru
	Jenis     string `json:"jenis"`       // "gizi", "kesehatan", "sosial"
	Tanggal   string `json:"tanggal"`     // Format: YYYY-MM-DD
	Deskripsi string `json:"deskripsi"`
	Hasil     string `json:"hasil"`
}

func (r *insertIntervensiRequest) validate() error {
	// ID Balita validation (wajib)
	if r.IdBalita == "" {
		return fmt.Errorf("id balita is required")
	}

	// Jenis validation: must be one of the allowed types
	if r.Jenis == "" {
		return fmt.Errorf("jenis intervensi is required")
	}
	allowedJenis := []string{"gizi", "kesehatan", "sosial"}
	jenisValid := slices.Contains(allowedJenis, r.Jenis)
	if !jenisValid {
		return fmt.Errorf("jenis must be one of: gizi, kesehatan, sosial")
	}

	// Tanggal validation: YYYY-MM-DD format
	if r.Tanggal == "" {
		return fmt.Errorf("tanggal intervensi is required")
	}
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(r.Tanggal) {
		return fmt.Errorf("tanggal must be in YYYY-MM-DD format")
	}

	// Parse and validate the date
	_, err := time.Parse("2006-01-02", r.Tanggal)
	if err != nil {
		return fmt.Errorf("invalid tanggal format")
	}

	// Check if date is not in the future
	intervensiDate, _ := time.Parse("2006-01-02", r.Tanggal)
	if intervensiDate.After(time.Now()) {
		return fmt.Errorf("tanggal intervensi cannot be in the future")
	}

	// Deskripsi validation
	if r.Deskripsi == "" {
		return fmt.Errorf("deskripsi is required")
	}
	if len(r.Deskripsi) < 10 || len(r.Deskripsi) > 1000 {
		return fmt.Errorf("deskripsi must be between 10-1000 characters")
	}

	// Hasil validation
	if r.Hasil == "" {
		return fmt.Errorf("hasil is required")
	}
	if len(r.Hasil) < 5 || len(r.Hasil) > 500 {
		return fmt.Errorf("hasil must be between 5-500 characters")
	}

	return nil
}

type insertIntervensiResponse struct {
	Id string `json:"id"`
}

// # IntervensiInsert handles inserting new intervensi data
//
// @Summary Insert new intervensi
// @Description Insert new intervensi data (Admin only)
// @Description
// @Description Creates a new intervensi record with:
// @Description - id_balita: ID of the balita being intervened
// @Description - jenis: type of intervention (gizi, kesehatan, sosial)
// @Description - tanggal: intervention date (YYYY-MM-DD format)
// @Description - deskripsi: detailed description of the intervention
// @Description - hasil: results or outcomes of the intervention
// @Description - Validates intervention type and date constraints
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param intervensi body insertIntervensiRequest true "Intervensi data"
// @Success 200 {object} object.Response{data=insertIntervensiResponse} "Intervensi inserted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/intervensi/insert [post]
func IntervensiInsert(w http.ResponseWriter, r *http.Request) {
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

	// Check if user is admin
	if role != "admin" {
		response := object.NewResponse(http.StatusForbidden, "Access denied. Admin role required", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Parse request body
	var req insertIntervensiRequest
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

	// Check if balita exists and not soft deleted
	var balitaExists int
	var namaBalita string
	checkBalitaQuery := `SELECT COUNT(*), b.nama 
        FROM balita b WHERE b.id = ? AND b.deleted_date IS NULL 
        GROUP BY b.nama`
	err = db.QueryRow(checkBalitaQuery, req.IdBalita).Scan(&balitaExists, &namaBalita)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusBadRequest, "Balita not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check balita existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if balitaExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Balita not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for duplicate intervensi (same balita, jenis, tanggal, and similar deskripsi)
	var duplicateExists int
	checkDuplicateQuery := `SELECT COUNT(*) FROM intervensi 
        WHERE id_balita = ? AND jenis = ? AND tanggal = ? AND deskripsi = ? AND deleted_date IS NULL`
	err = db.QueryRow(checkDuplicateQuery, req.IdBalita, req.Jenis, req.Tanggal, req.Deskripsi).Scan(&duplicateExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check duplicate intervensi", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if duplicateExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Intervensi with same balita '%s', type, date, and description already exists", namaBalita), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Insert intervensi
	insertQuery := `INSERT INTO intervensi 
        (id_balita, jenis, tanggal, deskripsi, hasil, created_id, created_date) 
        VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(insertQuery,
		req.IdBalita,    // <- Parameter baru
		req.Jenis,
		req.Tanggal,
		req.Deskripsi,
		req.Hasil,
		userId,
		currentTime,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert intervensi", nil)
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

	// Prepare success response with additional info
	message := fmt.Sprintf("Intervensi %s untuk balita '%s' berhasil ditambahkan untuk tanggal %s",
		req.Jenis, namaBalita, req.Tanggal)

	response := object.NewResponse(http.StatusOK, message, insertIntervensiResponse{
		Id: strconv.FormatInt(insertedId, 10),
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}