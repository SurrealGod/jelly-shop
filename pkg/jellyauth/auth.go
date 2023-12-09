package jellyauth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/SurrealGod/jelly-shop-code/config"
	"github.com/SurrealGod/jelly-shop-code/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type IJellyAuth interface {
	SignToken() string
}
type jellyAdmin struct {
	*jellyAuth
}

type jellyAuth struct {
	mapClaims *jellyMapClaims
	cfg       config.IJwtConfig
}

type jellyMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}
type IJellyAdmin interface {
	SignToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *jellyAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.SecretKey())
	return ss
}
func (a *jellyAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	ss, _ := token.SignedString(a.cfg.AdminKey())
	return ss
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*jellyMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jellyMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*jellyMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}
func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*jellyMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jellyMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method is invalid")
		}
		return cfg.AdminKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("token format is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token had expired")
		} else {
			return nil, fmt.Errorf("parse token failed: %v", err)
		}
	}

	if claims, ok := token.Claims.(*jellyMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("claims type is invalid")
	}
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &jellyAuth{
		cfg: cfg,
		mapClaims: &jellyMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "Jellyshop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewJellyAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IJellyAuth, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	default:
		return nil, fmt.Errorf("unknown token type")
	}

}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IJellyAuth {
	return &jellyAuth{
		cfg: cfg,
		mapClaims: &jellyMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "Jellyshop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IJellyAuth {
	return &jellyAuth{
		cfg: cfg,
		mapClaims: &jellyMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "Jellyshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}
func newAdminToken(cfg config.IJwtConfig) IJellyAdmin {
	return &jellyAdmin{
		jellyAuth: &jellyAuth{
			cfg: cfg,
			mapClaims: &jellyMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "jellyshop-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
