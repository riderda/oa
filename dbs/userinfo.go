package dbs

import (
	"github.com/jmoiron/sqlx"
	"time"
)

//由从sesion里的user中的id来获取更多信息的userinfo
type UserInfo struct{
	Id int `db:"id"`
	UserId int `db:"user id"` //对应user结构体里的Id
	FullName string `db:"full name"`
	JoinTime time.Time `db:"join time"`
	LeaveTime time.Time `db:"leave time"`
}

//根据user id查找用户真实信息
func GetFullName(userId int,Db *sqlx.DB)(string,error){
	row := Db.QueryRow("select `full name` from userinfo where `user id`=?",userId)
	var fullName string
	err := row.Scan(&fullName)
	if err != nil{
		return "", err
	}
	return fullName,nil
}
