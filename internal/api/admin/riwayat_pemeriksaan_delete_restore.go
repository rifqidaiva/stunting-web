package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type deleteRiwayatPemeriksaanRequest struct {
	Id string `json:"id"`
}

func (r *deleteRiwayatPemeriksaanRequest) validate() error {
	if r.Id == "" {
		return fmt.Errorf("riwayat pemeriksaan ID is required")
	}
	return nil
}

type deleteRiwayatPemeriksaanResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # RiwayatPemeriksaanDelete handles soft deleting riwayat pemeriksaan data
//
// @Summary Delete riwayat pemeriksaan data (soft delete)
// @Description Soft delete riwayat pemeriksaan data by setting deleted_date and deleted_id (Admin only)
// @Description
// @Description Performs soft delete operation:
// @Description - Sets deleted_date to current timestamp
// @Description - Sets deleted_id to current user ID
// @Description - Data remains in database but is excluded from queries
// @Description - Can be restored if needed in the future
// @Description - Provides detailed information about the deleted medical record
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param riwayat body deleteRiwayatPemeriksaanRequest true "Riwayat Pemeriksaan ID to delete"
// @Success 200 {object} object.Response{data=deleteRiwayatPemeriksaanResponse} "Riwayat pemeriksaan deleted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Riwayat pemeriksaan not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/riwayat-pemeriksaan/delete [delete]
func RiwayatPemeriksaanDelete(w http.ResponseWriter, r *http.Request) {
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
	var req deleteRiwayatPemeriksaanRequest
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

	// Check if riwayat pemeriksaan exists and not already soft deleted, also get current data
	var exists int
	var deletedDate sql.NullString
	var tanggal, statusGizi, namaBalita, jenisIntervensi string
	checkQuery := `SELECT COUNT(*), rp.deleted_date, rp.tanggal, rp.status_gizi, 
        b.nama as nama_balita, i.jenis as jenis_intervensi
        FROM riwayat_pemeriksaan rp
        LEFT JOIN balita b ON rp.id_balita = b.id
        LEFT JOIN intervensi i ON rp.id_intervensi = i.id
        WHERE rp.id = ? 
        GROUP BY rp.deleted_date, rp.tanggal, rp.status_gizi, b.nama, i.jenis`
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists, &deletedDate, &tanggal, &statusGizi, &namaBalita, &jenisIntervensi)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Riwayat pemeriksaan not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check riwayat pemeriksaan existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Riwayat pemeriksaan not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if already soft deleted
	if deletedDate.Valid {
		response := object.NewResponse(http.StatusBadRequest, "Riwayat pemeriksaan already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Additional business rule checks can be added here if needed
	// For example, preventing deletion of recent medical records or critical status records

	// Check if this is a critical medical record (e.g., stunting or gizi buruk status)
	var isCritical bool
	if statusGizi == "stunting" || statusGizi == "gizi buruk" {
		isCritical = true
	}

	// Count total riwayat pemeriksaan for this balita (to prevent deletion of last record)
	var totalRiwayatBalita int
	countRiwayatQuery := `SELECT COUNT(*) FROM riwayat_pemeriksaan rp
        LEFT JOIN balita b ON rp.id_balita = b.id
        WHERE b.nama = ? AND rp.deleted_date IS NULL`
	err = db.QueryRow(countRiwayatQuery, namaBalita).Scan(&totalRiwayatBalita)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check balita medical history count", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Warning if this is the only medical record for this balita
	var warningMessage string
	if totalRiwayatBalita == 1 {
		warningMessage = " (Warning: This is the only medical record for this balita)"
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Perform soft delete
	deleteQuery := `UPDATE riwayat_pemeriksaan SET 
        deleted_id = ?, 
        deleted_date = ? 
        WHERE id = ? AND deleted_date IS NULL`

	result, err := db.Exec(deleteQuery, userId, currentTime, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to delete riwayat pemeriksaan", nil)
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
		response := object.NewResponse(http.StatusNotFound, "Riwayat pemeriksaan not found or already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with detailed information
	message := fmt.Sprintf("Riwayat pemeriksaan balita '%s' untuk intervensi %s pada tanggal %s berhasil dihapus",
		namaBalita, jenisIntervensi, tanggal)

	// Add critical status warning if applicable
	if isCritical {
		message += fmt.Sprintf(" (Note: This was a critical medical record with status '%s')", statusGizi)
	}

	// Add warning message if applicable
	message += warningMessage

	response := object.NewResponse(http.StatusOK, "Riwayat pemeriksaan deleted successfully", deleteRiwayatPemeriksaanResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// # RiwayatPemeriksaanRestore handles restoring soft deleted riwayat pemeriksaan data
//
// @Summary Restore deleted riwayat pemeriksaan data
// @Description Restore soft deleted riwayat pemeriksaan data by clearing deleted_date and deleted_id (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param riwayat body deleteRiwayatPemeriksaanRequest true "Riwayat Pemeriksaan ID to restore"
// @Success 200 {object} object.Response{data=deleteRiwayatPemeriksaanResponse} "Riwayat pemeriksaan restored successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Riwayat pemeriksaan not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/riwayat-pemeriksaan/restore [post]
func RiwayatPemeriksaanRestore(w http.ResponseWriter, r *http.Request) {
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
	var req deleteRiwayatPemeriksaanRequest
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

	// Check if riwayat pemeriksaan exists and is soft deleted, also get current data
	var exists int
	var tanggal, statusGizi, namaBalita, jenisIntervensi string
	var idBalita, idIntervensi string
	checkQuery := `SELECT COUNT(*), rp.tanggal, rp.status_gizi, rp.id_balita, rp.id_intervensi,
        b.nama as nama_balita, i.jenis as jenis_intervensi
        FROM riwayat_pemeriksaan rp
        LEFT JOIN balita b ON rp.id_balita = b.id
        LEFT JOIN intervensi i ON rp.id_intervensi = i.id
        WHERE rp.id = ? AND rp.deleted_date IS NOT NULL
        GROUP BY rp.tanggal, rp.status_gizi, rp.id_balita, rp.id_intervensi, b.nama, i.jenis`
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists, &tanggal, &statusGizi, &idBalita, &idIntervensi, &namaBalita, &jenisIntervensi)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Riwayat pemeriksaan not found or not deleted", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check riwayat pemeriksaan existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Riwayat pemeriksaan not found or not deleted", nil)
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
		response := object.NewResponse(http.StatusBadRequest, "Cannot restore riwayat pemeriksaan. Related balita does not exist or is deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if related intervensi still exists and is not soft deleted
	var intervensiExists int
	checkIntervensiQuery := "SELECT COUNT(*) FROM intervensi WHERE id = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkIntervensiQuery, idIntervensi).Scan(&intervensiExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related intervensi", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if intervensiExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Cannot restore riwayat pemeriksaan. Related intervensi does not exist or is deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check for duplicates before restore (same balita, intervensi, and tanggal, not soft deleted)
	var duplicateExists int
	duplicateQuery := `SELECT COUNT(*) FROM riwayat_pemeriksaan 
        WHERE id_balita = ? AND id_intervensi = ? AND tanggal = ? AND id != ? AND deleted_date IS NULL`
	err = db.QueryRow(duplicateQuery, idBalita, idIntervensi, tanggal, req.Id).Scan(&duplicateExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check duplicate riwayat pemeriksaan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if duplicateExists > 0 {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot restore riwayat pemeriksaan. Another active medical record for balita '%s' on date '%s' in intervensi '%s' already exists",
				namaBalita, tanggal, jenisIntervensi), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Current timestamp for updated_date
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Restore riwayat pemeriksaan (clear soft delete fields)
	restoreQuery := `UPDATE riwayat_pemeriksaan SET 
        deleted_id = NULL, 
        deleted_date = NULL,
        updated_id = ?,
        updated_date = ?
        WHERE id = ? AND deleted_date IS NOT NULL`

	result, err := db.Exec(restoreQuery, userId, currentTime, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to restore riwayat pemeriksaan", nil)
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
		response := object.NewResponse(http.StatusNotFound, "Riwayat pemeriksaan not found or not deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with detailed information
	message := fmt.Sprintf("Riwayat pemeriksaan balita '%s' untuk intervensi %s pada tanggal %s berhasil dipulihkan",
		namaBalita, jenisIntervensi, tanggal)

	// Add status information
	message += fmt.Sprintf(" (Status gizi: %s)", statusGizi)

	response := object.NewResponse(http.StatusOK, "Riwayat pemeriksaan restored successfully", deleteRiwayatPemeriksaanResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
