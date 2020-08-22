package dbs
import (
	"github.com/jmoiron/sqlx"
)
type CourseControl interface {
	//GetCourseList(Db *sqlx.DB)(error)

}
type DbCourse struct{
	Id int `db:"id"`
	Type string `db:type`
	Url string `db:"url"`
	Title string `db:"title"`
	Content string `db:"content"`
}

/*
	@param courseType 课程类型
	@param Db 		  数据库连接
	@return 课程列表，错误
	根据类型查找，返回指定类型的课程数据
 */
func GetCourseListBy(courseType string,Db *sqlx.DB)([]DbCourse,error){
	list := make([]DbCourse,0)
	rows, err := Db.Query("select * from course where type=?",courseType)
	if err != nil {
		return list,err
	}

	for rows.Next(){
		dc := DbCourse{}
		err = rows.Scan(&dc.Id,&dc.Type,&dc.Url,&dc.Title,&dc.Content)
		if err != nil {
			return list,err
		}
		list = append(list,dc)
	}
	return list,nil
}
