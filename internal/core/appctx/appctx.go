package appctx

import "context"

type AuthContext struct {
	UserID       string
	DealershipID string
	Role         string
}

type authContextKey struct{}

func WithAuth(ctx context.Context, auth AuthContext) context.Context {
	return context.WithValue(ctx, authContextKey{}, auth)
}

func GetAuth(ctx context.Context) (AuthContext, bool) {
	auth, ok := ctx.Value(authContextKey{}).(AuthContext)
	return auth, ok
}

func MustGetAuth(ctx context.Context) AuthContext {
	auth, ok := GetAuth(ctx)
	if !ok {
		panic("auth context not found")
	}
	return auth
}
