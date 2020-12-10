package main

import (
	"github.com/golang/protobuf/proto"
	"mxs/bin/protoc/sample"
	"mxs/log"
)

func main() {
	person := &sample.Person{
		Name:   "lsill",
		Age:    24,
		Emails: []string{"924867750@qq.com"},
		Phones: [] *sample.PhoneNumber{
			&sample.PhoneNumber{
				Number: "17611163990",
				Type:   sample.PhoneType_MOBILE,
			},
		},
	}
	data, err := proto.Marshal(person)
	if err != nil {
		log.Error("marshal err :%v", err)
	}
	newdata := &sample.Person{}
	err = proto.Unmarshal(data, newdata)
	if err != nil {
		log.Error("unmarshal err %v", err)
	}
	log.Debug("%v", newdata)
}
