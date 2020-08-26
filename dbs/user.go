package dbs

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)
type UserControl interface{
	AddUser(Db *sqlx.DB)(error)
	DeleteUser(Db *sqlx.DB)(error)
	FindUser(Db *sqlx.DB)(error)
}
type User struct {
	Id int `db:"id""`
	Username string `db:"username"`
	Password string `db:"password"`
}
func (u *User)AddUser(Db *sqlx.DB)(error){
	result, err := Db.Exec("insert into user(username,password) values(?,?)",u.Username,u.Password)
	if err != nil{
		return err
	}
	id , _ := result.LastInsertId()
	fmt.Println(id)
	return nil
}

func (u *User)FindUser(Db *sqlx.DB)(error){
	row := Db.QueryRow("select * from user where username=? and password=?",&u.Username,&u.Password)
	err := row.Scan(&u.Id,&u.Username,&u.Password)
	if err != nil{
		return err
	}
	//fmt.Println(u)
	return nil
}



