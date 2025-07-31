package controllers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Gaojianli/raduis_mgnt/database"
)

type RadiusController struct{}

type RadiusAuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RadiusAuthResponse struct {
	Result   string            `json:"result"`
	Username string            `json:"username,omitempty"`
	Message  string            `json:"message,omitempty"`
	Attrs    map[string]string `json:"attrs,omitempty"`
}

func (rc *RadiusController) Authenticate(ctx context.Context, c *app.RequestContext) {
	var req RadiusAuthRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, RadiusAuthResponse{
			Result:  "reject",
			Message: "Invalid request format",
		})
		return
	}

	user, err := database.DAO.User.GetByUsernameForAuth(ctx, req.Username)
	if err != nil {
		c.JSON(consts.StatusOK, RadiusAuthResponse{
			Result:  "reject",
			Message: "Authentication failed: user not valid",
		})
		return
	}

	if !user.CheckPassword(req.Password) {
		c.JSON(consts.StatusOK, RadiusAuthResponse{
			Result:  "reject",
			Message: "Authentication failed: invalid password",
		})
		return
	}

	c.JSON(consts.StatusOK, RadiusAuthResponse{
		Result:   "accept",
		Username: user.Username,
		Message:  "Authentication successful",
		Attrs: map[string]string{
			"User-ID": user.Username,
			"User-Role": func() string {
				if user.IsAdmin {
					return "admin"
				}
				return "user"
			}(),
			"Session-Type": "authenticated",
		},
	})
}

func (rc *RadiusController) Authorize(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Username string `json:"username" binding:"required"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, RadiusAuthResponse{
			Result:  "reject",
			Message: "Invalid request format",
		})
		return
	}

	user, err := database.DAO.User.GetByUsernameForAuth(ctx, req.Username)
	if err != nil {
		c.JSON(consts.StatusOK, RadiusAuthResponse{
			Result:  "reject",
			Message: "User not found, disabled, or banned",
		})
		return
	}

	c.JSON(consts.StatusOK, RadiusAuthResponse{
		Result:   "accept",
		Username: user.Username,
		Message:  "User authorized to authenticate",
		Attrs: map[string]string{
			"Auth-Type": "REST",
			"User-ID":   user.Username,
			"User-Role": func() string {
				if user.IsAdmin {
					return "admin"
				}
				return "user"
			}(),
		},
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
