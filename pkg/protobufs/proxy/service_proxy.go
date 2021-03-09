package proxy

import (
	"github.com/phpyandong/im/conf"
	"container/list"
	"github.com/phpyandong/im/pkg/protobufs/proto"
	"github.com/phpyandong/im/pkg/protobufs/util"
	"github.com/phpyandong/im/pkg/protobufs/protobuf"
	"github.com/pkg/errors"
	"fmt"
)

type ServiceProxy struct{
	Config *conf.ServiceConf

}
func (serviceProxy *ServiceProxy) Invoke(methodName string,lookup string,params *list.List)( *dto.ResponseProto,error){
	var err error
	requestProto := &dto.RequestProto{
		Lookup : lookup,
		Method: methodName,
		ParamList:params,
	}
	serial := int32(util.GetSeq())
	sendProto := protobuf.NewSendProtobuf(requestProto,serviceProxy.Config,serial)
	var sendProtoData []byte
	if sendProtoData ,err =sendProto.ToBytes();err != nil{
		return nil,errors.WithMessage(err,fmt.Sprintf("Invoke sendProto:%+v",sendProtoData))
	}

	//todo 发送数据到socket
}