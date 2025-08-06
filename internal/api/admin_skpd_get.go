package api

import (
	"database/sql"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type skpdResponse struct {
	Id           string `json:"id"`
	Skpd         string `json:"skpd"`
	Jenis        string `json:"jenis"`
	PetugasCount int    `json:"petugas_count"` // jumlah petugas kesehatan terkait
	CreatedDate  string `json:"created_date"`
	UpdatedDate  string `json:"updated_date,omitempty"`
}

type getAllSkpdResponse struct {
	Data  []skpdResponse `json:"data"`
	Total int            `json:"total"`
}

type getSkpdByIdResponse struct {
	Data skpdResponse `json:"data"`
}

// # AdminSkpdGet handles getting SKPD data
//
// @Summary Get SKPD data
// @Description Get SKPD data based on query parameter (Admin only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all SKPD with total count
// @Description - With id parameter: Returns specific SKPD data
// @Description
// @Description SKPD data includes: skpd name, jenis (type), petugas count, creation/update dates
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "SKPD ID"
// @Success 200 {object} object.Response{data=getAllSkpdResponse} "SKPD data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "SKPD not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/skpd/get [get]
func AdminSkpdGet(w http.ResponseWriter, r *http.Request) {
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
		// Get specific SKPD by ID
		skpd, err := getSkpdById(db, idParam)
		if err != nil {
			if err == sql.ErrNoRows {
				response := object.NewResponse(http.StatusNotFound, "SKPD not found", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get SKPD", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "SKPD retrieved successfully", getSkpdByIdResponse{
			Data: skpd,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Get all SKPD
		skpdList, total, err := getAllSkpd(db)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get SKPD list", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "All SKPD retrieved successfully", getAllSkpdResponse{
			Data:  skpdList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Helper function to get SKPD by ID
func getSkpdById(db *sql.DB, id string) (skpdResponse, error) {
	var skpd skpdResponse
	var updatedDate sql.NullString

	query := `
        SELECT 
            s.id, s.skpd, s.jenis, s.created_date, s.updated_date,
            COALESCE(COUNT(pk.id), 0) as petugas_count
        FROM skpd s
        LEFT JOIN petugas_kesehatan pk ON s.id = pk.id_skpd AND pk.deleted_date IS NULL
        WHERE s.id = ? AND s.deleted_date IS NULL
        GROUP BY s.id, s.skpd, s.jenis, s.created_date, s.updated_date
    `

	err := db.QueryRow(query, id).Scan(
		&skpd.Id,
		&skpd.Skpd,
		&skpd.Jenis,
		&skpd.CreatedDate,
		&updatedDate,
		&skpd.PetugasCount,
	)

	if err != nil {
		return skpd, err
	}

	// Handle nullable updated_date
	if updatedDate.Valid {
		skpd.UpdatedDate = updatedDate.String
	}

	return skpd, nil
}

// Helper function to get all SKPD
func getAllSkpd(db *sql.DB) ([]skpdResponse, int, error) {
	var skpdList []skpdResponse

	query := `
        SELECT 
            s.id, s.skpd, s.jenis, s.created_date, s.updated_date,
            COALESCE(COUNT(pk.id), 0) as petugas_count
        FROM skpd s
        LEFT JOIN petugas_kesehatan pk ON s.id = pk.id_skpd AND pk.deleted_date IS NULL
        WHERE s.deleted_date IS NULL
        GROUP BY s.id, s.skpd, s.jenis, s.created_date, s.updated_date
        ORDER BY s.jenis ASC, s.skpd ASC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var skpd skpdResponse
		var updatedDate sql.NullString

		err := rows.Scan(
			&skpd.Id,
			&skpd.Skpd,
			&skpd.Jenis,
			&skpd.CreatedDate,
			&updatedDate,
			&skpd.PetugasCount,
		)

		if err != nil {
			return nil, 0, err
		}

		// Handle nullable updated_date
		if updatedDate.Valid {
			skpd.UpdatedDate = updatedDate.String
		}

		skpdList = append(skpdList, skpd)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM skpd WHERE deleted_date IS NULL"
	err = db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return skpdList, total, nil
}
