package dao

import "gorm.io/gorm"

type DAOManager struct {
	User UserDAO
}

func NewDAOManager(db *gorm.DB) *DAOManager {
	return &DAOManager{
		User: NewUserDAO(db),
	}
}