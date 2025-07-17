package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// GeoJsonGet retrieves GeoJSON data by ID.
// It responds with the GeoJSON data in JSON format.
func GeoJsonGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	db, err := object.ConnectDb()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	geojsonStr, err := getGeoJsonData(db, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Parse the GeoJSON string into a map
	var geojsonObj map[string]any
	if err := json.Unmarshal([]byte(geojsonStr), &geojsonObj); err != nil {
		http.Error(w, "Failed to parse GeoJSON", http.StatusInternalServerError)
		return
	}

	response := object.NewResponse(http.StatusOK, "GeoJSON data retrieved successfully", geojsonObj)
	err = response.WriteJson(w)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func getGeoJsonData(db *sql.DB, id string) (string, error) {
	var geojson string
	query := "SELECT geojson FROM stunting_geojson WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&geojson)
	if err != nil {
		return "", err
	}
	return geojson, nil
}
