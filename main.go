package main

import (
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"oa/Helper"
	"oa/back"
	"oa/dbs"
	forward2 "oa/forward"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var sessionMgr *Helper.SessionMgr = nil

//登陆检查
func checkLogin(w http.ResponseWriter, r *http.Request)(string,error){

	if r.Method == "GET"{
		return "",errors.New("method must is post")
	}else if r.Method == "POST"{
		//获取cookie里面的值
		sessionID := r.Header.Get("authentication")
		//判断是否存在Cookiename
		/**
			dasdasd,Cookiename=asdasdasdasdasd,dasdasdasd 大概是这个样子
			.*[,]{1}Cookiename=([^,]+).*
		 **/
		//不存在检验字段的时候
		if len(sessionID) == 0 {
			return sessionID,errors.New("no found sessionid")
		}

		if _, exist := sessionMgr.GetSessionVal(sessionID,"user");exist{
			return sessionID,nil
		}else{
			return sessionID,errors.New("is no login")
		}
	}
	return "",errors.New("other case")
}

//获取token
func index(w http.ResponseWriter, r *http.Request){

	back := back.BackToken{}
	sessionID := sessionMgr.StartSession(w,r)
	sessionMgr.SetSessionVal(sessionID, "token", sessionMgr.NewToken())
	token, _ := sessionMgr.GetSessionVal(sessionID,"token")
	back.Token = token.(string)

	w.Header().Set("Content-type","application/json")
	w.Header().Add("authentication",sessionID)
	json, err := json.MarshalIndent(&back,"","\t")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Write(json)

}

//登陆检验
func login(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET" {
		w.WriteHeader(201)
		return
	} else if r.Method == "POST" {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		w.Header().Set("Content-Type","application/json")
		r.ParseForm()
		//可以使用template.HTMLEscapeString()来避免用户进行js注入
		user := dbs.User{
			Id: 0,
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}
		fmt.Println(user.Username," ",user.Password)
		//获取Authentication的值
		sessionID := r.Header.Get("authentication")

		//如果不存在sessionid，则赋予新的sesionid，但结果肯定是不匹配
		if len(sessionID) == 0 {
			sessionID = sessionMgr.StartSession(w,r)
		}


		//对比token值
		token, _ := sessionMgr.GetSessionVal(sessionID,"token")
		if r.FormValue("token") != token{
			w.WriteHeader(201)
			//更新token
			sessionMgr.SetSessionVal(sessionID,"token",sessionMgr.NewToken())
			return
		}

		//在数据库中得到对应数据
		Db, err := dbs.GetConnect()
		defer Db.Close()
		if err != nil {
			panic("error")
		}
		//登陆验证
		err = user.FindUser(Db)

		//如果登陆失败，返回失败的信息
		if err != nil {
			backinfo := back.BackInfo{
				Sessionid: sessionID,
				Status: "fail",
			}
			json, _ := json.MarshalIndent(&backinfo,"","\t")
			w.Write(json)

			//更新token
			sessionMgr.SetSessionVal(sessionID,"token",sessionMgr.NewToken())
			return
		}

		//登陆成功后
		//更新token
		sessionMgr.SetSessionVal(sessionID,"token",sessionMgr.NewToken())

		//踢除重复登录的
		var onlineSessionIDList = sessionMgr.GetSessionIDList()

		for _, onlineSessionID := range onlineSessionIDList {
			if userInfo, ok := sessionMgr.GetSessionVal(onlineSessionID, "user"); ok {
				if value, ok := userInfo.(dbs.User); ok {
					if value.Id == user.Id {
							sessionMgr.EndSessionBy(onlineSessionID)
					}
				}
			}
		}

		//设置变量值，用该sessionid保存登陆的账号密码的信息
		sessionMgr.SetSessionVal(sessionID, "user", user)

		backinfo := back.BackInfo{
			Sessionid: sessionID,
			Status: "success",
		}
		json, _ := json.MarshalIndent(&backinfo,"","\t")
		w.Write(json)
	}
}

//上传图片
func uploadImage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")

	//登陆检测
	_, err := checkLogin(w,r)
	if err != nil{
		data := back.ImageUrl{
			Status: "not login sessionID",
		}
		json, _ := json.MarshalIndent(&data,"","\t")
		w.Write(json)
		return
	}


	//解析请求
	file,fileHeader,err := r.FormFile("image")
	if err != nil{
		log.Println("upload happen error ",err.Error())
		w.WriteHeader(500)
		return
	}
	defer file.Close()

	//数据库连接
	Db, err := dbs.GetConnect()
	if err != nil {
		log.Println("database is error ",err.Error())
		w.WriteHeader(500)
		return
	}
	defer Db.Close()

	/**应该多步进行文件校验
	1、文件后缀名判断（这个很蠢，但还是需要）
	2、文件类型判断（从请求里可以获取，但是可能被篡改）
	3、重新渲染，文件名和文件名后缀由服务器来决定
	**/

	//1、后缀名获取并判断
	//白名单

	followExt := []string{"jpg","png","jpeg"}
	reg := regexp.MustCompile(`.*[\.]{1}([^\.]+)`)

	ext := reg.FindStringSubmatch(fileHeader.Filename)
	if len(ext) != 2{
		log.Println("ext is error ", ext)
		w.WriteHeader(500)
		return
	}
	//后缀名为ext[1]
	in_follow := false
	for _, follow := range followExt{
		if follow == ext[1]{
			in_follow = true
		}
	}
	if in_follow == false {
		log.Println("ext is error ", ext)
		w.WriteHeader(500)
		return
	}

	//2、文件类型判断，从content-type里获取
	contentType := fileHeader.Header["Content-Type"][0]
	if contentType != "image/jpg" && contentType != "image/png" && contentType != "image/jpeg"{
		log.Println("contentType is error ", contentType)
		w.WriteHeader(500)
		return
	}

	//3、重新渲染图片
	//检查路径是否存在，没有则创建
	_,err = os.Stat("./upload/image/")
	if err != nil {
		os.MkdirAll("./upload/image",os.ModePerm)
	}

	//重新生成文件名
	uname := uuid.Must(uuid.NewV4())
	f, err := os.OpenFile("./upload/image/"+uname.String()+"."+ext[1],os.O_WRONLY|os.O_CREATE,0666)
	if err != nil{
		fmt.Println(err.Error())
		w.WriteHeader(500)
		return
	}
	defer f.Close()
	io.Copy(f,file)

	//数据库操作
	image := dbs.DbImage{
		Name: fileHeader.Filename,
		UploadTime: time.Now(),
		Path: "/upload/image/"+uname.String()+"."+ext[1],
		UploadAuthor: "",
	}

	_, err = image.Insert(Db)
	if err != nil{
		log.Println("insert error ",err.Error())
		w.WriteHeader(500)
		return
	}

	//返回上传图片的url
	data := back.ImageUrl{
		Url: "/static/"+uname.String()+"."+ext[1],
		Status: "OK",
	}
	json, _ := json.MarshalIndent(&data,"","\t")
	w.Write(json)
}

//获取文章
func getCourseLimit(w http.ResponseWriter, r *http.Request){
	data := back.CourseListLimit{}
	w.Header().Set("Content-Type","application/json")

	sessionID, err := checkLogin(w,r)
	//如果没有登陆
	if err != nil{
		info := back.QueryInfo{
			Length: 0,
			Error: "not login sessionID",
			PageNum: 0,
		}
		data.Info = info
		data.Courselist.LoginInfo = back.BackInfo{
			Sessionid: sessionID,
			Status: "fail",
		}
		json, _ := json.MarshalIndent(&data,"","\t")
		w.Write(json)
		return
	}




	//解析请求
	r.ParseForm()
	CourseType := r.FormValue("type")
	TempLimit := r.FormValue("limit")

	//请求的参数错误
	Limit, err := strconv.Atoi(TempLimit)
	if err != nil{
		data.Info = back.QueryInfo{
			Length: 0,
			Error: "limit is not number",
			PageNum: 0,
		}
		data.Courselist.LoginInfo = back.BackInfo{
			Sessionid: sessionID,
			Status: "success",
		}
		json, _ := json.MarshalIndent(&data,"","\t")
		w.Write(json)
		return
	}
	if len(CourseType) == 0{
		data.Info = back.QueryInfo{
			Length: 0,
			Error: "must send the type",
			PageNum: 0,
		}
		data.Courselist.LoginInfo = back.BackInfo{
			Sessionid: sessionID,
			Status: "success",
		}
		json, _ := json.MarshalIndent(&data,"","\t")
		w.Write(json)
		return
	}

	//数据库错误
	Db, err := dbs.GetConnect()
	if err != nil {
		log.Println("database is error",err.Error())
		w.WriteHeader(500)
		return
	}
	defer Db.Close()

	//查询错误（内部错误）
	data.Info.Length,err = dbs.GetLengthOfCourse(CourseType,Db)
	if err != nil {
		log.Println("query is error",err.Error())
		w.WriteHeader(500)
		return
	}



	//能到这里就意味着，登陆了，数据库正常，请求参数正确
	//开始查询数据
	list, err := dbs.GetCourseLimitBy(CourseType,Limit,Db)
	if err != nil{
		log.Println("query is error",err.Error())
		w.WriteHeader(500)
		return
	}
	for _, course := range list{
		temp := back.Course{}
		temp.Id, temp.Type, temp.Url, temp.Title, temp.Content = course.Id, course.Type, course.Url, course.Title, course.Content
		data.Courselist.List = append(data.Courselist.List,temp)
	}
	data.Courselist.LoginInfo.Status = "success"
	data.Courselist.LoginInfo.Sessionid = sessionID
	data.Info.PageNum = Limit
	json, _ := json.MarshalIndent(&data,"","\t")
	w.Write(json)
}

//搜索博客
func getBlogLimit(w http.ResponseWriter, r *http.Request){
	data := back.BlogListLimit{}
	w.Header().Set("Content-Type","application/json")

	//先检查登陆
	_, err := checkLogin(w,r)
	//没有登陆
	if err != nil {
		data.Info = back.QueryInfo{
			Length: 0,
			PageNum: 0,
			Error: "not login sessionID",
		}
		json, _ := json.MarshalIndent(&data,"","\t")
		w.Write(json)
		return
	}

	//解析请求参数
	r.ParseForm()

	Search := r.FormValue("search")
	TempLimit := r.FormValue("limit")

	//参数错误
	if len(Search) == 0 {
		data.Info = back.QueryInfo{
			Length: 0,
			PageNum: 0,
			Error: "search is null",
		}
		json, _ := json.MarshalIndent(&data,"","\t")
		w.Write(json)
		return
	}
	Limit,err := strconv.Atoi(TempLimit)
	if err != nil{
		data.Info = back.QueryInfo{
			Length: 0,
			PageNum: 0,
			Error: "Limit is not number",
		}
		json, _ := json.MarshalIndent(&data,"","\t")
		w.Write(json)
		return
	}

	//数据库
	Db, err := dbs.GetConnect()
	if err != nil {
		log.Println("database is error ",err.Error())
		w.WriteHeader(500)
		return
	}
	defer Db.Close()

	var wg sync.WaitGroup
	data.BlogListByKeyword = make([]back.Blog,0)
	//从搜索里获取关键词
	ch := dbs.Seg.CutForSearch(Search,true)
	for keyword := range ch{
		data.Keyword += keyword+";"
		wg.Add(2)
		//估计查询量比较大，使用并发查找
		go func(){
			//长度更新
			defer wg.Done()
			templength, err := dbs.GetLengthOfBlog(keyword,Db)
			if err != nil{
				log.Println("query is error",err.Error())
				w.WriteHeader(500)
				return
			}
			data.Info.Length += templength
		}()

		go func(){
			//内容更新
			defer wg.Done()
			dbBlogList, err := dbs.GetBlogLimitBy(keyword,Limit,Db)
			if err != nil{
				log.Println("query is error",err.Error())
				w.WriteHeader(500)
				return
			}

			for _, blog := range dbBlogList{
				temp := back.Blog{
					Id: blog.Id,
					Keyword: blog.Keyword,
					Title: blog.Title,
					Content: blog.Content,
					Summary: blog.Summary,
					Author: blog.Author,
					Record: blog.Record,
					PublicStatus: blog.PublicStatus,
					PublicTime: blog.PublicTime,
					IsShow: blog.IsShow,
				}
				data.BlogListByKeyword = append(data.BlogListByKeyword,temp)
			}

		}()
	}
	//只能阻塞，但不能影响write
	wg.Wait()
	data.Info.PageNum = Limit
	json, _ := json.MarshalIndent(&data,"","\t")
	w.Write(json)
}

//博客的增改
func execBlog(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	data := back.Normal{}
	var err error
	//验证
	sessionID, err := checkLogin(w,r)
	if err != nil{
		data.Status = "not login sessionID"
		json, _ := json.MarshalIndent(&data,"","\t")
		w.Write(json)
		return
	}

	//解析参数
	body,err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("read error ",err.Error())
		w.WriteHeader(500)
		return
	}
	forward := forward2.ClientBlog{}
	if err = json.Unmarshal(body,&forward);err != nil{
		log.Println("unmarshal error ",err.Error())
		w.WriteHeader(500)
		return
	}

	//获取命令参数
	if len(forward.Command) == 0 || (forward.Command != "add" && forward.Command != "update"){
		data.Status = "command must is add or update"
		json, _ := json.MarshalIndent(&data,"","\t")
		w.Write(json)
		return
	}

	//数据库连接
	Db,err := dbs.GetConnect()
	if err != nil{
		log.Println("database is error ",err.Error())
		w.WriteHeader(500)
		return
	}
	defer Db.Close()

	//author需要获取sessionid，然后根据里面user的id值找到userinfo表里fullname
	tempUser, _ := sessionMgr.GetSessionVal(sessionID,"user")
	user := tempUser.(dbs.User)
	fullName, err := dbs.GetFullName(user.Id,Db)
	if err != nil{
		log.Println("query error ",err.Error())
		w.WriteHeader(500)
		return
	}


	//确保status只有发布和未发布两个选项
	if forward.Status != "release" {
		forward.Status = ""
	}
	//初始化数据，add和update共用的数据
	db := dbs.DbBlog{
		Keyword: forward.Keyword,
		Title: forward.Title,
		Content: forward.Content,
		Author: fullName,
		PublicStatus: forward.Status,
		PublicTime: time.Now(),
	}
	//简介为内容的前一百个字
	if len(forward.Content) >= 100 {
		db.Summary = db.Content[:100]
	}else{
		db.Summary = db.Content
	}

	//返回id值
	var blogId int64
	//判断是新增还是修改
	if forward.Command == "add" {
		//写入数据库
		blogId, err = db.AddBlog(Db)
		if err != nil {
			log.Println("insert error ",err.Error())
			w.WriteHeader(500)
			return
		}
		data.BlogId = blogId
	}else{
		//修改博客操作
		db.Id = forward.Id
		_, err = db.UpdateBlog(Db)
		if err != nil {
			log.Println("update error ",err.Error())
			w.WriteHeader(500)
			return
		}
		data.BlogId = int64(db.Id)
	}
	data.Status = "OK"

	json, _ := json.MarshalIndent(&data,"","\t")
	w.Write(json)
}

//博客的删除
//博客的单篇查找

//测试用
func uploadtest(w http.ResponseWriter, r *http.Request){
	t,_ := template.ParseFiles("upload.html")
	t.Execute(w,nil)
}
func main(){
	sessionMgr = Helper.NewSessionMgr("Cookiename",3600)
	server := http.Server{
		Addr: ":8000",
	}
	http.HandleFunc("/se-login",login)
	http.HandleFunc("/se-token",index)
	http.HandleFunc("/se-course",getCourseLimit)
	http.HandleFunc("/se-blog",getBlogLimit)
	http.HandleFunc("/se-upload",uploadImage)
	http.HandleFunc("/se-blog-exec",execBlog)
	//http.HandleFunc("/test",uploadtest)
	images := http.FileServer(http.Dir("./upload/image/"))
	http.Handle("/static/",http.StripPrefix("/static/",images))

	server.ListenAndServe()
}

