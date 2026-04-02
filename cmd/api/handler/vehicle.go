package handler

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"torque/cmd/api/grpcerr"
	"torque/internal/core/pagination"
	vehicledto "torque/internal/modules/vehicle/application/dto"
	vehicleusecase "torque/internal/modules/vehicle/application/usecase"
	vehicledomain "torque/internal/modules/vehicle/domain"
	vehiclev1 "torque/gen/proto/vehicle/v1"
)

type VehicleHandler struct {
	vehiclev1.UnimplementedVehicleServiceServer
	create *vehicleusecase.CreateVehicleUseCase
	get    *vehicleusecase.GetVehicleUseCase
	list   *vehicleusecase.ListVehiclesUseCase
	update *vehicleusecase.UpdateVehicleUseCase
	delete *vehicleusecase.DeleteVehicleUseCase
}

func NewVehicleHandler(
	create *vehicleusecase.CreateVehicleUseCase,
	get *vehicleusecase.GetVehicleUseCase,
	list *vehicleusecase.ListVehiclesUseCase,
	update *vehicleusecase.UpdateVehicleUseCase,
	delete *vehicleusecase.DeleteVehicleUseCase,
) *VehicleHandler {
	return &VehicleHandler{create: create, get: get, list: list, update: update, delete: delete}
}

func (h *VehicleHandler) CreateVehicle(ctx context.Context, req *vehiclev1.CreateVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	var customerID *uuid.UUID
	if req.CustomerId != "" {
		parsed, err := uuid.Parse(req.CustomerId)
		if err != nil {
			return nil, grpcerr.From(err)
		}
		customerID = &parsed
	}

	output, err := h.create.Execute(ctx, vehicledto.CreateVehicleInput{
		CustomerID: customerID,
		VIN:        req.Vin,
		Plate:      req.Plate,
		Model:      req.Model,
		Year:       int(req.Year),
		Color:      req.Color,
	})
	if err != nil {
		return nil, grpcerr.From(err)
	}

	return toProto(output), nil
}

func (h *VehicleHandler) GetVehicle(ctx context.Context, req *vehiclev1.GetVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, grpcerr.From(err)
	}

	output, err := h.get.Execute(ctx, vehicledomain.VehicleID(id))
	if err != nil {
		return nil, grpcerr.From(err)
	}

	return toProto(output), nil
}

func (h *VehicleHandler) ListVehicles(ctx context.Context, req *vehiclev1.ListVehiclesRequest) (*vehiclev1.ListVehiclesResponse, error) {
	result, err := h.list.Execute(ctx, pagination.Page{
		Page:    int(req.Page),
		PerPage: int(req.PerPage),
	})
	if err != nil {
		return nil, grpcerr.From(err)
	}

	vehicles := make([]*vehiclev1.VehicleResponse, len(result.Data))
	for i, v := range result.Data {
		vehicles[i] = toProto(v)
	}

	return &vehiclev1.ListVehiclesResponse{
		Vehicles: vehicles,
		Meta: &vehiclev1.PageMeta{
			Page:       int32(result.Meta.Page),
			PerPage:    int32(result.Meta.PerPage),
			Total:      int32(result.Meta.Total),
			TotalPages: int32(result.Meta.TotalPages),
		},
	}, nil
}

func (h *VehicleHandler) UpdateVehicle(ctx context.Context, req *vehiclev1.UpdateVehicleRequest) (*vehiclev1.VehicleResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, grpcerr.From(err)
	}

	output, err := h.update.Execute(ctx, vehicledomain.VehicleID(id), vehicledto.UpdateVehicleInput{
		Plate: req.Plate,
		Model: req.Model,
		Year:  int(req.Year),
		Color: req.Color,
	})
	if err != nil {
		return nil, grpcerr.From(err)
	}

	return toProto(output), nil
}

func (h *VehicleHandler) DeleteVehicle(ctx context.Context, req *vehiclev1.DeleteVehicleRequest) (*emptypb.Empty, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, grpcerr.From(err)
	}

	if err := h.delete.Execute(ctx, vehicledomain.VehicleID(id)); err != nil {
		return nil, grpcerr.From(err)
	}

	return &emptypb.Empty{}, nil
}

func toProto(v *vehicledto.VehicleOutput) *vehiclev1.VehicleResponse {
	var customerID *wrapperspb.StringValue
	if v.CustomerID != nil {
		customerID = wrapperspb.String(*v.CustomerID)
	}

	return &vehiclev1.VehicleResponse{
		Id:         v.ID,
		CustomerId: customerID,
		Vin:        v.VIN,
		Plate:      v.Plate,
		Model:      v.Model,
		Year:       int32(v.Year),
		Color:      v.Color,
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  v.UpdatedAt,
	}
}
