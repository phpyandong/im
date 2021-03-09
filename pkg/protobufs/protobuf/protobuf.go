package protobuf

import (
	"github.com/phpyandong/im/conf"
	"github.com/phpyandong/im/pkg/protobufs/serialize"
	"encoding/binary"
	"bytes"
	"reflect"
	"github.com/phpyandong/im/pkg/protobufs/proto"
	"fmt"
)
const (
	HEAD_STACK_LENGTH int  = 14
	Version       = 1
	Length      = 4
	SerialId     = 4
	ServerId      = 1
	PType       = 1
	CompressType  = 1
	SerializeType = 1
	OS      = 1

)
type Protobuf struct {
	Version       byte        //版本号			0	 	1byte : 2
	Length        int32       //协议的长度	1-4		4byte : 111
	Serial        int32       //序列号		5-8		4byte :
	ServiceID     byte        //服务编号		9		1byte
	PType         byte        //消息体类型	10		1byte: 1response 2 request
	CompressType  byte        //压缩算法 11		1byte: 序列化规则
	SerialzerType byte        //序列化类型 12	1byte: 1 json 6 自定义
	Os            byte        //平台 0 ios 1 java 2 js 3 php 4 golang
	Content       []byte      //发送内容
	Object        interface{} //传输的对象
}

func NewRequestProtobuf(requestContent []byte, config *conf.ServiceConf, serialId int32) *Protobuf {
	return &Protobuf{
		Version:       2,
		Serial:        serialId,
		ServiceID:     config.ServiceId,
		PType:         2,
		CompressType:  0,
		SerialzerType: 6,
		Os:            4,
		Content:       requestContent,
	}
}

func NewSendProtobuf(req *dto.RequestProto, config *conf.ServiceConf, serialId int32) *Protobuf {
	return &Protobuf{
		Version:       2,
		Serial:        serialId,
		ServiceID:     config.ServiceId,
		PType:         2,
		CompressType:  0,
		SerialzerType: 6,
		Os:            4,
		Object:		req,
		//Content:       requestContent,
	}
}

func (protobuf *Protobuf) ToBytes() ([]byte, error) {
	var binaryEndian = binary.LittleEndian
	var serializer serialize.InterfaceSerializer
	var objectData []byte
	var protobufLen int
	var err error
	if serializer, err = serialize.GetSerializerByType(protobuf.SerialzerType); err != nil {
		return nil, err
	}
	if objectData, err = serializer.Serialize(protobuf.Object); err != nil {
		return nil, err
	}
	protobuf.Content = objectData
	protobufLen = HEAD_STACK_LENGTH + len(objectData)
	protobuf.Length = int32(protobufLen)
	buf := new(bytes.Buffer)
	buf.Grow(protobufLen)
	buf.WriteByte(byte(protobuf.Version))
	binary.Write(buf,binaryEndian,protobuf.Length)
	binary.Write(buf,binaryEndian,protobuf.Serial)
	buf.WriteByte(byte(protobuf.ServiceID))
	buf.WriteByte(byte(protobuf.PType))
	buf.WriteByte(byte(protobuf.CompressType))
	buf.WriteByte(byte(protobuf.SerialzerType))
	buf.WriteByte(byte(protobuf.Os))
	buf.Write(objectData)
	fmt.Printf("expect2 %+v\n",protobuf)
	return buf.Bytes(), nil
}

func (protobuf *Protobuf) ToProtobuf(mes []byte)(error){
	buf := bytes.NewBuffer(mes)
	var err error
	var binaryEndian = binary.LittleEndian
	binary.Read(buf,binaryEndian,&(protobuf.Version))
	err = binary.Read(buf,binaryEndian,&(protobuf.Length))
	binary.Read(buf,binaryEndian,&(protobuf.Serial))
	binary.Read(buf,binaryEndian,&(protobuf.ServiceID))
	binary.Read(buf,binaryEndian,&(protobuf.PType))
	binary.Read(buf,binaryEndian,&(protobuf.CompressType))
	binary.Read(buf,binaryEndian,&(protobuf.SerialzerType))
	binary.Read(buf,binaryEndian,&(protobuf.Os))
	protobuf.Content = buf.Bytes()
	var serializer serialize.InterfaceSerializer
	if serializer,err = serialize.GetSerializerByType(protobuf.SerialzerType); err != nil{
		return err
	}
	dotObj,err := serializer.Deserialize(protobuf.Content,reflect.TypeOf(&dto.RequestProto{}))
	if err != nil {
		return err
	}
	protobuf.Object = dotObj
	//fmt.Println("dotObj:",dotObj)

	//fmt.Printf("protobuf: %+v \n dot :%+v\n",protobuf,dotObj)
	//
	////var objectType reflect.Type
	////if objectType,err =
	return nil
}
