package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/unknwon/com"
	"os"
	"path"
	"time"
	_ "github.com/mattn/go-sqlite3"
)

const (
	_DB_NAME = "data/beeblog.db"
	_SQLITE_DRIVER = "sqlite3"
)

type Category struct {
	Id int64
	Title string
	Created time.Time `orm:"index"`
	Views int64 `orm:"index"`
	TopicTime time.Time `orm:"index"`
	TopicCount int64
	TopicLastUserId int64
}

type Topic struct {
	Id int64
	Uid int64
	Title string
	Content string `orm:"size(5000)"`
	Attachment string
	Created time.Time `orm:"index"`
	Update time.Time `orm:"index"`
	Views int64 `orm:"index"`
	Author string
	ReplyTime time.Time `orm:"index"`
	ReplyCount int64
	ReplyLastUserId int
}

func RegisterDB()  {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	orm.RegisterModel(new(Category), new(Topic))
	orm.RegisterDriver(_SQLITE_DRIVER, orm.DRSqlite)
	orm.RegisterDataBase("default", _SQLITE_DRIVER, _DB_NAME, 10)
}