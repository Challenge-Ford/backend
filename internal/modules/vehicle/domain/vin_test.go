package vehicledomain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVIN_Validate(t *testing.T) {
	tests := []struct {
		name    string
		vin     VIN
		wantErr error
	}{
		// Happy path
		{
			name:    "valid VIN with letters and numbers",
			vin:     "1HGBH41JXMN109186",
			wantErr: nil,
		},
		{
			name:    "valid VIN all uppercase",
			vin:     "WVWZZZ3CZWE123456",
			wantErr: nil,
		},

		// Forbidden characters I, O, Q
		{
			name:    "contains letter I",
			vin:     "1HGBH41JXMN10918I",
			wantErr: ErrInvalidVIN,
		},
		{
			name:    "contains letter O",
			vin:     "1HGBH41JXMN109O86",
			wantErr: ErrInvalidVIN,
		},
		{
			name:    "contains letter Q",
			vin:     "1HGBH41JXQN109186",
			wantErr: ErrInvalidVIN,
		},

		// Length
		{
			name:    "too short (16 chars)",
			vin:     "1HGBH41JXMN10918",
			wantErr: ErrInvalidVIN,
		},
		{
			name:    "too long (18 chars)",
			vin:     "1HGBH41JXMN1091860",
			wantErr: ErrInvalidVIN,
		},

		// Edge cases
		{
			name:    "empty string",
			vin:     "",
			wantErr: ErrInvalidVIN,
		},
		{
			name:    "lowercase letters",
			vin:     "1hgbh41jxmn109186",
			wantErr: ErrInvalidVIN,
		},
		{
			name:    "contains special characters",
			vin:     "1HGBH41JX-N109186",
			wantErr: ErrInvalidVIN,
		},
		{
			name:    "contains spaces",
			vin:     "1HGBH41JX MN109186",
			wantErr: ErrInvalidVIN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.vin.Validate()

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.wantErr))
			}
		})
	}
}
