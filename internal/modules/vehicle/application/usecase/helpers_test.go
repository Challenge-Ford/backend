package vehicleusecase_test

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"torque/internal/core/appctx"
	vehicledomain "torque/internal/modules/vehicle/domain"
)

func authCtx() context.Context {
	return appctx.WithAuth(context.Background(), appctx.AuthContext{
		UserID: uuid.New(),
		Role:   "admin",
	})
}

func newValidate() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("vin", func(fl validator.FieldLevel) bool {
		return vehicledomain.VIN(fl.Field().String()).Validate() == nil
	})
	v.RegisterValidation("plate", func(fl validator.FieldLevel) bool {
		return vehicledomain.Plate(fl.Field().String()).Validate() == nil
	})
	return v
}

func sampleModelYear() *vehicledomain.VehicleModelYear {
	return &vehicledomain.VehicleModelYear{
		ID:      vehicledomain.NewVehicleModelYearID(),
		ModelID: vehicledomain.NewVehicleModelID(),
		Year:    2024,
		Model:   &vehicledomain.VehicleModel{Name: "Corolla", Type: "sedan"},
	}
}
