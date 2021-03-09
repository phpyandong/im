// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jsonrpc implements a JSON-RPC 1.0 ClientCodec and ServerCodec
// for the rpc package.
// For JSON-RPC 2.0 support, see https://godoc.org/?q=json-rpc+2.0
package nrpc

import (
	"encoding/json"
	"io"
	"net"
	"net/rpc"
	"sync"
	"github.com/phpyandong/im/pkg/protobufs/protobuf"
	"github.com/phpyandong/im/pkg/protobufs/proto"
	"github.com/phpyandong/im/pkg/protobufs"
	"github.com/phpyandong/im/conf"
	"bufio"
	"fmt"
)

type clientCodec struct {
	conn    io.ReadWriteCloser  // todo net.Conn 不光可读可写可关闭，还包含其他
	//dec *json.Decoder // for reading JSON values
	//enc *json.Encoder // for writing JSON values
	c   io.Closer
	proto *protobuf.Protobuf
	dec  *protobufs.Transprot
	enc  *protobufs.Transprot

	// temporary work space
	req  clientRequest
	resp clientResponse

	// JSON-RPC responses include the request id but not the request method.
	// Package rpc expects both.
	// We save the request method in pending when sending a request
	// and then look it up by request ID when filling out the rpc Response.
	mutex   sync.Mutex        // protects pending
	pending map[uint64]string // map request id to method name
}
type proClientCodec struct {
	rwc    io.ReadWriteCloser

	dec     *protobufs.Transprot
	enc     *protobufs.Transprot

	encBuf *bufio.Writer
}


func (c *proClientCodec) ReadResponseBody(interface{}) error {
	fmt.Println("ReadResponseBody: ")
	return nil
}

func (c *proClientCodec) WriteRequest(r *rpc.Request, body interface{}) (err error) {

	conf := conf.NewConf()
	req := &dto.RequestProto{
		Lookup:    "xx",
		Method:    r.ServiceMethod,
		ParamList: nil,
	}
	proto := protobuf.NewSendProtobuf(req, conf.Service, int32(r.Seq))
	protodata, err := proto.ToBytes()
	if err != nil {
		return err
	}
	_, err = c.enc.Encode(protodata)
	fmt.Println("==================")
	return nil
}


// NewClientCodec returns a new rpc.ClientCodec using JSON-RPC on conn.
func NewClientCodec(conn io.ReadWriteCloser) rpc.ClientCodec {
	encBuf := bufio.NewWriter(conn)

	return &proClientCodec{
		//dec:     json.NewDecoder(conn),
		//enc:     json.NewEncoder(conn),
		rwc:       conn,
		//pending: make(map[uint64]string),
		dec: &protobufs.Transprot{EncBuf:conn},
		enc: &protobufs.Transprot{EncBuf:conn},
		encBuf:encBuf,

	}
}

type clientRequest struct {
	Method string         `json:"method"`
	Params [1]interface{} `json:"params"`
	Id     uint64         `json:"id"`
}

//func (c *clientCodec) WriteRequest(r *rpc.Request, param interface{}) error {
//	fmt.Println("========writeRequest============")
//	c.mutex.Lock()
//	c.pending[r.Seq] = r.ServiceMethod
//	c.mutex.Unlock()
//	c.req.Method = r.ServiceMethod
//	c.req.Params[0] = param
//	c.req.Id = r.Seq
//	req := &dto.RequestProto{
//		Lookup:"xx",
//		Method: r.ServiceMethod,
//		ParamList:nil,
//	}
//	conf := conf.NewConf()
//	proto := protobuf.NewSendProtobuf(req,conf.Service,int32(r.Seq))
//	protodata,err := proto.ToBytes()
//	if err != nil {
//		return err
//	}
//	_,err = c.trans.Encode(protodata)
//
//	fmt.Println(c.trans.Data)
//
//	return c.enc.Encode(&c.req)
//}

type clientResponse struct {
	Id     uint64           `json:"id"`
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
}

func (r *clientResponse) reset() {
	r.Id = 0
	r.Result = nil
	r.Error = nil
}

func (c *proClientCodec) ReadResponseHeader(r *rpc.Response) error {
	//c.resp.reset()
	//if err := c.dec.Decode(&c.resp); err != nil {
	//	return err
	//}
	//
	//c.mutex.Lock()
	//r.ServiceMethod = c.pending[c.resp.Id]
	//delete(c.pending, c.resp.Id)
	//c.mutex.Unlock()
	//
	//r.Error = ""
	//r.Seq = c.resp.Id
	//if c.resp.Error != nil || c.resp.Result == nil {
	//	x, ok := c.resp.Error.(string)
	//	if !ok {
	//		return fmt.Errorf("invalid error %v", c.resp.Error)
	//	}
	//	if x == "" {
	//		x = "unspecified error"
	//	}
	//	r.Error = x
	//}
	return nil
}

func (c *clientCodec) ReadResponseBody(x interface{}) error {
	if x == nil {
		return nil
	}
	return json.Unmarshal(*c.resp.Result, x)
}

func (c *proClientCodec) Close() error {
	//return c.c.Close()
	return nil
}

// NewClient returns a new rpc.Client to handle requests to the
// set of services at the other end of the connection.
func NewClient(conn io.ReadWriteCloser) *rpc.Client {
	return rpc.NewClientWithCodec(NewClientCodec(conn))
}

// Dial connects to a JSON-RPC server at the specified network address.
func Dial(network, address string) (*rpc.Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return NewClient(conn), err
}
