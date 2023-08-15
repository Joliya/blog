/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 15:14
 * @Description:
 */

package orm

import (
	"blog/pkg/utils"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"

	// GORM MySQL
	"gorm.io/gorm"
)

type Config struct {
	Name            string
	Addr            string
	UserName        string
	Password        string
	ShowLog         bool
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifeTime time.Duration
	SlowThreshold   time.Duration
}

func NewMySQL(c *Config) (db *gorm.DB) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=%t&loc=%s",
		c.UserName,
		c.Password,
		c.Addr,
		c.Name,
		true,
		//"Asia/Shanghai"),
		"Local")
	sqlDB, err := sql.Open("mysql", dsn)
	if utils.IsNotNil(err) {
		log.Panicf("open mysql failed. database name: %s, err: %+v", c.Name, err)
	}
	sqlDB.SetMaxOpenConns(c.MaxOpenConn)
	sqlDB.SetMaxIdleConns(c.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime)

	db, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), gormConfig(c))
	if utils.IsNotNil(err) {
		log.Panicf("databse connection failed. database name: %s, err: %+v", c.Name, err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4")
	return db
}

func gormConfig(c *Config) *gorm.Config {
	config := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁止外键约束
	}
	if c.ShowLog {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	// 只打印慢查询
	if c.SlowThreshold > 0 {
		config.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: c.SlowThreshold,
				Colorful:      true,
				LogLevel:      logger.Warn,
			},
		)
	}
	config.NamingStrategy = schema.NamingStrategy{
		SingularTable: true, // 表名禁止自动复数
	}
	return config
}
