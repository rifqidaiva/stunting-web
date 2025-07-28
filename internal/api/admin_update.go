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
//
// AdminUpdate swagger documentation
//
// @Summary Update a sufferer
// @Description	Updates a sufferer by ID. Only the fields you want to update need to be included in the request body (partial update).
// @Description	The 'id' field is required to identify the sufferer to update.
// @Description	The 'reported_by_id' field must NOT be included in the request body.
// @Description	Other fields (such as name, nik, date_of_birth, coordinates, status) can be included as needed.
// @Tags Admin
// @Accept json
// @Produce json
// @Param sufferer body object.Sufferer true "Sufferer data to update (only fields to be updated, id and reported_by_id are required)"
// @Success 200 {object} object.Response{data=nil} "Sufferer updated successfully"
// @Failure 400 {object} object.Response{data=nil} "Invalid request body or missing ID"
// @Failure 404 {object} object.Response{data=nil} "Sufferer ID not found"
// @Failure 500 {object} object.Response{data=nil} "Internal server error"
// @Router /api/admin/update [put]
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
