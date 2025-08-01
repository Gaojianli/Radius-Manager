package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/Gaojianli/raduis_mgnt/models"
)

type UserDAO interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByUsernameOrEmail(ctx context.Context, username, email string) (*models.User, error)
	GetByUsernameForAuth(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]models.User, int64, error)
	CountAdmins(ctx context.Context) (int64, error)
	UpdatePassword(ctx context.Context, id uint, password, salt string) error
	UpdateBanned(ctx context.Context, id uint, banned bool) error
	GetTotalCount(ctx context.Context) (int64, error)
	GetActiveCount(ctx context.Context) (int64, error)
	GetBannedCount(ctx context.Context) (int64, error)
}

type userDAOImpl struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &userDAOImpl{db: db}
}

func (d *userDAOImpl) Create(ctx context.Context, user *models.User) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *userDAOImpl) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := d.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAOImpl) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := d.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAOImpl) GetByUsernameOrEmail(ctx context.Context, username, email string) (*models.User, error) {
	var user models.User
	err := d.db.WithContext(ctx).Where("username = ? OR email = ?", username, email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func (d *userDAOImpl) Update(ctx context.Context, user *models.User) error {
	return d.db.WithContext(ctx).Save(user).Error
}

func (d *userDAOImpl) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (d *userDAOImpl) List(ctx context.Context, offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := d.db.WithContext(ctx).Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := d.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (d *userDAOImpl) CountAdmins(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&models.User{}).Where("is_admin = ?", true).Count(&count).Error
	return count, err
}

func (d *userDAOImpl) UpdatePassword(ctx context.Context, id uint, password, salt string) error {
	return d.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password": password,
		"salt":     salt,
	}).Error
}


func (d *userDAOImpl) GetByUsernameForAuth(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := d.db.WithContext(ctx).Where("username = ? AND banned = ?", username, false).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *userDAOImpl) UpdateBanned(ctx context.Context, id uint, banned bool) error {
	return d.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Update("banned", banned).Error
}

func (d *userDAOImpl) GetTotalCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	return count, err
}

func (d *userDAOImpl) GetActiveCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&models.User{}).Where("banned = ?", false).Count(&count).Error
	return count, err
}

func (d *userDAOImpl) GetBannedCount(ctx context.Context) (int64, error) {
	var count int64
	err := d.db.WithContext(ctx).Model(&models.User{}).Where("banned = ?", true).Count(&count).Error
	return count, err
}
