package middleware

import (
	"context"
	"errors"
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type ctxKey int

const (
	accessTokenClaimsContextKey ctxKey = iota
)

type accessTokenClaims struct {
	Uuid        string
	Email       string
	Role        string
	DisplayName string
}

type UserClaims struct {
	Uuid        uuid.UUID
	Email       string
	Role        models.Role
	DisplayName string
}

func NewUserClaims(claims accessTokenClaims) (UserClaims, error) {
	uuid, err := uuid.Parse(claims.Uuid)
	if err != nil {
		return UserClaims{}, err
	}
	role := models.RoleFromString(claims.Role)
	if role == models.UnknownRole {
		return UserClaims{}, CustomErrors.ErrUnknownUserRole
	}
	return UserClaims{
		Uuid:        uuid,
		Email:       claims.Email,
		Role:        role,
		DisplayName: claims.DisplayName,
	}, nil
}

func GetUserClaimsFromContext(ctx context.Context) (UserClaims, error) {
	claims, ok := ctx.Value(accessTokenClaimsContextKey).(accessTokenClaims)
	if !ok {
		return UserClaims{}, CustomErrors.ErrForbidden
	}
	user, err := NewUserClaims(claims)
	if err != nil {
		return UserClaims{}, CustomErrors.ErrForbidden
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
