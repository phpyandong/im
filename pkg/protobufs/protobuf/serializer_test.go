package protobuf

import (
	"github.com/phpyandong/im/pkg/protobufs/proto"
	"testing"
	"container/list"
	"fmt"
)

func TestSerizlizer(t *testing.T){
	//serializer,_ := serialize.GetSerializerByType(1)
	params := &dto.Params{
		Key :"StringKey",
		Value:"value",
	}
	lis := list.New()
	lis.PushBack(params)
	req := &dto.RequestProto{
		Lookup:"xx",
		Method:"methodName",
		ParamList:lis,
	}
	protoExpect := &Protobuf{
		Version:2,
		Serial :12,
		ServiceID:11,
		PType 	:  2,
		SerialzerType:1, //序列化类型 12	1byte: 1 json 6 自定义
		Os			:4,
		Object :req,
	}
	fmt.Printf("expect:%+v \n",protoExpect)
	sendData,err := protoExpect.ToBytes()
	if err != nil {
		t.Errorf("toBytes err :%v",err)
	}
	//fmt.Println(sendData)
	actProto := &Protobuf{
	}
	err = actProto.ToProtobuf(sendData)
	if err != nil {
		t.Errorf("toObj err :%v\n",err)
	}
	fmt.Printf("actual:%+v\n ",actProto)

	fmt.Printf("protobuf act:%+v\n exp :%+v\n",actProto.Object , protoExpect.Object)

}