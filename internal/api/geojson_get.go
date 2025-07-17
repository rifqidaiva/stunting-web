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
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		response := object.NewResponse(http.StatusBadRequest, "Missing id parameter", nil)
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

	geojsonStr, err := getGeoJsonData(db, id)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "Failed to retrieve GeoJSON data"
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
			msg = "GeoJSON not found"
		}
		response := object.NewResponse(status, msg, nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Parse the GeoJSON string into a map
	var geojsonObj map[string]any
	if err := json.Unmarshal([]byte(geojsonStr), &geojsonObj); err != nil {
		response := object.NewResponse(http.StatusInternalServerError, "Failed to parse GeoJSON", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "GeoJSON data retrieved successfully", geojsonObj)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
