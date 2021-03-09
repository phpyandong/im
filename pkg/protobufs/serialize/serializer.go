package serialize

import (
	"reflect"
	"encoding/json"
	"github.com/phpyandong/im/pkg/protobufs/proto"
)
type SerializerType struct {

}
type InterfaceSerializer interface {
	Serialize(obj interface{}) ([]byte, error)
	Deserialize(data []byte, t reflect.Type) (interface{}, error)
}
type JsonSerializer struct {
}

func (js *JsonSerializer) Serialize(strutObj interface{}) ([]byte, error) {
	return json.Marshal(strutObj)
}

func (js *JsonSerializer) Deserialize(data []byte, t reflect.Type) (interface{}, error) {
	request := &dto.RequestProto{}
	return request,json.Unmarshal(data,request)
}

func NewJsonSerializer ()( *JsonSerializer,error){
	return &JsonSerializer{},nil
}
/**
获取序列化类型对象
 */
func GetSerializerByType(types byte)( InterfaceSerializer,error){
	return NewJsonSerializer()
}
