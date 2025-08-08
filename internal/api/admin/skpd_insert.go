package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type insertSkpdRequest struct {
	Skpd  string `json:"skpd"`
	Jenis string `json:"jenis"` // "puskesmas", "kelurahan", "skpd"
}

func (r *insertSkpdRequest) validate() error {
	// SKPD name validation: 2-100 characters
	if r.Skpd == "" {
		return fmt.Errorf("nama SKPD is required")
	}
	if len(r.Skpd) < 2 || len(r.Skpd) > 100 {
		return fmt.Errorf("nama SKPD must be between 2-100 characters")
	}
	// Validasi format nama SKPD (huruf, angka, spasi, dan tanda baca umum)
	skpdRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\.\-\,\(\)\/]{2,100}$`)
	if !skpdRegex.MatchString(r.Skpd) {
		return fmt.Errorf("nama SKPD contains invalid characters")
	}

	// Jenis validation: must be one of the allowed values
	if r.Jenis == "" {
		return fmt.Errorf("jenis SKPD is required")
	}
	allowedJenis := []string{"puskesmas", "kelurahan", "skpd"}
	isValidJenis := slices.Contains(allowedJenis, r.Jenis)
	if !isValidJenis {
		return fmt.Errorf("jenis SKPD must be one of: %v", allowedJenis)
	}

	return nil
}

type insertSkpdResponse struct {
	Id string `json:"id"`
}

// # InsertSkpd handles inserting new SKPD data
//
// @Summary Insert new SKPD
// @Description Insert new SKPD data (Admin only)
// @Description
// @Description Inserts SKPD record with data including:
// @Description - skpd (nama SKPD), jenis (puskesmas/kelurahan/skpd)
// @Description - Validates uniqueness of SKPD name within the same jenis
// @Description - Supports different types of SKPD organizations
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param skpd body insertSkpdRequest true "SKPD data"
// @Success 200 {object} object.Response{data=insertSkpdResponse} "SKPD inserted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/skpd/insert [post]
func AdminSkpdInsert(w http.ResponseWriter, r *http.Request) {
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
	var req insertSkpdRequest
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

	// Check if SKPD name already exists within the same jenis (not soft deleted)
	var exists int
	checkQuery := "SELECT COUNT(*) FROM skpd WHERE skpd = ? AND jenis = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkQuery, req.Skpd, req.Jenis).Scan(&exists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check SKPD existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("SKPD '%s' with jenis '%s' already exists", req.Skpd, req.Jenis), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for similar SKPD names to warn about potential duplicates
	var similarCount int
	similarQuery := "SELECT COUNT(*) FROM skpd WHERE LOWER(skpd) LIKE LOWER(?) AND deleted_date IS NULL"
	similarPattern := "%" + req.Skpd + "%"
	err = db.QueryRow(similarQuery, similarPattern).Scan(&similarCount)
	if err != nil {
		// Not critical, continue with insertion
		similarCount = 0
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Insert SKPD
	insertQuery := `INSERT INTO skpd 
        (skpd, jenis, created_id, created_date) 
        VALUES (?, ?, ?, ?)`

	result, err := db.Exec(insertQuery,
		req.Skpd,
		req.Jenis,
		userId,
		currentTime,
	)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to insert SKPD", nil)
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

	// Prepare response message with additional info if needed
	message := "SKPD inserted successfully"
	if similarCount > 0 {
		message += fmt.Sprintf(" (Note: Found %d similar SKPD names)", similarCount)
	}

	response := object.NewResponse(http.StatusOK, message, insertSkpdResponse{
		Id: strconv.FormatInt(insertedId, 10),
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
