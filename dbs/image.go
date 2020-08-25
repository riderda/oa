package dbs

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type ImageControl interface {
	Insert(Db *sqlx.DB)(int64,error)

}
type DbImage struct{
	Id int `db:"id"`
	Name string `db:"name"`
	UploadTime time.Time `db:"uploadtime"`
	Path string `db:"path"`
	UploadAuthor string `db:"uploadauthor"`
}

func (dbi *DbImage)Insert(Db *sqlx.DB)(int64,error){
	uploadTime := dbi.UploadTime.Format("2006-01-02 15:04:05")
	result, err :=Db.Exec("insert into image(name,`upload time`,path,`upload author`) values(?,?,?,?)",dbi.Name,uploadTime,dbi.Path,dbi.UploadAuthor)
	if err != nil{
		return 0,err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0,nil
	}
	return affected,nil
}

