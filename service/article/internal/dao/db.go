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
	_selArticleBaseInfoByAId	  =  "SELECT * FROM `art_baseinfo` WHERE `aid`=?"
	_insertArticle                = "INSERT article (`content`) VALUES (?)"
	_insertArticleInfo            = "INSERT art_baseinfo (`aid`, `author`, `uid`, `title`, `description`, `date`) VALUES (?,?,?,?,?,?)"
	_searchArticleBaseInfoByTitle = "SELECT * FROM `art_baseinfo` WHERE `title`=? LIMIT 0,100"
	_searchArticleBaseInfoByUId	  = "SELECT * FROM `art_baseinfo` WHERE `uid`=? LIMIT 0,100"
	_searchArticleBaseInfoByDate  = "SELECT * FROM `art_baseinfo` WHERE `date`>? OR `date`<?  LIMIT 0,100"
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

func (d *Dao) RawArticle(ctx context.Context, id int64) (art *model.Article, err error) {
	art = new(model.Article)
	err = d.db.QueryRow(ctx, _selArticle, id).Scan(&art.ID, &art.Content)
	if err != nil {
		log.Errorf("RawArticle Failed (%v)", err)
	}
	return
}

func (d *Dao) PostArticle(c context.Context, info *artapi.ArticleBaseInfo, content string) error {
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
	_ = d.cache.Do(c, func(ctx context.Context) {
		_ = d.AddCacheArticle(c, info.Aid, &model.Article{
			ID:      info.Aid,
			Content: content,
		})
		_ = d.AddCacheABI(c, info.Aid, info)

	})
	return nil
}

// add to cache
func (d *Dao) ArticleBaseInfosByTitle(c context.Context, title string) (infos []artapi.ArticleBaseInfo, err error)  {
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
	d.cache.Do(c, func(ctx context.Context) {
		for _, temp := range infos {
			d.AddCacheABI(c, temp.Aid, &temp)
		}
	})
	return
}

func (d *Dao) ArticleBaseInfosByUId(c context.Context, uid int64) (infos []artapi.ArticleBaseInfo, err error) {
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
			// FIXME: processing error?
			log.Errorf("Rows Scan Fail (%v)", err)
			return
		}
		infos = append(infos, temp)
	}
	// FIXME: cache 是否会因为重复的 id 而出错
	d.cache.Do(c, func(ctx context.Context) {
		for _, temp := range infos {
			d.AddCacheABI(c, temp.Aid, &temp)
		}
	})

	return
}

func (d *Dao) RawArticleBaseInfoByAId(c context.Context, aid int64) (info *artapi.ArticleBaseInfo, err error) {
	info = &artapi.ArticleBaseInfo{}
	err = d.db.QueryRow(c, _selArticleBaseInfoByAId, aid).Scan(&info.Aid, &info.Author, &info.Uid, &info.Title, &info.Desc, &info.Date)
	if err != nil {
		info = nil
		log.Errorf("dao.ArticleBaseInfoByAId Failed (%v)", err)
		return
	}
	d.cache.Do(c, func(ctx context.Context) {
		d.AddCacheABI(c, aid, info)
	})
	return
}

func (d *Dao) RawArticleBaseInfoByDate(c context.Context, beg int64, end int64) (infos []artapi.ArticleBaseInfo, err error) {
	var (
		rows *sql.Rows
	)

	rows, err = d.db.Query(c, _searchArticleBaseInfoByDate, end, beg)
	if err != nil {
		log.Errorf("dao.ArticleBaseInfosByUId Failed (%v)", err)
		return
	}

	for rows.Next() {
		temp := artapi.ArticleBaseInfo{}
		err = rows.Scan(&temp.Aid, &temp.Author, &temp.Uid, &temp.Title, &temp.Desc, &temp.Date)
		if err != nil {
			// FIXME: processing error?
			log.Errorf("Rows Scan Fail (%v)", err)
			return
		}
		infos = append(infos, temp)
	}
	// FIXME: cache 是否会因为重复的 id 而出错
	d.cache.Do(c, func(ctx context.Context) {
		for _, temp := range infos {
			d.AddCacheABI(c, temp.Aid, &temp)
		}
	})

	return
}