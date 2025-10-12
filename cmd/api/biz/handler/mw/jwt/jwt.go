package jwt

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/hertz-contrib/jwt"

	"github.com/nnieie/golanglab5/cmd/api/biz/model/base"
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
		TokenLookup: "header: Access-Token, cookie: access_token",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*base.User); ok {
				return jwt.MapClaims{identityKey: v.ID}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			id, ok := claims[identityKey].(float64)
			if !ok {
				return int64(0)
			}
			return int64(id)
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
			id, ok := claims[identityKey].(float64)
			if !ok {
				return int64(0)
			}
			return int64(id)
		},
	})
	if err != nil {
		panic("JWT Error:" + err.Error())
	}
}

func ExtractUserID(c *app.RequestContext) (int64, error) {
	v, ok := c.Get(identityKey)
	if !ok {
		return 0, errno.AuthorizationFailedErr
	}

	userID, ok := v.(int64)
	if !ok {
		return 0, errno.AuthorizationFailedErr
	}

	return userID, nil
}
