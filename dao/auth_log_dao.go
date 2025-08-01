package dao

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/Gaojianli/raduis_mgnt/models"
)

type AuthLogDAO interface {
	Create(ctx context.Context, authLog *models.AuthLog) error
	GetSuccessCountByUsername(ctx context.Context, username string) (int64, error)
	GetTotalSuccessCount(ctx context.Context) (int64, error)
	GetSuccessCountByDateRange(ctx context.Context, start, end time.Time) (int64, error)
	List(ctx context.Context, offset, limit int, username string) ([]models.AuthLog, int64, error)
}

type authLogDAOImpl struct {
	db *gorm.DB
}

func NewAuthLogDAO(db *gorm.DB) AuthLogDAO {
	return &authLogDAOImpl{db: db}
}

func (d *authLogDAOImpl) Create(ctx context.Context, authLog *models.AuthLog) error {
	return d.db.WithContext(ctx).Create(authLog).Error
}

func (d *authLogDAOImpl) GetSuccessCountByUsername(ctx context.Context, username string) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&models.AuthLog{}).
		Where("username = ? AND success = ?", username, true).
		Count(&count).Error
	return count, err
}

func (d *authLogDAOImpl) GetTotalSuccessCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&models.AuthLog{}).
		Where("success = ?", true).
		Count(&count).Error
	return count, err
}

func (d *authLogDAOImpl) GetSuccessCountByDateRange(ctx context.Context, start, end time.Time) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&models.AuthLog{}).
		Where("success = ? AND auth_type = ? AND created_at BETWEEN ? AND ?", true, "authorize", start, end).
		Count(&count).Error
	return count, err
}

func (d *authLogDAOImpl) List(ctx context.Context, offset, limit int, username string) ([]models.AuthLog, int64, error) {
	var logs []models.AuthLog
	var total int64

	query := d.db.WithContext(ctx).Model(&models.AuthLog{})
	
	// 如果提供了用户名筛选
	if username != "" {
		query = query.Where("username = ?", username)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据，按创建时间倒序
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}