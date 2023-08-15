/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 15:53
 * @Description:
 */

package dao

import (
	"blog/internal/model"
	"blog/pkg/utils"
	"database/sql"
	"log"
)

func GetCategories() (categories []model.Category, err error) {
	rs, err := db.Query("select id, name from category")
	if utils.IsNotNil(err) {
		return nil, err
	}
	defer rs.Close()
	for rs.Next() {
		var category model.Category
		rs.Scan(&category.Id, &category.Name)
		categories = append(categories, category)
	}
	return
}

func GetCategoryIdsByName(name string) (categoryIds []string, err error) {
	rs, err := db.Query("select id from category where name like ?", "%"+name+"%")
	if utils.IsNotNil(err) {
		return nil, err
	}
	defer rs.Close()
	for rs.Next() {
		var categoryId string
		rs.Scan(&categoryId)
		categoryIds = append(categoryIds, categoryId)
	}
	return
}

func GetCategory(id int) (category model.Category) {
	row := db.QueryRow("select id, name from category where id=?", id)
	row.Scan(&category.Id, &category.Name)
	return
}

func DeleteCategory(category model.Category) (id int, err error) {
	id = category.Id
	log.Println(category)
	_, err = db.Exec("delete from category where id=?", id)
	if utils.IsNotNil(err) {
		log.Printf("category id: %d delete err: %v", id, err)
		return
	}
	return
}

func SaveCategory(category model.Category) (id int, err error) {
	var rs sql.Result
	if category.Id > 0 {
		id = category.Id
		log.Println(category)
		_, err = db.Exec("update category set name=? where id=?", category.Name, category.Id)
		if utils.IsNotNil(err) {
			log.Printf("category id: %d update err : %v", id, err)
			return
		}
	} else {
		rs, err = db.Exec("insert into category (`name`) values (?)", category.Name)
		if utils.IsNotNil(err) {
			log.Printf("category name: %s insert err: %v", category.Name, err)
			return
		}
		id64, _ := rs.LastInsertId()
		id = int(id64)
	}
	return
}
