package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type  Msg struct {
	MsgUser string
	MsgContext string
	Creattime time.Time
}
type Message struct{
	M []Msg
	Sum int
}
var Ms Message
var UMs Message
func PubilishMsg(username string,msg string){
	data,_:=json.Marshal(Msg{MsgContext:msg ,MsgUser: username,Creattime: time.Now()})
	n, err := Rdb.Publish("Admin",data ).Result()
	SaveMsg(username,data)
	Rdb.Publish(username,data)
	if err != nil{
		fmt.Printf("cuowu",err.Error())
		return
	}
	fmt.Printf("%d clients received the message\n", n)

}
func GetMsg() {
	pubsub := Rdb.Subscribe("Admin")
	defer pubsub.Close()
	for msg := range pubsub.Channel() {
		fmt.Printf("channel=%s message=%s\n,%s", msg.Channel, msg.Payload,msg.Pattern)
		var data Msg
		json.Unmarshal([]byte(msg.Payload),&data)
		Ms.M = append(Ms.M, data)
		fmt.Println(data,"msgs")
	}
}
func ToAdmin(msg Msg){
	data,_:=json.Marshal(msg)
	n, err := Rdb.Publish("Admin",data ).Result()
	if err != nil{
		fmt.Printf("cuowu",err.Error())
		return
	}
	fmt.Printf("%d clients received the message\n", n)

}
//保存消息
func SaveMsg(username string,m []byte)error{
	result,err:=Rdb.LPush(username,m).Result()
	Rdb.LPush("admin",m).Result()
	if err!=nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(username,result)
	return err
}
//显示消息
func GetMsgs(username string)([]Msg,error,int){
	result:=Rdb.LRange(username,0,-1).Val()
	sum:=Rdb.LLen(username).Val()
	var showresult Msg
	var resultstruct []Msg
	for i:=0;i< len(result);i++{
		_=json.Unmarshal([]byte(result[i]),&showresult)
		resultstruct=append(resultstruct,showresult)
	}
	fmt.Println(resultstruct)
	return resultstruct,err,int(sum)
}