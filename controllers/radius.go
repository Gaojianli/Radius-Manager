package controllers

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Gaojianli/raduis_mgnt/database"
	"github.com/Gaojianli/raduis_mgnt/models"
)

type RadiusController struct{}

type RadiusAuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RadiusAuthResponse struct {
	StatusCode int    `json:"REST-HTTP-Status-Code,omitempty"`
	Reply      string `json:"reply:Reply-Message,omitempty"`
}

type RadiusAuthorizeResponse struct {
	Reply   string `json:"reply:Reply-Message,omitempty"`
}

func (rc *RadiusController) Authenticate(ctx context.Context, c *app.RequestContext) {
	var req RadiusAuthRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, RadiusAuthResponse{
			StatusCode: 400,
			Reply:      "Invalid request format",
		})
		return
	}

	user, err := database.DAO.User.GetByUsernameForAuth(ctx, req.Username)
	if err != nil {
		c.JSON(consts.StatusNotFound, RadiusAuthResponse{
			StatusCode: 404,
			Reply:      "Authentication failed: user not found or disabled",
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		// 记录密码错误的认证日志
		authLog := &models.AuthLog{
			Username:  req.Username,
			AuthType:  "authenticate",
			Success:   false,
			IPAddress: c.ClientIP(),
			UserAgent: string(c.GetHeader("User-Agent")),
			CreatedAt: time.Now(),
		}
		
		// 异步记录日志，不影响响应速度
		go func() {
			database.DAO.AuthLog.Create(context.Background(), authLog)
		}()

		c.JSON(consts.StatusForbidden, RadiusAuthResponse{
			StatusCode: 403,
			Reply:      "Authentication failed: invalid password",
		})
		return
	}

	// 记录认证成功的日志
	authLog := &models.AuthLog{
		Username:  user.Username,
		AuthType:  "authenticate",
		Success:   true,
		IPAddress: c.ClientIP(),
		UserAgent: string(c.GetHeader("User-Agent")),
		CreatedAt: time.Now(),
	}
	
	// 异步记录日志，不影响响应速度
	go func() {
		database.DAO.AuthLog.Create(context.Background(), authLog)
	}()

	c.JSON(consts.StatusOK, RadiusAuthResponse{
		StatusCode: 200,
		Reply:      "Welcome, " + user.Username,
	})
}

func (rc *RadiusController) Authorize(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, RadiusAuthorizeResponse{
			Reply: "Invalid request format",
		})
		return
	}

	_, err := database.DAO.User.GetByUsernameForAuth(ctx, req.Username)
	success := err == nil

	if !success {
		c.JSON(consts.StatusForbidden, RadiusAuthorizeResponse{
			Reply: "User not found, disabled, or banned",
		})
		return
	}

	c.JSON(consts.StatusOK, RadiusAuthorizeResponse{
	})
}

func (rc *RadiusController) Accounting(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Username       string `json:"username" binding:"required"`
		AccountingType string `json:"acct_type"`
		SessionID      string `json:"session_id"`
		SessionTime    string `json:"session_time"`
		InputOctets    string `json:"input_octets"`
		OutputOctets   string `json:"output_octets"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"result":  "error",
			"message": "Invalid request format",
		})
		return
	}

	_, err := database.DAO.User.GetByUsername(ctx, req.Username)
	if err != nil {
		c.JSON(consts.StatusOK, map[string]interface{}{
			"result":  "ok",
			"message": "Accounting logged",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"result":  "ok",
		"message": "Accounting logged successfully",
	})
}
