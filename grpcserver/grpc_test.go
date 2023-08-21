package grpcserver

import (
	"context"
	pb "ibdatabase/grpcserver/proto"
	"io"
	"log"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
)

var addr string = "localhost:9393"

func TestGRPCServer(t *testing.T) {
	tls := true
	opts := []grpc.DialOption{}

	if tls {
		certFile := "../ssl/ca.crt"
		cres, err := credentials.NewClientTLSFromFile(certFile, "")

		if err != nil {
			log.Fatal("failed cert", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(cres))
	}

	conn, err := grpc.Dial(addr, opts...)

	if err != nil {
		log.Fatal("failed", err)
	}

	defer conn.Close()
	clt := pb.NewStockServiceClient(conn)

	strm, err := clt.GetStocks(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatal("failed", err)
	}

	for {
		res, err := strm.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Something happened: %v\n", err)
		}

		log.Println(res.StockSymbol)
	}
}
