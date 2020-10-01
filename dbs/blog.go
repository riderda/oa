package dbs
import (
	"github.com/jmoiron/sqlx"
	"github.com/wangbin/jiebago"
	"time"
)

type DbBlog struct{
	Id int `db:"id"`
	Keyword string `db:"keyword"`
	Title string `db:"title"`
	Content string `db:"content"`
	Summary string `db:"summary"`
	Author string `db:"author"`
	Record int `db:"record"`
	PublicStatus string `db:"public status"`
	PublicTime time.Time `db:"public time"`
	IsShow string `db:"is show"`
}

//结巴分词器对象
var Seg jiebago.Segmenter
func init(){
	Seg.LoadDictionary("dict.txt")
}

//分页获取博客
//需要将返回的数据里的content删除。下个版本更新
//已更新 2020-10-01
func GetBlogLimitBy(Keyword string, Limit int, Db *sqlx.DB)([]DbBlog, error){
	list := make([]DbBlog,0)
	rows, err := Db.Query("select id,keyword,title,summary,author,record,`public status`,`public time` from blog where title like ? and `public status`='release' limit ?,10 ","%"+Keyword+"%",Limit*10)
	if err != nil {
		return list,err
	}
	for rows.Next(){
		db := DbBlog{}
		//时间需要处理一下，先保存为string，然后再转time
		var datetime string
		err = rows.Scan(&db.Id,&db.Keyword,&db.Title,&db.Summary,&db.Author,&db.Record,&db.PublicStatus,&datetime)
		db.PublicTime, _ = time.Parse("2006-01-02 15:04:05", datetime)
		if err != nil{
			return list,err
		}
		list = append(list,db)
	}
	return list,nil
}

//返回搜索的博客的总量
func GetLengthOfBlog(Keyword string,Db *sqlx.DB)(int, error){
	var length int
	row := Db.QueryRow("select count(*) from blog where title like ? and `public status`='release'","%"+Keyword+"%")
	err := row.Scan(&length)
	if err != nil {
		return 0,err
	}
	return length,nil
}

//添加文章
func (db *DbBlog)AddBlog(Db *sqlx.DB)(int64,error){
	publicTime := db.PublicTime.Format("2006-01-02 15:04:05")
	result, err :=Db.Exec("insert into blog(keyword,title,content,summary,author,record,`public status`,`public time`) values(?,?,?,?,?,?,?,?)",
		db.Keyword,db.Title,db.Content,db.Summary,db.Author,db.Record,db.PublicStatus,publicTime)
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0,err
	}
	return lastInsertId,nil
}

//更新文章
func (db *DbBlog)UpdateBlog(Db *sqlx.DB)(int64,error){
	publicTime := db.PublicTime.Format("2006-01-02 15:04:05")
	result, err := Db.Exec("update blog set keyword=?, title=?, content=?, summary=?, author=?, `public status`=?, `public time`=? where id=?",
		db.Keyword,db.Title,db.Content,db.Summary,db.Author,db.PublicStatus,publicTime,db.Id)
	if err != nil {
		return 0, err
	}
	lastInsetId, err := result.LastInsertId()
	if err != nil {
		return 0,err
	}
	return lastInsetId,nil
}

func (db *DbBlog)SearchBlog(Db *sqlx.DB)(error){
	var datetime string
	row := Db.QueryRow("select keyword,title,content,summary,author,record,`public status`,`public time` from blog where id=? and `public status`='release'",db.Id)
	//row := Db.QueryRow("select keyword,title,content,summary,author,record,`public status` from blog where id=? and `public status`='release'",db.Id)
	err := row.Scan(&db.Keyword,&db.Title,&db.Content,&db.Summary,&db.Author,&db.Record,&db.PublicStatus,&datetime)
	db.PublicTime, _ = time.Parse("2006-01-02 15:04:05", datetime)
	if err != nil{
		return err
	}
	return nil
}