package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

type Category struct{
	gorm.Model
	TypeName string
}
//查寻分类是否存在
func CheckCategory(name string)error{
	var cate Category
	DB.Select("id").Where("type_name = ?",name).First(&cate)
	if cate.ID>0{

		err:=fmt.Errorf("用户已存在")
		return err
	}
	return nil
}
//查询分类下书籍数目

//新增分类
func CreatCategory(data *Category)error{

	err:=DB.Create(&data).Error
	if err!=nil{
		return err
	}
	return nil
}
//查询分类列表
func GetCategoey(PageSize int,Pagenum int)[]Category{
	var cate []Category
	err:=DB.Limit(PageSize).Offset((Pagenum-1)*PageSize).Find(&cate).Error
	if err!=nil{
		return nil
	}
	return cate
}
//删除分类
func DeleteCategory(id int)int{
	var cate Category
	err :=DB.Where("id = ?",id).Delete(&cate).Error
	if err!=nil{
		log.Fatal("删除错误")
		return -1
	}
	ID,_:=strconv.Atoi(fmt.Sprint(cate.ID))
	return ID
}
//编辑分类
func EditCategory(id int,u *Category){
	var maps=make(map[string]interface{})
	var cate Category
	maps["type_name"]=u.TypeName
	err:=DB.Model(&cate).Where("id = ?",id).Update(maps).Error
	if err!=nil{
		log.Fatal(err)
	}
}
