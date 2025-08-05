package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type deleteKeluargaRequest struct {
    Id string `json:"id"`
}

func (r *deleteKeluargaRequest) validate() error {
    if r.Id == "" {
        return fmt.Errorf("keluarga ID is required")
    }
    return nil
}

type deleteKeluargaResponse struct {
    Id      string `json:"id"`
    Message string `json:"message"`
}

// # DeleteKeluarga handles soft deleting keluarga data
//
// @Summary Delete keluarga data (soft delete)
// @Description Soft delete keluarga data by setting deleted_date and deleted_id (Admin only)
// @Description
// @Description Performs soft delete operation:
// @Description - Sets deleted_date to current timestamp
// @Description - Sets deleted_id to current user ID
// @Description - Data remains in database but is excluded from queries
// @Description - Can be restored if needed in the future
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param keluarga body deleteKeluargaRequest true "Keluarga ID to delete"
// @Success 200 {object} object.Response{data=deleteKeluargaResponse} "Keluarga deleted successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Keluarga not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/keluarga/delete [delete]
func AdminKeluargaDelete(w http.ResponseWriter, r *http.Request) {
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
    var req deleteKeluargaRequest
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

    // Check if keluarga exists and not already soft deleted
    var exists int
    var deletedDate sql.NullString
    checkQuery := "SELECT COUNT(*), deleted_date FROM keluarga WHERE id = ? GROUP BY deleted_date"
    err = db.QueryRow(checkQuery, req.Id).Scan(&exists, &deletedDate)
    if err != nil {
        if err == sql.ErrNoRows {
            response := object.NewResponse(http.StatusNotFound, "Keluarga not found", nil)
            if err := response.WriteJson(w); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            return
        }
        response := object.NewResponse(http.StatusInternalServerError, "Failed to check keluarga existence", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    if exists == 0 {
        response := object.NewResponse(http.StatusNotFound, "Keluarga not found", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Check if already soft deleted
    if deletedDate.Valid {
        response := object.NewResponse(http.StatusBadRequest, "Keluarga already deleted", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Check if keluarga has related balita records (prevent deletion if has children)
    var balitaCount int
    checkBalitaQuery := "SELECT COUNT(*) FROM balita WHERE id_keluarga = ? AND deleted_date IS NULL"
    err = db.QueryRow(checkBalitaQuery, req.Id).Scan(&balitaCount)
    if err != nil {
        response := object.NewResponse(http.StatusInternalServerError, "Failed to check related balita", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    if balitaCount > 0 {
        response := object.NewResponse(http.StatusBadRequest, 
            fmt.Sprintf("Cannot delete keluarga. There are %d active balita records related to this keluarga", balitaCount), nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Check if keluarga has related laporan masyarakat records (prevent deletion if has reports)
    var laporanCount int
    checkLaporanQuery := `SELECT COUNT(*) FROM laporan_masyarakat lm 
        JOIN balita b ON lm.id_balita = b.id 
        WHERE b.id_keluarga = ? AND lm.deleted_date IS NULL AND b.deleted_date IS NULL`
    err = db.QueryRow(checkLaporanQuery, req.Id).Scan(&laporanCount)
    if err != nil {
        response := object.NewResponse(http.StatusInternalServerError, "Failed to check related laporan", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    if laporanCount > 0 {
        response := object.NewResponse(http.StatusBadRequest, 
            fmt.Sprintf("Cannot delete keluarga. There are %d active laporan records related to this keluarga", laporanCount), nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Current timestamp
    currentTime := time.Now().Format("2006-01-02 15:04:05")

    // Perform soft delete
    deleteQuery := `UPDATE keluarga SET 
        deleted_id = ?, 
        deleted_date = ? 
        WHERE id = ? AND deleted_date IS NULL`

    result, err := db.Exec(deleteQuery, userId, currentTime, req.Id)
    if err != nil {
        response := object.NewResponse(http.StatusInternalServerError, "Failed to delete keluarga", nil)
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
        response := object.NewResponse(http.StatusNotFound, "Keluarga not found or already deleted", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    response := object.NewResponse(http.StatusOK, "Keluarga deleted successfully", deleteKeluargaResponse{
        Id:      req.Id,
        Message: "Data keluarga berhasil dihapus",
    })
    if err := response.WriteJson(w); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// # RestoreKeluarga handles restoring soft deleted keluarga data (optional feature)
//
// @Summary Restore deleted keluarga data
// @Description Restore soft deleted keluarga data by clearing deleted_date and deleted_id (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param keluarga body deleteKeluargaRequest true "Keluarga ID to restore"
// @Success 200 {object} object.Response{data=deleteKeluargaResponse} "Keluarga restored successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Keluarga not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/keluarga/restore [post]
func AdminKeluargaRestore(w http.ResponseWriter, r *http.Request) {
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
    var req deleteKeluargaRequest
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

    // Check if keluarga exists and is soft deleted
    var exists int
    checkQuery := "SELECT COUNT(*) FROM keluarga WHERE id = ? AND deleted_date IS NOT NULL"
    err = db.QueryRow(checkQuery, req.Id).Scan(&exists)
    if err != nil {
        response := object.NewResponse(http.StatusInternalServerError, "Failed to check keluarga existence", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    if exists == 0 {
        response := object.NewResponse(http.StatusNotFound, "Keluarga not found or not deleted", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Current timestamp for updated_date
    currentTime := time.Now().Format("2006-01-02 15:04:05")

    // Restore keluarga (clear soft delete fields)
    restoreQuery := `UPDATE keluarga SET 
        deleted_id = NULL, 
        deleted_date = NULL,
        updated_id = ?,
        updated_date = ?
        WHERE id = ? AND deleted_date IS NOT NULL`

    result, err := db.Exec(restoreQuery, userId, currentTime, req.Id)
    if err != nil {
        response := object.NewResponse(http.StatusInternalServerError, "Failed to restore keluarga", nil)
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
        response := object.NewResponse(http.StatusNotFound, "Keluarga not found or not deleted", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    response := object.NewResponse(http.StatusOK, "Keluarga restored successfully", deleteKeluargaResponse{
        Id:      req.Id,
        Message: "Data keluarga berhasil dipulihkan",
    })
    if err := response.WriteJson(w); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}