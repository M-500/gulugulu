// Package gormx 统一管理 GORM 数据库客户端的创建和连接池配置。
package gormx

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLConfig 描述 MySQL 连接和连接池参数。
// 默认值标签由 go-zero conf 在加载业务配置时填充。
type MySQLConfig struct {
	DataSource   string
	MaxIdleConns int `json:",default=10"`
	MaxOpenConns int `json:",default=100"`
	ConnMaxLife  int `json:",default=3600"`
}

// NewMySQL 创建配置好连接池的 GORM 客户端。
func NewMySQL(config MySQLConfig) (*gorm.DB, error) {
	if err := validateMySQLConfig(config); err != nil {
		return nil, err
	}
	db, err := gorm.Open(mysql.Open(config.DataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("初始化GORM MySQL客户端失败: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取GORM底层连接池失败: %w", err)
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLife) * time.Second)
	return db, nil
}

// Close 关闭 GORM 持有的底层 database/sql 连接池。
func Close(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取GORM底层连接池失败: %w", err)
	}
	return sqlDB.Close()
}

func validateMySQLConfig(config MySQLConfig) error {
	if strings.TrimSpace(config.DataSource) == "" {
		return fmt.Errorf("MySQL DataSource不能为空")
	}
	if config.MaxIdleConns < 0 || config.MaxOpenConns <= 0 || config.ConnMaxLife <= 0 {
		return fmt.Errorf("MySQL连接池参数不合法")
	}
	if config.MaxIdleConns > config.MaxOpenConns {
		return fmt.Errorf("MySQL MaxIdleConns不能大于MaxOpenConns")
	}
	return nil
}
