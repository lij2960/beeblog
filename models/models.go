package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/unknwon/com"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	_DB_NAME       = "data/beeblog.db"
	_SQLITE_DRIVER = "sqlite3"
)

type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Category        string
	Labels          string
	Content         string `orm:"size(5000)"`
	Attachment      string
	Created         time.Time `orm:"index"`
	Update          time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time `orm:"index"`
	ReplyCount      int64
	ReplyLastUserId int
}

type Comment struct {
	Id       int64
	Tid      int64
	Nickname string
	Content  string    `orm:"size(1000)"`
	Created  time.Time `orm:"index"`
}

func RegisterDB() {
	if !com.IsExist(_DB_NAME) {
		os.MkdirAll(path.Dir(_DB_NAME), os.ModePerm)
		os.Create(_DB_NAME)
	}
	orm.RegisterModel(new(Category), new(Topic), new(Comment))
	orm.RegisterDriver(_SQLITE_DRIVER, orm.DRSqlite)
	orm.RegisterDataBase("default", _SQLITE_DRIVER, _DB_NAME, 10)
}

func AddReply(tid, nickname, content string) error {
	tidInt, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	comment := &Comment{
		Tid:      tidInt,
		Nickname: nickname,
		Content:  content,
		Created:  time.Now(),
	}
	_, err = o.Insert(comment)
	if err != nil {
		return err
	}
	topic := &Topic{Id: tidInt}
	if o.Read(topic) == nil {
		topic.ReplyTime = time.Now()
		topic.ReplyCount++
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	return err
}

func AddCategory(name string) error {
	o := orm.NewOrm()
	/**
	此处与视频不同，设置了索引，则字段不能为空，不知道是不是版本的原因
	 */
	cate := &Category{
		Title:     name,
		Created:   time.Now(),
		TopicTime: time.Now(),
	}
	qs := o.QueryTable("category")
	err := qs.Filter("title", name).One(cate)
	if err == nil {
		return err
	}
	_, err = o.Insert(cate)
	if err != nil {
		return err
	}
	return nil
}

func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

func GetAllComments(tid string) ([]*Comment, error) {
	tidInt, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	comments := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tidInt).All(&comments)
	return comments, err
}

func GetAllTopics(cate, label string, isDesc bool) ([]*Topic, error) {
	o := orm.NewOrm()
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	var err error
	if isDesc {
		if len(cate) > 0 {
			qs = qs.Filter("category", cate)
		}
		if len(label) > 0 {
			qs = qs.Filter("labels__contains", "$" + label + "#")
		}
		_, err = qs.OrderBy("-created").All(&topics)
	} else {
		_, err = qs.All(&topics)
	}
	return topics, err
}

func DelCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	return err
}

func DeleteTopic(id string) error {
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id: tid}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		_, err = o.Delete(topic)
		if err != nil {
			return err
		}
	}
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err = qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}
	return err
}

func DeleteComment(rid string) error {
	ridInt, err := strconv.ParseInt(rid, 10, 64)
	if err != nil {
		return err
	}
	var tid int64
	o := orm.NewOrm()
	comment := &Comment{Id: ridInt}
	if o.Read(comment) == nil {
		tid = comment.Tid
		_, err = o.Delete(comment)
		if err != nil {
			return err
		}
	}
	comments := make([]*Comment, 0)
	qs := o.QueryTable("comment")
	_, err = qs.Filter("tid", tid).OrderBy("-created").All(&comments)
	if err != nil {
		return nil
	}
	topic := &Topic{Id: tid}
	if o.Read(topic) == nil {
		topic.ReplyTime = comments[0].Created
		topic.ReplyCount = int64(len(comments))
		_, err = o.Update(topic)
	}
	return err
}

func GetTopic(id string) (*Topic, error) {
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	topic := new(Topic)
	qs := o.QueryTable("topic")
	err = qs.Filter("id", tid).One(topic)
	if err != nil {
		return nil, err
	}

	topic.Views++
	_, err = o.Update(topic)
	topic.Labels = strings.Replace(strings.Replace(
		topic.Labels, "#", " ", -1),
		"$", "", -1)
	return topic, err
}

func AddTopic(title, category, label, content string) error {
	//处理标签
	label = "$" + strings.Join(
		strings.Split(label, " "), "#$") + "#"
	o := orm.NewOrm()
	topic := &Topic{
		Title:     title,
		Category:  category,
		Labels:    label,
		Content:   content,
		Created:   time.Now(),
		Update:    time.Now(),
		ReplyTime: time.Now(),
	}
	_, err := o.Insert(topic)
	if err != nil {
		return err
	}
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		//如果不存在，简单的忽略更新操作
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return err
}

func ModifyTopic(id, title, category, label, content string) error {
	label = "$" + strings.Join(
		strings.Split(label, " "), "#$") + "#"
	tid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	var oldCate string
	o := orm.NewOrm()
	topic := &Topic{Id: tid}
	if o.Read(topic) == nil {
		oldCate = topic.Category
		topic.Title = title
		topic.Category = category
		topic.Labels = label
		topic.Content = content
		topic.Update = time.Now()
		_, err = o.Update(topic)
		if err != nil {
			return err
		}
	}
	//更新分类统计
	if len(oldCate) > 0 {
		cate := new(Category)
		qs := o.QueryTable("category")
		err := qs.Filter("title", oldCate).One(cate)
		if err == nil {
			cate.TopicCount--
			_, err = o.Update(cate)
		}
	}
	cate := new(Category)
	qs := o.QueryTable("category")
	err = qs.Filter("title", category).One(cate)
	if err == nil {
		cate.TopicCount++
		_, err = o.Update(cate)
	}
	return nil
}
