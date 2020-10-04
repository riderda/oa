package dbs

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type FileControl interface {
	Insert(Db sqlx.DB)(int64,error)
}

type DbFile struct {
	Id int `db:"id""`
	Time time.Time `db:"time"`
	Size float64 `db:"size"`
	Brief string `db:"brief"`
	Name string `db:"name"`
	ImgUrl string `db:"imgurl"`
	Url string `db:"url"`
	Uploader string `db:"uploader"`
}

func (df *DbFile) Insert(Db *sqlx.DB)(int64,error){
	time := df.Time.Format("2006-01-02 15:04:05")
	result, err := Db.Exec("insert into file(time,size,brief,name,imgurl,url,uploader) values(?,?,?,?,?,?,?)",time,df.Size,df.Brief,df.Name,df.ImgUrl,df.Url,df.Uploader)
	if err != nil{
		return 0,nil
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil{
		return 0,nil
	}
	return lastInsertId,nil
}

func GetNameByUrl(Db *sqlx.DB,url string)(string,error){
	result := Db.QueryRow("select name from file where url=?",url)
	var name string
	err := result.Scan(&name)
	if err != nil{
		return "",err
	}
	return name,nil
}