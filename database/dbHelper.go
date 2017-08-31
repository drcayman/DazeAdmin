package database
import (
	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
	"fmt"
)
var engine *xorm.Engine
type User struct {
	Id int64
	Username string `xorm:"unique"`
	Password string
	Created time.Time `xorm:"created"`
	Expired time.Time
	Group string

}
func (User) TableName() string {
	return "DF_User"
}
func GetUserCount() int64 {
	var user User
	b,_:=engine.Count(&user)
	return b
}
func GetAllUser() []User {
	var user []User=make([]User,0)
	engine.Find(&user)
	return user
}
func GetUserById(id int)(User,bool){
	var user User
	b,_:=engine.Where("id = ?", id).Get(&user)
	return user,b
}
func AddUser(u User)(error){
	_,err:=engine.Insert(u)
	return err
}
func DeleteById(id int) (bool,string){
	v,b:=GetUserById(id)
	if b==false{
		return b,""
	}
	_,err:=engine.Delete(&v)
	return err==nil,v.Username
}
func EditUserById(id int,u User){
	_,err:=engine.Where("id = ?",id).Update(u)
	if err!=nil{
		fmt.Println("编辑此用户失败，可能是用户名已存在！错误代码：",err.Error())
	}else{
		fmt.Println("编辑此用户成功！")
	}
}
func LoadDatabase(driver string,connectionString string){
	var err error
	engine,err=xorm.NewEngine(driver,connectionString)
	if err!=nil{
		log.Fatal("数据库连接失败！原因：",err)
	}
	err=engine.Sync2(new(User))
	if err!=nil{
		log.Fatal("数据库加载失败！原因：",err)
	}
	count:=GetUserCount()
	log.Println("数据库连接成功！用户数：",count)
}