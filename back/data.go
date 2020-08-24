package back

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

