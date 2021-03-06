package repository

import (
	"database/sql"
	"goBlog/logger"
	"time"
)

//单一post
type Post struct {
	Tid       int //id
	Cid       int //category id
	Uid       int //user id
	Title     string
	Content   string
	Tags      string
	Type      int
	Status    int
	Views     int
	Replys    int
	Created   time.Time
	Updated   time.Time
	Lastreply time.Time
}

//带回复
type Article struct {
	Post     *Post
	Comments []*Comment
}

//发布主题
func AddPost(cid, uid int, title, content string) (int, error) {
	s := "insert into post (cid,uid,title,content) VALUES " +
		"($1,$2,$3,$4) RETURNING id"
	return add(s, cid, uid, title, content)
}

//删除主题
func DelPost(tid int) (int64, error) {
	s := "delete from post where tid = $1"
	return del(s, tid)
}

//编辑文章
func EditPost(tid int, title, content string) (int64, error) {
	s := "UPDATE post SET title = $1, content = $2, updated = now(), lastreply = updated WHERE id = $3"
	return update(s, title, content, tid)
}

//#0-正常，1-不可回复2不可查看
func setPostStatus(tid, status int) (int64, error) {
	if tid <= 0 || status < 0 || status > 2 {
		return -1, ErrParama
	}
	s := " UPDATE post SET status = $1 WHERE id = $2"
	return update(s, status, tid)
}

//查看一篇文章status
func getPostStatus(tid int) (int, error) {
	status := 2
	err := db.QueryRow("SELECT  status FROM post WHERE tid = $1", tid).Scan(&status)
	return status, err
}

//根据tid获取单篇文章
func GetPost(tid int) (*Post, error) {
	s := "SELECT cid,uid,title,content,tags,type,status," +
		"views,replys,created,updated,lastreply FROM post WHERE tid = $1"

	p := &Post{Tid: tid}
	err := db.QueryRow(s, tid).Scan(&p.Cid, &p.Uid, &p.Title, &p.Content, &p.Tags, &p.Type, &p.Status,
		&p.Views, &p.Replys, &p.Created, &p.Updated, &p.Lastreply)
	return p, err
}

//获得文章带回复
func GetArticle(tid int) (a *Article, err error) {
	c := make(chan error, 1)
	var p *Post

	//开一个线程去获取文章
	go func() {
		p, err = GetPost(tid)
		c <- err
	}()

	//主线程获取评论
	var commnets = make([]*Comment, 0, 30)
	if commnets, err = GetComments(tid, 1, 30); err != nil {
		return
	}

	//等待线程结束
	if err = <-c; err != nil {
		return
	}

	a = &Article{Post: p, Comments: commnets}
	return a, nil
}

//获得某一cid的文章列表
//如果cid<0 则表示不分区
//只获取文章前一部分(120)
//created hot3 hot7 其余是按照最后回复排序
func getPostsList(cid, page, pagesize int, order string) (posts []*Post, err error) {
	//查询数据
	var where string
	switch order {
	case "created":
		where = "ORDER BY created DESC"
	case "hot3":
		//最近3天的热帖
		where = "AND DATEDIFF(NOW(),lastreply)< 3 ORDER BY replys DESC,lastreply DESC"
	case "hot7":
		//最近7天的热帖
		where = "AND DATEDIFF(NOW(),lastreply)< 7 ORDER BY replys DESC,lastreply DESC"
	default:
		//新帖
		where = "ORDER BY lastreply DESC"
	}

	offset := (page - 1) * pagesize

	var whereCid string
	if cid < 0 {
		whereCid = "where $1 = $1 "
	} else {
		whereCid = "where cid = $1 "
	}

	s := "SELECT id,uid,title,left(content,120)," +
		"tags,type,status,views," +
		"replys,created,updated,lastreply" +
		" FROM post " + whereCid + where + " LIMIT $2 OFFSET $3"

	var rows *sql.Rows
	logger.D("get posts sql: %s", s)
	if rows, err = db.Query(s, -1, pagesize, offset); err != nil {
		return nil, err
	}
	defer rows.Close()

	posts = make([]*Post, 0, pagesize)

	for rows.Next() {
		p := &Post{Cid: cid}
		if err = rows.Scan(
			&p.Tid, &p.Uid, &p.Title, &p.Content,
			&p.Tags, &p.Type, &p.Status, &p.Views,
			&p.Replys, &p.Created, &p.Updated, &p.Lastreply); err != nil {
			return
		}
		posts = append(posts, p)
	}

	err = rows.Err()
	return
}

// 按发帖时间排序
func GetPostListOrderCreated(page, pagesize int) ([]*Post, error) {
	return getPostsList(-1, page, pagesize, "created")
}

func GetPostListByCidOrderCreated(cid, page, pagesize int) ([]*Post, error) {
	return getPostsList(cid, page, pagesize, "created")
}

//最近3天热贴
func GetPostListHot3(page, pagesize int) ([]*Post, error) {
	return getPostsList(-1, page, pagesize, "hot3")
}

func GetPostListByCidHot3(cid, page, pagesize int) ([]*Post, error) {
	return getPostsList(cid, page, pagesize, "hot3")
}

//最近7天热贴
func GetPostListHot7(page, pagesize int) ([]*Post, error) {
	return getPostsList(-1, page, pagesize, "hot7")
}

func GetPostListByCidHot7(cid, page, pagesize int) ([]*Post, error) {
	return getPostsList(cid, page, pagesize, "hot7")
}

//最后回复
func GetPostListOrderLastReply(page, pagesize int) ([]*Post, error) {
	return getPostsList(-1, page, pagesize, "lastreply")
}

func GetPostListByCidOrderLastReply(cid, page, pagesize int) ([]*Post, error) {
	return getPostsList(cid, page, pagesize, "lastreply")
}
