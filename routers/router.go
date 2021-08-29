package routers

import (
	v1 "bookmanager/api/v1"
	"bookmanager/middleware"
	"bookmanager/models"
	"bookmanager/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)
func InitRouter(){
	//gin.SetMode(utils.AppMode)
	r:=gin.Default()
	r.LoadHTMLGlob("statics/dist/*.html")
	//r.LoadHTMLFiles("statics/dist/static/*/*")
	r.Static("statics","./statics")
	r.Static("/static","./statics/dist/static")
	//r.Static("static/admin/css","./statics/static")

	r.GET("index",func(c*gin.Context){
		c.HTML(http.StatusOK,"index.html",nil)
	})
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://8.130.51.87:3030/","http://8.130.51.87:3000/"},
		AllowMethods:     []string{"PUT", "GET","DELETE","POST"},
		AllowHeaders:     []string{"Origin","Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 24 * time.Hour,
	}))

	router :=r.Group("api/v1")
	router.Use(middleware.JwtAuth())
	{
		//用户模块模块
		router.PUT("user/:id",v1.EditUser)
		router.DELETE("user/:id",middleware.AdminAuth(),v1.DeleteUser)
		router.GET("user/books",v1.Getborrow)
		router.POST("user/book/borrow/:id",v1.BorrowBook)
		router.POST("user/book/return/:id",v1.ReturnBook)
		//分类模块模块
		router.POST("category/add",v1.AddCategory)
		router.PUT("category/:id",v1.EditCategory)
		router.DELETE("category/:id",v1.DeleteCategory)
		//普通管理员模块
		router.POST("book/add",middleware.AdminAuth(),v1.AddBook)
		router.PUT("book/:id",middleware.AdminAuth(),v1.EditBook)
		router.PUT("book/:id/up",middleware.AdminAuth(),v1.Upbook)
		router.PUT("book/:id/down",middleware.AdminAuth(),v1.Downbook)
		router.DELETE("book/:id",middleware.AdminAuth(),v1.DeleteBook)
		router.PUT("admin/add/:id",middleware.AdminAuth(),v1.AddAdmin)
		router.PUT("admin/creat",middleware.AdminAuth(),v1.CreatAdmin)
		router.PUT("admin/delete/:id",middleware.AdminAuth(),v1.DeleteAdmin)
		//搜索模块
		router.POST("user/book/search",v1.SearchUserBook)
		//评论
		router.POST("comment/:bookid",v1.SendComment)
		router.POST("toadmin",v1.ToAdmin)
	}
	router2 :=r.Group("api/v1")
	{
		router2.POST("user/search",v1.SearchUser)
		router2.POST("book/search",v1.SearchBook)
		router2.POST("admin/search",v1.SearchAdmin)
		router2.POST("user/add",v1.AddUser)
		router2.GET("/users",v1.GetUsers)
		router.GET("user",v1.GetUser)
		router2.GET("/admins",v1.GetAdmins)
		router2.GET("/category",v1.GetCategory)
		router2.GET("books",v1.GetBooks)
		router2.GET("books/up",v1.GetUpbooks)
		router2.GET("book/id",v1.GetBook)
		router2.POST("login",v1.Login)
		router2.GET("comment/:bookid",v1.GetComment)
		router2.GET("admin/msg",v1.GetMsg)
	}
	go func() {
		models.GetMsg()
	}()
	r.Run(utils.HttpPort)
}
