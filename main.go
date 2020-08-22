package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"sesiontest/Helper"
	"sesiontest/back"
	"sesiontest/dbs"
)

var sessionMgr *Helper.SessionMgr = nil

func checkLogin(w http.ResponseWriter, r *http.Request)(string,error){

	if r.Method == "GET"{
		return "",errors.New("method must is post")
	}else if r.Method == "POST"{
		//获取cookie里面的值
		sessionID := r.Header.Get("Cookie")
		//判断是否存在Cookiename
		/**
			dasdasd,Cookiename=asdasdasdasdasd,dasdasdasd 大概是这个样子
			.*[,]{1}Cookiename=([^,]+).*
		 **/
		re := regexp.MustCompile(`.*[;]{1} Cookiename=([^;]+).*`)
		result := re.FindStringSubmatch(sessionID)
		//匹配到结果，0是本身，1是匹配结果
		if len(result) == 2 {
			sessionID = result[1]
			_, exist := sessionMgr.GetSessionVal(sessionID,"user")
			if exist{
				return sessionID,nil
			}else{
				return "",errors.New("not found user")
			}
		}else{
			return "",errors.New("not found cookie")
		}
	}
	return "",errors.New("other case")
}
func index(w http.ResponseWriter, r *http.Request){

	back := back.BackToken{

	}
	sessionID := sessionMgr.StartSession(w,r)
	sessionMgr.SetSessionVal(sessionID, "token", sessionMgr.NewToken())
	token, _ := sessionMgr.GetSessionVal(sessionID,"token")
	back.Token = token.(string)

	w.Header().Set("Content-type","application/json")
	json, err := json.MarshalIndent(&back,"","\t")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Write(json)

}
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
		//获取cookiename的值
		sessionID := r.Header.Get("Cookie")

		re := regexp.MustCompile(`.*[;]{1} Cookiename=([^;]+).*`)
		result := re.FindStringSubmatch(sessionID)
		if len(result) == 2 {
			sessionID = result[1]
		}else{
			sessionID = sessionMgr.StartSession(w,r)
		}


		//fmt.Println("得到的session值:",sessionID)
		//在数据库中得到对应数据

		Db, err := dbs.GetConnect()
		defer Db.Close()
		if err != nil {
			panic("error")
		}

		//对比token值
		token, _ := sessionMgr.GetSessionVal(sessionID,"token")
		if r.FormValue("token") != token{
			w.WriteHeader(201)
			//更新token
			sessionMgr.SetSessionVal(sessionID,"token",sessionMgr.NewToken())
			return
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
func getCourse(w http.ResponseWriter, r *http.Request){
	data := back.CourseList{}
	w.Header().Set("Content-Type","application/json")

	//先校验登陆状态
	sesionID,err := checkLogin(w,r)
	if err != nil {
		back := back.BackInfo{
			Sessionid: sesionID,
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
	//添加错误原因
	if err != nil {
		log.Println("database is error")
		w.WriteHeader(500)
	}
	defer Db.Close()

	//添加错误原因
	list,err :=dbs.GetCourseListBy(courseType,Db)
	if err != nil {
		log.Println("query err:",err.Error(),"query type is : ",courseType)
		w.WriteHeader(500)
	}
	//写入数据
	for _,course := range list{
		temp := back.Course{}
		temp.Id, temp.Type, temp.Url, temp.Title, temp.Content = course.Id, course.Type, course.Url, course.Title, course.Content
		data.List = append(data.List,temp)
	}
	data.LoginInfo.Status = "success"
	data.LoginInfo.Sessionid = sesionID
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
	server.ListenAndServe()
}

