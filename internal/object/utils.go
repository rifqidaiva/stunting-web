package object

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseCoordinates converts a string representation of coordinates into a [2]float64 array.
// It supports both "longitude, latitude" (e.g., "110.124, -7.123") and "POINT(110.124 -7.123)" formats.
func ParseCoordinates(coordinates string) [2]float64 {
	coordinates = strings.TrimSpace(coordinates)
	// Check for POINT format
	if strings.HasPrefix(coordinates, "POINT(") && strings.HasSuffix(coordinates, ")") {
		re := regexp.MustCompile(`POINT\(\s*([-\d.]+)\s+([-\d.]+)\s*\)`)
		matches := re.FindStringSubmatch(coordinates)
		if len(matches) == 3 {
			lon, err1 := strconv.ParseFloat(matches[1], 64)
			lat, err2 := strconv.ParseFloat(matches[2], 64)
			if err1 == nil && err2 == nil {
				return [2]float64{lon, lat}
			}
		}
		return [2]float64{0, 0}
	}

	// Default: "longitude, latitude"
	parts := strings.Split(coordinates, ",")
	if len(parts) != 2 {
		return [2]float64{0, 0}
	}

	lon, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	lat, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil {
		return [2]float64{0, 0}
	}

	return [2]float64{lon, lat}
}

// FormatCoordinates converts a [2]float64 array into a string representation of coordinates.
// The output format is "POINT(long lat)".
func FormatCoordinates(coordinates [2]float64) string {
	return "POINT(" +
		strconv.FormatFloat(coordinates[0], 'f', -1, 64) + " " +
		strconv.FormatFloat(coordinates[1], 'f', -1, 64) + ")"
}

// ToGeoJSON converts a slice of Sufferer objects into GeoJSON format.
func ToGeoJSON(sufferers []Sufferer) map[string]any {
	features := make([]map[string]any, len(sufferers))

	for i, sufferer := range sufferers {
		features[i] = map[string]any{
			"type": "Feature",
			"geometry": map[string]any{
				"type":        "Point",
				"coordinates": []float64{sufferer.Coordinates[0], sufferer.Coordinates[1]},
			},
			"properties": map[string]any{
				"id":             sufferer.Id,
				"name":           sufferer.Name,
				"nik":            sufferer.Nik,
				"date_of_birth":  sufferer.DateOfBirth,
				"status":         sufferer.Status,
				"reported_by_id": sufferer.ReportedById,
			},
		}
	}

	return map[string]any{
		"type":     "FeatureCollection",
		"features": features,
	}
}
