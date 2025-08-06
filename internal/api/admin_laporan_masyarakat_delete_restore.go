package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type deleteLaporanMasyarakatRequest struct {
	Id string `json:"id"`
}

func (r *deleteLaporanMasyarakatRequest) validate() error {
	if r.Id == "" {
		return fmt.Errorf("laporan masyarakat ID is required")
	}
	return nil
}

type deleteLaporanMasyarakatResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

// # DeleteLaporanMasyarakat handles soft deleting laporan masyarakat data
//
// @Summary Delete laporan masyarakat data (soft delete)
// @Description Soft delete laporan masyarakat data by setting deleted_date and deleted_id (Admin only)
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
// @Param laporan body deleteLaporanMasyarakatRequest true "Laporan Masyarakat ID to delete"
// @Success 200 {object} object.Response{data=deleteLaporanMasyarakatResponse} "Laporan masyarakat deleted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Laporan masyarakat not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/laporan-masyarakat/delete [delete]
func AdminLaporanMasyarakatDelete(w http.ResponseWriter, r *http.Request) {
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
	var req deleteLaporanMasyarakatRequest
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

	// Check if laporan masyarakat exists and not already soft deleted
	var exists int
	var deletedDate sql.NullString
	checkQuery := "SELECT COUNT(*), deleted_date FROM laporan_masyarakat WHERE id = ? GROUP BY deleted_date"
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists, &deletedDate)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusNotFound, "Laporan masyarakat not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check laporan masyarakat existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Laporan masyarakat not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if already soft deleted
	if deletedDate.Valid {
		response := object.NewResponse(http.StatusBadRequest, "Laporan masyarakat already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get laporan details for additional information
	var idBalita, statusLaporan, jenisLaporan string
	var idMasyarakat sql.NullString
	detailQuery := `
        SELECT lm.id_balita, lm.id_masyarakat, sl.status
        FROM laporan_masyarakat lm
        LEFT JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        WHERE lm.id = ?
    `
	err = db.QueryRow(detailQuery, req.Id).Scan(&idBalita, &idMasyarakat, &statusLaporan)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get laporan details", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Determine laporan type
	if idMasyarakat.Valid {
		jenisLaporan = "masyarakat"
	} else {
		jenisLaporan = "admin"
	}

	// Check if laporan has related riwayat pemeriksaan records (warn user if balita has medical history)
	var riwayatCount int
	checkRiwayatQuery := "SELECT COUNT(*) FROM riwayat_pemeriksaan WHERE id_balita = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkRiwayatQuery, idBalita).Scan(&riwayatCount)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related riwayat pemeriksaan", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prevent deletion if laporan is in processed status (optional business rule)
	processedStatuses := []string{"Diproses dan data sesuai", "Sudah ditindaklanjuti", "Sudah perbaikan gizi"}
	if slices.Contains(processedStatuses, statusLaporan) {
		response := object.NewResponse(http.StatusBadRequest,
			fmt.Sprintf("Cannot delete laporan masyarakat. Status '%s' indicates this report has been processed", statusLaporan), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Additional check: if laporan from masyarakat and has been responded, warn before deletion
	if jenisLaporan == "masyarakat" && statusLaporan != "Belum diproses" {
		// Could add a confirmation parameter here for frontend to handle
		// For now, we'll proceed with warning in response message
	}

	// Current timestamp
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Perform soft delete
	deleteQuery := `UPDATE laporan_masyarakat SET 
        deleted_id = ?, 
        deleted_date = ? 
        WHERE id = ? AND deleted_date IS NULL`

	result, err := db.Exec(deleteQuery, userId, currentTime, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to delete laporan masyarakat", nil)
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
		response := object.NewResponse(http.StatusNotFound, "Laporan masyarakat not found or already deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Prepare response message with additional information
	message := "Data laporan masyarakat berhasil dihapus"
	if riwayatCount > 0 {
		message += fmt.Sprintf(" (Note: Balita terkait memiliki %d riwayat pemeriksaan)", riwayatCount)
	}
	if jenisLaporan == "masyarakat" && statusLaporan != "Belum diproses" {
		message += " (Warning: This was a processed community report)"
	}

	response := object.NewResponse(http.StatusOK, "Laporan masyarakat deleted successfully", deleteLaporanMasyarakatResponse{
		Id:      req.Id,
		Message: message,
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// # RestoreLaporanMasyarakat handles restoring soft deleted laporan masyarakat data (optional feature)
//
// @Summary Restore deleted laporan masyarakat data
// @Description Restore soft deleted laporan masyarakat data by clearing deleted_date and deleted_id (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param laporan body deleteLaporanMasyarakatRequest true "Laporan Masyarakat ID to restore"
// @Success 200 {object} object.Response{data=deleteLaporanMasyarakatResponse} "Laporan masyarakat restored successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Laporan masyarakat not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/laporan-masyarakat/restore [post]
func AdminLaporanMasyarakatRestore(w http.ResponseWriter, r *http.Request) {
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
	var req deleteLaporanMasyarakatRequest
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

	// Check if laporan masyarakat exists and is soft deleted
	var exists int
	checkQuery := "SELECT COUNT(*) FROM laporan_masyarakat WHERE id = ? AND deleted_date IS NOT NULL"
	err = db.QueryRow(checkQuery, req.Id).Scan(&exists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check laporan masyarakat existence", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if exists == 0 {
		response := object.NewResponse(http.StatusNotFound, "Laporan masyarakat not found or not deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if related balita still exists and is not soft deleted
	var balitaExists int
	var laporanIdBalita string
	getBalitaQuery := "SELECT id_balita FROM laporan_masyarakat WHERE id = ?"
	err = db.QueryRow(getBalitaQuery, req.Id).Scan(&laporanIdBalita)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get laporan balita", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	checkBalitaQuery := "SELECT COUNT(*) FROM balita WHERE id = ? AND deleted_date IS NULL"
	err = db.QueryRow(checkBalitaQuery, laporanIdBalita).Scan(&balitaExists)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to check related balita", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if balitaExists == 0 {
		response := object.NewResponse(http.StatusBadRequest, "Cannot restore laporan masyarakat. Related balita does not exist or is deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if masyarakat still exists (if laporan from masyarakat)
	var idMasyarakat sql.NullString
	getMasyarakatQuery := "SELECT id_masyarakat FROM laporan_masyarakat WHERE id = ?"
	err = db.QueryRow(getMasyarakatQuery, req.Id).Scan(&idMasyarakat)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to get laporan masyarakat", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if idMasyarakat.Valid {
		var masyarakatExists int
		checkMasyarakatQuery := "SELECT COUNT(*) FROM masyarakat WHERE id = ?"
		err = db.QueryRow(checkMasyarakatQuery, idMasyarakat.String).Scan(&masyarakatExists)
		if err != nil {
			response := object.NewResponse(http.StatusInternalServerError, "Failed to check related masyarakat", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if masyarakatExists == 0 {
			response := object.NewResponse(http.StatusBadRequest, "Cannot restore laporan masyarakat. Related masyarakat does not exist", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}

	// Current timestamp for updated_date
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// Restore laporan masyarakat (clear soft delete fields)
	restoreQuery := `UPDATE laporan_masyarakat SET 
        deleted_id = NULL, 
        deleted_date = NULL,
        updated_id = ?,
        updated_date = ?
        WHERE id = ? AND deleted_date IS NOT NULL`

	result, err := db.Exec(restoreQuery, userId, currentTime, req.Id)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to restore laporan masyarakat", nil)
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
		response := object.NewResponse(http.StatusNotFound, "Laporan masyarakat not found or not deleted", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Laporan masyarakat restored successfully", deleteLaporanMasyarakatResponse{
		Id:      req.Id,
		Message: "Data laporan masyarakat berhasil dipulihkan",
	})
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
