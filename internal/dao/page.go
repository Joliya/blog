/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 16:10
 * @Description:
 */

package dao

import (
	"blog/internal/model"
	"blog/pkg/logger"
	"blog/pkg/utils"
	"database/sql"
	"go.uber.org/zap"
	"strconv"
)

type PageParams struct {
	PerPage int
	Page    int
}

func GetPages(params PageParams) (pages []model.Page, err error) {
	querySql := "select id, title, content from page"
	querySql += " order by id asc"
	if params.PerPage > 0 && params.Page > 0 {
		offset := (params.Page - 1) * params.PerPage
		querySql += " limit " + strconv.Itoa(offset) + "," + strconv.Itoa(params.PerPage)
	}
	rows, err := db.Query(querySql)
	if utils.IsNotNil(err) {
		logger.Error("GetPagesError", zap.Error(err))
		return
	}
	if utils.IsNotNil(rows.Err()) {
		logger.Error("GetPagesError", zap.Error(rows.Err()))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var page model.Page
		_ = rows.Scan(&page.Id, &page.Title, &page.Content)
		pages = append(pages, page)
	}
	return
}

func GetPage(id string) (page model.Page) {
	row := db.QueryRow("select id, title, content from page where id=? limit 1", id)
	if utils.IsNotNil(row.Err()) {
		logger.Error("GetPageError", zap.Error(row.Err()))
		return
	}
	err := row.Scan(&page.Id, &page.Title, &page.Content)
	if utils.IsNotNil(err) {
		logger.Error("GetPageError", zap.Error(err))
		return
	}
	return
}

func PageDelete(page model.Page) (id int, err error) {
	id = page.Id
	_, err = db.Exec("delete from page where id=?", id)
	if utils.IsNotNil(err) {
		logger.Error("PageDelete", zap.Error(err))
		return
	}
	return
}

func PageSave(page model.Page) (id int, err error) {
	var rs sql.Result
	if page.Id > 0 {
		id = page.Id
		_, err = db.Exec("update page set title=?, content=? where id=?", page.Title, page.Content, page.Id)
		if utils.IsNotNil(err) {
			logger.Error("PageSage", zap.Error(err))
			return
		}
	} else {
		rs, err = db.Exec("insert into page (titile, content) values (?,?)", page.Title, page.Content)
		if utils.IsNotNil(err) {
			logger.Error("PageSave", zap.Error(err))
			return
		}
		id64, _ := rs.LastInsertId()
		id = int(id64)
	}
	return
}
