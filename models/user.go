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
	RealName string `json:"realname" form:"realname"`
	RemainSum int	`json:"remainsum" form:"remainsum"`
	Role int`json:"role"`
	//0为超级管理员，1为管理员，2为普通用户
	Head string`json:"head" form:"head"`
	//头像
	Phone string`json:"phone" form:"phone"`
	Email string`json:"email" form:"email"`
	Books []Book`json:"books" gorm:"many2many:user_book"`
	BookSum int`json:"booksum"`
	Introduce string `json:"introduce" json:"introduce"`

}
type UserMsg struct{
	Username string`json:"username" form:"username"`
	//Password string`json:"password" form:"password"`
	RealName string `json:"realname" form:"realname"`
	RemainSum int	`json:"remainsum" form:"remainsum"`
	Role int`json:"role"`
	//0为超级管理员，1为管理员，2为普通用户
	Head string`json:"head" form:"head"`
	//头像
	Phone string`json:"phone" form:"phone"`
	Email string`json:"email" form:"email"`
	Books []Book`json:"books" gorm"many2many:user_book"`
	BookSum int`json:"booksum"`
	Id uint `json:"id"`
	Introduce string `json:"introduce" json:"introduce"`
}
func TakeUserMsg(user User )UserMsg{
	var msg UserMsg
	msg.Role=user.Role
	msg.Phone=user.Phone
	msg.Email=user.Email
	msg.Username=user.Username
	msg.RealName=user.RealName
	msg.Books=user.Books
	msg.BookSum=user.BookSum
	msg.Introduce=user.Introduce
	msg.RemainSum= user.RemainSum
	msg.Id=user.ID
	msg.Head=user.Head
	return msg
}
//验证管理员权限
func CheckAdmin(username string)int{
	var user User
	err:=DB.Where("username = ?",username).First(&user).Error
	if err!=nil{
		fmt.Println("权限验证错误",err)
		return -1
	}
	return user.Role
}
//登陆验证
func CheckLogin(username string,password string)(User,error){
	var user User
	err:=DB.Where("username = ?",username).First(&user).Error
	fmt.Println(user)
	if err!=nil{
		return User{},err
	}
	if ScryptPW(password)!=user.Password{
		err=errors.New("密码错误")
		return User{},err
	}
	return user,nil
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
	err:=DB.Create(&data).Error
	if err!=nil{
		log.Fatal(err)
		return err
	}
	return nil
}
//查询用户列表
func GteUsers(PageSize int,Pagenum int)([]User,int){
	var users []User
	var sum int
	DB.Where("role = ?",2).Find(&users).Count(&sum)
	err:=DB.Limit(PageSize).Offset((Pagenum-1)*PageSize).Where("role = ?",2).Find(&users).Error
	if err!=nil{
		return nil,-1
	}
	return users, sum
}
//搜索用户
func SearchU(username string,pagesize int,pagenum int)([]User,error,int){
	var u []User
	var sum int
	username="%"+username+"%"
	err:=DB.Limit(pagesize).Offset((pagenum-1)*pagesize).Where("role = ? AND username LIKE ?",2,username).Find(&u).Count(&sum).Error
	if err!=nil{
		fmt.Println(err)
		return []User{},err,0
	}
	return u,err,sum
}
//搜索管理员
func SearchA(username string,pagesize int,pagenum int)([]User,error,int){
	var u []User
	var sum int
	username="%"+username+"%"
	err:=DB.Limit(pagesize).Offset((pagenum-1)*pagesize).Where("role = ? AND username LIKE ?",1,username).Find(&u).Count(&sum).Error
	if err!=nil{
		fmt.Println(err)
		return []User{},err,0
	}
	return u,err,sum
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
func EditUser(id int,u *User)User{
	var user User
	err:=DB.Debug().Model(&user).Where("id = ?",id).Updates(User{
		Introduce: u.Introduce,
		RealName: u.RealName,
		Phone: u.Phone,
		Email: u.Email,
		Head: u.Head,
	}).Error
	if err!=nil{
		fmt.Sprint(err)
	}
	return user
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
func GteAdmins(PageSize int,Pagenum int)([]User,int){
	var users []User
	var sum int
	DB.Where("role = ?",1).Find(&users).Count(&sum)
	err:=DB.Limit(PageSize).Offset((Pagenum-1)*PageSize).Where("role = ?",1).Find(&users).Error
	if err!=nil{
		fmt.Println(err)
		return nil,-1
	}
	return users, sum
}
//查询用户所借书籍
func Borrowedbooks(PageSize int,Pagenum int,username string)([]Book,int){
	var u User
	var books []Book
	DB.Where("username = ?",username).First(&u)
	err:=DB.Debug().Limit(PageSize).Offset((Pagenum-1)*PageSize).Model(&u).Association("Books").Find(&books).Error
	if err!=nil{
		fmt.Println(err.Error())
	}
	tmpb:=books
	sum:=DB.Model(&u).Association("Books").Find(&tmpb).Count()
	return books,sum
}
//用户借书
func Borrowbook(username string,bookname int)(Book,error){
	var u User
	var b Book
	DB.Where("username = ?",username).First(&u)
	DB.Debug().Where("id = ?",bookname).First(&b)
	status,_:=CheckBorrowed(u,b)
	if status==1{
		return Book{},errors.New("该书正在借阅")
	}else{
		if b.Sum>0{
			err=DB.Model(&u).Association("Books").Append(&b).Error
			if err!=nil{
				fmt.Println(err.Error())
				return Book{},nil
			}
			DB.Model(&u).Update("book_sum",u.BookSum+1)
			err=DB.Model(&b).Select("sum","borrow_sum").Updates(map[string]interface{}{
				"sum":b.Sum-1,
				"borrow_sum":b.BorrowSum+1,
			}).Error
			if err!=nil{
				fmt.Println(err.Error())
			}
			msg:="借阅了"+b.BookName+"bookid:"+strconv.Itoa(int(b.ID))
			fmt.Println(msg)
			PubilishMsg(u.Username,msg)
			return b,err
		}else{
			fmt.Println("余量不足")
			return Book{},errors.New("余量不足")
		}
	}


}
//用户还书
func ReturnBook(username string,bookname int)(Book,error){
	var u User
	var b Book
	DB.Where("username = ?",username).First(&u)
	DB.Debug().Where("id = ?",bookname).First(&b)
		err=DB.Debug().Model(&u).Association("Books").Delete(&b).Error
		if err!=nil{
			fmt.Println(err.Error())
			return Book{},nil
		}
	     DB.Model(&u).Select("book_sum").Updates(map[string]interface{}{"book_sum":u.BookSum-1})
		fmt.Println(u.Books)
		err=DB.Model(&b).Select("sum","borrow_sum").Updates(map[string]interface{}{
			"sum":b.Sum+1,
			"borrow_sum":b.BorrowSum-1,
		}).Error
		if err!=nil{
			fmt.Println(err.Error())
		}
	msg:="还了"+b.BookName+"bookid:"+strconv.Itoa(int(b.ID))
	PubilishMsg(u.Username,msg)
		return b,err

}
func CheckBorrowed(u User,b Book)(int,error){
	var books []Book
	err=DB.Model(&u).Association("Books").Find(&books).Error
	for i:=0;i< len(books);i++{
		if b.ID==books[i].ID{
			return 1,err
		}
	}
	return 0,err
}
func SearchUrl(username string)(string){
	var u User
	err:=DB.Where("username = ?",username).Find(&u).Error
	if err!=nil{
		fmt.Println(err)
		return ""
	}
	return u.Head
}