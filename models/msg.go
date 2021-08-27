package models

import (
	"encoding/json"
	"fmt"
)

type  Msg struct {
	MsgUser string
	MsgContext string
}
var Msgs []Msg
func PubilishMsg(username string,msg string){
	data,_:=json.Marshal(Msg{MsgContext:msg ,MsgUser: username})
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
		Msgs = append(Msgs, data)
		fmt.Println(data,"msgs")
	}
}