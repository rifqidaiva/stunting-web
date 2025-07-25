package object

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
}

// NewResponse creates a new Response object.
// Pass nil for data if you want to omit it.
func NewResponse(statusCode int, message string, data any) *Response {
	return &Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

// WriteJson writes the Response as JSON to the http.ResponseWriter
func (r *Response) WriteJson(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)

	return json.NewEncoder(w).Encode(r)
}

type GeoJson struct {
	Type     string       `json:"type"`
	Features []GeoFeature `json:"features"`
}

type GeoFeature struct {
	Type       string        `json:"type"`
	Properties GeoProperties `json:"properties"`
	Geometry   GeoGeometry   `json:"geometry"`
}

type GeoProperties struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Nik          string `json:"nik"`
	DateOfBirth  string `json:"date_of_birth"`
	Status       string `json:"status"`
	ReportedById int    `json:"reported_by_id"`
}
type GeoGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
