package api

import (
	"net/http"

	"github.com/rifqidaiva/stunting-web/internal/object"
)

// AdminDelete handles the deletion of a sufferer from the database.
//
// MARK: TODO
//   - Implement JWT authentication and authorization for admin routes.
func AdminDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Get the sufferer ID from the URL parameters
	id := r.URL.Query().Get("id")
	if id == "" {
		response := object.NewResponse(http.StatusBadRequest, "Missing sufferer ID", nil)
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

	stmt, err := db.Prepare("DELETE FROM sufferer WHERE id = ?")
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
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
		response := object.NewResponse(http.StatusNotFound, "Sufferer ID not found", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "Sufferer deleted successfully", nil)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
