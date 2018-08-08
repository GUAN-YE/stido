package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/duguying/blog/utils"
	// "github.com/gogather/com/log"
	"strconv"
	"strings"
	"time"
)

type Article struct {
	Id       int
	Title    string
	Uri      string
	Keywords string
	Abstract string
	Content  string
	Author   string
	Time     time.Time `orm:"column(created_at)"`
	Count    int
	Status   int
}

const (
	ART_STATUS_DRAFT   = 0
	ART_STATUS_PUBLISH = 1
)

func (this *Article) TableName() string {
	return "articles"
}

func init() {
	orm.RegisterModel(new(Article))
}

// 添加文章
func AddArticle(title string, content string, keywords string, abstract string, status int, author string) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")

	sql := "insert into articles(title, uri, keywords, abstract, content, author, status) values(?, ?, ?, ?, ?, ?, ?)"
	res, err := o.Raw(sql, title, strings.Replace(title, "/", "-", -1), keywords, abstract, content, author, status).Exec()
	if nil != err {
		return 0, err
	} else {
		return res.LastInsertId()
	}
}

// 通过id获取文章-cached
func GetArticle(id int) (Article, error) {
	var err error
	var art Article

	err = utils.GetCache("GetArticle.id."+fmt.Sprintf("%d", id), &art)
	if err != nil {
		o := orm.NewOrm()
		o.Using("default")
		art = Article{Id: id}
		err = o.Read(&art, "id")
		utils.SetCache("GetArticle.id."+fmt.Sprintf("%d", id), art, 600)
	}

	return art, err
}

// 通过uri获取文章-cached
func GetArticleByUri(uri string) (Article, error) {
	var err error
	var art Article

	err = utils.GetCache("GetArticleByUri.uri."+uri, &art)
	if err == nil {
		// get view count
		count, err := GetArticleViewCount(art.Id)
		if err == nil {
			art.Count = int(count)
		}

		return art, nil
	} else {
		o := orm.NewOrm()
		o.Using("default")
		art = Article{Uri: uri}
		err = o.Read(&art, "uri")
		utils.SetCache("GetArticleByUri.uri."+uri, art, 600)
	}

	return art, err
}

// 通过文章标题获取文章-cached
func GetArticleByTitle(title string) (Article, error) {
	var err error
	var art Article

	err = utils.GetCache("GetArticleByTitle.title."+title, &art)
	if err == nil {
		// get view count
		count, err := GetArticleViewCount(art.Id)
		if err == nil {
			art.Count = int(count)
		}

		return art, nil
	} else {
		o := orm.NewOrm()
		o.Using("default")
		art = Article{Title: title}
		err = o.Read(&art, "title")
		utils.SetCache("GetArticleByTitle.title."+title, art, 600)
	}

	return art, err
}

// 获取文章浏览量
func GetArticleViewCount(id int) (int, error) {
	var maps []orm.Params

	sql := `select count from articles where id=?`
	o := orm.NewOrm()
	num, err := o.Raw(sql, id).Values(&maps)
	if err == nil && num > 0 {
		count := maps[0]["count"].(string)

		return strconv.Atoi(count)
	} else {
		return 0, err
	}
}

// 更新阅览数统计
func UpdateCount(id int) error {
	o := orm.NewOrm()
	o.Using("default")
	art := Article{Id: id}
	err := o.Read(&art)

	o.QueryTable("articles").Filter("id", id).Update(orm.Params{
		"count": art.Count + 1,
	})

	return err
}

// 更新文章
func UpdateArticle(id int64, uri string, newArt Article) error {
	o := orm.NewOrm()
	o.Using("default")
	var art Article

	if 0 != id {
		art = Article{Id: int(id)}
	} else if "" != uri {
		art = Article{Uri: uri}
	}

	art.Title = newArt.Title
	art.Keywords = newArt.Keywords
	art.Abstract = newArt.Abstract
	art.Content = newArt.Content
	art.Status = newArt.Status

	getArt, _ := GetArticle(int(id))
	utils.DelCache("GetArticleByUri.uri." + getArt.Uri)
	utils.DelCache("GetArticle.id." + fmt.Sprintf("%d", art.Id))

	_, err := o.Update(&art, "title", "keywords", "abstract", "content", "status")
	return err
}

// 通过uri删除文章
func DeleteArticle(id int64, uri string) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")
	var art Article

	if 0 != id {
		art.Id = int(id)
	} else if "" != uri {
		art.Uri = uri
	}

	getArt, _ := GetArticle(int(id))
	utils.DelCache("GetArticleByUri.uri." + getArt.Uri)
	utils.DelCache("GetArticle.id." + fmt.Sprintf("%d", art.Id))

	return o.Delete(&art)
}

// 按月份统计文章数-cached
// select DATE_FORMAT(created_at,'%Y年%m月') as date,count(*) as number ,year(created_at) as year, month(created_at) as month from article group by date order by year desc, month desc
func CountByMonth() ([]orm.Params, error) {
	var maps []orm.Params

	err := utils.GetCache("CountByMonth", &maps)
	if nil != err {
		sql := "select DATE_FORMAT(created_at,'%Y年%m月') as date,count(*) as number ,year(created_at) as year, month(created_at) as month from articles where status=? group by date order by year desc, month desc"
		o := orm.NewOrm()
		num, err := o.Raw(sql, ART_STATUS_PUBLISH).Values(&maps)
		if err == nil && num > 0 {
			utils.SetCache("CountByMonth", maps, 3600)
			return maps, nil
		} else {
			return nil, err
		}
	} else {
		return maps, err
	}

}

// 获取某月的文章列表-cached
// select * from article where year(created_at)=2014 and month(created_at)=8
// year 年
// month 月
// page 页码
// numPerPage 每页条数
// 返回值:
// []orm.Params 文章
// bool 是否有下一页
// int 总页数
// error 错误
func ListByMonth(year int, month int, page int, numPerPage int) ([]orm.Params, bool, int, error) {
	if year < 0 {
		year = 1970
	}

	if month < 0 || month > 12 {
		month = 1
	}

	if page < 1 {
		page = 1
	}

	if numPerPage < 1 {
		numPerPage = 10
	}

	var maps, maps2 []orm.Params
	o := orm.NewOrm()
	var err error

	// get data - cached
	err = utils.GetCache(fmt.Sprintf("ListByMonth.list.%d.%d.%d", year, month, page), &maps)
	if nil != err {
		sql1 := "select * from articles where status=? and year(created_at)=? and month(created_at)=? order by created_at desc limit ?,?"
		_, err = o.Raw(sql1, ART_STATUS_PUBLISH, year, month, numPerPage*(page-1), numPerPage).Values(&maps)
		utils.SetCache(fmt.Sprintf("ListByMonth.list.%d.%d.%d", year, month, page), maps, 3600)
	}

	err = utils.GetCache(fmt.Sprintf("ListByMonth.count.%d.%d", year, month), &maps2)
	if nil != err {
		sql2 := "select count(*)as number from articles where status=? and year(created_at)=? and month(created_at)=?"
		_, err = o.Raw(sql2, ART_STATUS_PUBLISH, year, month).Values(&maps2)
		utils.SetCache(fmt.Sprintf("ListByMonth.count.%d.%d", year, month), maps2, 3600)
	}

	// calculate pages
	number, _ := strconv.Atoi(maps2[0]["number"].(string))
	var addFlag int
	if 0 == (number % numPerPage) {
		addFlag = 0
	} else {
		addFlag = 1
	}
	pages := number/numPerPage + addFlag

	var flagNextPage bool
	if pages == page {
		flagNextPage = false
	} else {
		flagNextPage = true
	}

	if err == nil {
		return maps, flagNextPage, pages, nil
	} else {
		return nil, false, pages, err
	}

}

// 文章分页列表
// select * from article order by created_at desc limit 0,6
// page 页码
// numPerPage 每页条数
// 返回值:
// []orm.Params 文章
// bool 是否有下一页
// int 总页数
// error 错误
func ListPage(page int, numPerPage int) ([]orm.Params, bool, int, error) {
	// pagePerNum := 6
	sql1 := "select * from articles where status = ? order by created_at desc limit ?," + fmt.Sprintf("%d", numPerPage)
	sql2 := "select count(*) as number from articles where status = ?"
	var maps, maps2 []orm.Params
	o := orm.NewOrm()
	num, err := o.Raw(sql1, ART_STATUS_PUBLISH, numPerPage*(page-1)).Values(&maps)
	if err != nil {
		fmt.Println("execute sql1 error:")
		fmt.Println(err)
		return nil, false, 0, err
	}

	err = utils.GetCache("ArticleNumber", &maps2)
	if nil != err {
		_, err = o.Raw(sql2, ART_STATUS_PUBLISH).Values(&maps2)
		if err != nil {
			fmt.Println("execute sql2 error:")
			fmt.Println(err)
			return nil, false, 0, err
		}
		utils.SetCache("ArticleNumber", maps2, 3600)
	}

	number, err := strconv.Atoi(maps2[0]["number"].(string))

	var addFlag int
	if 0 == (number % numPerPage) {
		addFlag = 0
	} else {
		addFlag = 1
	}

	pages := number/numPerPage + addFlag

	var flagNextPage bool
	if pages == page {
		flagNextPage = false
	} else {
		flagNextPage = true
	}

	if err == nil && num > 0 {
		return maps, flagNextPage, pages, nil
	} else {
		return nil, false, pages, err
	}
}

// 同关键词文章列表
// select * from article where keywords like '%keyword%'
// 返回值:
// []orm.Params 文章
// bool 是否有下一页
// error 错误
func ListByKeyword(keyword string, page int, numPerPage int) ([]orm.Params, bool, int, error) {
	// numPerPage := 6
	sql1 := "select * from articles where keywords like ? order by created_at desc limit ?,?"
	sql2 := "select count(*) as number from articles where keywords like ?"
	var maps, maps2 []orm.Params
	o := orm.NewOrm()
	num, err := o.Raw(sql1, fmt.Sprintf("%%%s%%", keyword), numPerPage*(page-1), numPerPage).Values(&maps)
	o.Raw(sql2, fmt.Sprintf("%%%s%%", keyword)).Values(&maps2)

	number, _ := strconv.Atoi(maps2[0]["number"].(string))

	var addFlag int
	if 0 == (number % numPerPage) {
		addFlag = 0
	} else {
		addFlag = 1
	}

	pages := number/numPerPage + addFlag

	var flagNextPage bool
	if pages == page {
		flagNextPage = false
	} else {
		flagNextPage = true
	}

	if err == nil && num > 0 {
		return maps, flagNextPage, pages, nil
	} else {
		return nil, false, pages, err
	}
}

// 最热文章列表 - cached
// select * from article order by count desc limit 10
func HottestArticleList() ([]orm.Params, error) {
	var maps []orm.Params

	// get data - cached
	err := utils.GetCache("HottestArticleList", &maps)
	if nil != err {
		sql := "select id,uri,title,count from articles order by count desc limit 20"
		o := orm.NewOrm()
		_, err = o.Raw(sql).Values(&maps)

		utils.SetCache("HottestArticleList", maps, 3600)
	}

	return maps, err
}

// 列出文章 for admin
func ArticleListForAdmin(page int, numPerPage int) ([]orm.Params, bool, int, error) {
	sql1 := "select id,uri,title,count,status,created_at from articles order by created_at desc limit ?," + fmt.Sprintf("%d", numPerPage)
	sql2 := "select count(*) as number from articles"
	var maps, maps2 []orm.Params
	o := orm.NewOrm()
	num, err := o.Raw(sql1, numPerPage*(page-1)).Values(&maps)
	if err != nil {
		fmt.Println("execute sql1 error:")
		fmt.Println(err)
		return nil, false, 0, err
	}

	err = utils.GetCache("ArticleNumber", &maps2)
	if nil != err {
		_, err = o.Raw(sql2).Values(&maps2)
		if err != nil {
			fmt.Println("execute sql2 error:")
			fmt.Println(err)
			return nil, false, 0, err
		}
		utils.SetCache("ArticleNumber", maps2, 3600)
	}

	number, err := strconv.Atoi(maps2[0]["number"].(string))

	var addFlag int
	if 0 == (number % numPerPage) {
		addFlag = 0
	} else {
		addFlag = 1
	}

	pages := number/numPerPage + addFlag

	var flagNextPage bool
	if pages == page {
		flagNextPage = false
	} else {
		flagNextPage = true
	}

	if err == nil && num > 0 {
		return maps, flagNextPage, pages, nil
	} else {
		return nil, false, pages, err
	}
}
