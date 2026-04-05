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

var validInput = vehicledto.CreateVehicleInput{
	ModelYearID: uuid.New(),
	VIN:         "1HGBH41JXMN109186",
	Plate:       "ABC1234",
	Color:       "#FF0000",
}

func TestCreateVehicle_Execute(t *testing.T) {
	modelYear := sampleModelYear()

	t.Run("creates vehicle successfully", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		modelRepo.EXPECT().GetModelYearByID(ctx, vehicledomain.VehicleModelYearID(validInput.ModelYearID)).Return(modelYear, nil)
		repo.EXPECT().GetByVIN(ctx, vehicledomain.VIN(validInput.VIN)).Return(nil, nil)
		repo.EXPECT().GetByPlate(ctx, vehicledomain.Plate(validInput.Plate)).Return(nil, nil)
		repo.EXPECT().Save(ctx, mock.MatchedBy(func(v *vehicledomain.Vehicle) bool {
			return v.VIN == vehicledomain.VIN(validInput.VIN) && v.Plate == vehicledomain.Plate(validInput.Plate)
		})).Return(nil)

		uc := vehicleusecase.NewCreateVehicle(repo, modelRepo, newValidate())
		out, err := uc.Execute(ctx, validInput)

		require.NoError(t, err)
		assert.Equal(t, validInput.VIN, out.VIN)
		assert.Equal(t, validInput.Plate, out.Plate)
	})

	t.Run("returns not found when model year does not exist", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		modelRepo.EXPECT().GetModelYearByID(ctx, vehicledomain.VehicleModelYearID(validInput.ModelYearID)).Return(nil, nil)

		uc := vehicleusecase.NewCreateVehicle(repo, modelRepo, newValidate())
		_, err := uc.Execute(ctx, validInput)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindNotFound, appErr.Kind)
	})

	t.Run("returns conflict when VIN already exists", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		modelRepo.EXPECT().GetModelYearByID(ctx, vehicledomain.VehicleModelYearID(validInput.ModelYearID)).Return(modelYear, nil)
		repo.EXPECT().GetByVIN(ctx, vehicledomain.VIN(validInput.VIN)).Return(&vehicledomain.Vehicle{VIN: vehicledomain.VIN(validInput.VIN)}, nil)

		uc := vehicleusecase.NewCreateVehicle(repo, modelRepo, newValidate())
		_, err := uc.Execute(ctx, validInput)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindConflict, appErr.Kind)
	})

	t.Run("returns conflict when plate already exists", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		modelRepo.EXPECT().GetModelYearByID(ctx, vehicledomain.VehicleModelYearID(validInput.ModelYearID)).Return(modelYear, nil)
		repo.EXPECT().GetByVIN(ctx, vehicledomain.VIN(validInput.VIN)).Return(nil, nil)
		repo.EXPECT().GetByPlate(ctx, vehicledomain.Plate(validInput.Plate)).Return(&vehicledomain.Vehicle{Plate: vehicledomain.Plate(validInput.Plate)}, nil)

		uc := vehicleusecase.NewCreateVehicle(repo, modelRepo, newValidate())
		_, err := uc.Execute(ctx, validInput)

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindConflict, appErr.Kind)
	})

	t.Run("returns validation error for invalid VIN", func(t *testing.T) {
		repo := mockvehicle.NewMockRepository(t)
		modelRepo := mockvehicle.NewMockModelRepository(t)
		ctx := authCtx()

		uc := vehicleusecase.NewCreateVehicle(repo, modelRepo, newValidate())
		_, err := uc.Execute(ctx, vehicledto.CreateVehicleInput{
			ModelYearID: uuid.New(),
			VIN:         "INVALID",
			Plate:       "ABC1234",
			Color:       "#FF0000",
		})

		require.Error(t, err)
		var appErr *apperr.Error
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, apperr.KindValidation, appErr.Kind)
	})
}
