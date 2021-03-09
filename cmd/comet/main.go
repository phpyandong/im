package main

import "github.com/phpyandong/im/comet"

func main(){
	server := comet.NewServer()
	server.Run()
}