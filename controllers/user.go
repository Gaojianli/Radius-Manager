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

type StatsResponse struct {
	TotalUsers  int64 `json:"total_users"`
	ActiveUsers int64 `json:"active_users"`
	BannedUsers int64 `json:"banned_users"`
	AuthCount   int64 `json:"auth_count"`
}

type AuthLogsResponse struct {
	Data struct {
		Logs       []AuthLogResponse `json:"logs"`
		Pagination struct {
			Page  int   `json:"page"`
			Limit int   `json:"limit"`
			Total int64 `json:"total"`
		} `json:"pagination"`
	} `json:"data"`
}

type AuthLogResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AuthType  string `json:"auth_type"`
	Success   bool   `json:"success"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	CreatedAt int64  `json:"created_at"`
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

	// 更新用户状态
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

func (uc *UserController) GetStats(ctx context.Context, c *app.RequestContext) {
	currentUser, err := middleware.GetCurrentUser(ctx, c)
	if err != nil {
		c.JSON(consts.StatusUnauthorized, map[string]interface{}{
			"code":    consts.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	var stats StatsResponse

	if currentUser.IsAdmin {
		// 管理员可以查看所有统计数据
		totalUsers, err := database.DAO.User.GetTotalCount(ctx)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    consts.StatusInternalServerError,
				"message": "Failed to get total users count",
			})
			return
		}

		activeUsers, err := database.DAO.User.GetActiveCount(ctx)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    consts.StatusInternalServerError,
				"message": "Failed to get active users count",
			})
			return
		}

		bannedUsers, err := database.DAO.User.GetBannedCount(ctx)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    consts.StatusInternalServerError,
				"message": "Failed to get banned users count",
			})
			return
		}

		totalAuthCount, err := database.DAO.AuthLog.GetTotalSuccessCount(ctx)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    consts.StatusInternalServerError,
				"message": "Failed to get total auth count",
			})
			return
		}

		stats = StatsResponse{
			TotalUsers:  totalUsers,
			ActiveUsers: activeUsers,
			BannedUsers: bannedUsers,
			AuthCount:   totalAuthCount,
		}
	} else {
		// 普通用户只能查看自己的授权次数
		authCount, err := database.DAO.AuthLog.GetSuccessCountByUsername(ctx, currentUser.Username)
		if err != nil {
			c.JSON(consts.StatusInternalServerError, map[string]interface{}{
				"code":    consts.StatusInternalServerError,
				"message": "Failed to get auth count",
			})
			return
		}

		stats = StatsResponse{
			AuthCount: authCount,
		}
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": consts.StatusOK,
		"data": stats,
	})
}

// 管理员专用统计接口
func (uc *UserController) GetAdminStats(ctx context.Context, c *app.RequestContext) {
	// 管理员可以查看所有统计数据
	totalUsers, err := database.DAO.User.GetTotalCount(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to get total users count",
		})
		return
	}

	activeUsers, err := database.DAO.User.GetActiveCount(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to get active users count",
		})
		return
	}

	bannedUsers, err := database.DAO.User.GetBannedCount(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to get banned users count",
		})
		return
	}

	totalAuthCount, err := database.DAO.AuthLog.GetTotalSuccessCount(ctx)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to get total auth count",
		})
		return
	}

	stats := StatsResponse{
		TotalUsers:  totalUsers,
		ActiveUsers: activeUsers,
		BannedUsers: bannedUsers,
		AuthCount:   totalAuthCount,
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": consts.StatusOK,
		"data": stats,
	})
}

func (uc *UserController) GetAuthLogs(ctx context.Context, c *app.RequestContext) {
	// 获取分页参数
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")
	username := c.Query("username") // 可选的用户名筛选

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 || limitInt > 100 {
		limitInt = 20
	}

	offset := (pageInt - 1) * limitInt

	// 查询日志
	logs, total, err := database.DAO.AuthLog.List(ctx, offset, limitInt, username)
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]interface{}{
			"code":    consts.StatusInternalServerError,
			"message": "Failed to get auth logs",
		})
		return
	}

	// 转换为响应格式
	logResponses := make([]AuthLogResponse, len(logs))
	for i, log := range logs {
		logResponses[i] = AuthLogResponse{
			ID:        log.ID,
			Username:  log.Username,
			AuthType:  log.AuthType,
			Success:   log.Success,
			IPAddress: log.IPAddress,
			UserAgent: log.UserAgent,
			CreatedAt: log.CreatedAt.Unix(),
		}
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"code": consts.StatusOK,
		"data": map[string]interface{}{
			"logs": logResponses,
			"pagination": map[string]interface{}{
				"page":  pageInt,
				"limit": limitInt,
				"total": total,
			},
		},
	})
}
