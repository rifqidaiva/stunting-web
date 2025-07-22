package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// AdminUpdate handles the update of a sufferer in the database.
// It allows updating multiple fields dynamically based on the request body.
//
// MARK: TODO
//   - Implement JWT authentication and authorization for admin routes.
func AdminUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var req object.Sufferer
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err := req.ValidateFields("Id")
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	nonEmptyFields := req.NonEmptyFields()
	if len(nonEmptyFields) == 0 {
		response := object.NewResponse(http.StatusBadRequest, "No updatable fields provided", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := req.ValidateFields(nonEmptyFields...); err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Build the SQL query dynamically based on non-empty fields
	query := "UPDATE sufferer SET "
	var setClauses []string
	var args []any

	for _, field := range nonEmptyFields {
		switch field {
		case "Name":
			setClauses = append(setClauses, "name = ?")
			args = append(args, req.Name)
		case "Nik":
			setClauses = append(setClauses, "nik = ?")
			args = append(args, req.Nik)
		case "DateOfBirth":
			setClauses = append(setClauses, "date_of_birth = ?")
			args = append(args, req.DateOfBirth)
		case "Coordinates":
			wkt := object.FormatCoordinates(req.Coordinates)
			setClauses = append(setClauses, "coordinates = ST_GeomFromText(?)")
			args = append(args, wkt)
		case "Status":
			setClauses = append(setClauses, "status = ?")
			args = append(args, req.Status)
		case "Id":
			continue
		default:
			// Unknown field
			response := object.NewResponse(http.StatusBadRequest, "Unknown field: "+field, nil)
			if err := response.WriteJson(w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
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

	query += strings.Join(setClauses, ", ") + " WHERE id = ?"
	args = append(args, req.Id)
	stmt, err := db.Prepare(query)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(args...)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if rowsAffected == 0 {
		response := object.NewResponse(http.StatusNotFound, "Sufferer ID not found or no changes made", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Sufferer updated successfully", nil)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
