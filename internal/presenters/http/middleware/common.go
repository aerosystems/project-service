package middleware

import (
	"context"
	"errors"
	"github.com/aerosystems/project-service/internal/models"
	"net/http"
	"strings"
)

type ctxKey int

const (
	userContextKey ctxKey = iota
)

type User struct {
	Uuid        string
	Email       string
	Role        string
	DisplayName string
}

func GetUserFromContext(ctx context.Context) (User, error) {
	user, ok := ctx.Value(userContextKey).(User)
	if !ok {
		return User{}, errors.New("user not found in context")
	}
	return user, nil
}

func isAccess(roles []models.Role, role string) bool {
	for _, r := range roles {
		if r.String() == role {
			return true
		}
	}
	return false
}

func getAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", errors.New("missing Authorization header")
	}
	return authHeader, nil
}

func getTokenFromHeader(r *http.Request) (string, error) {
	authHeader, err := getAuthHeader(r)
	if err != nil {
		return "", err
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}
	return tokenParts[1], nil
}
