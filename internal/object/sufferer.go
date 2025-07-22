package object

import "fmt"

// Sufferer represents a person suffering from stunting.
type Sufferer struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Nik          string     `json:"nik"`
	DateOfBirth  string     `json:"date_of_birth"`
	Coordinates  [2]float64 `json:"coordinates"` // [longitude, latitude]
	Status       string     `json:"status"`
	ReportedById string     `json:"reported_by_id"`
}

// ValidateFields checks if the required fields of the Sufferer are set.
//   s.ValidateFields("Id", "Name", "Nik", "DateOfBirth", "Coordinates", "Status", "ReportedById")
// Returns an error if any required field is missing or invalid.
func (s *Sufferer) ValidateFields(fields ...string) error {
	for _, field := range fields {
		switch field {
		case "Id":
			if s.Id == "" {
				return fmt.Errorf("ID is required")
			}
		case "Name":
			if s.Name == "" {
				return fmt.Errorf("name is required")
			}
		case "Nik":
			if s.Nik == "" {
				return fmt.Errorf("NIK is required")
			}
		case "DateOfBirth":
			if s.DateOfBirth == "" {
				return fmt.Errorf("date of birth is required")
			}
		case "Coordinates":
			if len(s.Coordinates) != 2 {
				return fmt.Errorf("coordinates must be a valid [longitude, latitude] pair")
			}
			if s.Coordinates[0] == 0 && s.Coordinates[1] == 0 {
				return fmt.Errorf("coordinates cannot be [0, 0]")
			}
		case "Status":
			if s.Status == "" {
				return fmt.Errorf("status is required")
			}
		case "ReportedById":
			if s.ReportedById == "" {
				return fmt.Errorf("reported by ID is required")
			}
		default:
			return fmt.Errorf("unknown field: %s", field)
		}
	}
	return nil
}

// EmptyFields returns a list of fields that are empty in the Sufferer struct.
// It checks all fields and returns a slice of strings with the names of the empty fields.
func (s *Sufferer) EmptyFields() []string {
	var emptyFields []string
	if s.Id == "" {
		emptyFields = append(emptyFields, "Id")
	}
	if s.Name == "" {
		emptyFields = append(emptyFields, "Name")
	}
	if s.Nik == "" {
		emptyFields = append(emptyFields, "Nik")
	}
	if s.DateOfBirth == "" {
		emptyFields = append(emptyFields, "DateOfBirth")
	}
	if len(s.Coordinates) != 2 || (s.Coordinates[0] == 0 && s.Coordinates[1] == 0) {
		emptyFields = append(emptyFields, "Coordinates")
	}
	if s.Status == "" {
		emptyFields = append(emptyFields, "Status")
	}
	if s.ReportedById == "" {
		emptyFields = append(emptyFields, "ReportedById")
	}
	return emptyFields
}

// NonEmptyFields returns a list of fields that are not empty in the Sufferer struct.
// It checks all fields and returns a slice of strings with the names of the non-empty fields
func (s *Sufferer) NonEmptyFields() []string {
	var nonEmptyFields []string
	if s.Id != "" {
		nonEmptyFields = append(nonEmptyFields, "Id")
	}
	if s.Name != "" {
		nonEmptyFields = append(nonEmptyFields, "Name")
	}
	if s.Nik != "" {
		nonEmptyFields = append(nonEmptyFields, "Nik")
	}
	if s.DateOfBirth != "" {
		nonEmptyFields = append(nonEmptyFields, "DateOfBirth")
	}
	if len(s.Coordinates) == 2 && (s.Coordinates[0] != 0 || s.Coordinates[1] != 0) {
		nonEmptyFields = append(nonEmptyFields, "Coordinates")
	}
	if s.Status != "" {
		nonEmptyFields = append(nonEmptyFields, "Status")
	}
	if s.ReportedById != "" {
		nonEmptyFields = append(nonEmptyFields, "ReportedById")
	}
	return nonEmptyFields
}
