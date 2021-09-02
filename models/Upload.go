package models

import (
	"bookmanager/utils"
	"context"
	"fmt"
	"mime/multipart"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
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
	//putExtra:=storage.PutExtra{}
	formUploader:=storage.NewFormUploader(&cfg)
	ret:=storage.PutRet{}

	err:=formUploader.PutWithoutKey(context.Background(),&ret,upToken,file,filesize,nil)

	url:=Imgurl+ret.Key
	fmt.Println("imgurl is",Imgurl)
	return url,err
}