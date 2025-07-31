package controllers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/Gaojianli/raduis_mgnt/database"
	"github.com/Gaojianli/raduis_mgnt/middleware"
	"github.com/Gaojianli/raduis_mgnt/models"
)

type UserController struct{}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	IsAdmin  *bool  `json:"is_admin"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type AdminChangePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (uc *UserController) AdminCreateUser(ctx context.Context, c *app.RequestContext) {
	var req CreateUserRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    consts.StatusBadRequest,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	_, err := database.DAO.User.GetByUsernameOrEmail(ctx, req.Username, req.Email)
	if err == nil {
		c.JSON(consts.StatusConflict, map[string]interface{}{
			"code":    consts.StatusConflict,
			"message": "Username or email already exists",
		})
		return
	}

	isAdmin := false
	if req.IsAdmin != nil {
		isAdmin = *req.IsAdmin
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		IsAdmin:  isAdmin,
		Status:   true,
	}

	if err := database.DAO.User.Create(ctx, &user); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(consts.StatusCreated, map[string]interface{}{
		"code":    consts.StatusCreated,
		"message": "User created successfully",
		"data":    user.ToResponse(),
	})
}

func (uc *UserController) GetProfile(ctx context.Context, c *app.RequestContext) {
	currentUser, err := middleware.GetCurrentUser(ctx, c)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    consts.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	user, err := database.DAO.User.GetByID(ctx, currentUser.UserID)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    consts.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": consts.StatusOK,
		"data": user.ToResponse(),
	})
}

func (uc *UserController) ChangePassword(ctx context.Context, c *app.RequestContext) {
	currentUser, err := middleware.GetCurrentUser(ctx, c)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    consts.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    consts.StatusBadRequest,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	user, err := database.DAO.User.GetByID(ctx, currentUser.UserID)
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    consts.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	if !user.CheckPassword(req.OldPassword) {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    consts.StatusBadRequest,
			"message": "Invalid old password",
		})
		return
	}

	if err := user.HashPassword(req.NewPassword); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to hash password",
		})
		return
	}

	if err := database.DAO.User.Update(ctx, user); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to update password",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    consts.StatusOK,
		"message": "Password changed successfully",
	})
}

func (uc *UserController) GetUsers(ctx context.Context, c *app.RequestContext) {
	page, _ := strconv.Atoi(string(c.Query("page")))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(string(c.Query("limit")))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	users, total, err := database.DAO.User.List(ctx, offset, limit)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to fetch users",
		})
		return
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, user.ToResponse())
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": consts.StatusOK,
		"data": map[string]interface{}{
			"users": userResponses,
			"pagination": map[string]interface{}{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		},
	})
}

func (uc *UserController) AdminChangePassword(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    consts.StatusBadRequest,
			"message": "Invalid user ID",
		})
		return
	}

	var req AdminChangePasswordRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    consts.StatusBadRequest,
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	user, err := database.DAO.User.GetByID(ctx, uint(userID))
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    consts.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	if err := user.HashPassword(req.NewPassword); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to hash password",
		})
		return
	}

	if err := database.DAO.User.Update(ctx, user); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to update password",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    consts.StatusOK,
		"message": "Password changed successfully",
	})
}

func (uc *UserController) ToggleUserStatus(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    consts.StatusBadRequest,
			"message": "Invalid user ID",
		})
		return
	}

	user, err := database.DAO.User.GetByID(ctx, uint(userID))
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    consts.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	user.Status = !user.Status
	if err := database.DAO.User.Update(ctx, user); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to update user status",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    consts.StatusOK,
		"message": "User status updated successfully",
		"data":    user.ToResponse(),
	})
}

func (uc *UserController) AdminDeleteUser(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    consts.StatusBadRequest,
			"message": "Invalid user ID",
		})
		return
	}

	currentUser, err := middleware.GetCurrentUser(ctx, c)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    consts.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	if currentUser.UserID == uint(userID) {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    consts.StatusBadRequest,
			"message": "Cannot delete yourself",
		})
		return
	}

	_, err = database.DAO.User.GetByID(ctx, uint(userID))
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    consts.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	if err := database.DAO.User.Delete(ctx, uint(userID)); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to delete user",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    consts.StatusOK,
		"message": "User deleted successfully",
	})
}

func (uc *UserController) AdminToggleBanUser(ctx context.Context, c *app.RequestContext) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    consts.StatusBadRequest,
			"message": "Invalid user ID",
		})
		return
	}

	currentUser, err := middleware.GetCurrentUser(ctx, c)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    consts.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	if currentUser.UserID == uint(userID) {
		c.JSON(consts.StatusBadRequest, map[string]interface{}{
			"code":    consts.StatusBadRequest,
			"message": "Cannot ban yourself",
		})
		return
	}

	user, err := database.DAO.User.GetByID(ctx, uint(userID))
	if err != nil {
		c.JSON(consts.StatusNotFound, map[string]interface{}{
			"code":    consts.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	newBannedStatus := !user.Banned
	if err := database.DAO.User.UpdateBanned(ctx, uint(userID), newBannedStatus); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to update user ban status",
		})
		return
	}

	user.Banned = newBannedStatus
	action := "banned"
	if !newBannedStatus {
		action = "unbanned"
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code":    consts.StatusOK,
		"message": fmt.Sprintf("User %s successfully", action),
		"data":    user.ToResponse(),
	})
}
