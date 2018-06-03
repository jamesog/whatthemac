// Package mac returns information about IEEE assigned MAC addresses.
package mac

import (
	"encoding/csv"
	"io"
)

// Registration holds registration information for each MAC address prefix.
type Registration struct {
	List         string // MA-L, MA-M, MA-S
	OUI          string
	Organization string
}

// Assignments is a map of MAC address prefix assignments to Registrations.
type Assignments map[string]Registration

// ReadAssignments reads the given CSV file and returns Assignments.
func ReadAssignments(file io.Reader) (Assignments, error) {
	a := make(Assignments)

	r := csv.NewReader(file)
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch {
		// Skip the header line
		case rec[0] == "Registry":
			continue
		// Skip sub-allocations contained in another list
		case rec[2] == "IEEE Registration Authority":
			continue
		}

		a[rec[1]] = Registration{
			List:         rec[0],
			OUI:          rec[1],
			Organization: rec[2],
		}
	}

	return a, nil
}
