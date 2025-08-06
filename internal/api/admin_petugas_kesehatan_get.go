package api

import (
	"database/sql"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type petugasKesehatanResponse struct {
	Id              string `json:"id"`
	IdPengguna      string `json:"id_pengguna"`
	IdSkpd          string `json:"id_skpd"`
	Email           string `json:"email"`
	Nama            string `json:"nama"`
	Skpd            string `json:"skpd"`
	JenisSkpd       string `json:"jenis_skpd"`
	IntervensiCount int    `json:"intervensi_count"` // jumlah intervensi terkait
	CreatedDate     string `json:"created_date"`
	UpdatedDate     string `json:"updated_date,omitempty"`
}

type getAllPetugasKesehatanResponse struct {
	Data  []petugasKesehatanResponse `json:"data"`
	Total int                        `json:"total"`
}

type getPetugasKesehatanByIdResponse struct {
	Data petugasKesehatanResponse `json:"data"`
}

// # AdminPetugasKesehatanGet handles getting petugas kesehatan data
//
// @Summary Get petugas kesehatan data
// @Description Get petugas kesehatan data based on query parameter (Admin only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all petugas kesehatan with total count
// @Description - With id parameter: Returns specific petugas kesehatan data
// @Description
// @Description Petugas kesehatan data includes: nama, email, SKPD info, intervensi count, creation/update dates
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Petugas Kesehatan ID"
// @Success 200 {object} object.Response{data=getAllPetugasKesehatanResponse} "Petugas kesehatan data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Petugas kesehatan not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/petugas-kesehatan [get]
func AdminPetugasKesehatanGet(w http.ResponseWriter, r *http.Request) {
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
		// Get specific petugas kesehatan by ID
		petugas, err := getPetugasKesehatanById(db, idParam)
		if err != nil {
			if err == sql.ErrNoRows {
				response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get petugas kesehatan", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Petugas kesehatan retrieved successfully", getPetugasKesehatanByIdResponse{
			Data: petugas,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Get all petugas kesehatan
		petugasList, total, err := getAllPetugasKesehatan(db)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get petugas kesehatan list", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "All petugas kesehatan retrieved successfully", getAllPetugasKesehatanResponse{
			Data:  petugasList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Helper function to get petugas kesehatan by ID
func getPetugasKesehatanById(db *sql.DB, id string) (petugasKesehatanResponse, error) {
	var petugas petugasKesehatanResponse
	var updatedDate sql.NullString

	query := `
        SELECT 
            pk.id, pk.id_pengguna, pk.id_skpd, pk.nama, pk.created_date, pk.updated_date,
            p.email,
            s.skpd, s.jenis as jenis_skpd,
            COALESCE(COUNT(ip.id), 0) as intervensi_count
        FROM petugas_kesehatan pk
        LEFT JOIN pengguna p ON pk.id_pengguna = p.id
        LEFT JOIN skpd s ON pk.id_skpd = s.id AND s.deleted_date IS NULL
        LEFT JOIN intervensi_petugas ip ON pk.id = ip.id_petugas_kesehatan
        WHERE pk.id = ? AND pk.deleted_date IS NULL
        GROUP BY pk.id, pk.id_pengguna, pk.id_skpd, pk.nama, pk.created_date, pk.updated_date,
                 p.email, s.skpd, s.jenis
    `

	err := db.QueryRow(query, id).Scan(
		&petugas.Id,
		&petugas.IdPengguna,
		&petugas.IdSkpd,
		&petugas.Nama,
		&petugas.CreatedDate,
		&updatedDate,
		&petugas.Email,
		&petugas.Skpd,
		&petugas.JenisSkpd,
		&petugas.IntervensiCount,
	)

	if err != nil {
		return petugas, err
	}

	// Handle nullable updated_date
	if updatedDate.Valid {
		petugas.UpdatedDate = updatedDate.String
	}

	return petugas, nil
}

// Helper function to get all petugas kesehatan
func getAllPetugasKesehatan(db *sql.DB) ([]petugasKesehatanResponse, int, error) {
	var petugasList []petugasKesehatanResponse

	query := `
        SELECT 
            pk.id, pk.id_pengguna, pk.id_skpd, pk.nama, pk.created_date, pk.updated_date,
            p.email,
            s.skpd, s.jenis as jenis_skpd,
            COALESCE(COUNT(ip.id), 0) as intervensi_count
        FROM petugas_kesehatan pk
        LEFT JOIN pengguna p ON pk.id_pengguna = p.id
        LEFT JOIN skpd s ON pk.id_skpd = s.id AND s.deleted_date IS NULL
        LEFT JOIN intervensi_petugas ip ON pk.id = ip.id_petugas_kesehatan
        WHERE pk.deleted_date IS NULL
        GROUP BY pk.id, pk.id_pengguna, pk.id_skpd, pk.nama, pk.created_date, pk.updated_date,
                 p.email, s.skpd, s.jenis
        ORDER BY s.jenis ASC, s.skpd ASC, pk.nama ASC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var petugas petugasKesehatanResponse
		var updatedDate sql.NullString

		err := rows.Scan(
			&petugas.Id,
			&petugas.IdPengguna,
			&petugas.IdSkpd,
			&petugas.Nama,
			&petugas.CreatedDate,
			&updatedDate,
			&petugas.Email,
			&petugas.Skpd,
			&petugas.JenisSkpd,
			&petugas.IntervensiCount,
		)

		if err != nil {
			return nil, 0, err
		}

		// Handle nullable updated_date
		if updatedDate.Valid {
			petugas.UpdatedDate = updatedDate.String
		}

		petugasList = append(petugasList, petugas)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM petugas_kesehatan WHERE deleted_date IS NULL"
	err = db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return petugasList, total, nil
}
