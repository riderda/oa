package forward

//客户端发送的
type ClientBlog struct{
	Id int `json:"id"`
	Keyword string `json:"keyword"`
	Title string `json:"title"`
	Content string `json:"content"`
	Status string `json:"status"`
	Command string `json:"command"`
}