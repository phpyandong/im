package protobufs

import (
	"io"
	"errors"
	"encoding/binary"
	"bufio"
)

type Transprot struct {
	encBuf *bufio.Writer
	Data []byte
}
func (t *Transprot) Write(p []byte) (int, error) {
	t.Data = append(t.Data, p...)
	return len(p), nil
}
var ProtoBeginBytes = []byte{18, 17, 13, 10, 9}
var ProtoEndBytes = []byte{9, 10, 13, 17, 18}
var HeadTailLength = len(ProtoBeginBytes) + len(ProtoEndBytes)
/**
| header |  dataLength  | data                 |  end  |
| header |  dataLength  | EventType| RetData   |  end  |
| 5 bytes|   4 bytes    | 4 bytes  | N-4 bytes |5 bytes|
*/
const (
	LenHead = 5
	LenData = 4

	LenEnd  = 5
)
type transportResp struct {
	//tmpIndex int
	head 		[]byte
	Serial 		int32
	ResultData 	[]byte
	tail		[]byte
}
/**
| header |  dataLength  | data                 |  end  |
| header |  dataLength  | EventType| RetData   |  end  |
| 5 bytes|   4 bytes    | 4 bytes  | N-4 bytes |5 bytes|
*/
func (transprot *Transprot) Encode(data []byte)([]byte,error){
	totalLen := HeadTailLength + LenData + len(data)
	transportData := make([]byte,totalLen)
	//数据头部
	copy(transportData[0:LenHead],ProtoBeginBytes)
	//数据长度
	binary.BigEndian.PutUint32(transportData[LenHead:LenHead+LenData],uint32(totalLen))
	//数据
	copy(transportData[LenHead+LenData:totalLen-LenEnd],data)
	//数据尾部
	copy(transportData[totalLen-LenEnd:totalLen],ProtoEndBytes)
	if _, err := transprot.Write(transportData); err != nil {
		return nil,err
	}

	return transportData,nil
}

/**
| header |  header  	| data                 |  end  |
| begin  |  serial 		| EventType| RetData   |  end  |
| 5 bytes|   4 bytes    | 4 bytes  | N-4 bytes |5 bytes|
| 01234  |56789			|
*/
func (resp *transportResp) Decode(r io.Reader) error{
		var err error
		_,err = io.ReadFull(r,resp.head)
		if err != nil {
			return err
		}
		for i:=0;i<5;i++{
			if ProtoBeginBytes[i] !=  resp.head[i] {
				return errors.New("head error")
			}
		}
		dataLength := binary.LittleEndian.Uint32(resp.head[1:5])//0 1 2 3,4
		data := make([]byte,int(dataLength))

		copy(data[0:5],resp.head)
		_,err = io.ReadFull(r,data[5:])
		resp.Serial = int32(binary.LittleEndian.Uint32(data[5:9]))
		resp.ResultData = data
		//读取io流中的tail
		_,err = io.ReadFull(r,resp.tail)
		for i:=0;i<5 ;i++  {
			if resp.tail[i] != ProtoEndBytes[i]{
				return errors.New("tail err")
			}
		}
		return nil
}

func NewTransportResp() *transportResp{
	return &transportResp{
		head :make([]byte,5),
		tail :make([]byte,5),
	}
}