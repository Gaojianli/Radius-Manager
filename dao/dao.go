package dao

import "gorm.io/gorm"

type DAOManager struct {
	User    UserDAO
	AuthLog AuthLogDAO
}

func NewDAOManager(db *gorm.DB) *DAOManager {
	return &DAOManager{
		User:    NewUserDAO(db),
		AuthLog: NewAuthLogDAO(db),
	}
}