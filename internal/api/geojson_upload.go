package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/rifqidaiva/stunting-web/internal/object"
)

// GeoJsonUpload handles the upload of GeoJSON files.
// It processes the uploaded file and stores it in the database.
func GeoJsonUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // Limit to 10 MB
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse form: %v", err), http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("geojson")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get file from form: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if fileHeader.Filename == "" {
		http.Error(w, "File name is missing", http.StatusBadRequest)
		return
	}

	if fileHeader.Header.Get("Content-Type") != "application/json" &&
		fileHeader.Header.Get("Content-Type") != "application/geo+json" {
		http.Error(w, "Invalid file type. Only GeoJSON files are allowed.", http.StatusBadRequest)
		return
	}

	const maxFileSize = 10 << 20 // 10 MB
	if fileHeader.Size > maxFileSize {
		http.Error(w, "File size exceeds 10 MB limit", http.StatusBadRequest)
		return
	}

	// Process the GeoJSON file and store it in the database
	err = processGeoJsonFile(file, fileHeader.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process GeoJSON file: %v", err), http.StatusInternalServerError)
		return
	}

	response := object.NewResponse(http.StatusOK, "GeoJSON file uploaded successfully", nil)
	err = response.WriteJson(w)
	if err != nil {
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

	id := uuid.New().String() // Generate UUID in Go
	fmt.Printf("Uploading GeoJSON file: '%s' with generated ID: %s\n", filename, id)

	// Insert the whole GeoJSON as one row
	stmt, err := db.Prepare("INSERT INTO stunting_geojson (id, name, geojson) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	geojsonStr := string(data)
	_, err = stmt.Exec(id, filename, geojsonStr)
	if err != nil {
		return fmt.Errorf("failed to insert geojson: %v", err)
	}

	return nil
}
