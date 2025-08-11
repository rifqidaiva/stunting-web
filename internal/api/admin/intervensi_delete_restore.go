package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type deleteIntervensiRequest struct {
	Id string `json:"id"`
}

func (r *deleteIntervensiRequest) validate() error {
	if r.Id == "" {
		return fmt.Errorf("intervensi ID is required")
	}
	return nil
}

type deleteIntervensiResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # IntervensiDelete handles soft deleting intervensi data
//
// @Summary Delete intervensi data (soft delete)
// @Description Soft delete intervensi data by setting deleted_date and deleted_id (Admin only)
// @Description
// @Description Performs soft delete operation:
// @Description - Sets deleted_date to current timestamp
// @Description - Sets deleted_id to current user ID
// @Description - Data remains in database but is excluded from queries
// @Description - Can be restored if needed in the future
// @Description - Checks for related records before deletion
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param intervensi body deleteIntervensiRequest true "Intervensi ID to delete"
// @Success 200 {object} object.Response{data=deleteIntervensiResponse} "Intervensi deleted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Intervensi not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/intervensi/delete [delete]
func IntervensiDelete(w http.ResponseWriter, r *http.Request) {
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
	var req deleteIntervensiRequest
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

	// Check if intervensi exists and not already soft deleted
	var exists int
	var deletedDate sql.NullString
	var jenis, tanggal, deskripsi, namaBalita string
	checkQuery := `SELECT COUNT(*), i.deleted_date, i.jenis, i.tanggal, i.deskripsi, b.nama as nama_balita
        FROM intervensi i
        LEFT JOIN balita b ON i.id_balita = b.id
        WHERE i.id = ? 
        GROUP BY i.deleted_date, i.jenis, i.tanggal, i.deskripsi, b.nama`
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists, &deletedDate, &jenis, &tanggal, &deskripsi, &namaBalita)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Intervensi not found", nil)
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

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Intervensi not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if already soft deleted
	if deletedDate.Valid {
		response := object.NewResponse(http.StatusBadRequest, "Intervensi already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if intervensi has related petugas records (prevent deletion if has assignments)
	var petugasCount int
	checkPetugasQuery := "SELECT COUNT(*) FROM intervensi_petugas WHERE id_intervensi = ?"
	err = db.QueryRow(checkPetugasQuery, req.Id).Scan(&petugasCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related petugas", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if petugasCount > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot delete intervensi '%s'. There are %d petugas assigned to this intervensi", jenis, petugasCount), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if intervensi has related riwayat pemeriksaan records (prevent deletion if has medical records)
	var riwayatCount int
	checkRiwayatQuery := "SELECT COUNT(*) FROM riwayat_pemeriksaan WHERE id_intervensi = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkRiwayatQuery, req.Id).Scan(&riwayatCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related riwayat pemeriksaan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if riwayatCount > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot delete intervensi '%s'. There are %d active riwayat pemeriksaan records related to this intervensi", jenis, riwayatCount), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Begin transaction (to handle both intervensi and related junction tables if needed)
	tx, err := db.Begin()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to begin transaction", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer tx.Rollback()

	// Perform soft delete on intervensi table
	deleteIntervensiQuery := `UPDATE intervensi SET 
        deleted_id = ?, 
        deleted_date = ? 
        WHERE id = ? AND deleted_date IS NULL`

	result, err := tx.Exec(deleteIntervensiQuery, userId, currentTime, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to delete intervensi", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check delete result", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if rowsAffected == 0 {
		response := object.NewResponse(http.StatusNotFound, "Intervensi not found or already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to commit transaction", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with additional information
	message := fmt.Sprintf("Intervensi %s untuk balita '%s' tanggal %s berhasil dihapus",
		jenis, namaBalita, tanggal)

	response := object.NewResponse(http.StatusOK, "Intervensi deleted successfully", deleteIntervensiResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// # IntervensiRestore handles restoring soft deleted intervensi data
//
// @Summary Restore deleted intervensi data
// @Description Restore soft deleted intervensi data by clearing deleted_date and deleted_id (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param intervensi body deleteIntervensiRequest true "Intervensi ID to restore"
// @Success 200 {object} object.Response{data=deleteIntervensiResponse} "Intervensi restored successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Intervensi not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/intervensi/restore [post]
func IntervensiRestore(w http.ResponseWriter, r *http.Request) {
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
	var req deleteIntervensiRequest
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

	// Check if intervensi exists and is soft deleted
	var exists int
	var jenis, tanggal, deskripsi, idBalita, namaBalita string
	checkQuery := `SELECT COUNT(*), i.jenis, i.tanggal, i.deskripsi, i.id_balita, b.nama as nama_balita
        FROM intervensi i
        LEFT JOIN balita b ON i.id_balita = b.id
        WHERE i.id = ? AND i.deleted_date IS NOT NULL
        GROUP BY i.jenis, i.tanggal, i.deskripsi, i.id_balita, b.nama`
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists, &jenis, &tanggal, &deskripsi, &idBalita, &namaBalita)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Intervensi not found or not deleted", nil)
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

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Intervensi not found or not deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for duplicates before restore (same balita, jenis, tanggal, and deskripsi, not soft deleted)
	var duplicateExists int
	duplicateQuery := `SELECT COUNT(*) FROM intervensi 
        WHERE id_balita = ? AND jenis = ? AND tanggal = ? AND deskripsi = ? AND id != ? AND deleted_date IS NULL`
	err = db.QueryRow(duplicateQuery, idBalita, jenis, tanggal, deskripsi, req.Id).Scan(&duplicateExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check duplicate intervensi", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if duplicateExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot restore intervensi. Another active intervensi for balita '%s' with type '%s', date '%s', and similar description already exists",
				namaBalita, jenis, tanggal), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp for updated_date
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Restore intervensi (clear soft delete fields)
	restoreQuery := `UPDATE intervensi SET 
        deleted_id = NULL, 
        deleted_date = NULL,
        updated_id = ?,
        updated_date = ?
        WHERE id = ? AND deleted_date IS NOT NULL`

	result, err := db.Exec(restoreQuery, userId, currentTime, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to restore intervensi", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check restore result", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if rowsAffected == 0 {
		response := object.NewResponse(http.StatusNotFound, "Intervensi not found or not deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if related balita still exists and is not soft deleted
	var balitaExists int
	checkBalitaQuery := "SELECT COUNT(*) FROM balita WHERE id = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkBalitaQuery, idBalita).Scan(&balitaExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related balita", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if balitaExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Cannot restore intervensi. Related balita does not exist or is deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with additional information
	message := fmt.Sprintf("Intervensi %s untuk balita '%s' tanggal %s berhasil dipulihkan",
		jenis, namaBalita, tanggal)

	response := object.NewResponse(http.StatusOK, "Intervensi restored successfully", deleteIntervensiResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}