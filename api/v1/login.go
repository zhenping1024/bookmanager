package v1

import (
	"bookmanager/middleware"
	"bookmanager/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c*gin.Context){
	var u models.User
	c.ShouldBind(&u)
	fmt.Println("login",u)
	var token string
	err:=models.CheckLogin(u.Username,u.Password)
	if err==nil{
		token=middleware.InitJWT(u.Username)
	}
	c.JSON(http.StatusOK,gin.H{
		"token":token,
		"error":err,
	})
}
