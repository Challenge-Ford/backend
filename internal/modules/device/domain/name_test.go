package devicedomain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceName_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   DeviceName
		wantErr error
	}{
		// Happy path
		{
			name:    "valid device name with single digit",
			input:   "TRQ-1",
			wantErr: nil,
		},
		{
			name:    "valid device name with multiple digits",
			input:   "TRQ-001",
			wantErr: nil,
		},
		{
			name:    "valid device name with many digits",
			input:   "TRQ-12345",
			wantErr: nil,
		},

		// Invalid formats
		{
			name:    "lowercase prefix",
			input:   "trq-001",
			wantErr: ErrInvalidDeviceName,
		},
		{
			name:    "missing dash",
			input:   "TRQ001",
			wantErr: ErrInvalidDeviceName,
		},
		{
			name:    "letters after dash",
			input:   "TRQ-ABC",
			wantErr: ErrInvalidDeviceName,
		},
		{
			name:    "no prefix at all",
			input:   "device-001",
			wantErr: ErrInvalidDeviceName,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: ErrInvalidDeviceName,
		},
		{
			name:    "extra dash at end",
			input:   "TRQ-001-",
			wantErr: ErrInvalidDeviceName,
		},
		{
			name:    "space in name",
			input:   "TRQ-00 1",
			wantErr: ErrInvalidDeviceName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.wantErr))
			}
		})
	}
}
