package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type deleteSkpdRequest struct {
	Id string `json:"id"`
}

func (r *deleteSkpdRequest) validate() error {
	if r.Id == "" {
		return fmt.Errorf("SKPD ID is required")
	}
	return nil
}

type deleteSkpdResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # DeleteSkpd handles soft deleting SKPD data
//
// @Summary Delete SKPD data (soft delete)
// @Description Soft delete SKPD data by setting deleted_date and deleted_id (Admin only)
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
// @Param skpd body deleteSkpdRequest true "SKPD ID to delete"
// @Success 200 {object} object.Response{data=deleteSkpdResponse} "SKPD deleted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "SKPD not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/skpd/delete [delete]
func AdminSkpdDelete(w http.ResponseWriter, r *http.Request) {
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
	var req deleteSkpdRequest
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

	// Check if SKPD exists and not already soft deleted
	var exists int
	var deletedDate sql.NullString
	checkQuery := "SELECT COUNT(*), deleted_date FROM skpd WHERE id = ? GROUP BY deleted_date"
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists, &deletedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "SKPD not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check SKPD existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "SKPD not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if already soft deleted
	if deletedDate.Valid {
		response := object.NewResponse(http.StatusBadRequest, "SKPD already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get SKPD details for additional information
	var skpdName, jenisSkpd string
	detailQuery := "SELECT skpd, jenis FROM skpd WHERE id = ?"
	err = db.QueryRow(detailQuery, req.Id).Scan(&skpdName, &jenisSkpd)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get SKPD details", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if SKPD has related petugas kesehatan records (prevent deletion if has staff)
	var petugasCount int
	checkPetugasQuery := "SELECT COUNT(*) FROM petugas_kesehatan WHERE id_skpd = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkPetugasQuery, req.Id).Scan(&petugasCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related petugas kesehatan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if petugasCount > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot delete SKPD '%s'. There are %d active petugas kesehatan records related to this SKPD", skpdName, petugasCount), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if SKPD has related intervensi records through petugas kesehatan (informational)
	var intervensiCount int
	checkIntervensiQuery := `SELECT COUNT(*) FROM intervensi_petugas ip 
        JOIN petugas_kesehatan pk ON ip.id_petugas_kesehatan = pk.id 
        WHERE pk.id_skpd = ? AND pk.deleted_date IS NULL`
	err = db.QueryRow(checkIntervensiQuery, req.Id).Scan(&intervensiCount)
	if err != nil {
		// Not critical, continue with deletion
		intervensiCount = 0
	}

	// Check if there are any soft deleted petugas kesehatan that could be restored later
	var deletedPetugasCount int
	checkDeletedPetugasQuery := "SELECT COUNT(*) FROM petugas_kesehatan WHERE id_skpd = ? AND deleted_date IS NOT NULL"
	err = db.QueryRow(checkDeletedPetugasQuery, req.Id).Scan(&deletedPetugasCount)
	if err != nil {
		// Not critical, continue with deletion
		deletedPetugasCount = 0
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Perform soft delete
	deleteQuery := `UPDATE skpd SET 
        deleted_id = ?, 
        deleted_date = ? 
        WHERE id = ? AND deleted_date IS NULL`

	result, err := db.Exec(deleteQuery, userId, currentTime, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to delete SKPD", nil)
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
		response := object.NewResponse(http.StatusNotFound, "SKPD not found or already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with additional information
	message := fmt.Sprintf("Data SKPD '%s' (jenis: %s) berhasil dihapus", skpdName, jenisSkpd)
	if deletedPetugasCount > 0 {
		message += fmt.Sprintf(" (Note: This SKPD had %d deleted petugas kesehatan)", deletedPetugasCount)
	}
	if intervensiCount > 0 {
		message += fmt.Sprintf(" (Note: Related to %d historical intervensi records)", intervensiCount)
	}

	response := object.NewResponse(http.StatusOK, "SKPD deleted successfully", deleteSkpdResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// # RestoreSkpd handles restoring soft deleted SKPD data (optional feature)
//
// @Summary Restore deleted SKPD data
// @Description Restore soft deleted SKPD data by clearing deleted_date and deleted_id (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param skpd body deleteSkpdRequest true "SKPD ID to restore"
// @Success 200 {object} object.Response{data=deleteSkpdResponse} "SKPD restored successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "SKPD not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/skpd/restore [post]
func AdminSkpdRestore(w http.ResponseWriter, r *http.Request) {
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
	var req deleteSkpdRequest
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

	// Check if SKPD exists and is soft deleted
	var exists int
	checkQuery := "SELECT COUNT(*) FROM skpd WHERE id = ? AND deleted_date IS NOT NULL"
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check SKPD existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "SKPD not found or not deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get SKPD details for validation and response
	var skpdName, jenisSkpd string
	detailQuery := "SELECT skpd, jenis FROM skpd WHERE id = ?"
	err = db.QueryRow(detailQuery, req.Id).Scan(&skpdName, &jenisSkpd)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get SKPD details", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for naming conflicts before restore (same name and jenis, not soft deleted)
	var conflictExists int
	conflictQuery := "SELECT COUNT(*) FROM skpd WHERE skpd = ? AND jenis = ? AND id != ? AND deleted_date IS NULL"
	err = db.QueryRow(conflictQuery, skpdName, jenisSkpd, req.Id).Scan(&conflictExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check naming conflicts", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if conflictExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot restore SKPD. Another SKPD with name '%s' and jenis '%s' already exists", skpdName, jenisSkpd), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for soft deleted petugas kesehatan that could be restored with this SKPD
	var deletedPetugasCount int
	checkDeletedPetugasQuery := "SELECT COUNT(*) FROM petugas_kesehatan WHERE id_skpd = ? AND deleted_date IS NOT NULL"
	err = db.QueryRow(checkDeletedPetugasQuery, req.Id).Scan(&deletedPetugasCount)
	if err != nil {
		// Not critical, continue with restore
		deletedPetugasCount = 0
	}

	// Current timestamp for updated_date
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Restore SKPD (clear soft delete fields)
	restoreQuery := `UPDATE skpd SET 
        deleted_id = NULL, 
        deleted_date = NULL,
        updated_id = ?,
        updated_date = ?
        WHERE id = ? AND deleted_date IS NOT NULL`

	result, err := db.Exec(restoreQuery, userId, currentTime, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to restore SKPD", nil)
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
		response := object.NewResponse(http.StatusNotFound, "SKPD not found or not deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with additional information
	message := fmt.Sprintf("Data SKPD '%s' (jenis: %s) berhasil dipulihkan", skpdName, jenisSkpd)
	if deletedPetugasCount > 0 {
		message += fmt.Sprintf(" (Note: This SKPD has %d deleted petugas kesehatan that can be restored separately)", deletedPetugasCount)
	}

	response := object.NewResponse(http.StatusOK, "SKPD restored successfully", deleteSkpdResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}