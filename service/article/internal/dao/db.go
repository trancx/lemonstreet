package dao

import (
	"article/api/artapi"
	"context"
	"github.com/prometheus/common/log"

	"article/internal/model"
	"github.com/bilibili/kratos/pkg/conf/paladin"
	"github.com/bilibili/kratos/pkg/database/sql"
)

const (
	_selArticle                   = "SELECT * FROM article WHERE aid=?"
	_insertArticle                = "INSERT article (`content`) VALUES (?)"
	_insertArticleInfo            = "INSERT art_baseinfo (`aid`, `author`, `uid`, `title`, `description`, `date`) VALUES (?,?,?,?,?,?)"
	_searchArticleBaseInfoByTitle = "SELECT * FROM `art_baseinfo` WHERE `title`=?"
	_searchArticleBaseInfoByUId	  = "SELECT * FROM `art_baseinfo` WHERE `uid`=?"
	)

func NewDB() (db *sql.DB, err error) {
	var cfg struct {
		Client *sql.Config
	}
	if err = paladin.Get("db.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	db = sql.NewMySQL(cfg.Client)
	db.Ping(context.Background())
	return
}

func (d *dao) RawArticle(ctx context.Context, id int64) (art *model.Article, err error) {
	art = new(model.Article)
	err = d.db.QueryRow(ctx, _selArticle, id).Scan(&art.ID, &art.Content)
	if err != nil {
		log.Errorf("RawArticle Failed (%v)", err)
	}
	return
}

func (d *dao) PostArticle(c context.Context, info *artapi.ArticleBaseInfo, content string) error {
	var (
		id int64
	)
	res, err := d.db.Exec(c, _insertArticle, content)
	if err != nil {
		log.Errorf("PostArticle Failed (%v)", err)
		return err
	}
	if id, err = res.LastInsertId(); err != nil {
		log.Errorf("PostArticle LastInsertId Failed (%v)", err)
		return err
	}
	info.Aid = id
	//aid`, `author`, `uid`, `title`, `description`, `date`
	res, err = d.db.Exec(c, _insertArticleInfo, info.Aid, info.Author, info.Uid, info.Title, info.Desc, info.Date)
	if err != nil {
		log.Errorf("PostArticleInfo Failed (%v)", err)
	}
	return err
}

// add to cache
func (d *dao) ArticleBaseInfosByTitle(c context.Context, title string) (infos []artapi.ArticleBaseInfo, err error)  {
	var (
		rows *sql.Rows
	)

	rows, err = d.db.Query(c, _searchArticleBaseInfoByTitle, title)
	if err != nil {
		log.Errorf("dao.ArticleBaseInfosByTitle Failed (%v)", err)
		return
	}

	for rows.Next() {
		temp := artapi.ArticleBaseInfo{}
		err = rows.Scan(&temp.Aid, &temp.Author, &temp.Uid, &temp.Title, &temp.Desc, &temp.Date)
		if err != nil {
			// FIXME: middle error?
			log.Errorf("Rows Scan Fail (%v)", err)
			return
		}
		infos = append(infos, temp)
	}

	return
}

func (d *dao) ArticleBaseInfosByUId(c context.Context, uid int64) (infos []artapi.ArticleBaseInfo, err error) {
	var (
		rows *sql.Rows
	)

	rows, err = d.db.Query(c, _searchArticleBaseInfoByUId, uid)
	if err != nil {
		log.Errorf("dao.ArticleBaseInfosByUId Failed (%v)", err)
		return
	}

	for rows.Next() {
		temp := artapi.ArticleBaseInfo{}
		err = rows.Scan(&temp.Aid, &temp.Author, &temp.Uid, &temp.Title, &temp.Desc, &temp.Date)
		if err != nil {
			// FIXME: middle error?
			log.Errorf("Rows Scan Fail (%v)", err)
			return
		}
		infos = append(infos, temp)
	}

	return
}