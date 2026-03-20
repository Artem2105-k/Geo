package rpcserver

import "ar.konovalov202_gmail.com/rpc/general"

type RPCServer struct {
	GeoSer general.GeoProvider
}

type SearchArgs struct {
	Query string
}

type GeosodeArgs struct {
	Lat string
	Lon string
}

func (s *RPCServer) AddressSearch(args *SearchArgs, reply *[]*general.Address) error {
	res, err := s.GeoSer.AddressSearch(args.Query)
	if err != nil {
		return err
	}
	*reply = res
	return nil
}

func (s *RPCServer) GeoCode(args *GeosodeArgs, reply *[]*general.Address) error {
	res, err := s.GeoSer.GeoCode(args.Lat, args.Lon)
	if err != nil {
		return err
	}
	*reply = res
	return nil
}
