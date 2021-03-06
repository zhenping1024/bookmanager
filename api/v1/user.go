package v1

import (
	"bookmanager/middleware"
	"bookmanager/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//获取用户信息
func GetUser( c *gin.Context){
	authString := c.Request.Header.Get("Authorization")
	kv := strings.Split(authString, " ")
	tokenString := kv[1]
	token, err := middleware.ParseJwt(tokenString)
	if err!=nil{
		fmt.Println(err)
	}
	u,e:=models.GetUser(token.Username)
	c.JSON(http.StatusOK,gin.H{
		"status":e,
		"data":u,
	})
}
//添加用户
func AddUser(c*gin.Context){
	var u models.User
	u.Username=c.PostForm("username")
	u.Password=c.PostForm("password")
	u.Email=c.PostForm("email")
	u.Phone=c.PostForm("phone")
	u.Role=2
	fmt.Println(u)
	err:=models.CheckUser(u.Username)
	if err!=nil{
		fmt.Println(err)
		c.JSON(http.StatusOK,gin.H{
			"status":"用户已存在",
			"data":models.TakeUserMsg(u),
			"message":err.Error(),
		})
	}else{
		models.CreatUser(&u)
		c.JSON(http.StatusOK,gin.H{
			"status":"创建成功",
			"data":models.TakeUserMsg(u),
			"message":nil,
		})
	}

}
//查询用户列表
func GetUsers(c*gin.Context){
	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	if pagesize==0{
		pagesize = -1
	}
	if pagenum == 0{
		pagenum =-1
	}
	users,sum:=models.GteUsers(pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":users,
		"pagesize":pagesize,
		"pagenum":pagenum,
		"num": sum,
	})
}
//编辑用户资料
func EditUser(c*gin.Context){
	var u models.User
	id,_:=strconv.Atoi(c.Param("id"))
	u.Email=c.PostForm("email")
	u.Phone=c.PostForm("phone")
	u.RealName=c.PostForm("realname")
	u.Introduce=c.PostForm("introduce")
	er:=models.CheckUser(u.Username)
	if er!=nil{
		c.JSON(200,gin.H{
			"err":er,
			"message":"用户名重复",
		})
		return
	}
	file,fileHeader,e:=c.Request.FormFile("imag")
	if e!=nil{
		fmt.Println(e)
	}else{
		filesize:=fileHeader.Size
		url,ee:=models.UpLoadFile(file,filesize)
		if ee!=nil{
			fmt.Println("上传错误")
		}
		u.Head=url
		fmt.Println(url,"urlshi")
	}
	fmt.Println(u.Head,u)
	u=models.EditUser(id,&u)
	c.JSON(http.StatusOK,gin.H{
		"status":models.TakeUserMsg(u),
		"message":nil,
	})
}
//删除用户
func DeleteUser(c*gin.Context){
	id,_:=strconv.Atoi(c.Param("id"))
	code:=models.DeleteUser(id)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":code,
		"message":nil,
	})
}
//创建新管理员
func CreatAdmin(c*gin.Context){
	var u1 models.User
	var u2 models.User
	authString := c.Request.Header.Get("Authorization")
	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		c.Abort()
		return
	}
	tokenString := kv[1]
	token, err := middleware.ParseJwt(tokenString)
	err=models.DB.Where("username = ?",token.Username).First(&u1).Error
	fmt.Println("管理员是",u1.Role,u1.Username)
	if err!=nil{
		c.JSON(http.StatusOK,gin.H{
			"status":"错误",
			"message":err,
		})
	}else if u1.Role!=0{
		c.JSON(http.StatusOK,gin.H{
			"status":"错误",
			"message":errors.New("无权限"),
		})
	}
	u2.Username=c.PostForm("username")
	u2.Password=c.PostForm("password")
	u2.Email=c.PostForm("email")
	u2.Phone=c.PostForm("phone")
	u2.Role=1
	file,fileHeader,e:=c.Request.FormFile("imag")
	if e!=nil{
		fmt.Println(e.Error())
	}else{
		filesize:=fileHeader.Size
		url,ee:=models.UpLoadFile(file,filesize)
		if ee!=nil{
			fmt.Println("上传错误")
		}
		u2.Head=url
		fmt.Println("save",u2)
		fmt.Println("上传成功")
	}
	err=models.CheckUser(u2.Username)
	if err!=nil{
		fmt.Println(err)
		c.JSON(http.StatusOK,gin.H{
			"status":"用户已存在",
			"data":models.TakeUserMsg(u2),
			"message":err.Error(),
		})
	}else{
		models.CreatUser(&u2)
		c.JSON(http.StatusOK,gin.H{
			"status":"管理员创建成功",
			"data":models.TakeUserMsg(u2),
			"message":nil,
		})
	}
}
//新增管理员
func AddAdmin(c*gin.Context){
	var u models.User
	id,_:=strconv.Atoi(c.Param("id"))

	err:=models.AddAdmin(id,&u)
	c.JSON(http.StatusOK,gin.H{
		"status":u,
		"Error":err,
	})
}
//撤销管理员
func DeleteAdmin(c*gin.Context){
	var u models.User
	id,_:=strconv.Atoi(c.Param("id"))

	err:=models.DeleteAdmin(id,&u)
	c.JSON(http.StatusOK,gin.H{
		"status":u,
		"Error":err,
	})
}
func GetAdmins(c*gin.Context){
	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	if pagesize==0{
		pagesize = -1
	}
	if pagenum == 0{
		pagenum =-1
	}
	users,sum:=models.GteAdmins(pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":users,
		"num": sum,
		"pagesize":pagesize,
		"pagenum":pagenum,
	})
}
//用户借书列表
func Getborrow(c*gin.Context){
	authString := c.Request.Header.Get("Authorization")
	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		c.Abort()
		return
	}
	tokenString := kv[1]
	token, _ := middleware.ParseJwt(tokenString)
	//err=models.DB.Where("username = ?",token.Username).First(&u1).Error
	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	if pagesize==0{
		pagesize = -1
	}
	if pagenum == 0{
		pagenum =-1
	}
	sum,books:=models.Borrowedbooks(pagesize,pagenum,token.Username)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":books,
		"num": sum,
		"pagesize":pagesize,
		"pagenum":pagenum,
	})
}
//用户借书
func BorrowBook(c*gin.Context){
	id,err:=strconv.Atoi(c.Param("id"))
	fmt.Println(id,"cuowu",err)
	authString := c.Request.Header.Get("Authorization")
	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		c.Abort()
		return
	}
	tokenString := kv[1]
	token, _ := middleware.ParseJwt(tokenString)
	b,e:=models.Borrowbook(token.Username,id)
	c.JSON(http.StatusOK,gin.H{
		"status":fmt.Sprint(e),
		"data":b,
	})
}
//用户还书
func ReturnBook(c*gin.Context){
	id,err:=strconv.Atoi(c.Param("id"))
	fmt.Println("错误",err)
	authString := c.Request.Header.Get("Authorization")
	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		c.Abort()
		return
	}
	tokenString := kv[1]
	token, _ := middleware.ParseJwt(tokenString)
	b,e:=models.ReturnBook(token.Username,id)
	c.JSON(http.StatusOK,gin.H{
		"status":e,
		"data":b,
	})
}
//搜索普通用户
func SearchUser(c*gin.Context){
	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	if pagesize==0{
		pagesize = -1
	}
	if pagenum == 0{
		pagenum =-1
	}

	username:=c.Query("username")
	u,err,sum:=models.SearchU(username,pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":fmt.Sprint(err),
		"data":u,
		"datasum": sum,
	})
}
//搜索普通管理员
func SearchAdmin(c*gin.Context){

	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	if pagesize==0{
		pagesize = -1
	}
	if pagenum == 0{
		pagenum =-1
	}
	username:=c.Query("username")
	u,err,sum:=models.SearchA(username,pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":fmt.Sprint(err),
		"data":u,
		"datasum": sum,
	})
}
//书名搜索书籍
func SearchBook(c*gin.Context){
	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	if pagesize==0{
		pagesize = -1
	}
	if pagenum == 0{
		pagenum =-1
	}

	bookname:=c.Query("bookname")
	fmt.Println(bookname)
	u,err,sum:=models.SearchBook(bookname,pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":fmt.Sprint(err),
		"data":u,
		"datasum": sum,
	})
}
//书名搜索已借书籍
func SearchUserBook(c*gin.Context){
	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	if pagesize==0{
		pagesize = -1
	}
	if pagenum == 0{
		pagenum =-1
	}

	bookname:=c.Query("bookname")
	authString := c.Request.Header.Get("Authorization")
	kv := strings.Split(authString, " ")
	tokenString := kv[1]

	// Parse token
	token, err := middleware.ParseJwt(tokenString)
	if err!=nil{
		fmt.Println(err)
	}
	u,err,sum:=models.SearchBorrowedBook(token.Username,bookname,pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":fmt.Sprint(err),
		"data":u,
		"datasum": sum,
	})
}
//管理员接受消息
func GetMsg(c*gin.Context){

	authString := c.Request.Header.Get("Authorization")

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		//result := models.UnauthorizedResult()
		c.JSON(200, gin.H{
			"status":"错误",
			"message":"Token错误",
		})
		c.Abort()
		return
	}
	tokenString := kv[1]
	// Parse token
	token, _:= middleware.ParseJwt(tokenString)
	role:=models.CheckAdmin(token.Username)
	if role>1{
		m,_,_:=models.GetMsgs(token.Username)
		c.JSON(http.StatusOK,gin.H{
			"message":m,
			"status":1,
		})
	}else{
		fmt.Println(models.Ms.Sum, len(models.Ms.M))
		if models.Ms.Sum< len(models.Ms.M){
			models.Ms.Sum= len(models.Ms.M)
			c.JSON(http.StatusOK,gin.H{
				"message":models.Ms.M,
				"status":1,
			})
		}else{
			c.JSON(http.StatusOK,gin.H{
				"message":models.Ms.M,
				"status":0,
			})
		}
	}

}
//置信管理员
func ToAdmin(c*gin.Context){
	var msg models.Msg
	authString := c.Request.Header.Get("Authorization")
	kv := strings.Split(authString, " ")
	tokenString := kv[1]
	token, err := middleware.ParseJwt(tokenString)
	if err!=nil{
		fmt.Println(err)
	}

	context:=c.PostForm("comment")
	msg.MsgUser=token.Username
	msg.MsgContext=context
	msg.Creattime=time.Now()
	models.ToAdmin(msg)
	c.JSON(http.StatusOK,gin.H{
		"status":fmt.Sprint(err),
		"data":msg,
	})
}