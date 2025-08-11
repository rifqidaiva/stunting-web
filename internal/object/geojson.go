package object

import (
	"fmt"
	"strconv"
	"strings"
)

// GeoJSONFeature represents a single GeoJSON feature
type GeoJSONFeature struct {
	Type       string                 `json:"type"`
	Geometry   GeoJSONGeometry        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// GeoJSONGeometry represents the geometry part of a GeoJSON feature
type GeoJSONGeometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

// GeoJSONFeatureCollection represents a collection of GeoJSON features
type GeoJSONFeatureCollection struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}

// ToWKT converts a [2]float64 array into a Well-Known Text (WKT) representation of coordinates.
// The output format is "POINT(long lat)".
func ToWKT(coordinates [2]float64) string {
	return "POINT(" +
		strconv.FormatFloat(coordinates[0], 'f', -1, 64) + " " +
		strconv.FormatFloat(coordinates[1], 'f', -1, 64) + ")"
}

// ParseWKT converts a WKT representation of coordinates into a [2]float64 array.
// It supports the "POINT(long lat)" format.
func ParseWKT(wkt string) [2]float64 {
	if len(wkt) < 7 || wkt[:6] != "POINT(" || wkt[len(wkt)-1] != ')' {
		return [2]float64{0, 0}
	}

	// Remove "POINT(" prefix and ")" suffix
	coords := wkt[6 : len(wkt)-1]
	parts := strings.Split(coords, " ")
	if len(parts) != 2 {
		return [2]float64{0, 0}
	}

	lon, err1 := strconv.ParseFloat(parts[0], 64)
	lat, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		return [2]float64{0, 0}
	}

	return [2]float64{lon, lat}
}

// WKTToGeoJSON converts WKT geometry to GeoJSON geometry
func WKTToGeoJSON(wkt string) (GeoJSONGeometry, error) {
	var geometry GeoJSONGeometry

	wkt = strings.TrimSpace(wkt)

	if strings.HasPrefix(wkt, "POINT") {
		return parseWKTPoint(wkt)
	} else if strings.HasPrefix(wkt, "MULTIPOLYGON") {
		return parseWKTMultiPolygon(wkt)
	} else if strings.HasPrefix(wkt, "POLYGON") {
		return parseWKTPolygon(wkt)
	}

	return geometry, fmt.Errorf("unsupported WKT geometry type")
}

// parseWKTPoint converts WKT POINT to GeoJSON Point
func parseWKTPoint(wkt string) (GeoJSONGeometry, error) {
	coords := ParseWKT(wkt)

	return GeoJSONGeometry{
		Type:        "Point",
		Coordinates: []float64{coords[0], coords[1]},
	}, nil
}

// parseWKTPolygon converts WKT POLYGON to GeoJSON Polygon
func parseWKTPolygon(wkt string) (GeoJSONGeometry, error) {
	// Remove "POLYGON(" prefix and ")" suffix
	if !strings.HasPrefix(wkt, "POLYGON((") || !strings.HasSuffix(wkt, "))") {
		return GeoJSONGeometry{}, fmt.Errorf("invalid POLYGON format")
	}

	coordStr := wkt[9 : len(wkt)-2] // Remove "POLYGON((" and "))"
	rings := parsePolygonRings(coordStr)

	return GeoJSONGeometry{
		Type:        "Polygon",
		Coordinates: rings,
	}, nil
}

// parseWKTMultiPolygon converts WKT MULTIPOLYGON to GeoJSON MultiPolygon
func parseWKTMultiPolygon(wkt string) (GeoJSONGeometry, error) {
	// Remove "MULTIPOLYGON(" prefix and ")" suffix
	if !strings.HasPrefix(wkt, "MULTIPOLYGON(") || !strings.HasSuffix(wkt, ")") {
		return GeoJSONGeometry{}, fmt.Errorf("invalid MULTIPOLYGON format")
	}

	coordStr := wkt[13 : len(wkt)-1] // Remove "MULTIPOLYGON(" and ")"
	polygons := parseMultiPolygonCoordinates(coordStr)

	return GeoJSONGeometry{
		Type:        "MultiPolygon",
		Coordinates: polygons,
	}, nil
}

// parseMultiPolygonCoordinates parses the coordinate string for MultiPolygon
func parseMultiPolygonCoordinates(coordStr string) [][][][2]float64 {
	var polygons [][][][2]float64

	// Split by ")),((" to separate individual polygons
	polygonStrs := strings.Split(coordStr, ")),(")

	for i, polygonStr := range polygonStrs {
		// Clean up the polygon string
		if i == 0 {
			polygonStr = strings.TrimPrefix(polygonStr, "((")
		}
		if i == len(polygonStrs)-1 {
			polygonStr = strings.TrimSuffix(polygonStr, "))")
		}
		if i > 0 && i < len(polygonStrs)-1 {
			// Middle polygons, no trimming needed
		}

		rings := parsePolygonRings(polygonStr)
		polygons = append(polygons, rings)
	}

	return polygons
}

// parsePolygonRings parses the coordinate string for Polygon rings
func parsePolygonRings(coordStr string) [][][2]float64 {
	var rings [][][2]float64

	// Split by "),(" to separate rings
	ringStrs := strings.Split(coordStr, "),(")

	for i, ringStr := range ringStrs {
		// Clean up the ring string
		if i == 0 {
			ringStr = strings.TrimPrefix(ringStr, "(")
		}
		if i == len(ringStrs)-1 {
			ringStr = strings.TrimSuffix(ringStr, ")")
		}

		coordinates := parseCoordinateString(ringStr)
		rings = append(rings, coordinates)
	}

	return rings
}

// parseCoordinateString parses a string of coordinates into [][2]float64
func parseCoordinateString(coordStr string) [][2]float64 {
	var coordinates [][2]float64

	// Split by comma to get individual coordinate pairs
	pairs := strings.Split(coordStr, ",")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		coords := strings.Fields(pair)

		if len(coords) >= 2 {
			lon, err1 := strconv.ParseFloat(coords[0], 64)
			lat, err2 := strconv.ParseFloat(coords[1], 64)

			if err1 == nil && err2 == nil {
				coordinates = append(coordinates, [2]float64{lon, lat})
			}
		}
	}

	return coordinates
}

// CreateGeoJSONFeature creates a GeoJSON feature with geometry and properties
func CreateGeoJSONFeature(wkt string, properties map[string]interface{}) (GeoJSONFeature, error) {
	geometry, err := WKTToGeoJSON(wkt)
	if err != nil {
		return GeoJSONFeature{}, err
	}

	return GeoJSONFeature{
		Type:       "Feature",
		Geometry:   geometry,
		Properties: properties,
	}, nil
}

// CreateGeoJSONFeatureCollection creates a GeoJSON feature collection
func CreateGeoJSONFeatureCollection(features []GeoJSONFeature) GeoJSONFeatureCollection {
	return GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}
}
