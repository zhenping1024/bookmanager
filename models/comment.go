package models

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
)

type Comment struct{
	gorm.Model
	Commentuser string
	Context		string
	Head  string
}

//发布评论
func PublishComment(bookid int,c Comment)(Comment,error){
	//_,err:=Rdb.LPush(strconv.Itoa(bookid),c).Result()
	//return err
	data,_:=json.Marshal(c)
	result,err:=Rdb.LPush(strconv.Itoa(bookid),data).Result()
	if err!=nil {
		fmt.Println(err)
		return Comment{},err
	}
	fmt.Println(bookid,result)
	return c,err
}
//显示评论
func GetComment(bookid int,pagesize int,pagenum int)([]Comment,error,int){
	result:=Rdb.LRange(strconv.Itoa(bookid),int64(pagesize*(pagenum-1)),int64(pagesize*pagenum-1)).Val()
	sum:=Rdb.LLen(strconv.Itoa(bookid)).Val()
	var showresult Comment
	var resultstruct []Comment
	for i:=0;i< len(result);i++{
		_=json.Unmarshal([]byte(result[i]),&showresult)
		resultstruct=append(resultstruct,showresult)
	}
	fmt.Println(resultstruct)
	return resultstruct,err,int(sum)
}
//删除评论