// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nrpc

import (
	"errors"
	"io"
	"github.com/phpyandong/im/pkg/rpc"
	"sync"
	"github.com/phpyandong/im/pkg/protobufs"
	"bufio"
	"encoding/json"
)

var errMissingParams = errors.New("jsonrpc: request body missing params")

type serverCodec struct {
	dec *protobufs.Transprot // for reading JSON values
	enc *protobufs.Transprot // for writing JSON values
	trans  *protobufs.Transprot
	c   io.Closer

	// temporary work space
	req serverRequest

	// JSON-RPC clients can use arbitrary json values as request IDs.
	// Package rpc expects uint64 request IDs.
	// We assign uint64 sequence numbers to incoming requests
	// but save the original request ID in the pending map.
	// When rpc responds, we use the sequence number in
	// the response to find the original request ID.
	mutex   sync.Mutex // protects seq, pending
	seq     uint64
	pending map[uint64]*json.RawMessage
}

type nserverCodec struct {
	rwc    io.ReadWriteCloser

	dec     *protobufs.Transprot
	enc     *protobufs.Transprot

	encBuf *bufio.Writer
}
// NewServerCodec returns a new rpc.ServerCodec using JSON-RPC on conn.
func NewServerCodec(conn io.ReadWriteCloser) rpc.ServerCodec{
	encBuf := bufio.NewWriter(conn)

	return &nserverCodec{
		rwc:       conn,
		//pending: make(map[uint64]string),
		dec: &protobufs.Transprot{EncBuf:conn},
		enc: &protobufs.Transprot{EncBuf:conn},
		encBuf:encBuf,
	}
}

type serverRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
	Id     *json.RawMessage `json:"id"`
}

func (r *serverRequest) reset() {
	r.Method = ""
	r.Params = nil
	r.Id = nil
}

type serverResponse struct {
	Id     *json.RawMessage `json:"id"`
	Result interface{}      `json:"result"`
	Error  interface{}      `json:"error"`
}

func (c *nserverCodec) ReadRequestHeader(r *rpc.Request) error {
	c.dec.Decode()
	//fmt.Println("xxhead:%v",r)

	//c.req.reset()
	//if err := c.dec.Decode(&c.req); err != nil {
	//	return err
	//}
	//r.ServiceMethod = c.req.Method
	//
	//// JSON request id can be any JSON value;
	//// RPC package expects uint64.  Translate to
	//// internal uint64 and save JSON on the side.
	//c.mutex.Lock()
	//c.seq++
	//c.pending[c.seq] = c.req.Id
	//c.req.Id = nil
	//r.Seq = c.seq
	//c.mutex.Unlock()

	return nil
}

func (c *nserverCodec) ReadRequestBody(x interface{}) error {
	c.dec.Decode()

	//fmt.Println("xxbody :%v",x)
	return nil
	//if x == nil {
	//	return nil
	//}
	//if c.req.Params == nil {
	//	return errMissingParams
	//}
	//// JSON params is array value.
	//// RPC params is struct.
	//// Unmarshal into array containing struct for now.
	//// Should think about making RPC more general.
	//var params [1]interface{}
	//params[0] = x
	//return json.Unmarshal(*c.req.Params, &params)
}

var null = json.RawMessage([]byte("null"))

func (c *nserverCodec) WriteResponse(r *rpc.Response, x interface{}) error {
	//c.mutex.Lock()
	//b, ok := c.pending[r.Seq]
	//if !ok {
	//	c.mutex.Unlock()
	//	return errors.New("invalid sequence number in response")
	//}
	//delete(c.pending, r.Seq)
	//c.mutex.Unlock()
	//
	//if b == nil {
	//	// Invalid request so no id. Use JSON null.
	//	b = &null
	//}
	//resp := serverResponse{Id: b}
	//if r.Error == "" {
	//	resp.Result = x
	//} else {
	//	resp.Error = r.Error
	//}

	//return c.enc.Encode(resp)
	return nil
}

func (c *nserverCodec) Close() error {
	return nil
}

// ServeConn runs the JSON-RPC server on a single connection.
// ServeConn blocks, serving the connection until the client hangs up.
// The caller typically invokes ServeConn in a go statement.
func ServeConn(conn io.ReadWriteCloser) {
	rpc.ServeCodec(NewServerCodec(conn))
}
