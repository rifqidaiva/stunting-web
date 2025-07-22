package api

import (
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// AdminGetGeoJson handles the retrieval of GeoJSON data for all sufferers from the database.
//
// MARK: TODO
//   - Implement JWT authentication and authorization for admin routes.
func AdminGetGeoJson(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	db, err := object.ConnectDb()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, nik, date_of_birth, ST_AsText(coordinates), status FROM sufferer")
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer rows.Close()

	var sufferers []object.Sufferer
	for rows.Next() {
		var sufferer object.Sufferer
		var coordinates string
		if err := rows.Scan(&sufferer.Id, &sufferer.Name, &sufferer.Nik, &sufferer.DateOfBirth, &coordinates, &sufferer.Status); err != nil {
			response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		sufferer.Coordinates = object.ParseCoordinates(coordinates)
		sufferers = append(sufferers, sufferer)
	}

	// Convert sufferers to GeoJSON format
	geoJSON := object.ToGeoJSON(sufferers)
	response := object.NewResponse(http.StatusOK, "Success", geoJSON)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}