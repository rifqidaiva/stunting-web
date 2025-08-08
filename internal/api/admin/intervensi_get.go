package admin

import (
	"database/sql"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type intervensiResponse struct {
	Id           string `json:"id"`
	Jenis        string `json:"jenis"`
	Tanggal      string `json:"tanggal"`
	Deskripsi    string `json:"deskripsi"`
	Hasil        string `json:"hasil"`
	PetugasCount int    `json:"petugas_count"` // jumlah petugas yang di-assign
	RiwayatCount int    `json:"riwayat_count"` // jumlah riwayat pemeriksaan terkait
	CreatedDate  string `json:"created_date"`
	UpdatedDate  string `json:"updated_date,omitempty"`
	CreatedBy    string `json:"created_by,omitempty"` // nama admin yang membuat
	UpdatedBy    string `json:"updated_by,omitempty"` // nama admin yang mengupdate
}

type getAllIntervensiResponse struct {
	Data  []intervensiResponse `json:"data"`
	Total int                  `json:"total"`
}

type getIntervensiByIdResponse struct {
	Data intervensiResponse `json:"data"`
}

// # IntervensiGet handles getting intervensi data
//
// @Summary Get intervensi data
// @Description Get intervensi data based on query parameter (Admin only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all intervensi with total count
// @Description - With id parameter: Returns specific intervensi data
// @Description
// @Description Intervensi data includes: jenis, tanggal, deskripsi, hasil, petugas count, riwayat count, creation/update info
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Intervensi ID"
// @Success 200 {object} object.Response{data=getAllIntervensiResponse} "Intervensi data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Intervensi not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/intervensi/get [get]
func IntervensiGet(w http.ResponseWriter, r *http.Request) {
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
		// Get specific intervensi by ID
		intervensi, err := getIntervensiById(db, idParam)
		if err != nil {
			if err == sql.ErrNoRows {
				response := object.NewResponse(http.StatusNotFound, "Intervensi not found", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get intervensi", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Intervensi retrieved successfully", getIntervensiByIdResponse{
			Data: intervensi,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Get all intervensi
		intervensiList, total, err := getAllIntervensi(db)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get intervensi list", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "All intervensi retrieved successfully", getAllIntervensiResponse{
			Data:  intervensiList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Helper function to get intervensi by ID
func getIntervensiById(db *sql.DB, id string) (intervensiResponse, error) {
	var intervensi intervensiResponse
	var updatedDate sql.NullString
	var createdBy, updatedBy sql.NullString

	query := `
        SELECT 
            i.id, i.jenis, i.tanggal, i.deskripsi, i.hasil,
            i.created_date, i.updated_date,
            COALESCE(COUNT(DISTINCT ip.id), 0) as petugas_count,
            COALESCE(COUNT(DISTINCT rp.id), 0) as riwayat_count,
            pc.email as created_by,
            pu.email as updated_by
        FROM intervensi i
        LEFT JOIN intervensi_petugas ip ON i.id = ip.id_intervensi
        LEFT JOIN riwayat_pemeriksaan rp ON i.id = rp.id_intervensi AND rp.deleted_date IS NULL
        LEFT JOIN pengguna pc ON i.created_id = pc.id
        LEFT JOIN pengguna pu ON i.updated_id = pu.id
        WHERE i.id = ? AND i.deleted_date IS NULL
        GROUP BY i.id, i.jenis, i.tanggal, i.deskripsi, i.hasil,
                 i.created_date, i.updated_date, pc.email, pu.email
    `

	err := db.QueryRow(query, id).Scan(
		&intervensi.Id,
		&intervensi.Jenis,
		&intervensi.Tanggal,
		&intervensi.Deskripsi,
		&intervensi.Hasil,
		&intervensi.CreatedDate,
		&updatedDate,
		&intervensi.PetugasCount,
		&intervensi.RiwayatCount,
		&createdBy,
		&updatedBy,
	)

	if err != nil {
		return intervensi, err
	}

	// Handle nullable fields
	if updatedDate.Valid {
		intervensi.UpdatedDate = updatedDate.String
	}
	if createdBy.Valid {
		intervensi.CreatedBy = createdBy.String
	}
	if updatedBy.Valid {
		intervensi.UpdatedBy = updatedBy.String
	}

	return intervensi, nil
}

// Helper function to get all intervensi
func getAllIntervensi(db *sql.DB) ([]intervensiResponse, int, error) {
	var intervensiList []intervensiResponse

	query := `
        SELECT 
            i.id, i.jenis, i.tanggal, i.deskripsi, i.hasil,
            i.created_date, i.updated_date,
            COALESCE(COUNT(DISTINCT ip.id), 0) as petugas_count,
            COALESCE(COUNT(DISTINCT rp.id), 0) as riwayat_count,
            pc.email as created_by,
            pu.email as updated_by
        FROM intervensi i
        LEFT JOIN intervensi_petugas ip ON i.id = ip.id_intervensi
        LEFT JOIN riwayat_pemeriksaan rp ON i.id = rp.id_intervensi AND rp.deleted_date IS NULL
        LEFT JOIN pengguna pc ON i.created_id = pc.id
        LEFT JOIN pengguna pu ON i.updated_id = pu.id
        WHERE i.deleted_date IS NULL
        GROUP BY i.id, i.jenis, i.tanggal, i.deskripsi, i.hasil,
                 i.created_date, i.updated_date, pc.email, pu.email
        ORDER BY i.tanggal DESC, i.created_date DESC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var intervensi intervensiResponse
		var updatedDate sql.NullString
		var createdBy, updatedBy sql.NullString

		err := rows.Scan(
			&intervensi.Id,
			&intervensi.Jenis,
			&intervensi.Tanggal,
			&intervensi.Deskripsi,
			&intervensi.Hasil,
			&intervensi.CreatedDate,
			&updatedDate,
			&intervensi.PetugasCount,
			&intervensi.RiwayatCount,
			&createdBy,
			&updatedBy,
		)

		if err != nil {
			return nil, 0, err
		}

		// Handle nullable fields
		if updatedDate.Valid {
			intervensi.UpdatedDate = updatedDate.String
		}
		if createdBy.Valid {
			intervensi.CreatedBy = createdBy.String
		}
		if updatedBy.Valid {
			intervensi.UpdatedBy = updatedBy.String
		}

		intervensiList = append(intervensiList, intervensi)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM intervensi WHERE deleted_date IS NULL"
	err = db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return intervensiList, total, nil
}
