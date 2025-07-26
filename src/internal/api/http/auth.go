package http

import (
	"context"
	"net/http"
	"t03/internal/domain"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func UserIDFromCtx(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDKey).(string)
	return v, ok
}

type UserAuthenticator struct {
	AuthService domain.UserService
}

func NewUserAuthenticator(authService domain.UserService) *UserAuthenticator {
	return &UserAuthenticator{AuthService: authService}
}

func (ua *UserAuthenticator) Protect(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		userID, err := ua.AuthService.AuthenticateBasic(auth)
		if err != nil {
			http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}
