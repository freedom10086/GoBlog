package models

import (
	"log"
	"time"
)

type Category struct {
	Cid         int
	Title       string
	Description string
	Posts       int
	ToadyPosts  int
	LastPost    time.Time
	Created     time.Time
}

//新增category
func AddCate(name, description string) error {
	res, err := db.Exec(
		"call cate_add(?,?);",
		name,
		description)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if rowCnt < 1 {
		return ErrNoAff
	}
	return err
}

//删除category 里面的帖子怎么办？
//所以最好不要删除category
func DelCate(cid int) error {
	res, err := db.Exec(
		"call cate_del(?);",
		cid)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	log.Println("aff", rowCnt)
	if err != nil {
		return err
	} else if rowCnt < 1 {
		return ErrNoAff
	}
	return err
}

//获得category
func GetCate(cid int) (*Category, error) {

	cate := &Category{Cid: cid}
	err := db.QueryRow(
		"SELECT  `title`, `description`,`posts`,`todayposts`,`lastpost`,`created` FROM `category` WHERE `cid` = ?",
		cid).Scan(&cate.Title, &cate.Description, &cate.Posts, &cate.ToadyPosts, &cate.LastPost, &cate.Created)
	if err != nil {
		return nil, err
	}
	return cate, err
}

//获得所有category
func GetCates() ([]*Category, error) {
	//查询数据
	rows, err := db.Query("SELECT `title`, `description`,`posts`,`todayposts`,`lastpost`,`created` FROM `category`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cates := make([]*Category, 0)

	for rows.Next() {
		cate := &Category{}
		err = rows.Scan(&cate.Title, &cate.Description, &cate.Posts, &cate.ToadyPosts, &cate.LastPost, &cate.Created)
		if err != nil {
			log.Fatal(err)
			continue
		}
		cates = append(cates, cate)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return cates, err
}

//修改category
func ModifyCate(cid int, name, description string) error {
	res, err := db.Exec(
		"call cate_edit(?,?,?)",
		cid, name, description)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if rowCnt < 1 {
		return ErrNoAff
	}
	return err
}