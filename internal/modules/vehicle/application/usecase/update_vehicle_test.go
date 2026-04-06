package vehicleusecase_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"torque/internal/core/apperr"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicleusecase "torque/internal/modules/vehicle/application/usecase"
	vehicledomain "torque/internal/modules/vehicle/domain"
	mockvehicle "torque/mocks/vehicle/domain"
)

func TestUpdateVehicle_Execute(t *testing.T) {
	vehicleID := vehicledomain.NewVehicleID()
	existing := &vehicledomain.Vehicle{
		ID:    vehicleID,
		VIN:   "1HGBH41JXMN109186",
		Plate: "ABC1234",
		Color: "#FF0000",
	}

	t.Run("updates plate and color successfully", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)
		repo.EXPECT().Save(ctx, mock.MatchedBy(func(v *vehicledomain.Vehicle) bool {
			return v.Plate == "XYZ9E72" && v.Color == "#00FF00"
		})).Return(nil)

		uc := vehicleusecase.NewUpdateVehicle(repo, modelRepo)
		out, err := uc.Execute(ctx, vehicleID, vehicledto.UpdateVehicleInput{
			Plate: "XYZ9E72",
			Color: "#00FF00",
		})

		require.NoError(t, err)
		assert.Equal(t, "XYZ9E72", out.Plate)
	})

	t.Run("updates model year successfully", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()
		newModelYearID := uuid.New()
		newModelYear := sampleModelYear()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)
		modelRepo.EXPECT().GetModelYearByID(ctx, vehicledomain.VehicleModelYearID(newModelYearID)).Return(newModelYear, nil)
		repo.EXPECT().Save(ctx, mock.MatchedBy(func(v *vehicledomain.Vehicle) bool {
			return v.ModelYearID == vehicledomain.VehicleModelYearID(newModelYearID)
		})).Return(nil)

		uc := vehicleusecase.NewUpdateVehicle(repo, modelRepo)
		out, err := uc.Execute(ctx, vehicleID, vehicledto.UpdateVehicleInput{ModelYearID: &newModelYearID})

		require.NoError(t, err)
		assert.NotNil(t, out)
	})

	t.Run("returns not found when vehicle does not exist", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(nil, nil)

		uc := vehicleusecase.NewUpdateVehicle(repo, modelRepo)
		_, err := uc.Execute(ctx, vehicleID, vehicledto.UpdateVehicleInput{Plate: "XYZ9E72"})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns not found when model year does not exist", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()
		newModelYearID := uuid.New()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)
		modelRepo.EXPECT().GetModelYearByID(ctx, vehicledomain.VehicleModelYearID(newModelYearID)).Return(nil, nil)

		uc := vehicleusecase.NewUpdateVehicle(repo, modelRepo)
		_, err := uc.Execute(ctx, vehicleID, vehicledto.UpdateVehicleInput{ModelYearID: &newModelYearID})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns bad request for invalid plate", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)

		uc := vehicleusecase.NewUpdateVehicle(repo, modelRepo)
		_, err := uc.Execute(ctx, vehicleID, vehicledto.UpdateVehicleInput{Plate: "INVALID"})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindBadRequest, appErr.Kind)
	})

	t.Run("returns bad request for invalid color", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)

		uc := vehicleusecase.NewUpdateVehicle(repo, modelRepo)
		_, err := uc.Execute(ctx, vehicleID, vehicledto.UpdateVehicleInput{Color: "not-a-color"})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindBadRequest, appErr.Kind)
	})

	t.Run("returns internal error when GetByID fails", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(nil, assert.AnError)

		uc := vehicleusecase.NewUpdateVehicle(repo, modelRepo)
		_, err := uc.Execute(ctx, vehicleID, vehicledto.UpdateVehicleInput{Plate: "XYZ9E72"})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})

	t.Run("returns internal error when Save fails", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		repo.EXPECT().GetByID(ctx, vehicleID).Return(existing, nil)
		repo.EXPECT().Save(ctx, mock.Anything).Return(assert.AnError)

		uc := vehicleusecase.NewUpdateVehicle(repo, modelRepo)
		_, err := uc.Execute(ctx, vehicleID, vehicledto.UpdateVehicleInput{Plate: "XYZ9E72"})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindInternal, appErr.Kind)
	})
}
