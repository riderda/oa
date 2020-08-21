package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sesiontest/Helper"
	"sesiontest/dbs"
	"strings"
)
type BackToken struct{
	Token string `json:"token"`
}
type BackInfo struct{
	Sessionid string `json:"sessionid"`
	Status string `json:"status"`
}
var sessionMgr *Helper.SessionMgr = nil

func index(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(404)
		return
	}
	back := BackToken{

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

		//fmt.Println("本该的session值:",sessionID)
		if strings.Contains(sessionID,"Cookiename="){
			//获取之前的session
			sessionID = strings.Split(sessionID,"Cookiename=")[1]
		}else{
			//之前没有session
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

		//如果登陆失败，直接重定向到登陆界面
		if err != nil {
			backinfo := BackInfo{
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
			if userInfo, ok := sessionMgr.GetSessionVal(onlineSessionID, "UserInfo"); ok {
				if value, ok := userInfo.(dbs.User); ok {
					if value.Id == user.Id {
							sessionMgr.EndSessionBy(onlineSessionID)
					}
				}
			}
		}

		//设置变量值
		sessionMgr.SetSessionVal(sessionID, "UserInfo", user)

		backinfo := BackInfo{
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
func main(){
	sessionMgr = Helper.NewSessionMgr("Cookiename",3600)
	server := http.Server{
		Addr: ":8000",
	}
	http.HandleFunc("/login",login)
	http.HandleFunc("/",index)
	http.HandleFunc("/test",test)
	server.ListenAndServe()
}

