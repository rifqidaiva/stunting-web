package api

import (
	"database/sql"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type keluargaResponse struct {
    Id          string    `json:"id"`
    NomorKk     string    `json:"nomor_kk"`
    NamaAyah    string    `json:"nama_ayah"`
    NamaIbu     string    `json:"nama_ibu"`
    NikAyah     string    `json:"nik_ayah"`
    NikIbu      string    `json:"nik_ibu"`
    Alamat      string    `json:"alamat"`
    Rt          string    `json:"rt"`
    Rw          string    `json:"rw"`
    IdKelurahan string    `json:"id_kelurahan"`
    Kelurahan   string    `json:"kelurahan"`
    Kecamatan   string    `json:"kecamatan"`
    Koordinat   [2]float64 `json:"koordinat"`
    CreatedDate string    `json:"created_date"`
    UpdatedDate string    `json:"updated_date,omitempty"`
}

type getAllKeluargaResponse struct {
    Data  []keluargaResponse `json:"data"`
    Total int                `json:"total"`
}

type getKeluargaByIdResponse struct {
    Data keluargaResponse `json:"data"`
}

// # AdminKeluargaGet handles getting keluarga data
//
// @Summary Get keluarga data
// @Description Get keluarga data based on query parameter (Admin only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all keluarga with total count
// @Description - With id parameter: Returns specific keluarga data
// @Description
// @Description Keluarga data includes: nomor_kk, nama_ayah, nama_ibu, nik_ayah, nik_ibu, alamat, rt, rw, kelurahan, kecamatan, koordinat
// @Tags admin
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Keluarga ID"
// @Success 200 {object} object.Response{data=getAllKeluargaResponse} "Keluarga data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 404 {object} object.Response{data=nil} "Keluarga not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/keluarga/get [get]
func AdminKeluargaGet(w http.ResponseWriter, r *http.Request) {
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
        // Get specific keluarga by ID
        keluarga, err := getKeluargaById(db, idParam)
        if err != nil {
            if err == sql.ErrNoRows {
                response := object.NewResponse(http.StatusNotFound, "Keluarga not found", nil)
                if err := response.WriteJson(w); err != nil {
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                }
                return
            }
            response := object.NewResponse(http.StatusInternalServerError, "Failed to get keluarga", nil)
            if err := response.WriteJson(w); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            return
        }

        response := object.NewResponse(http.StatusOK, "Keluarga retrieved successfully", getKeluargaByIdResponse{
            Data: keluarga,
        })
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else {
        // Get all keluarga
        keluargaList, total, err := getAllKeluarga(db)
        if err != nil {
            response := object.NewResponse(http.StatusInternalServerError, "Failed to get keluarga list", nil)
            if err := response.WriteJson(w); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            return
        }

        response := object.NewResponse(http.StatusOK, "All keluarga retrieved successfully", getAllKeluargaResponse{
            Data:  keluargaList,
            Total: total,
        })
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    }
}

// Helper function to get keluarga by ID
func getKeluargaById(db *sql.DB, id string) (keluargaResponse, error) {
    var keluarga keluargaResponse
    var koordinatWKT string
    var updatedDate sql.NullString

    query := `
        SELECT 
            k.id, k.nomor_kk, k.nama_ayah, k.nama_ibu, k.nik_ayah, k.nik_ibu,
            k.alamat, k.rt, k.rw, k.id_kelurahan,
            kel.kelurahan, kec.kecamatan,
            ST_AsText(k.koordinat) as koordinat_wkt,
            k.created_date, k.updated_date
        FROM keluarga k
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        WHERE k.id = ? AND k.deleted_date IS NULL
    `

    err := db.QueryRow(query, id).Scan(
        &keluarga.Id,
        &keluarga.NomorKk,
        &keluarga.NamaAyah,
        &keluarga.NamaIbu,
        &keluarga.NikAyah,
        &keluarga.NikIbu,
        &keluarga.Alamat,
        &keluarga.Rt,
        &keluarga.Rw,
        &keluarga.IdKelurahan,
        &keluarga.Kelurahan,
        &keluarga.Kecamatan,
        &koordinatWKT,
        &keluarga.CreatedDate,
        &updatedDate,
    )

    if err != nil {
        return keluarga, err
    }

    // Parse WKT coordinates
    keluarga.Koordinat = object.ParseWKT(koordinatWKT)

    // Handle nullable updated_date
    if updatedDate.Valid {
        keluarga.UpdatedDate = updatedDate.String
    }

    return keluarga, nil
}

// Helper function to get all keluarga
func getAllKeluarga(db *sql.DB) ([]keluargaResponse, int, error) {
    var keluargaList []keluargaResponse

    query := `
        SELECT 
            k.id, k.nomor_kk, k.nama_ayah, k.nama_ibu, k.nik_ayah, k.nik_ibu,
            k.alamat, k.rt, k.rw, k.id_kelurahan,
            kel.kelurahan, kec.kecamatan,
            ST_AsText(k.koordinat) as koordinat_wkt,
            k.created_date, k.updated_date
        FROM keluarga k
        LEFT JOIN kelurahan kel ON k.id_kelurahan = kel.id
        LEFT JOIN kecamatan kec ON kel.id_kecamatan = kec.id
        WHERE k.deleted_date IS NULL
        ORDER BY k.created_date DESC
    `

    rows, err := db.Query(query)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    for rows.Next() {
        var keluarga keluargaResponse
        var koordinatWKT string
        var updatedDate sql.NullString

        err := rows.Scan(
            &keluarga.Id,
            &keluarga.NomorKk,
            &keluarga.NamaAyah,
            &keluarga.NamaIbu,
            &keluarga.NikAyah,
            &keluarga.NikIbu,
            &keluarga.Alamat,
            &keluarga.Rt,
            &keluarga.Rw,
            &keluarga.IdKelurahan,
            &keluarga.Kelurahan,
            &keluarga.Kecamatan,
            &koordinatWKT,
            &keluarga.CreatedDate,
            &updatedDate,
        )

        if err != nil {
            return nil, 0, err
        }

        // Parse WKT coordinates
        keluarga.Koordinat = object.ParseWKT(koordinatWKT)

        // Handle nullable updated_date
        if updatedDate.Valid {
            keluarga.UpdatedDate = updatedDate.String
        }

        keluargaList = append(keluargaList, keluarga)
    }

    if err = rows.Err(); err != nil {
        return nil, 0, err
    }

    // Get total count
    var total int
    countQuery := "SELECT COUNT(*) FROM keluarga WHERE deleted_date IS NULL"
    err = db.QueryRow(countQuery).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    return keluargaList, total, nil
}