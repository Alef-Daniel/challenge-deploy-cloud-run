package entity

import (
	"errors"
	"testing"
)

func TestValidateZipcode(t *testing.T) {
	tests := []struct {
		zipcode string
		want    error
	}{
		{"01001000", nil},
		{"0100100", ErrInvalidZipcode},
		{"010010000", ErrInvalidZipcode},
		{"0100100a", ErrInvalidZipcode},
		{"01001-000", ErrInvalidZipcode},
		{"", ErrInvalidZipcode},
	}

	for _, tt := range tests {
		if got := ValidateZipcode(tt.zipcode); !errors.Is(got, tt.want) {
			t.Errorf("ValidateZipcode(%q) = %v, want %v", tt.zipcode, got, tt.want)
		}
	}
}
