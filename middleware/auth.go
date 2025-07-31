package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"

	"github.com/Gaojianli/raduis_mgnt/config"
	"github.com/Gaojianli/raduis_mgnt/database"
	"github.com/Gaojianli/raduis_mgnt/models"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}

var (
	JWTMiddleware *jwt.HertzJWTMiddleware
	identityKey   = "user_id"
)

func InitJWT() error {
	var err error
	JWTMiddleware, err = jwt.New(&jwt.HertzJWTMiddleware{
		Realm:       "radius_mgnt",
		Key:         []byte(config.AppConfig.JWTSecret),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					"username":  v.Username,
					"is_admin":  v.IsAdmin,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			userID, _ := claims[identityKey].(float64)
			username, _ := claims["username"].(string)
			isAdmin, _ := claims["is_admin"].(bool)

			return &Claims{
				UserID:   uint(userID),
				Username: username,
				IsAdmin:  isAdmin,
			}
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var loginReq struct {
				Username string `json:"username" binding:"required"`
				Password string `json:"password" binding:"required"`
			}

			if err := c.BindAndValidate(&loginReq); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}

			user, err := database.DAO.User.GetByUsernameForAuth(ctx, loginReq.Username)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if !user.CheckPassword(loginReq.Password) {
				return nil, jwt.ErrFailedAuthentication
			}

			return &user, nil
		},
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			if claims, ok := data.(*Claims); ok {
				user, err := database.DAO.User.GetByID(ctx, claims.UserID)
				if err != nil || !user.Status || user.Banned {
					return false
				}
				return true
			}
			return false
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			c.JSON(code, map[string]interface{}{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	return err
}

func RequireAdmin() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		claims := jwt.ExtractClaims(ctx, c)
		isAdmin, exists := claims["is_admin"]
		if !exists || !isAdmin.(bool) {
			c.JSON(consts.StatusForbidden, map[string]interface{}{
				"code":    consts.StatusForbidden,
				"message": "Admin access required",
			})
			c.Abort()
			return
		}
		c.Next(ctx)
	}
}

func GetCurrentUser(ctx context.Context, c *app.RequestContext) (*Claims, error) {
	claims := jwt.ExtractClaims(ctx, c)
	if claims == nil {
		return nil, errors.New("no claims found")
	}

	userID, ok := claims[identityKey].(float64)
	if !ok {
		return nil, errors.New("invalid user ID in claims")
	}

	username, _ := claims["username"].(string)
	isAdmin, _ := claims["is_admin"].(bool)

	return &Claims{
		UserID:   uint(userID),
		Username: username,
		IsAdmin:  isAdmin,
	}, nil
}
