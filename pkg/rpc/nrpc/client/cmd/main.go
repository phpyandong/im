package main

import (
"fmt"
	"github.com/phpyandong/im/pkg/rpc/nrpc/server"
	"github.com/phpyandong/im/pkg/rpc/nrpc/client"
)

func main() {
	args := server.Args{
		A:  1,
		B:  2,
		Op: "+",
	}

	var reply server.Reply

	err := client.Call(args, &reply)
	display(err, args, reply)

	args = server.Args{
		A: 1, B: 0, Op: "/",
	}
	err = client.Call(args, &reply)
	display(err, args, reply)
}

func display(err error, args server.Args, reply server.Reply) {
	if err != nil {
		fmt.Printf("err:%v\n", err)
	} else {
		// 如果err不为nil,这里的reply是上个调用的值
		// 因此可能会出现1.00 / 0.00=3.00
		fmt.Printf("%.2f %s %.2f=%.2f\n", args.A, args.Op, args.B, reply.Data)
	}
}

