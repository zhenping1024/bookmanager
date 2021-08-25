package v1

import (
	"bookmanager/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//添加分类
func AddCategory(c*gin.Context){
	var u models.Category
	c.ShouldBindJSON(&u)
	err:=models.CheckCategory(u.TypeName)
	if err!=nil{
		fmt.Println(err)
		c.JSON(http.StatusOK,gin.H{
			"status":"分类已存在",
			"data":u,
			"message":err.Error(),
		})
	}else{
		models.CreatCategory(&u)
		c.JSON(http.StatusOK,gin.H{
			"status":"创建成功",
			"data":u,
			"message":nil,
		})
	}

}
//查询分类下的文章

//查询分类列表
func GetCategory(c*gin.Context){
	pagesize,_:=strconv.Atoi(c.Query("pagesize"))
	pagenum,_:=strconv.Atoi(c.Query("pagenum"))
	if pagesize==0{
		pagesize = -1
	}
	if pagenum == 0{
		pagenum =-1
	}
	cate:=models.GetCategoey(pagesize,pagenum)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":cate,
		"message":nil,
	})
}
//编辑分类
func EditCategory(c*gin.Context){
	var u models.Category
	id,_:=strconv.Atoi(c.Param("id"))
	c.ShouldBind(&u)
	err:=models.CheckCategory(u.TypeName)
	if err==nil{
		models.EditCategory(id,&u)
	}
	c.JSON(http.StatusOK,gin.H{
		"status":u,
		"message":err,
	})
}

//删除分类
func DeleteCategory(c*gin.Context){
	id,_:=strconv.Atoi(c.Param("id"))
	code:=models.DeleteCategory(id)
	c.JSON(http.StatusOK,gin.H{
		"status":"成功",
		"data":code,
		"message":nil,
	})
}