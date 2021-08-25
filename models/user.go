package models

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"
	"log"
	"strconv"
)

type User struct{
	gorm.Model
	Username string`json:"username" form:"username"`
	Password string`json:"password" form:"password"`
	Role int`json:"role"`
	//0为超级管理员，1为管理员，2为普通用户
	Head string`json:"head" form:"head"`
	//头像
	Phone string`json:"phone" form:"phone"`
	Email string`json:"email" form:"email"`
	Books []Book`json:"books" gorm"many2many:user_book"`
	BookSum int`json:"booksum"`
}
//登陆验证
func CheckLogin(username string,password string)error{
	var user User
	err:=DB.Where("username = ?",username).First(&user).Error
	//log.Fatal(user)
	fmt.Println(user)
	if err!=nil{
		return err
	}
	if ScryptPW(password)!=user.Password{
		err=errors.New("密码错误")
		return err
	}
	return nil
}
//获取用户信息
func GetUser(name string)( User ,error){
	var u User
	err:=DB.Where("username = ?",name).First(&u).Error
	return u,err

}
//查寻用户是否存在
func CheckUser(name string)error{
	var u User
	DB.Select("id").Where("username = ?",name).First(&u)
	if u.ID>0{

		err:=errors.New("用户已存在")
		return err
	}
	return nil
}
//新增用户
func CreatUser(data *User)error{
	data.Password=ScryptPW(data.Password)
	fmt.Println(data)
	err:=DB.Create(&data).Error
	if err!=nil{
		log.Fatal(err)
		return err
	}
	return nil
}
//查询用户列表
func GteUsers(PageSize int,Pagenum int)[]User{
	var users []User
	err:=DB.Limit(PageSize).Offset((Pagenum-1)*PageSize).Find(&users).Error
	if err!=nil{
		return nil
	}
	return users
}
//密码加密
func ScryptPW(password string)string{
	const Keylen =10
	salt:=make([]byte ,8)
	salt=[]byte{12,32,4,6,66,22,222,11}
	HashPw,err:=scrypt.Key([]byte(password),salt,16384,8,1,Keylen)
	if err!=nil{
		log.Fatal(err)
	}
	PW:=base64.StdEncoding.EncodeToString(HashPw)
	return PW
}
//删除用户
func DeleteUser(id int)int{
	var user User
	err :=DB.Where("id = ?",id).Delete(&user).Error
	if err!=nil{
		log.Fatal("删除错误")
		return -1
	}
	ID,_:=strconv.Atoi(fmt.Sprint(user.ID))
	return ID
}
//编辑用户信息
func EditUser(id int,u *User){
	var maps=make(map[string]interface{})
	var user User
	maps["username"]=u.Username
	maps["phone"]=u.Phone
	maps["email"]=u.Email
	maps["head"]=u.Head
	err:=DB.Model(&user).Where("id = ?",id).Update(maps).Error
	if err!=nil{
		log.Fatal(err)
	}
}
//增加管理员
func AddAdmin(id int,u*User)error{
	var maps=make(map[string]interface{})
	var user User
	maps["role"]=1
	err:=DB.Model(&user).Where("id = ?",id).Update(maps).Error
	if err!=nil{
		return err
		log.Fatal(err)
	}
	return err
}

//解除管理员
func DeleteAdmin(id int,u*User)error{
	var maps=make(map[string]interface{})
	var user User
	maps["role"]=2
	err:=DB.Model(&user).Where("id = ?",id).Update(maps).Error
	if err!=nil{
		return err
		log.Fatal(err)
	}
	return err
}
//获取管理员列表
func GteAdmins(PageSize int,Pagenum int)[]User{
	var users []User
	err:=DB.Limit(PageSize).Offset((Pagenum-1)*PageSize).Where("role = ?",0).Or("role = ?",1).Find(&users).Error
	if err!=nil{
		return nil
	}
	return users
}
//查询用户所借书籍
func Borrowedbooks(){

}
//用户借书
func Borrowbook(){

}
//用户还书
func ReturnBook(){

}