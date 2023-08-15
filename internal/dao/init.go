/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 15:55
 * @Description:
 */

package dao

import (
	"blog/pkg/storage/mysql"
	"database/sql"
)

var (
	db *sql.DB
)

func Init(cfg *mysql.Config) *sql.DB {
	db = mysql.NewMySQL(cfg)
	return db
}

func GetDB() *sql.DB {
	return db
}
