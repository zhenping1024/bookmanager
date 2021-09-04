package models
import (
	"bookmanager/utils"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)
var DB *gorm.DB
var err error
func InitDb(){
	DB,err=gorm.Open(utils.Db,fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
	utils.DbUser,
	utils.DbPassWord,
	utils.DbHost,
	utils.DbPort,
	utils.DbName,
	))
	if err!=nil{
		fmt.Println("连接数据库失败",err)
	}
	DB.SingularTable(true)
	DB.AutoMigrate(&User{},&Book{})
	//DB.Close()
}
