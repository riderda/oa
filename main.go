package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"oa/Helper"
	"oa/back"
	"oa/dbs"
	"strconv"
	"sync"
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


func test(w http.ResponseWriter,r *http.Request){
	sessionID := sessionMgr.StartSession(w,r)
	sessionMgr.SetSessionVal(sessionID, "token", sessionMgr.NewToken())
	token, _ := sessionMgr.GetSessionVal(sessionID,"token")
	t, _ := template.ParseFiles("login.html")
	t.Execute(w,token)
}
func testLogin(w http.ResponseWriter, r *http.Request){
	sessionID := sessionMgr.StartSession(w,r)
	sessionMgr.SetSessionVal(sessionID,"token",sessionMgr.NewToken())
	r.Header.Add("Authentication",sessionID)
	login(w,r)
}
func getCourse(w http.ResponseWriter, r *http.Request){
	data := back.CourseList{}
	w.Header().Set("Content-Type","application/json")

	//先校验登陆状态
	sessionID,err := checkLogin(w,r)
	if err != nil {
		back := back.BackInfo{
			Sessionid: sessionID,
			Status: "fail",
		}
		data.LoginInfo = back
		json,_ := json.MarshalIndent(&data,"","\t")
		w.Write(json)
		return
	}

	//解析请求并返回数据
	r.ParseForm()
	courseType := r.FormValue("type")
	Db, err := dbs.GetConnect()
	//打印错误日志
	if err != nil {
		log.Println("database is error",err.Error())
		w.WriteHeader(500)
		return
	}
	defer Db.Close()

	//打印错误原因
	list,err :=dbs.GetCourseListBy(courseType,Db)
	if err != nil {
		log.Println("query err:",err.Error(),"query type is : ",courseType)
		w.WriteHeader(500)
		return
	}
	//写入数据
	for _,course := range list{
		temp := back.Course{}
		temp.Id, temp.Type, temp.Url, temp.Title, temp.Content = course.Id, course.Type, course.Url, course.Title, course.Content
		data.List = append(data.List,temp)
	}
	data.LoginInfo.Status = "success"
	data.LoginInfo.Sessionid = sessionID
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

//获取博客
func getBlogLimit(w http.ResponseWriter, r *http.Request){
	data := back.BlogListLimit{}
	w.Header().Set("Content-Type","application/json")

	//先检查登陆
	//_, err := checkLogin(w,r)
	////没有登陆
	//if err != nil {
	//	data.Info = back.QueryInfo{
	//		Length: 0,
	//		PageNum: 0,
	//		Error: "not login sessionID",
	//	}
	//	json, _ := json.MarshalIndent(&data,"","\t")
	//	w.Write(json)
	//	return
	//}

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
	data.BlogListByKeyword = make(map[string][]dbs.DbBlog)
	//从搜索里获取关键词
	ch := dbs.Seg.CutForSearch(Search,true)
	for keyword := range ch{
		fmt.Print(keyword+" ")
		wg.Add(1)
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

			//内容更新
			Bloglist, err := dbs.GetBlogLimitBy(keyword,Limit,Db)
			if err != nil{
				log.Println("query is error",err.Error())
				w.WriteHeader(500)
				return
			}
			data.BlogListByKeyword[keyword] = Bloglist

		}()
	}
	//只能阻塞，但不能影响write
	wg.Wait()
	data.Info.PageNum = Limit
	json, _ := json.MarshalIndent(&data,"","\t")
	w.Write(json)
}
func main(){
	sessionMgr = Helper.NewSessionMgr("Cookiename",3600)
	server := http.Server{
		Addr: ":8000",
	}
	http.HandleFunc("/se-login",login)
	http.HandleFunc("/se-token",index)
	http.HandleFunc("/test",test)
	http.HandleFunc("/se-course",getCourse)
	http.HandleFunc("/testLogin",testLogin)
	http.HandleFunc("/se-course-limit",getCourseLimit)
	http.HandleFunc("/se-blog-limit",getBlogLimit)
	server.ListenAndServe()
}

