package server

import "net"

func main() {

	addr := "0.0.0.0:8080"

	tcpAddr, err := net.ResolveTCPAddr("tcp",addr)

	if err != nil {
		log.Fatalf("net.ResovleTCPAddr fail:%s", addr)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatalf("listen %s fail: %s", addr, err)
	} else {

		log.Println("rpc listening", addr)
	}


	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("listener.Accept error:", err)
			continue
		}

		go handleConnection(conn)

	}

}