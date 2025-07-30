package object

import (
	"strconv"
	"strings"
)

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
