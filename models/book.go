package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

type Book struct{
gorm.Model
//Category Category	`gorm:foreignkey:Cid`
BookName string `json:"bookname" form:"bookname"`
BookPrice string `json:"price"`
BookImag string `json:"imag" form:"imag"`
//Cid int
BookType string`json:"booktype" form:"booktype"`
State string `json:"state" gorm:"default:'no'"`
Sum int `json:"sum" form;"sum"`
BorrowSum int`json:"borrowsum"`
Users []User`gorm:"many2many:user_book"`
Author string `json:"author"`
Introduce string `introduce`
}
//新增书籍
func CreatBook(data *Book)error{
	err:=DB.Create(&data).Error
	if err!=nil{
		return err
	}
	return nil
}
//上架图书
func UpBook(id int,u *Book){
	var maps=make(map[string]interface{})
	var book Book
	maps["state"]="yes"
	err:=DB.Model(&book).Where("id = ?",id).Update(maps).Error
	if err!=nil{
		log.Fatal(err)
	}
}
//下架图书
func DownBook(id int,u *Book){
	var maps=make(map[string]interface{})
	var book Book
	maps["state"]="no"
	err:=DB.Model(&book).Where("id = ?",id).Update(maps).Error
	if err!=nil{
		log.Fatal(err)
	}
}
//查询分类下所有书籍
func GetCateBook(id string,pagesize int,pagenum int)[]Book{
	var books []Book
	DB.Limit(pagesize).Offset((pagenum-1)*pagesize).Where("book_type = ?",id).Find(&books)
	return books
}
//查询单一书籍
func GetBookInfo (id int)Book{
	var book Book
	err:=DB.Where("id = ?",id).First(&book).Error
	if err!=nil{
		return Book{}
	}
	return book
}
//查询书籍列表
func GetBooks(PageSize int,Pagenum int)([]Book,int){
	var book []Book
	var sum int
	DB.Find(&book).Count(&sum)
	err:=DB.Limit(PageSize).Offset((Pagenum-1)*PageSize).Find(&book).Error
	if err!=nil{
		return nil,-1
	}
	return book, sum
}
//删除分类
func DeleteBook(id int)int{
	var book Book
	err :=DB.Where("id = ?",id).Delete(&book).Error
	if err!=nil{
		log.Fatal("删除错误")
		return -1
	}
	ID,_:=strconv.Atoi(fmt.Sprint(book.ID))
	return ID
}
//编辑书籍
func EditBook(id int,u *Book){
	var book Book
	//maps["state"]=u.State
	err:=DB.Model(&book).Where("id = ?",id).Updates(Book{
		BookName: u.BookName,
		Sum: u.Sum,
		BookImag: u.BookImag,
		BookType: u.BookType}).Error
	if err!=nil{
		log.Fatal(err)
	}
}
//书名搜索书籍
func SearchBook(bookname string,pagesize int,pagenum int)([]Book,error,int){
	var u []Book
	bookname="%"+bookname+"%"
	var sum int
	fmt.Println(bookname)
	DB.Where("book_name LIKE ?",bookname).Find(&u).Count(&sum)
	err:=DB.Limit(pagesize).Offset((pagenum-1)*pagesize).Where("book_name LIKE ?",bookname).Find(&u).Error
	fmt.Println("一共有",sum)
	if err!=nil{
		fmt.Println(err)
		return []Book{},err,0
	}
	return u,err,sum
}
//书名搜索用户已借阅
func SearchBorrowedBook(username string,bookname string,pagesize int,pagenum int)([]Book,error,int){
	var u User
	var books []Book
	var sum int
	bookname="%"+bookname+"%"
	DB.Where("username = ?",username).First(&u)
	sum=DB.Model(&u).Preload("Books").Where("book_name LIKE ?",bookname).Association("Books").Find(&books).Count()
	DB.Limit(pagesize).Offset((pagenum-1)*pagesize).Debug().Model(&u).Preload("Books").Where("book_name LIKE ?",bookname).Association("Books").Find(&books)
	fmt.Println("chadao",sum)
	return books,err,sum
}