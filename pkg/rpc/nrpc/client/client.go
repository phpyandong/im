package client

import (
"log"
	"github.com/phpyandong/im/pkg/rpc/nrpc/server"
	"github.com/phpyandong/im/pkg/rpc/nrpc"
)

func Call(args server.Args, reply *server.Reply) error {
	conn, err := nrpc.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatalln("dailing error: ", err)
		return err
	}

	defer conn.Close()

	// 调用远程的Calc的Compute方法
	err = conn.Call("Calc.Compute", args, &reply)
	return err
}
