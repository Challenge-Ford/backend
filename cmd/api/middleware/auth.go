package middleware

import (
	"context"

	"torque/internal/core/appctx"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Auth(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	userID, err := extractUUID(md, "x-user-id")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "missing or invalid x-user-id")
	}

	values := md.Get("x-dealership-id")
	dealershipID := ""
	if len(values) > 0 {
		dealershipID = values[0]
	}

	role := ""
	if r := md.Get("x-user-role"); len(r) > 0 {
		role = r[0]
	}

	ctx = appctx.WithAuth(ctx, appctx.AuthContext{
		UserID:       userID,
		DealershipID: dealershipID,
		Role:         role,
	})

	return handler(ctx, req)
}

func extractUUID(md metadata.MD, key string) (uuid.UUID, error) {
	values := md.Get(key)
	if len(values) == 0 {
		return uuid.Nil, status.Error(codes.Unauthenticated, "missing "+key)
	}
	return uuid.Parse(values[0])
}
