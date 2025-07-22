package api

import (
	"encoding/json"
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// AdminInsert handles the insertion of new sufferers into the database.
//
// MARK: TODO
//   - Implement JWT authentication and authorization for admin routes.
func AdminInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var sufferer object.Sufferer
	err := json.NewDecoder(r.Body).Decode(&sufferer)
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = sufferer.ValidateFields("Name", "Nik", "DateOfBirth", "Coordinates")
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
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

	// MARK: TODO
	// - Implement JWT authentication and authorization for admin routes.
	// - Get admin ID and insert to reported_by_id.
	// For now, we assume the request is authenticated and authorized.
	stmt, err := db.Prepare("INSERT INTO sufferer (id, name, nik, date_of_birth, coordinates, reported_by_id) VALUES (UUID(), ?, ?, ?, ST_GeomFromText(?), ?)")
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(sufferer.Name, sufferer.Nik, sufferer.DateOfBirth, object.FormatCoordinates(sufferer.Coordinates), "b0429037-66ad-11f0-a701-2811a8eb1247")
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Sufferer added successfully", nil)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
