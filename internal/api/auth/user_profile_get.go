package auth

import (
	"database/sql"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

type userProfileResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Data  any    `json:"data,omitempty"`
}

type masyarakatData struct {
	Id     string `json:"id"`
	Nama   string `json:"nama"`
	Alamat string `json:"alamat"`
}

type petugasKesehatanData struct {
	Id      string `json:"id"`
	IdSkpd  string `json:"id_skpd"`
	Nama    string `json:"nama"`
	Created string `json:"created_date"`
}

// # UserProfileGet handles getting user profile based on token
//
// @Summary Get user profile
// @Description Get user profile data based on JWT token and role
// @Description
// @Description Response data varies by role:
// @Description - masyarakat: {id, nama, alamat}
// @Description - petugas kesehatan: {id, id_skpd, nama, created_date}
// @Description - admin: {nama: "Administrator"}
// @Tags auth
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} object.Response{data=userProfileResponse{data=nil}} "User profile retrieved successfully"
// @Failure 401 {object} object.Response{data=nil} "Unauthorized"
// @Failure 403 {object} object.Response{data=nil} "Forbidden"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/auth/profile [get]
func UserProfileGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Extract token from Authorization header using utility function
	authHeader := r.Header.Get("Authorization")
	token, err := object.GetJWTFromHeader(authHeader)
	if err != nil {
		response := object.NewResponse(http.StatusUnauthorized, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Validate and parse the JWT token
	userId, role, err := object.ParseJWT(token)
	if err != nil {
		response := object.NewResponse(http.StatusUnauthorized, "Invalid or expired token", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	db, err := object.ConnectDb()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Database connection error", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	// Get basic user info
	var user object.Pengguna
	query := "SELECT id, email, role FROM pengguna WHERE id = ?"
	err = db.QueryRow(query, userId).Scan(&user.Id, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			response := object.NewResponse(http.StatusUnauthorized, "User not found", nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		response := object.NewResponse(http.StatusInternalServerError, "Database query error", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	userProfile := userProfileResponse{
		Id:    user.Id,
		Email: user.Email,
		Role:  user.Role,
	}

	// Get role-specific data
	switch role {
	case "masyarakat":
		var masyarakat masyarakatData
		masyarakatQuery := "SELECT id, nama, alamat FROM masyarakat WHERE id_pengguna = ?"
		err = db.QueryRow(masyarakatQuery, userId).Scan(&masyarakat.Id, &masyarakat.Nama, &masyarakat.Alamat)
		if err != nil {
			if err != sql.ErrNoRows {
				response := object.NewResponse(http.StatusInternalServerError, "Failed to get masyarakat data", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		} else {
			userProfile.Data = masyarakat
		}

	case "petugas kesehatan":
		var petugas petugasKesehatanData
		petugasQuery := "SELECT id, id_skpd, nama, created_date FROM petugas_kesehatan WHERE id_pengguna = ?"
		err = db.QueryRow(petugasQuery, userId).Scan(&petugas.Id, &petugas.IdSkpd, &petugas.Nama, &petugas.Created)
		if err != nil {
			if err != sql.ErrNoRows {
				response := object.NewResponse(http.StatusInternalServerError, "Failed to get petugas kesehatan data", nil)
				if err := response.WriteJson(w); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
		} else {
			userProfile.Data = petugas
		}

	case "admin":
		// Admin doesn't have additional data table, just basic info
		userProfile.Data = map[string]string{
			"nama": "Administrator",
		}

	default:
		response := object.NewResponse(http.StatusForbidden, "Invalid user role", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "User profile retrieved successfully", userProfile)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
