package vehicledomain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlate_Validate(t *testing.T) {
	tests := []struct {
		name    string
		plate   Plate
		wantErr error
	}{
		// Happy path — old format (ABC-1234 without dash)
		{
			name:    "valid old format plate",
			plate:   "ABC1234",
			wantErr: nil,
		},
		{
			name:    "valid old format plate with Z",
			plate:   "ZZZ9999",
			wantErr: nil,
		},

		// Happy path — Mercosul format (ABC1D23)
		{
			name:    "valid Mercosul plate",
			plate:   "ABC1D23",
			wantErr: nil,
		},
		{
			name:    "valid Mercosul plate with different digit position",
			plate:   "AAA0A00",
			wantErr: nil,
		},

		// Invalid formats
		{
			name:    "lowercase old format",
			plate:   "abc1234",
			wantErr: ErrInvalidPlate,
		},
		{
			name:    "lowercase Mercosul format",
			plate:   "abc1d23",
			wantErr: ErrInvalidPlate,
		},
		{
			name:    "with dash (not accepted by regex)",
			plate:   "ABC-1234",
			wantErr: ErrInvalidPlate,
		},
		{
			name:    "too short",
			plate:   "ABC123",
			wantErr: ErrInvalidPlate,
		},
		{
			name:    "too long",
			plate:   "ABC12345",
			wantErr: ErrInvalidPlate,
		},
		{
			name:    "all numbers",
			plate:   "1234567",
			wantErr: ErrInvalidPlate,
		},
		{
			name:    "all letters",
			plate:   "ABCDEFG",
			wantErr: ErrInvalidPlate,
		},
		{
			name:    "empty string",
			plate:   "",
			wantErr: ErrInvalidPlate,
		},
		{
			name:    "special characters",
			plate:   "ABC@123",
			wantErr: ErrInvalidPlate,
		},
		{
			name:    "Mercosul with digit in wrong position",
			plate:   "AB1CD23",
			wantErr: ErrInvalidPlate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.plate.Validate()

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.wantErr))
			}
		})
	}
}
