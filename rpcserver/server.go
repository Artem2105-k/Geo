// rpcserver/server.go

package rpcserver

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"studentgit.kata.academy/ar.konovalov202_gmail.com/rpc/general"
)

func StartRpcServer(geoSer general.GeoProvider) {
	rpc.Register(&RPCServer{GeoSer: geoSer})

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Error starting RPC server:", err)
	}
	log.Println("RPC Server started on :8081")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go jsonrpc.ServeConn(conn)
	}
}
