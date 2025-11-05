package models

import (
	"time"

	"gorm.io/gorm"
)

type AuthLog struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	Username   string    `json:"username" gorm:"not null;index"`
	AuthType   string    `json:"auth_type" gorm:"not null"` // "authorize" or "authenticate"
	Success    bool      `json:"success" gorm:"not null;index"`
	IPAddress  string    `json:"ip_address" gorm:"not null"`
	UserAgent  string    `json:"user_agent"`
	DeviceMAC  string    `json:"device_mac"`  // 设备 MAC 地址
	TargetSSID string    `json:"target_ssid"` // 目标 SSID
	CreatedAt  time.Time `json:"created_at"`
}

func (AuthLog) TableName() string {
	return "auth_logs"
}

// AutoMigrate 自动迁移
func MigrateAuthLog(db *gorm.DB) error {
	return db.AutoMigrate(&AuthLog{})
}
