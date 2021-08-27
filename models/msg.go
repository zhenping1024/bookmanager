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
func PubilishMsg(username string,msg string){
	data,_:=json.Marshal(Msg{MsgContext:msg ,MsgUser: username,Creattime: time.Now()})
	n, err := Rdb.Publish("Admin",data ).Result()
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