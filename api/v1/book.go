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

//新增图书
func AddBook(c*gin.Context){
	var u models.Book
	u.BookName=c.PostForm("bookname")
	u.BookType=c.PostForm("booktype")
	u.Author=c.PostForm("author")
	u.Introduce=c.PostForm("introduce")
	u.Sum,_=strconv.Atoi(c.PostForm("sum"))
	file,e:=c.FormFile("imag")
	if e!=nil{
		fmt.Println(e)
	}else{
		c.FormFile("imag")
		time_int:=time.Now().Unix()
		time_str:=strconv.FormatInt(time_int,10)
		filename:=time_str+u.BookName
		dst:=path.Join("./statics/image/bookimage",filename)
		//获取存储路径
		u.BookImag=dst
		if err := c.SaveUploadedFile(file, dst);
			err != nil {
			//自己完成信息提示
			return
		}
		fmt.Println("save",u)
		fmt.Println("上传成功")
	}
		u.BorrowSum=0
		models.CreatBook(&u)
		c.JSON(http.StatusOK,gin.H{
			"status":"创建成功",
			"data":u,
			"message":nil,
		})

}
//上架图书
func Upbook(c*gin.Context){
	var u models.Book
	id,_:=strconv.Atoi(c.Param("id"))
	c.ShouldBind(&u)

	models.UpBook(id,&u)
	c.JSON(http.StatusOK,gin.H{
		"status":u,
		"message":"上架成功",
	})
}
//下架图书
func Downbook(c*gin.Context){
	var u models.Book
	id,_:=strconv.Atoi(c.Param("id"))
	c.ShouldBind(&u)

	models.DownBook(id,&u)
	c.JSON(http.StatusOK,gin.H{
		"status":u,
		"message":"下架成功",
	})
}
//删除图书
func DeleteBook(c*gin.Context){
	id,_:=strconv.Atoi(c.Param("id"))
	code:=models.DeleteBook(id)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":code,
		"message":nil,
	})
}
//查询图书列表
func GetBooks(c*gin.Context){
	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	if pagesize==0{
		pagesize = -1
	}
	if pagenum == 0{
		pagenum =-1
	}
	cate,sum:=models.GetBooks(pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":cate,
		"num": sum,
		"pagesize":pagesize,
		"pagenum":pagenum,
	})
}
//用户获取以上架图书
func GetUpbooks(c*gin.Context){
	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	if pagesize==0{
		pagesize = -1
	}
	if pagenum == 0{
		pagenum =-1
	}
	cate,sum:=models.GetUpBooks(pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":cate,
		"num": sum,
		"pagesize":pagesize,
		"pagenum":pagenum,
	})
}
//查询单一书籍
func GetBook(c*gin.Context){
	var u models.Book
	id,_:=strconv.Atoi(c.Param("id"))
	u=models.GetBookInfo(id)
	c.JSON(http.StatusOK,gin.H{
		"status":u,
		"message":"获取成功",
	})
}
//编辑图书资料
func EditBook(c*gin.Context){
	var u models.Book
	id,_:=strconv.Atoi(c.Param("id"))
	u.BookName=c.PostForm("bookname")
	u.BookType=c.PostForm("booktype")
	u.Author=c.PostForm("author")
	u.Introduce=c.PostForm("introduce")
	u.Sum,_=strconv.Atoi(c.PostForm("sum"))
	//c.ShouldBind(&u)
	file,e:=c.FormFile("imag")
	if e!=nil{
		c.JSON(http.StatusOK,gin.H{
			"status":e,
			"message":"图片上传错误",
		})
		fmt.Println(e)
	}else{
		c.FormFile("imag")
		time_int:=time.Now().Unix()
		time_str:=strconv.FormatInt(time_int,10)
		filename:=time_str+u.BookName
		dst:=path.Join("./statics/image/bookimage",filename)
		//获取存储路径
		u.BookImag=dst
		if err := c.SaveUploadedFile(file, dst);
			err != nil {
			//自己完成信息提示
			return
		}
		fmt.Println("save",u)
		fmt.Println("上传成功")
	}
		models.EditBook(id,&u)
	c.JSON(http.StatusOK,gin.H{
		"status":u,
		"message":"编辑成功",
	})
}
//发表评论
func SendComment(c*gin.Context){
	var comment models.Comment
	authString := c.Request.Header.Get("Authorization")
	kv := strings.Split(authString, " ")
	tokenString := kv[1]
	token, err := middleware.ParseJwt(tokenString)
	if err!=nil{
		fmt.Println(err)
	}
	bookid,_:=strconv.Atoi(c.Param("bookid"))
	context:=c.PostForm("comment")
	comment.Commentuser=token.Username
	comment.Context=context
	comment.UpdatedAt=time.Now()
	comment.ID=uint(bookid)
	fmt.Println(comment)
	comment,err=models.PublishComment(bookid,comment)
	c.JSON(http.StatusOK,gin.H{
		"status":fmt.Sprint(err),
		"data":comment,
	})
}
//显示评论
func GetComment(c*gin.Context){
	bookid,_:=strconv.Atoi(c.Param("bookid"))
	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	fmt.Println("bookid",bookid)
	comments,err,sum:=models.GetComment(bookid,pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":fmt.Sprint(err),
		"data":comments,
		"sum":sum,
	})
}

