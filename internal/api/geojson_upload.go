package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/rifqidaiva/stunting-web/internal/object"
)

// GeoJsonUpload handles the upload of GeoJSON files.
// It expects a POST request with a file upload containing GeoJSON data.
// The uploaded file is processed and stored in the database.
func GeoJsonUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := object.NewResponse(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // Limit to 10 MB
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	file, fileHeader, err := r.FormFile("geojson")
	if err != nil {
		response := object.NewResponse(http.StatusBadRequest, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	if fileHeader.Filename == "" {
		response := object.NewResponse(http.StatusBadRequest, "File name is missing", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if fileHeader.Header.Get("Content-Type") != "application/json" &&
		fileHeader.Header.Get("Content-Type") != "application/geo+json" {
		response := object.NewResponse(http.StatusBadRequest, "Invalid file type. Only GeoJSON files are allowed.", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	const maxFileSize = 10 << 20 // 10 MB
	if fileHeader.Size > maxFileSize {
		response := object.NewResponse(http.StatusBadRequest, "File size exceeds 10 MB limit", nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Process the GeoJSON file and store it in the database
	err = processGeoJsonFile(file, fileHeader.Filename)
	if err != nil {
		response := object.NewResponse(http.StatusInternalServerError, err.Error(), nil)
		if err := response.WriteJson(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := object.NewResponse(http.StatusOK, "GeoJSON file uploaded successfully", nil)
	if err := response.WriteJson(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func processGeoJsonFile(file multipart.File, filename string) error {
	db, err := object.ConnectDb()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Read the uploaded file into memory
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	id := uuid.New().String()
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	fmt.Printf("Processing GeoJSON file: %s (ID: %s)\n", name, id)

	// Insert the whole GeoJSON as one row
	stmt, err := db.Prepare("INSERT INTO stunting_geojson (id, status, geojson) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	geojsonStr := string(data)
	_, err = stmt.Exec(id, name, geojsonStr)
	if err != nil {
		return fmt.Errorf("failed to insert geojson: %v", err)
	}

	return nil
}
