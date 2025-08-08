package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type insertBalitaRequest struct {
    IdKeluarga   string `json:"id_keluarga"`
    Nama         string `json:"nama"`
    TanggalLahir string `json:"tanggal_lahir"` // Format: YYYY-MM-DD
    JenisKelamin string `json:"jenis_kelamin"` // "L" or "P"
    BeratLahir   string `json:"berat_lahir"`   // in grams
    TinggiLahir  string `json:"tinggi_lahir"`  // in cm
}

func (r *insertBalitaRequest) validate() error {
    // ID Keluarga validation
    if r.IdKeluarga == "" {
        return fmt.Errorf("id keluarga is required")
    }

    // Nama validation: 2-50 characters, only letters and spaces
    if r.Nama == "" {
        return fmt.Errorf("nama is required")
    }
    namaRegex := regexp.MustCompile(`^[a-zA-Z\s]{2,50}$`)
    if !namaRegex.MatchString(r.Nama) {
        return fmt.Errorf("nama must be 2-50 characters and contain only letters and spaces")
    }

    // Tanggal Lahir validation: YYYY-MM-DD format
    if r.TanggalLahir == "" {
        return fmt.Errorf("tanggal lahir is required")
    }
    dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
    if !dateRegex.MatchString(r.TanggalLahir) {
        return fmt.Errorf("tanggal lahir must be in YYYY-MM-DD format")
    }

    // Parse and validate the date
    _, err := time.Parse("2006-01-02", r.TanggalLahir)
    if err != nil {
        return fmt.Errorf("invalid tanggal lahir format")
    }

    // Check if date is not in the future
    birthDate, _ := time.Parse("2006-01-02", r.TanggalLahir)
    if birthDate.After(time.Now()) {
        return fmt.Errorf("tanggal lahir cannot be in the future")
    }

    // Check if child is not older than 5 years (balita criteria)
    fiveYearsAgo := time.Now().AddDate(-5, 0, 0)
    if birthDate.Before(fiveYearsAgo) {
        return fmt.Errorf("child must be under 5 years old (balita criteria)")
    }

    // Jenis Kelamin validation: L or P
    if r.JenisKelamin == "" {
        return fmt.Errorf("jenis kelamin is required")
    }
    if r.JenisKelamin != "L" && r.JenisKelamin != "P" {
        return fmt.Errorf("jenis kelamin must be 'L' (Laki-laki) or 'P' (Perempuan)")
    }

    // Berat Lahir validation: numeric, reasonable range (500-6000 grams)
    if r.BeratLahir == "" {
        return fmt.Errorf("berat lahir is required")
    }
    beratLahir, err := strconv.Atoi(r.BeratLahir)
    if err != nil {
        return fmt.Errorf("berat lahir must be a valid number (in grams)")
    }
    if beratLahir < 500 || beratLahir > 6000 {
        return fmt.Errorf("berat lahir must be between 500-6000 grams")
    }

    // Tinggi Lahir validation: numeric, reasonable range (25-65 cm)
    if r.TinggiLahir == "" {
        return fmt.Errorf("tinggi lahir is required")
    }
    tinggiLahir, err := strconv.Atoi(r.TinggiLahir)
    if err != nil {
        return fmt.Errorf("tinggi lahir must be a valid number (in cm)")
    }
    if tinggiLahir < 25 || tinggiLahir > 65 {
        return fmt.Errorf("tinggi lahir must be between 25-65 cm")
    }

    return nil
}

type insertBalitaResponse struct {
    Id string `json:"id"`
}

// # BalitaInsert handles inserting new balita data
//
// @Summary Insert new balita
// @Description Insert new balita data (Admin only)
// @Description
// @Description Inserts balita record with data including:
// @Description - id_keluarga, nama, tanggal_lahir, jenis_kelamin
// @Description - berat_lahir (in grams), tinggi_lahir (in cm)
// @Description - Validates keluarga existence and balita age criteria (under 5 years)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param balita body insertBalitaRequest true "Balita data"
// @Success 200 {object} object.Response{data=insertBalitaResponse} "Balita inserted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/balita/insert [post]
func BalitaInsert(w http.ResponseWriter, r *http.Request) {
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
    var req insertBalitaRequest
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

    // Check if keluarga exists and not soft deleted
    var keluargaExists int
    checkKeluargaQuery := "SELECT COUNT(*) FROM keluarga WHERE id = ? AND deleted_date IS NULL"
    err = db.QueryRow(checkKeluargaQuery, req.IdKeluarga).Scan(&keluargaExists)
    if err != nil {
        response := object.NewResponse(http.StatusInternalServerError, "Failed to check keluarga existence", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    if keluargaExists == 0 {
        response := object.NewResponse(http.StatusBadRequest, "Keluarga not found", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Check for duplicate balita (same name, birth date, and keluarga)
    var duplicateExists int
    checkDuplicateQuery := `SELECT COUNT(*) FROM balita 
        WHERE id_keluarga = ? AND nama = ? AND tanggal_lahir = ? AND deleted_date IS NULL`
    err = db.QueryRow(checkDuplicateQuery, req.IdKeluarga, req.Nama, req.TanggalLahir).Scan(&duplicateExists)
    if err != nil {
        response := object.NewResponse(http.StatusInternalServerError, "Failed to check duplicate balita", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }
    if duplicateExists > 0 {
        response := object.NewResponse(http.StatusBadRequest, "Balita with same name and birth date already exists in this keluarga", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Current timestamp
    currentTime := time.Now().Format("2006-01-02 15:04:05")

    // Insert balita
    insertQuery := `INSERT INTO balita 
        (id_keluarga, nama, tanggal_lahir, jenis_kelamin, berat_lahir, tinggi_lahir, created_id, created_date) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

    result, err := db.Exec(insertQuery,
        req.IdKeluarga,
        req.Nama,
        req.TanggalLahir,
        req.JenisKelamin,
        req.BeratLahir,
        req.TinggiLahir,
        userId,
        currentTime,
    )
    if err != nil {
        response := object.NewResponse(http.StatusInternalServerError, "Failed to insert balita", nil)
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

    response := object.NewResponse(http.StatusOK, "Balita inserted successfully", insertBalitaResponse{
        Id: strconv.FormatInt(insertedId, 10),
    })
    if err := response.WriteJson(w); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}