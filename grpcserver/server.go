package grpcserver

import (
	"context"
	"fmt"
	pb "ibdatabase/grpcserver/proto"
	"ibdatabase/models"
	"ibdatabase/services"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type GrpcServer struct {
	Port      string
	EnableTLS bool
	CertFile  string
	CertKey   string
	Services  *services.Service
}

type StockServer struct {
	pb.StockServiceServer
}

var Services *services.Service

func (s GrpcServer) StartGRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.Port))

	if err != nil {
		log.Fatal("failed", err)
	}

	opts := []grpc.ServerOption{}
	if s.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(s.CertFile, s.CertKey)
		if err != nil {
			log.Fatal("failed cert", err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	svr := grpc.NewServer(opts...)
	pb.RegisterStockServiceServer(svr, &StockServer{})

	Services = s.Services

	log.Printf("grpc server started at port: %s", s.Port)

	if err = svr.Serve(lis); err != nil {
		log.Fatal("failed", err)
	}
}

func (s *StockServer) GetStocks(_ *emptypb.Empty, stream pb.StockService_GetStocksServer) error {
	// log.Println("GetStocks called")
	stks, err := Services.GetStocks()
	if err != nil {
		return err
	}

	for _, stk := range stks {
		stream.Send(getPBStockByStock(stk))
	}

	return nil
}

func (s *StockServer) SaveOrUpdateStock(ctx context.Context, in *pb.GrpcStock) (*pb.GrpcResponse, error) {
	daily := getStockByBPStock(in)

	err := Services.SaveOrUpdateStock(daily)
	if err != nil {
		return nil, err
	}

	return &pb.GrpcResponse{
		Error:   false,
		Message: "stock update.",
		Data:    nil,
	}, nil
}

func (s *StockServer) CreateStocks(stream pb.StockService_CreateStocksServer) error {
	var stocks []*models.Stock

	for {
		grpcStock, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			stream.SendAndClose(&pb.GrpcResponse{
				Error:   true,
				Message: err.Error(),
				Data:    nil,
			})
			return err
		}

		stocks = append(stocks, getStockByBPStock(grpcStock))
	}

	err := Services.CreateStocks(stocks)
	if err != nil {
		stream.SendAndClose(&pb.GrpcResponse{
			Error:   true,
			Message: err.Error(),
			Data:    nil,
		})
		return err
	}

	stream.SendAndClose(&pb.GrpcResponse{
		Error:   false,
		Message: "stocks added.",
		Data:    nil,
	})

	return nil
}

func (s *StockServer) RemoveStock(ctx context.Context, id *wrapperspb.Int64Value) (*pb.GrpcResponse, error) {
	err := Services.RemoveStock(uint(id.Value))

	if err != nil {
		return &pb.GrpcResponse{
			Error:   true,
			Message: err.Error(),
			Data:    nil,
		}, err
	}

	return &pb.GrpcResponse{
		Error:   false,
		Message: "stock removed",
		Data:    nil,
	}, nil
}

func (s *StockServer) CreateDailyStocks(stream pb.StockService_CreateDailyStocksServer) error {
	var dailys []*models.DailyHistoricalStock

	for {
		grpcDaily, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			stream.SendAndClose(&pb.GrpcResponse{
				Error:   true,
				Message: err.Error(),
				Data:    nil,
			})
			return err
		}

		// log.Println(getDailyStockByBPStock(grpcDaily))

		dailys = append(dailys, getDailyStockByBPStock(grpcDaily))
	}

	err := Services.CreateDailyStocks(dailys)

	if err != nil {
		stream.SendAndClose(&pb.GrpcResponse{
			Error:   true,
			Message: err.Error(),
			Data:    nil,
		})
		return err
	}

	stream.SendAndClose(&pb.GrpcResponse{
		Error:   false,
		Message: "Daily stocks added.",
		Data:    nil,
	})

	return nil
}

func (s *StockServer) RetrieveDailyByDuration(req *pb.GrpcStockDurationRequest, stream pb.StockService_RetrieveDailyByDurationServer) error {
	dailys, err := Services.RetrieveDailyByDuration(req.Symbols, req.Start.AsTime(), req.End.AsTime())

	if err != nil {
		return err
	}

	for _, daily := range dailys {
		err := stream.Send(getBPDailyStockByStock(daily))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StockServer) SaveOrUdpdateDailyStock(ctx context.Context, daily *pb.GrpcDailyHistoricalStock) (*pb.GrpcResponse, error) {
	err := Services.SaveOrUdpdateDailyStock(getDailyStockByBPStock(daily))
	if err != nil {
		return &pb.GrpcResponse{
			Error:   true,
			Message: err.Error(),
			Data:    nil,
		}, err
	}

	return &pb.GrpcResponse{
		Error:   false,
		Message: "Daily update.",
		Data:    nil,
	}, nil
}

func getPBStockByStock(stk *models.Stock) *pb.GrpcStock {
	return &pb.GrpcStock{
		Id:          uint64(stk.Id),
		StockSymbol: stk.StockSymbol,
	}
}

func getStockByBPStock(stk *pb.GrpcStock) *models.Stock {
	return &models.Stock{
		Id:          uint(stk.Id),
		StockSymbol: stk.StockSymbol,
	}
}

func getDailyStockByBPStock(stk *pb.GrpcDailyHistoricalStock) *models.DailyHistoricalStock {
	return &models.DailyHistoricalStock{
		Id:      uint(stk.Id),
		Close:   stk.Close,
		Open:    stk.Open,
		High:    stk.High,
		Low:     stk.Low,
		StockId: uint(stk.StockId),
		Volume:  int64(stk.Volume),
		Count:   uint(stk.Count),
		Wap:     stk.Wap,
		Date:    stk.Date.AsTime().Local(),
	}
}

func getBPDailyStockByStock(stk *models.DailyHistoricalStock) *pb.GrpcDailyHistoricalStock {
	return &pb.GrpcDailyHistoricalStock{
		Id:      uint64(stk.Id),
		Close:   stk.Close,
		Open:    stk.Open,
		High:    stk.High,
		Low:     stk.Low,
		StockId: uint64(stk.StockId),
		Volume:  uint64(stk.Volume),
		Count:   uint64(stk.Count),
		Wap:     stk.Wap,
		Date:    timestamppb.New(stk.Date),
	}
}
