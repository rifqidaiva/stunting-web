package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type assignIntervensiPetugasRequest struct {
	IdIntervensi       string `json:"id_intervensi"`
	IdPetugasKesehatan string `json:"id_petugas_kesehatan"`
}

func (r *assignIntervensiPetugasRequest) validate() error {
	if r.IdIntervensi == "" {
		return fmt.Errorf("id intervensi is required")
	}
	if r.IdPetugasKesehatan == "" {
		return fmt.Errorf("id petugas kesehatan is required")
	}
	return nil
}

type removeIntervensiPetugasRequest struct {
	Id string `json:"id"` // ID dari tabel intervensi_petugas
}

func (r *removeIntervensiPetugasRequest) validate() error {
	if r.Id == "" {
		return fmt.Errorf("assignment ID is required")
	}
	return nil
}

type assignIntervensiPetugasResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type removeIntervensiPetugasResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

type intervensiPetugasResponse struct {
	Id                  string `json:"id"`
	IdIntervensi        string `json:"id_intervensi"`
	JenisIntervensi     string `json:"jenis_intervensi"`
	TanggalIntervensi   string `json:"tanggal_intervensi"`
	DeskripsiIntervensi string `json:"deskripsi_intervensi"`
	IdPetugasKesehatan  string `json:"id_petugas_kesehatan"`
	NamaPetugas         string `json:"nama_petugas"`
	SkpdPetugas         string `json:"skpd_petugas"`
	JenisSkpd           string `json:"jenis_skpd"`
	EmailPetugas        string `json:"email_petugas"`
}

type getAllIntervensiPetugasResponse struct {
	Data  []intervensiPetugasResponse `json:"data"`
	Total int                         `json:"total"`
}

type getIntervensiPetugasByIdResponse struct {
	Data intervensiPetugasResponse `json:"data"`
}

// # IntervensiPetugasGet handles getting intervensi petugas assignments
//
// @Summary Get intervensi petugas assignments
// @Description Get intervensi petugas assignments based on query parameter (Admin only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all assignments with total count
// @Description - With id parameter: Returns specific assignment data
// @Description - With id_intervensi parameter: Returns all petugas assigned to specific intervensi
// @Description - With id_petugas_kesehatan parameter: Returns all intervensi assigned to specific petugas
// @Description
// @Description Assignment data includes: intervensi info, petugas info, SKPD info
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Assignment ID"
// @Param id_intervensi query string false "Intervensi ID"
// @Param id_petugas_kesehatan query string false "Petugas Kesehatan ID"
// @Success 200 {object} object.Response{data=getAllIntervensiPetugasResponse} "Assignment data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Assignment not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/intervensi-petugas/get [get]
func IntervensiPetugasGet(w http.ResponseWriter, r *http.Request) {
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

	// Check query parameters
	idParam := r.URL.Query().Get("id")
	idIntervensiParam := r.URL.Query().Get("id_intervensi")
	idPetugasParam := r.URL.Query().Get("id_petugas_kesehatan")

	if idParam != "" {
		// Get specific assignment by ID
		assignment, err := getIntervensiPetugasById(db, idParam)
		if err != nil {
			if err == sql.ErrNoRows {
				response := object.NewResponse(http.StatusNotFound, "Assignment not found", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get assignment", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Assignment retrieved successfully", getIntervensiPetugasByIdResponse{
			Data: assignment,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if idIntervensiParam != "" {
		// Get all petugas assigned to specific intervensi
		assignmentList, total, err := getIntervensiPetugasByIntervensi(db, idIntervensiParam)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get assignments by intervensi", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Assignments by intervensi retrieved successfully", getAllIntervensiPetugasResponse{
			Data:  assignmentList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if idPetugasParam != "" {
		// Get all intervensi assigned to specific petugas
		assignmentList, total, err := getIntervensiPetugasByPetugas(db, idPetugasParam)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get assignments by petugas", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "Assignments by petugas retrieved successfully", getAllIntervensiPetugasResponse{
			Data:  assignmentList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		// Get all assignments
		assignmentList, total, err := getAllIntervensiPetugas(db)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to get assignment list", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		response := object.NewResponse(http.StatusOK, "All assignments retrieved successfully", getAllIntervensiPetugasResponse{
			Data:  assignmentList,
			Total: total,
		})
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// # IntervensiPetugasAssign handles assigning petugas to intervensi
//
// @Summary Assign petugas to intervensi
// @Description Assign petugas kesehatan to specific intervensi (Admin only)
// @Description
// @Description Creates assignment between petugas and intervensi:
// @Description - Validates both intervensi and petugas existence
// @Description - Prevents duplicate assignments
// @Description - Ensures petugas is active and not soft deleted
// @Description - Ensures intervensi is active and not soft deleted
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param assignment body assignIntervensiPetugasRequest true "Assignment data"
// @Success 200 {object} object.Response{data=assignIntervensiPetugasResponse} "Petugas assigned successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/intervensi-petugas/assign [post]
func IntervensiPetugasAssign(w http.ResponseWriter, r *http.Request) {
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

	// Parse request body
	var req assignIntervensiPetugasRequest
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

	// Check if intervensi exists and not soft deleted
	var intervensiExists int
	var jenisIntervensi, tanggalIntervensi string
	checkIntervensiQuery := `SELECT COUNT(*), jenis, tanggal 
        FROM intervensi WHERE id = ? AND deleted_date IS NULL 
        GROUP BY jenis, tanggal`
	err = db.QueryRow(checkIntervensiQuery, req.IdIntervensi).Scan(&intervensiExists, &jenisIntervensi, &tanggalIntervensi)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusBadRequest, "Intervensi not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check intervensi existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if intervensiExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Intervensi not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if petugas kesehatan exists and not soft deleted
	var petugasExists int
	var namaPetugas, skpdPetugas string
	checkPetugasQuery := `SELECT COUNT(*), pk.nama, s.skpd 
        FROM petugas_kesehatan pk 
        LEFT JOIN skpd s ON pk.id_skpd = s.id
        WHERE pk.id = ? AND pk.deleted_date IS NULL 
        GROUP BY pk.nama, s.skpd`
	err = db.QueryRow(checkPetugasQuery, req.IdPetugasKesehatan).Scan(&petugasExists, &namaPetugas, &skpdPetugas)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusBadRequest, "Petugas kesehatan not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check petugas kesehatan existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if petugasExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Petugas kesehatan not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for duplicate assignment
	var duplicateExists int
	checkDuplicateQuery := `SELECT COUNT(*) FROM intervensi_petugas 
        WHERE id_intervensi = ? AND id_petugas_kesehatan = ?`
	err = db.QueryRow(checkDuplicateQuery, req.IdIntervensi, req.IdPetugasKesehatan).Scan(&duplicateExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check duplicate assignment", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if duplicateExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Petugas '%s' sudah di-assign ke intervensi %s pada tanggal %s",
				namaPetugas, jenisIntervensi, tanggalIntervensi), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Insert assignment
	insertQuery := `INSERT INTO intervensi_petugas (id_intervensi, id_petugas_kesehatan) VALUES (?, ?)`
	result, err := db.Exec(insertQuery, req.IdIntervensi, req.IdPetugasKesehatan)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to assign petugas to intervensi", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get the inserted ID
	insertedId, err := result.LastInsertId()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to retrieve assignment ID", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare success response with detailed information
	message := fmt.Sprintf("Petugas '%s' dari %s berhasil di-assign ke intervensi %s pada tanggal %s",
		namaPetugas, skpdPetugas, jenisIntervensi, tanggalIntervensi)

	response := object.NewResponse(http.StatusOK, "Petugas assigned to intervensi successfully", assignIntervensiPetugasResponse{
		Id:      fmt.Sprintf("%d", insertedId),
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// # IntervensiPetugasRemove handles removing petugas assignment from intervensi
//
// @Summary Remove petugas assignment from intervensi
// @Description Remove petugas kesehatan assignment from intervensi (Admin only)
// @Description
// @Description Removes assignment between petugas and intervensi:
// @Description - Validates assignment existence
// @Description - Provides detailed information about removed assignment
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param assignment body removeIntervensiPetugasRequest true "Assignment ID to remove"
// @Success 200 {object} object.Response{data=removeIntervensiPetugasResponse} "Assignment removed successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Assignment not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/intervensi-petugas/remove [delete]
func IntervensiPetugasRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	// Parse request body
	var req removeIntervensiPetugasRequest
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

	// Check if assignment exists and get details
	var exists int
	var namaPetugas, jenisIntervensi, tanggalIntervensi, skpdPetugas string
	checkQuery := `SELECT COUNT(*), pk.nama, i.jenis, i.tanggal, s.skpd
        FROM intervensi_petugas ip
        LEFT JOIN petugas_kesehatan pk ON ip.id_petugas_kesehatan = pk.id
        LEFT JOIN intervensi i ON ip.id_intervensi = i.id
        LEFT JOIN skpd s ON pk.id_skpd = s.id
        WHERE ip.id = ?
        GROUP BY pk.nama, i.jenis, i.tanggal, s.skpd`
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists, &namaPetugas, &jenisIntervensi, &tanggalIntervensi, &skpdPetugas)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Assignment not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check assignment existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Assignment not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if there are related riwayat pemeriksaan records for this assignment
	var riwayatCount int
	checkRiwayatQuery := `SELECT COUNT(*) FROM riwayat_pemeriksaan rp
        INNER JOIN intervensi_petugas ip ON rp.id_intervensi = ip.id_intervensi
        WHERE ip.id = ? AND rp.deleted_date IS NULL`
	err = db.QueryRow(checkRiwayatQuery, req.Id).Scan(&riwayatCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related medical records", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Remove assignment
	deleteQuery := `DELETE FROM intervensi_petugas WHERE id = ?`
	result, err := db.Exec(deleteQuery, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to remove assignment", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check remove result", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if rowsAffected == 0 {
		response := object.NewResponse(http.StatusNotFound, "Assignment not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare success response with detailed information
	message := fmt.Sprintf("Assignment petugas '%s' dari %s untuk intervensi %s pada tanggal %s berhasil dihapus",
		namaPetugas, skpdPetugas, jenisIntervensi, tanggalIntervensi)

	if riwayatCount > 0 {
		message += fmt.Sprintf(" (Note: Terdapat %d riwayat pemeriksaan terkait intervensi ini)", riwayatCount)
	}

	response := object.NewResponse(http.StatusOK, "Assignment removed successfully", removeIntervensiPetugasResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Helper function to get assignment by ID
func getIntervensiPetugasById(db *sql.DB, id string) (intervensiPetugasResponse, error) {
	var assignment intervensiPetugasResponse

	query := `
        SELECT 
            ip.id, ip.id_intervensi, ip.id_petugas_kesehatan,
            i.jenis, i.tanggal, i.deskripsi,
            pk.nama, s.skpd, s.jenis, p.email
        FROM intervensi_petugas ip
        LEFT JOIN intervensi i ON ip.id_intervensi = i.id
        LEFT JOIN petugas_kesehatan pk ON ip.id_petugas_kesehatan = pk.id
        LEFT JOIN skpd s ON pk.id_skpd = s.id
        LEFT JOIN pengguna p ON pk.id_pengguna = p.id
        WHERE ip.id = ?
    `

	err := db.QueryRow(query, id).Scan(
		&assignment.Id,
		&assignment.IdIntervensi,
		&assignment.IdPetugasKesehatan,
		&assignment.JenisIntervensi,
		&assignment.TanggalIntervensi,
		&assignment.DeskripsiIntervensi,
		&assignment.NamaPetugas,
		&assignment.SkpdPetugas,
		&assignment.JenisSkpd,
		&assignment.EmailPetugas,
	)

	return assignment, err
}

// Helper function to get all assignments
func getAllIntervensiPetugas(db *sql.DB) ([]intervensiPetugasResponse, int, error) {
	var assignmentList []intervensiPetugasResponse

	query := `
        SELECT 
            ip.id, ip.id_intervensi, ip.id_petugas_kesehatan,
            i.jenis, i.tanggal, i.deskripsi,
            pk.nama, s.skpd, s.jenis, p.email
        FROM intervensi_petugas ip
        LEFT JOIN intervensi i ON ip.id_intervensi = i.id AND i.deleted_date IS NULL
        LEFT JOIN petugas_kesehatan pk ON ip.id_petugas_kesehatan = pk.id AND pk.deleted_date IS NULL
        LEFT JOIN skpd s ON pk.id_skpd = s.id AND s.deleted_date IS NULL
        LEFT JOIN pengguna p ON pk.id_pengguna = p.id
        ORDER BY i.tanggal DESC, pk.nama ASC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var assignment intervensiPetugasResponse

		err := rows.Scan(
			&assignment.Id,
			&assignment.IdIntervensi,
			&assignment.IdPetugasKesehatan,
			&assignment.JenisIntervensi,
			&assignment.TanggalIntervensi,
			&assignment.DeskripsiIntervensi,
			&assignment.NamaPetugas,
			&assignment.SkpdPetugas,
			&assignment.JenisSkpd,
			&assignment.EmailPetugas,
		)

		if err != nil {
			return nil, 0, err
		}

		assignmentList = append(assignmentList, assignment)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM intervensi_petugas"
	err = db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return assignmentList, total, nil
}

// Helper function to get assignments by intervensi
func getIntervensiPetugasByIntervensi(db *sql.DB, idIntervensi string) ([]intervensiPetugasResponse, int, error) {
	var assignmentList []intervensiPetugasResponse

	query := `
        SELECT 
            ip.id, ip.id_intervensi, ip.id_petugas_kesehatan,
            i.jenis, i.tanggal, i.deskripsi,
            pk.nama, s.skpd, s.jenis, p.email
        FROM intervensi_petugas ip
        LEFT JOIN intervensi i ON ip.id_intervensi = i.id AND i.deleted_date IS NULL
        LEFT JOIN petugas_kesehatan pk ON ip.id_petugas_kesehatan = pk.id AND pk.deleted_date IS NULL
        LEFT JOIN skpd s ON pk.id_skpd = s.id AND s.deleted_date IS NULL
        LEFT JOIN pengguna p ON pk.id_pengguna = p.id
        WHERE ip.id_intervensi = ?
        ORDER BY pk.nama ASC
    `

	rows, err := db.Query(query, idIntervensi)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var assignment intervensiPetugasResponse

		err := rows.Scan(
			&assignment.Id,
			&assignment.IdIntervensi,
			&assignment.IdPetugasKesehatan,
			&assignment.JenisIntervensi,
			&assignment.TanggalIntervensi,
			&assignment.DeskripsiIntervensi,
			&assignment.NamaPetugas,
			&assignment.SkpdPetugas,
			&assignment.JenisSkpd,
			&assignment.EmailPetugas,
		)

		if err != nil {
			return nil, 0, err
		}

		assignmentList = append(assignmentList, assignment)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return assignmentList, len(assignmentList), nil
}

// Helper function to get assignments by petugas
func getIntervensiPetugasByPetugas(db *sql.DB, idPetugas string) ([]intervensiPetugasResponse, int, error) {
	var assignmentList []intervensiPetugasResponse

	query := `
        SELECT 
            ip.id, ip.id_intervensi, ip.id_petugas_kesehatan,
            i.jenis, i.tanggal, i.deskripsi,
            pk.nama, s.skpd, s.jenis, p.email
        FROM intervensi_petugas ip
        LEFT JOIN intervensi i ON ip.id_intervensi = i.id AND i.deleted_date IS NULL
        LEFT JOIN petugas_kesehatan pk ON ip.id_petugas_kesehatan = pk.id AND pk.deleted_date IS NULL
        LEFT JOIN skpd s ON pk.id_skpd = s.id AND s.deleted_date IS NULL
        LEFT JOIN pengguna p ON pk.id_pengguna = p.id
        WHERE ip.id_petugas_kesehatan = ?
        ORDER BY i.tanggal DESC
    `

	rows, err := db.Query(query, idPetugas)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var assignment intervensiPetugasResponse

		err := rows.Scan(
			&assignment.Id,
			&assignment.IdIntervensi,
			&assignment.IdPetugasKesehatan,
			&assignment.JenisIntervensi,
			&assignment.TanggalIntervensi,
			&assignment.DeskripsiIntervensi,
			&assignment.NamaPetugas,
			&assignment.SkpdPetugas,
			&assignment.JenisSkpd,
			&assignment.EmailPetugas,
		)

		if err != nil {
			return nil, 0, err
		}

		assignmentList = append(assignmentList, assignment)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return assignmentList, len(assignmentList), nil
}
