package vehicledomain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVehicleID_NewVehicleID(t *testing.T) {
	id := NewVehicleID()
	assert.NotEqual(t, VehicleID{}, id)
}

func TestVehicleID_String(t *testing.T) {
	original := NewVehicleID()
	s := original.String()
	assert.NotEmpty(t, s)

	parsed, err := uuid.Parse(s)
	require.NoError(t, err)
	assert.Equal(t, uuid.UUID(original), parsed)
}

func TestVehicleID_Value(t *testing.T) {
	id := NewVehicleID()
	val, err := id.Value()

	require.NoError(t, err)
	assert.IsType(t, "", val)
	assert.Equal(t, id.String(), val)
}

func TestVehicleID_Scan(t *testing.T) {
	original := NewVehicleID()

	t.Run("parses from string", func(t *testing.T) {
		var scanned VehicleID
		err := scanned.Scan(original.String())
		require.NoError(t, err)
		assert.Equal(t, original, scanned)
	})

	t.Run("parses from []byte", func(t *testing.T) {
		var scanned VehicleID
		err := scanned.Scan([]byte(original.String()))
		require.NoError(t, err)
		assert.Equal(t, original, scanned)
	})

	t.Run("errors on invalid string", func(t *testing.T) {
		var scanned VehicleID
		err := scanned.Scan("not-a-uuid")
		assert.Error(t, err)
	})

	t.Run("errors on unsupported type", func(t *testing.T) {
		var scanned VehicleID
		err := scanned.Scan(int64(123))
		assert.Error(t, err)
	})
}

func TestVehicleID_RoundTrip(t *testing.T) {
	original := NewVehicleID()
	val, _ := original.Value()

	var scanned VehicleID
	err := scanned.Scan(val)
	require.NoError(t, err)
	assert.Equal(t, original, scanned)
}

func TestVehicle_Delete(t *testing.T) {
	userID := uuid.New()
	vehicle := &Vehicle{
		ID:    NewVehicleID(),
		VIN:   "1HGBH41JXMN109186",
		Plate: "ABC1234",
		Color: "#FF0000",
	}

	vehicle.Delete(userID)

	assert.True(t, vehicle.DeletedAt.Valid)
	assert.NotNil(t, vehicle.DeletedBy)
	assert.Equal(t, userID, *vehicle.DeletedBy)
}

func TestVehicle_Delete_OverwritesDeletedBy(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()
	vehicle := &Vehicle{
		ID:    NewVehicleID(),
		VIN:   "1HGBH41JXMN109186",
		Plate: "ABC1234",
		Color: "#FF0000",
	}

	vehicle.Delete(userID1)
	vehicle.Delete(userID2)

	assert.True(t, vehicle.DeletedAt.Valid)
	assert.Equal(t, userID2, *vehicle.DeletedBy)
}
