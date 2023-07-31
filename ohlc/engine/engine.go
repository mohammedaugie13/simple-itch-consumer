package engine

import (
	"context"
	"log"
	"ohlc/models"
	"ohlc/models/pb"
)

type Server struct {
	Map *models.OHLCMap
	pb.UnimplementedOHLCServiceServer
}

func (s *Server) GetOHLC(ctx context.Context, message *pb.StockCode) (*pb.OHLC, error) {
	log.Printf("Recevied %v", message.StockCode)
	return &pb.OHLC{Value: "1000", Volume: "1000", AveragePrice: "1000", HighestPrice: "1000", LowestPrice: "1000", OpenPrice: "10000", ClosePrice: "1000", PreviousPrice: "1000"}, nil
}
