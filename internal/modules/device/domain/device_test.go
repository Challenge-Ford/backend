package devicedomain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceID_NewDeviceID(t *testing.T) {
	id := NewDeviceID()
	assert.NotEqual(t, DeviceID{}, id)
}

func TestDeviceID_String(t *testing.T) {
	original := NewDeviceID()
	s := original.String()
	assert.NotEmpty(t, s)

	parsed, err := uuid.Parse(s)
	require.NoError(t, err)
	assert.Equal(t, uuid.UUID(original), parsed)
}

func TestDeviceID_Value(t *testing.T) {
	id := NewDeviceID()
	val, err := id.Value()

	require.NoError(t, err)
	assert.IsType(t, "", val)
	assert.Equal(t, id.String(), val)
}

func TestDeviceID_Scan(t *testing.T) {
	original := NewDeviceID()

	t.Run("parses from string", func(t *testing.T) {
		var scanned DeviceID
		err := scanned.Scan(original.String())
		require.NoError(t, err)
		assert.Equal(t, original, scanned)
	})

	t.Run("parses from []byte", func(t *testing.T) {
		var scanned DeviceID
		err := scanned.Scan([]byte(original.String()))
		require.NoError(t, err)
		assert.Equal(t, original, scanned)
	})

	t.Run("errors on invalid string", func(t *testing.T) {
		var scanned DeviceID
		err := scanned.Scan("not-a-uuid")
		assert.Error(t, err)
	})

	t.Run("errors on unsupported type", func(t *testing.T) {
		var scanned DeviceID
		err := scanned.Scan(int64(123))
		assert.Error(t, err)
	})
}

func TestDeviceID_RoundTrip(t *testing.T) {
	original := NewDeviceID()
	val, _ := original.Value()

	var scanned DeviceID
	err := scanned.Scan(val)
	require.NoError(t, err)
	assert.Equal(t, original, scanned)
}

func TestDevice_Delete(t *testing.T) {
	userID := uuid.New()
	device := &Device{
		ID:            NewDeviceID(),
		Name:          "TRQ-001",
		CertificateCN: "TRQ-001",
		CertificateSN: "sn-123",
	}

	device.Delete(userID)

	assert.True(t, device.DeletedAt.Valid)
	assert.NotNil(t, device.DeletedBy)
	assert.Equal(t, userID, *device.DeletedBy)
}

func TestDevice_Delete_OverwritesDeletedBy(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()
	device := &Device{
		ID:            NewDeviceID(),
		Name:          "TRQ-001",
		CertificateCN: "TRQ-001",
		CertificateSN: "sn-123",
	}

	device.Delete(userID1)
	device.Delete(userID2)

	assert.True(t, device.DeletedAt.Valid)
	assert.Equal(t, userID2, *device.DeletedBy)
}
