package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type deletePetugasKesehatanRequest struct {
	Id string `json:"id"`
}

func (r *deletePetugasKesehatanRequest) validate() error {
	if r.Id == "" {
		return fmt.Errorf("petugas kesehatan ID is required")
	}
	return nil
}

type deletePetugasKesehatanResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # DeletePetugasKesehatan handles soft deleting petugas kesehatan data
//
// @Summary Delete petugas kesehatan data (soft delete)
// @Description Soft delete petugas kesehatan data by setting deleted_date and deleted_id (Admin only)
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
// @Param petugas body deletePetugasKesehatanRequest true "Petugas Kesehatan ID to delete"
// @Success 200 {object} object.Response{data=deletePetugasKesehatanResponse} "Petugas kesehatan deleted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Petugas kesehatan not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/petugas-kesehatan/delete [delete]
func AdminPetugasKesehatanDelete(w http.ResponseWriter, r *http.Request) {
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
	var req deletePetugasKesehatanRequest
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

	// Check if petugas kesehatan exists and not already soft deleted
	var exists int
	var deletedDate sql.NullString
	var idPengguna, nama, skpdName, jenisSkpd string
	checkQuery := `SELECT COUNT(*), pk.deleted_date, pk.id_pengguna, pk.nama, s.skpd, s.jenis
        FROM petugas_kesehatan pk
        LEFT JOIN skpd s ON pk.id_skpd = s.id
        WHERE pk.id = ? 
        GROUP BY pk.deleted_date, pk.id_pengguna, pk.nama, s.skpd, s.jenis`
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists, &deletedDate, &idPengguna, &nama, &skpdName, &jenisSkpd)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found", nil)
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

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if already soft deleted
	if deletedDate.Valid {
		response := object.NewResponse(http.StatusBadRequest, "Petugas kesehatan already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if petugas kesehatan has related intervensi records (prevent deletion if has interventions)
	var intervensiCount int
	checkIntervensiQuery := "SELECT COUNT(*) FROM intervensi_petugas WHERE id_petugas_kesehatan = ?"
	err = db.QueryRow(checkIntervensiQuery, req.Id).Scan(&intervensiCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related intervensi", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if intervensiCount > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot delete petugas kesehatan '%s'. There are %d intervensi records related to this petugas", nama, intervensiCount), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if petugas kesehatan has related riwayat pemeriksaan records (prevent deletion if has medical records)
	var riwayatCount int
	checkRiwayatQuery := `SELECT COUNT(*) FROM riwayat_pemeriksaan rp 
        JOIN intervensi i ON rp.id_intervensi = i.id 
        JOIN intervensi_petugas ip ON i.id = ip.id_intervensi 
        WHERE ip.id_petugas_kesehatan = ? AND rp.deleted_date IS NULL`
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
			fmt.Sprintf("Cannot delete petugas kesehatan '%s'. There are %d riwayat pemeriksaan records related to this petugas", nama, riwayatCount), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Begin transaction (to handle both petugas_kesehatan and pengguna tables)
	tx, err := db.Begin()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to begin transaction", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer tx.Rollback()

	// Perform soft delete on petugas_kesehatan table
	deletePetugasQuery := `UPDATE petugas_kesehatan SET 
        deleted_id = ?, 
        deleted_date = ? 
        WHERE id = ? AND deleted_date IS NULL`

	result, err := tx.Exec(deletePetugasQuery, userId, currentTime, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to delete petugas kesehatan", nil)
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
		response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found or already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Optional: You might want to deactivate the related pengguna account
	// This depends on your business requirements
	// For now, we'll keep the pengguna account active but add a comment field if needed

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
	message := fmt.Sprintf("Data petugas kesehatan '%s' dari %s '%s' berhasil dihapus", nama, jenisSkpd, skpdName)

	response := object.NewResponse(http.StatusOK, "Petugas kesehatan deleted successfully", deletePetugasKesehatanResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// # RestorePetugasKesehatan handles restoring soft deleted petugas kesehatan data
//
// @Summary Restore deleted petugas kesehatan data
// @Description Restore soft deleted petugas kesehatan data by clearing deleted_date and deleted_id (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param petugas body deletePetugasKesehatanRequest true "Petugas Kesehatan ID to restore"
// @Success 200 {object} object.Response{data=deletePetugasKesehatanResponse} "Petugas kesehatan restored successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Petugas kesehatan not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/petugas-kesehatan/restore [post]
func AdminPetugasKesehatanRestore(w http.ResponseWriter, r *http.Request) {
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
	var req deletePetugasKesehatanRequest
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

	// Check if petugas kesehatan exists and is soft deleted
	var exists int
	var nama, email, idSkpd, skpdName, jenisSkpd string
	checkQuery := `SELECT COUNT(*), pk.nama, p.email, pk.id_skpd, s.skpd, s.jenis
        FROM petugas_kesehatan pk
        LEFT JOIN pengguna p ON pk.id_pengguna = p.id
        LEFT JOIN skpd s ON pk.id_skpd = s.id
        WHERE pk.id = ? AND pk.deleted_date IS NOT NULL
        GROUP BY pk.nama, p.email, pk.id_skpd, s.skpd, s.jenis`
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists, &nama, &email, &idSkpd, &skpdName, &jenisSkpd)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found or not deleted", nil)
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

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found or not deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if related SKPD still exists and is not soft deleted
	var skpdExists int
	checkSkpdQuery := "SELECT COUNT(*) FROM skpd WHERE id = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkSkpdQuery, idSkpd).Scan(&skpdExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related SKPD", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if skpdExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Cannot restore petugas kesehatan. Related SKPD does not exist or is deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for naming conflicts before restore (same name in same SKPD, not soft deleted)
	var nameConflictExists int
	nameConflictQuery := "SELECT COUNT(*) FROM petugas_kesehatan WHERE nama = ? AND id_skpd = ? AND id != ? AND deleted_date IS NULL"
	err = db.QueryRow(nameConflictQuery, nama, idSkpd, req.Id).Scan(&nameConflictExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check naming conflicts", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if nameConflictExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot restore petugas kesehatan. Another petugas with name '%s' already exists in %s '%s'", nama, jenisSkpd, skpdName), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for email conflicts before restore
	var emailConflictExists int
	emailConflictQuery := `SELECT COUNT(*) FROM pengguna p1 
        JOIN petugas_kesehatan pk1 ON p1.id = pk1.id_pengguna 
        WHERE p1.email = ? AND pk1.id != ? AND pk1.deleted_date IS NULL`
	err = db.QueryRow(emailConflictQuery, email, req.Id).Scan(&emailConflictExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check email conflicts", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if emailConflictExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot restore petugas kesehatan. Another active petugas with email '%s' already exists", email), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp for updated_date
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Restore petugas kesehatan (clear soft delete fields)
	restoreQuery := `UPDATE petugas_kesehatan SET 
        deleted_id = NULL, 
        deleted_date = NULL,
        updated_id = ?,
        updated_date = ?
        WHERE id = ? AND deleted_date IS NOT NULL`

	result, err := db.Exec(restoreQuery, userId, currentTime, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to restore petugas kesehatan", nil)
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
		response := object.NewResponse(http.StatusNotFound, "Petugas kesehatan not found or not deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with additional information
	message := fmt.Sprintf("Data petugas kesehatan '%s' dari %s '%s' berhasil dipulihkan", nama, jenisSkpd, skpdName)

	response := object.NewResponse(http.StatusOK, "Petugas kesehatan restored successfully", deletePetugasKesehatanResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
