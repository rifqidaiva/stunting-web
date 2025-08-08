package admin

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type balitaResponse struct {
	Id           string `json:"id"`
	IdKeluarga   string `json:"id_keluarga"`
	NomorKk      string `json:"nomor_kk"`
	NamaAyah     string `json:"nama_ayah"`
	NamaIbu      string `json:"nama_ibu"`
	Nama         string `json:"nama"`
	TanggalLahir string `json:"tanggal_lahir"`
	JenisKelamin string `json:"jenis_kelamin"`
	BeratLahir   string `json:"berat_lahir"`
	TinggiLahir  string `json:"tinggi_lahir"`
	Umur         string `json:"umur"` // calculated field in months
	Kelurahan    string `json:"kelurahan"`
	Kecamatan    string `json:"kecamatan"`
	CreatedDate  string `json:"created_date"`
	UpdatedDate  string `json:"updated_date,omitempty"`
}

type getAllBalitaResponse struct {
	Data  []balitaResponse `json:"data"`
	Total int              `json:"total"`
}

type getBalitaByIdResponse struct {
	Data balitaResponse `json:"data"`
}

// # AdminBalitaGet handles getting balita data
//
// @Summary Get balita data
// @Description Get balita data based on query parameter (Admin only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all balita with total count
// @Description - With id parameter: Returns specific balita data
// @Description
// @Description Balita data includes: nama, tanggal_lahir, jenis_kelamin, berat_lahir, tinggi_lahir, umur, keluarga info, location info
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Balita ID"
// @Success 200 {object} object.Response{data=getAllBalitaResponse} "Balita data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Balita not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/balita/get [get]
func AdminBalitaGet(w http.ResponseWriter, r *http.Request) {
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

	// Check if ID parameter is provided
	idParam := r.URL.Query().Get("id")
	if idParam != "" {
		// Get specific balita by ID
		balita, err := getBalitaById(db, idParam)
		if err != nil {
			if err == sql.ErrNoRows {
				response := object.NewResponse(http.StatusNotFound, "Balita not found", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get balita", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Balita retrieved successfully", getBalitaByIdResponse{
			Data: balita,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Get all balita
		balitaList, total, err := getAllBalita(db)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get balita list", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "All balita retrieved successfully", getAllBalitaResponse{
			Data:  balitaList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Helper function to get balita by ID
func getBalitaById(db *sql.DB, id string) (balitaResponse, error) {
	var balita balitaResponse
	var updatedDate sql.NullString

	query := `
        SELECT 
            b.id, b.id_keluarga, b.nama, b.tanggal_lahir, b.jenis_kelamin,
            b.berat_lahir, b.tinggi_lahir, b.created_date, b.updated_date,
            k.nomor_kk, k.nama_ayah, k.nama_ibu,
            kel.kelurahan, kec.kecamatan
        FROM balita b
        LEFT JOIN keluarga k ON b.id_keluarga = k.id
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        WHERE b.id = ? AND b.deleted_date IS NULL
    `

	err := db.QueryRow(query, id).Scan(
		&balita.Id,
		&balita.IdKeluarga,
		&balita.Nama,
		&balita.TanggalLahir,
		&balita.JenisKelamin,
		&balita.BeratLahir,
		&balita.TinggiLahir,
		&balita.CreatedDate,
		&updatedDate,
		&balita.NomorKk,
		&balita.NamaAyah,
		&balita.NamaIbu,
		&balita.Kelurahan,
		&balita.Kecamatan,
	)

	if err != nil {
		return balita, err
	}

	// Calculate age in months
	balita.Umur = calculateAgeInMonths(balita.TanggalLahir)

	// Handle nullable updated_date
	if updatedDate.Valid {
		balita.UpdatedDate = updatedDate.String
	}

	return balita, nil
}

// Helper function to get all balita
func getAllBalita(db *sql.DB) ([]balitaResponse, int, error) {
	var balitaList []balitaResponse

	query := `
        SELECT 
            b.id, b.id_keluarga, b.nama, b.tanggal_lahir, b.jenis_kelamin,
            b.berat_lahir, b.tinggi_lahir, b.created_date, b.updated_date,
            k.nomor_kk, k.nama_ayah, k.nama_ibu,
            kel.kelurahan, kec.kecamatan
        FROM balita b
        LEFT JOIN keluarga k ON b.id_keluarga = k.id AND k.deleted_date IS NULL
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        WHERE b.deleted_date IS NULL
        ORDER BY b.created_date DESC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var balita balitaResponse
		var updatedDate sql.NullString

		err := rows.Scan(
			&balita.Id,
			&balita.IdKeluarga,
			&balita.Nama,
			&balita.TanggalLahir,
			&balita.JenisKelamin,
			&balita.BeratLahir,
			&balita.TinggiLahir,
			&balita.CreatedDate,
			&updatedDate,
			&balita.NomorKk,
			&balita.NamaAyah,
			&balita.NamaIbu,
			&balita.Kelurahan,
			&balita.Kecamatan,
		)

		if err != nil {
			return nil, 0, err
		}

		// Calculate age in months
		balita.Umur = calculateAgeInMonths(balita.TanggalLahir)

		// Handle nullable updated_date
		if updatedDate.Valid {
			balita.UpdatedDate = updatedDate.String
		}

		balitaList = append(balitaList, balita)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM balita WHERE deleted_date IS NULL"
	err = db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return balitaList, total, nil
}

// Helper function to calculate age in months
func calculateAgeInMonths(birthDateStr string) string {
	birthDate, err := time.Parse("2006-01-02", birthDateStr)
	if err != nil {
		return "0 bulan"
	}

	now := time.Now()
	years := now.Year() - birthDate.Year()
	months := int(now.Month()) - int(birthDate.Month())

	// Adjust if the current day is before the birth day
	if now.Day() < birthDate.Day() {
		months--
	}

	// Convert to total months
	totalMonths := years*12 + months

	if totalMonths < 0 {
		return "0 bulan"
	}

	if totalMonths < 12 {
		return strconv.Itoa(totalMonths) + " bulan"
	} else {
		years := totalMonths / 12
		remainingMonths := totalMonths % 12
		if remainingMonths == 0 {
			return strconv.Itoa(years) + " tahun"
		}
		return strconv.Itoa(years) + " tahun " + strconv.Itoa(remainingMonths) + " bulan"
	}
}
