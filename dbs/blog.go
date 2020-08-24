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
	Author string `db:"author"`
	Record int `db:"record"`
	PublicStatus string `db:"public status"`
	PublicTime time.Time `db:"public time"`
}

//结巴分词器对象
var Seg jiebago.Segmenter
func init(){
	Seg.LoadDictionary("dict.txt")
}

func GetBlogLimitBy(Keyword string, Limit int, Db *sqlx.DB)([]DbBlog, error){
	list := make([]DbBlog,0)
	rows, err := Db.Query("select * from blog where title like ? limit ?,10","%"+Keyword+"%",Limit*10)
	if err != nil {
		return list,err
	}
	for rows.Next(){
		db := DbBlog{}
		//时间需要处理一下，先保存为string，然后再转time
		var datetime string
		err = rows.Scan(&db.Id,&db.Keyword,&db.Title,&db.Content,&db.Author,&db.Record,&db.PublicStatus,&datetime)
		db.PublicTime, _ = time.Parse("2006-01-02 15:04:05", datetime)
		if err != nil{
			return list,err
		}
		list = append(list,db)
	}
	return list,nil
}

func GetLengthOfBlog(Keyword string,Db *sqlx.DB)(int, error){
	var length int
	row := Db.QueryRow("select count(*) from blog where title like ?","%"+Keyword+"%")
	err := row.Scan(&length)
	if err != nil {
		return 0,err
	}
	return length,nil
}