package vehicledomain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColor_Validate(t *testing.T) {
	tests := []struct {
		name    string
		color   Color
		wantErr error
	}{
		// Happy path — 6-digit hex
		{
			name:    "valid 6-digit uppercase hex",
			color:   "#FF0000",
			wantErr: nil,
		},
		{
			name:    "valid 6-digit lowercase hex",
			color:   "#ff0000",
			wantErr: nil,
		},
		{
			name:    "valid 6-digit mixed case hex",
			color:   "#Ff00Aa",
			wantErr: nil,
		},

		// Happy path — 3-digit hex (shorthand)
		{
			name:    "valid 3-digit uppercase hex",
			color:   "#F00",
			wantErr: nil,
		},
		{
			name:    "valid 3-digit lowercase hex",
			color:   "#f00",
			wantErr: nil,
		},

		// Invalid formats
		{
			name:    "missing hash prefix",
			color:   "FF0000",
			wantErr: ErrInvalidColor,
		},
		{
			name:    "invalid hex characters",
			color:   "#GGGGGG",
			wantErr: ErrInvalidColor,
		},
		{
			name:    "too short (2 digits)",
			color:   "#FF",
			wantErr: ErrInvalidColor,
		},
		{
			name:    "too long (8 digits)",
			color:   "#FF0000FF",
			wantErr: ErrInvalidColor,
		},
		{
			name:    "5 digits (not 3 or 6)",
			color:   "#FF000",
			wantErr: ErrInvalidColor,
		},
		{
			name:    "empty string",
			color:   "",
			wantErr: ErrInvalidColor,
		},
		{
			name:    "just hash",
			color:   "#",
			wantErr: ErrInvalidColor,
		},
		{
			name:    "css named color (not supported)",
			color:   "red",
			wantErr: ErrInvalidColor,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.color.Validate()

			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.True(t, errors.Is(err, tt.wantErr))
			}
		})
	}
}
