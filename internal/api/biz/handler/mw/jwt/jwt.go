package jwt

import (
	"context"
	"strconv"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"

	"github.com/nnieie/golanglab5/internal/api/biz/model/base"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/errno"
)

var (
	AccessTokenJwtMiddleware  *jwt.HertzJWTMiddleware
	RefreshTokenJwtMiddleware *jwt.HertzJWTMiddleware
	identityKey               = constants.IdentityKey
)

func InitJwt() {
	var err error
	AccessTokenJwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "access",
		Key:         []byte(constants.AccessTokenSecretKey),
		Timeout:     time.Hour,
		IdentityKey: identityKey,
		TokenLookup: "header: Access-Token, cookie: access_token, query: token",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*base.User); ok {
				return jwt.MapClaims{identityKey: v.ID}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			if id, ok := claims[identityKey].(string); ok {
				return id
			}
			if id, ok := claims[identityKey].(float64); ok {
				return strconv.FormatInt(int64(id), 10)
			}
			return ""
		},
	})
	if err != nil {
		panic("JWT Error:" + err.Error())
	}

	RefreshTokenJwtMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "refresh",
		Key:         []byte(constants.RefreshTokenSecretKey),
		Timeout:     time.Hour * 24 * 7,
		IdentityKey: identityKey,
		TokenLookup: "header: Refresh-Token, cookie: refresh_token",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*base.User); ok {
				return jwt.MapClaims{identityKey: v.ID}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			if id, ok := claims[identityKey].(string); ok {
				return id
			}
			if id, ok := claims[identityKey].(float64); ok {
				return strconv.FormatInt(int64(id), 10)
			}
			return ""
		},
	})
	if err != nil {
		panic("JWT Error:" + err.Error())
	}
}

func ExtractUserID(c *app.RequestContext) (string, error) {
	v, ok := c.Get(identityKey)
	if !ok {
		return "", errno.AuthorizationFailedErr
	}

	if userID, ok := v.(string); ok {
		return userID, nil
	}
	if userID, ok := v.(float64); ok {
		return strconv.FormatInt(int64(userID), 10), nil
	}

	return "", errno.AuthorizationFailedErr
}
