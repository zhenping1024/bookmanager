package v1

import (
	"bookmanager/middleware"
	"bookmanager/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

//获取用户信息
func GetUser( c *gin.Context){
	authString := c.Request.Header.Get("Authorization")
	kv := strings.Split(authString, " ")
	tokenString := kv[1]

	// Parse token
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
	//u.Password=c.ShouldBindJSON(&u)
	fmt.Println(u)
	err:=models.CheckUser(u.Username)
	if err!=nil{
		fmt.Println(err)
		c.JSON(http.StatusOK,gin.H{
			"status":"用户已存在",
			"data":u,
			"message":err.Error(),
		})
	}else{
		models.CreatUser(&u)
		c.JSON(http.StatusOK,gin.H{
			"status":"创建成功",
			"data":u,
			"message":nil,
		})
	}

}

//查询用户
//func FindUser(c*gin.Context){
//
//}

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
	users:=models.GteUsers(pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":users,
		"message":nil,
	})
}

//编辑用户资料
func EditUser(c*gin.Context){
	var u models.User
	var dst string
	id,_:=strconv.Atoi(c.Param("id"))
	u.Username=c.PostForm("username")
	u.Email=c.PostForm("email")
	u.Phone=c.PostForm("phone")
	//c.ShouldBind(&u)
	err:=models.CheckUser(u.Username)
	file,e:=c.FormFile("imag")
	if e!=nil{
		//c.JSON(200,gin.H{
		//	"err":e.Error(),
		//})
		fmt.Println(e)
	}else{
		c.FormFile("imag")
		time_int:=time.Now().Unix()
		time_str:=strconv.FormatInt(time_int,10)
		filename:=time_str+file.Filename
		dst=path.Join("./statics/image/userimag",filename)
		//获取存储路径
		u.Head=dst
		if err := c.SaveUploadedFile(file, dst);
			err != nil {
			//自己完成信息提示
				fmt.Println("上传失败",err)
			return
		}
		fmt.Println("save",u)
		fmt.Println("上传成功")
	}
	fmt.Println(u.Head,u)
		models.EditUser(id,&u)
	c.JSON(http.StatusOK,gin.H{
		"status":u,
		"message":err,
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
//新增管理员
func AddAdmin(c*gin.Context){
	var u models.User
	id,_:=strconv.Atoi(c.Param("id"))

	err:=models.AddAdmin(id,&u)
	c.JSON(http.StatusOK,gin.H{
		"status":u,
		"message":err,
	})
}
//撤销管理员
func DeleteAdmin(c*gin.Context){
	var u models.User
	id,_:=strconv.Atoi(c.Param("id"))

	err:=models.DeleteAdmin(id,&u)
	c.JSON(http.StatusOK,gin.H{
		"status":u,
		"message":err,
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
	users:=models.GteAdmins(pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":users,
		"message":nil,
	})
}
//用户借书列表
func Getborrow(c*gin.Context){

}
//用户借书
func BorrowBook(c*gin.Context){

}
//用户还书
func ReturnBook(c*gin.Context){

}