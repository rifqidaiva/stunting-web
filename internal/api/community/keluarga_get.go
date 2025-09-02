package community

import (
	"database/sql"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type keluargaResponse struct {
    Id          string     `json:"id"`
    NomorKk     string     `json:"nomor_kk"`
    NamaAyah    string     `json:"nama_ayah"`
    NamaIbu     string     `json:"nama_ibu"`
    NikAyah     string     `json:"nik_ayah"`
    NikIbu      string     `json:"nik_ibu"`
    Alamat      string     `json:"alamat"`
    Rt          string     `json:"rt"`
    Rw          string     `json:"rw"`
    IdKelurahan string     `json:"id_kelurahan"`
    Kelurahan   string     `json:"kelurahan"`
    Kecamatan   string     `json:"kecamatan"`
    Koordinat   [2]float64 `json:"koordinat"`
    CreatedDate string     `json:"created_date"`
    UpdatedDate string     `json:"updated_date,omitempty"`
    
    // Status untuk masyarakat
    JumlahBalita      int    `json:"jumlah_balita"`
    JumlahLaporan     int    `json:"jumlah_laporan"`
    StatusLaporanAktif string `json:"status_laporan_aktif,omitempty"`
    CanEdit           bool   `json:"can_edit"`
}

type getAllKeluargaResponse struct {
    Data  []keluargaResponse `json:"data"`
    Total int                `json:"total"`
}

type getKeluargaByIdResponse struct {
    Data keluargaResponse `json:"data"`
}

// # KeluargaGet handles getting keluarga data for masyarakat
//
// @Summary Get keluarga data (Community)
// @Description Get keluarga data for community/masyarakat users (own data only)
// @Description
// @Description Response data varies by parameter:
// @Description - Without id parameter: Returns all keluarga created by the user
// @Description - With id parameter: Returns specific keluarga data (if owned by user)
// @Description
// @Description Data includes family information, balita count, laporan status, and edit permissions.
// @Description Users can only access keluarga data they have created themselves.
// @Tags community
// @Accept json
// @Produce json
// @Security Bearer
// @Param id query string false "Keluarga ID"
// @Success 200 {object} object.Response{data=getAllKeluargaResponse} "Keluarga data retrieved successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden - Masyarakat role required"
// @Failure 404 {object} object.Response{data=nil} "Keluarga not found or not owned by user"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/community/keluarga/get [get]
func KeluargaGet(w http.ResponseWriter, r *http.Request) {
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

    userId, role, err := object.ParseJWT(token)
    if err != nil {
        response := object.NewResponse(http.StatusUnauthorized, "Invalid or expired token", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Check if user is masyarakat
    if role != "masyarakat" {
        response := object.NewResponse(http.StatusForbidden, "Access denied. Masyarakat role required", nil)
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

    // Verify user exists and get masyarakat ID
    var masyarakatId string
    checkUserQuery := "SELECT m.id FROM masyarakat m JOIN pengguna p ON m.id_pengguna = p.id WHERE p.id = ?"
    err = db.QueryRow(checkUserQuery, userId).Scan(&masyarakatId)
    if err != nil {
        response := object.NewResponse(http.StatusUnauthorized, "Masyarakat profile not found", nil)
        if err := response.WriteJson(w); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    // Check if ID parameter is provided
    idParam := r.URL.Query().Get("id")
    if idParam != "" {
        // Get specific keluarga by ID (only if owned by user)
        keluarga, err := getKeluargaByIdForUser(db, idParam, userId)
        if err != nil {
            if err == sql.ErrNoRows {
                response := object.NewResponse(http.StatusNotFound, "Keluarga not found or not owned by you", nil)
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
        // Get all keluarga for this user
        keluargaList, total, err := getAllKeluargaForUser(db, userId)
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

// Helper function to get keluarga by ID for specific user (ownership check)
func getKeluargaByIdForUser(db *sql.DB, id string, userId string) (keluargaResponse, error) {
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
        WHERE k.id = ? AND k.created_id = ? AND k.deleted_date IS NULL
    `

    err := db.QueryRow(query, id, userId).Scan(
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

    // Get additional information for masyarakat
    err = getKeluargaAdditionalInfo(db, &keluarga)
    if err != nil {
        // Log error but don't fail the request
        // Additional info is not critical
    }

    return keluarga, nil
}

// Helper function to get all keluarga for specific user
func getAllKeluargaForUser(db *sql.DB, userId string) ([]keluargaResponse, int, error) {
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
        WHERE k.created_id = ? AND k.deleted_date IS NULL
        ORDER BY k.created_date DESC
    `

    rows, err := db.Query(query, userId)
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

        // Get additional information for masyarakat
        err = getKeluargaAdditionalInfo(db, &keluarga)
        if err != nil {
            // Log error but don't fail the request
            // Set default values
            keluarga.JumlahBalita = 0
            keluarga.JumlahLaporan = 0
            keluarga.CanEdit = true
        }

        keluargaList = append(keluargaList, keluarga)
    }

    if err = rows.Err(); err != nil {
        return nil, 0, err
    }

    // Get total count for this user
    var total int
    countQuery := "SELECT COUNT(*) FROM keluarga WHERE created_id = ? AND deleted_date IS NULL"
    err = db.QueryRow(countQuery, userId).Scan(&total)
    if err != nil {
        return nil, 0, err
    }

    return keluargaList, total, nil
}

// Helper function to get additional information for keluarga (balita count, laporan status, etc.)
func getKeluargaAdditionalInfo(db *sql.DB, keluarga *keluargaResponse) error {
    // Get balita count
    var balitaCount int
    balitaQuery := "SELECT COUNT(*) FROM balita WHERE id_keluarga = ? AND deleted_date IS NULL"
    err := db.QueryRow(balitaQuery, keluarga.Id).Scan(&balitaCount)
    if err != nil {
        balitaCount = 0
    }
    keluarga.JumlahBalita = balitaCount

    // Get laporan count and active status
    var laporanCount int
    var activeStatus sql.NullString
    laporanQuery := `
        SELECT COUNT(*), 
            GROUP_CONCAT(
                CASE 
                    WHEN sl.status IN ('Belum diproses', 'Diproses dan data sesuai', 'Belum ditindaklanjuti') 
                    THEN sl.status 
                    ELSE NULL 
                END
            ) as active_statuses
        FROM laporan_masyarakat lm
        JOIN balita b ON lm.id_balita = b.id
        JOIN status_laporan sl ON lm.id_status_laporan = sl.id
        WHERE b.id_keluarga = ? 
        AND lm.deleted_date IS NULL 
        AND b.deleted_date IS NULL
    `
    err = db.QueryRow(laporanQuery, keluarga.Id).Scan(&laporanCount, &activeStatus)
    if err != nil {
        laporanCount = 0
    }
    keluarga.JumlahLaporan = laporanCount

    // Set active status if any
    if activeStatus.Valid && activeStatus.String != "" {
        keluarga.StatusLaporanAktif = activeStatus.String
    }

    // Determine if keluarga can be edited
    // Can't edit if there are active reports
    canEdit := true
    if activeStatus.Valid && activeStatus.String != "" {
        canEdit = false
    }
    keluarga.CanEdit = canEdit

    return nil
}