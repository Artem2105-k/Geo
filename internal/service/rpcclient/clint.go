package rpcclient

import (
	"errors"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"

	"studentgit.kata.academy/ar.konovalov202_gmail.com/rpc/general"
)

var ErrConnectionFailed = errors.New("connection failed")

type RPCClient struct {
	client *rpc.Client
}

func NewRPCClient(addr string) (*RPCClient, error) {
	conn, err := jsonrpc.Dial("tcp", addr)
	if err != nil {
		log.Printf("Failed to connect to RPC server: %v", err)
		return nil, ErrConnectionFailed
	}
	log.Println("Connected to RPC server at", addr)
	return &RPCClient{client: conn}, nil
}

func (c *RPCClient) AddressSearch(input string) ([]*general.Address, error) {
	var result []*general.Address
	err := c.client.Call("RPCServer.AddressSearch", &SearchArgs{Query: input}, &result)
	return result, err
}

func (c *RPCClient) GeoCode(lat, lng string) ([]*general.Address, error) {
	var result []*general.Address
	err := c.client.Call("RPCServer.GeoCode", &GeocodeArgs{Lat: lat, Lng: lng}, &result)
	return result, err
}

type SearchArgs struct {
	Query string
}

type GeocodeArgs struct {
	Lat string
	Lng string
}
