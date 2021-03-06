package middleware

import (
	"bookmanager/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)
var TwtKey =[]byte("bookmanager")

type MyClaims struct{
	Username string `json:"username"`
	jwt.StandardClaims
}
//初始化创建Token
func InitJWT( uname string)(s string){
	mySigningKey:=[]byte("usertest")
	c:=MyClaims{
		Username: uname,
		StandardClaims:jwt.StandardClaims{
			NotBefore: time.Now().Unix()-60,
			ExpiresAt: time.Now().Unix()+60*60*2,
			Issuer: "user2",
		},
	}
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,c)
	s,e:=token.SignedString(mySigningKey)
	if e!=nil{
		fmt.Println(e.Error())
		return
	}else{
		fmt.Println(s)
		return s
	}
}
//解析token
func ParseJwt(s string)(claims *MyClaims,err error){
	mySigningKey:=[]byte("usertest")
	t,err:=jwt.ParseWithClaims(s,&MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey,nil
	})
	if err!=nil{
		fmt.Println(err.Error())
		return nil,err
	}else{
		fmt.Println(t.Claims.(*MyClaims).Username)
		return t.Claims.(*MyClaims),nil
	}
}
//登陆验证
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		authString := c.Request.Header.Get("Authorization")

		kv := strings.Split(authString, " ")
		if len(kv) != 2 || kv[0] != "Bearer" {
			//result := models.UnauthorizedResult()
			//c.JSON(200, result)
			c.Abort()
			return
		}

		tokenString := kv[1]

		// Parse token
		token, err := ParseJwt(tokenString)

		if err != nil {
			//result := models.UnauthorizedResult()
			//c.JSON(200, result)
			c.Abort()
			return
		}
		if time.Now().Unix()>token.ExpiresAt{
			fmt.Println("token已过期")
			c.JSON(http.StatusOK,gin.H{
				"err":err,
			})
			c.Abort()
			return
		}
		c.Set("username", token.Username)
		c.Next()
	}
}
//管理员权限中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

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
		token, _:= ParseJwt(tokenString)
		role:=models.CheckAdmin(token.Username)
		if role>1{
			c.JSON(http.StatusOK,gin.H{
				"status":"错误",
				"message":"无权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}





func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}