/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 14:33
 * @Description:
 */

package mysql

import (
	"blog/pkg/utils"
	"database/sql"
	"github.com/pkg/errors"
	"time"
)

type Config struct {
	DSN             string
	MaxOpenConn     int
	MaxIdleConn     int
	ConnMaxLifeTime int
}

func NewMySQL(c *Config) (db *sql.DB) {
	db, err := connect(c, c.DSN)
	if utils.IsNotNil(err) {
		panic(err)
	}
	return
}

func connect(c *Config, dataSourceName string) (*sql.DB, error) {
	d, err := sql.Open("mysql", dataSourceName)
	if utils.IsNotNil(err) {
		err = errors.WithStack(err)
		return nil, err
	}
	d.SetMaxOpenConns(c.MaxOpenConn)
	d.SetMaxIdleConns(c.MaxIdleConn)
	d.SetConnMaxLifetime(time.Duration(c.ConnMaxLifeTime))
	return d, nil
}
