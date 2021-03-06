package models

import (
	"bookmanager/utils"
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"log"
	"mime/multipart"
)

var AK =utils.AccessKey
var SK =utils.SecretKey
var Bucket = utils.Bucket
var Imgurl = utils.QiniuSever

func UpLoadFile(file multipart.File,filesize int64)(string,error){
	putPolicy:=storage.PutPolicy{
		Scope: Bucket,
	}
	mac:=qbox.NewMac(AK,SK)
	upToken:=putPolicy.UploadToken(mac)

	cfg:=storage.Config{
		Zone: &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS: false,

	}
	formUploader:=storage.NewFormUploader(&cfg)
	ret:=storage.PutRet{}
	err:=formUploader.PutWithoutKey(context.Background(),&ret,upToken,file,filesize,nil)

	url:=Imgurl+ret.Key
	return url,err
}
func UpLoadBook(file multipart.File,filesize int64,id int)(string,error){
	putPolicy:=storage.PutPolicy{
		Scope: Bucket,
	}
	mac:=qbox.NewMac(AK,SK)
	upToken:=putPolicy.UploadToken(mac)

	cfg:=storage.Config{
		Zone: &storage.ZoneHuanan,
		UseCdnDomains: false,
		UseHTTPS: false,

	}
	formUploader:=storage.NewFormUploader(&cfg)
	ret:=storage.PutRet{}

	err:=formUploader.PutWithoutKey(context.Background(),&ret,upToken,file,filesize,nil)

	url:=Imgurl+ret.Key
	SaveBook(id,url)
	return url,err
}
func SaveBook(id int,url string){
	var maps=make(map[string]interface{})
	var book Book
	maps["bookurl"]=url
	err:=DB.Model(&book).Where("id = ?",id).Update(maps).Error
	if err!=nil{
		log.Fatal(err)
	}
}