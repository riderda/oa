package back

import (
	"time"
)


//校验登陆的token
type BackToken struct{
	Token string `json:"token"`
}
//校验登陆状态
type BackInfo struct{
	Sessionid string `json:"sessionid"`
	Status string `json:"status"`
}
//查询数据库的信息
type QueryInfo struct{
	Length int `json:"length"`
	PageNum int `json:"pagenum"`
	Error string `json:"error"`
}
//查询返回的课程结构
type Course struct{
	Id int `json:"id"`
	Type string `json:type`
	Url string `json:"url"`
	Title string `json:"title"`
	Content string `json:"content"`
}
//查询返回的博客结构
type Blog struct{
	Id int `json:"id"`
	Keyword string `json:"keyword"`
	Title string `json:"title"`
	Content string `json:"content"`
	Summary string `json:"summary"`
	Author string `json:"author"`
	Record int `json:"record"`
	PublicStatus string `json:"public status"`
	PublicTime time.Time `json:"public time"`
	IsShow string `json:"isshow"`
}
//集合课程的列表
type CourseList struct{
	LoginInfo BackInfo `json:"logininfo"`
	List []Course `json:"list"`
}

//带查询信息的课程列表
type CourseListLimit struct{
	Info QueryInfo `json:"queryinfo"`
	Courselist CourseList `json:"infolist"`
}

//带查询信息的博客列表
type BlogListLimit struct{
	Info QueryInfo `json:"queryinfo"`
	Keyword string `json:keyword`
	BlogListByKeyword []Blog `json:"bloglist"`
}

//图片上传后返回的url
type ImageUrl struct{
	Url string `json:"url"`
	Status string `json:"status"`
}

//查询上传文件的单个文件信息
type File struct {

}

//博客增改删通用返回数据
type Normal struct{
	BlogId int64 `json:"blogid"`
	Status string `json:"status"`
}