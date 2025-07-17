package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// GeoJsonGet handles the retrieval of GeoJSON data.
// It expects a GET request with an optional ID query parameter.
// If ID is provided, it retrieves the specific GeoJSON data.
// If no ID is provided, it retrieves all GeoJSON data.
func GeoJsonGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	id := r.URL.Query().Get("id")

	db, err := object.ConnectDb()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer db.Close()

	if id != "" {
		geojsonObj, err := getGeoJsonData(db, id)
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

		response := object.NewResponse(http.StatusOK, "GeoJSON data retrieved successfully", geojsonObj)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get all
	geojsonList, err := getAllGeoJsonData(db)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "All GeoJSON data retrieved successfully", geojsonList)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getGeoJsonData(db *sql.DB, id string) (map[string]any, error) {
	var name, geojsonStr string
	query := "SELECT name, geojson FROM stunting_geojson WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&name, &geojsonStr)
	if err != nil {
		return nil, err
	}
	var geojsonObj map[string]any
	if err := json.Unmarshal([]byte(geojsonStr), &geojsonObj); err != nil {
		return nil, err
	}
	geojsonObj["id"] = id
	geojsonObj["name"] = name
	return geojsonObj, nil
}

func getAllGeoJsonData(db *sql.DB) ([]map[string]any, error) {
	rows, err := db.Query("SELECT id, name, geojson FROM stunting_geojson")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]any
	for rows.Next() {
		var id, name, geojsonStr string
		if err := rows.Scan(&id, &name, &geojsonStr); err != nil {
			return nil, err
		}
		var geojsonObj map[string]any
		if err := json.Unmarshal([]byte(geojsonStr), &geojsonObj); err != nil {
			return nil, err
		}

		geojsonObj["id"] = id
		geojsonObj["name"] = name
		result = append(result, geojsonObj)
	}
	return result, nil
}
