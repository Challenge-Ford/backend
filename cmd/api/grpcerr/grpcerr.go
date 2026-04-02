package grpcerr

import (
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"torque/internal/core/apperr"
)

func From(err error) error {
	if err == nil {
		return nil
	}

	var appErr *apperr.Error
	if !errors.As(err, &appErr) {
		return status.Error(codes.Internal, err.Error())
	}

	switch appErr.Kind {
	case apperr.KindValidation:
		return validationStatus(appErr)
	case apperr.KindBadRequest:
		return status.Error(codes.InvalidArgument, appErr.Message)
	case apperr.KindNotFound:
		return status.Error(codes.NotFound, appErr.Message)
	case apperr.KindConflict:
		return status.Error(codes.AlreadyExists, appErr.Message)
	case apperr.KindForbidden:
		return status.Error(codes.PermissionDenied, appErr.Message)
	case apperr.KindUnauthorized:
		return status.Error(codes.Unauthenticated, appErr.Message)
	default:
		return status.Error(codes.Internal, appErr.Message)
	}
}

func validationStatus(appErr *apperr.Error) error {
	violations := make([]*errdetails.BadRequest_FieldViolation, len(appErr.ValidationErrors))
	for i, v := range appErr.ValidationErrors {
		violations[i] = &errdetails.BadRequest_FieldViolation{
			Field:       v.Field,
			Description: v.Message,
		}
	}

	st, err := status.New(codes.InvalidArgument, appErr.Message).
		WithDetails(&errdetails.BadRequest{FieldViolations: violations})
	if err != nil {
		return status.Error(codes.InvalidArgument, appErr.Message)
	}
	return st.Err()
}
