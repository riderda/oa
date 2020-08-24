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

func GetCourseLimitBy(CourseType string,Limit int,Db *sqlx.DB)([]DbCourse,error){
	list := make([]DbCourse,0)
	rows, err := Db.Query("select * from course where type=? order by id limit ?,10",CourseType,Limit*10)
	if err != nil {
		return list, err
	}

	for rows.Next(){
		dc := DbCourse{}
		err = rows.Scan(&dc.Id,&dc.Type,&dc.Url,&dc.Title,&dc.Content)
		if err != nil {
			return list, err
		}
		list = append(list,dc)
	}
	return list,nil
}

//获取查询类型的文章的数量，但其实并没有太大必要
//后来想想，还是有必要的，可以减少main的篇幅，尤其是对一些结构比较离谱进行长度计算时
func GetLengthOfCourse(CourseType string,Db *sqlx.DB)(int, error){
	var length int
	row := Db.QueryRow("select count(*) from course where type=?",CourseType)

	if err := row.Scan(&length); err != nil{
		return 0,err
	}
	return length,nil
}