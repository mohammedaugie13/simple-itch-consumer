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
	details, ok := s.Map.Get(message.StockCode)
	if !ok {
		return &pb.OHLC{Value: "0", Volume: "0", AveragePrice: "0", HighestPrice: "0", LowestPrice: "0", OpenPrice: "0", ClosePrice: "0", PreviousPrice: "0"}, nil
	}
	return &pb.OHLC{Value: details.Value.String(), Volume: details.Volume.String(), AveragePrice: details.AveragePrice.String(), HighestPrice: details.HighestPrice.String(), LowestPrice: details.LowestPrice.String(), OpenPrice: details.OpenPrice.String(), ClosePrice: details.ClosePrice.String(), PreviousPrice: details.PreviousPrice.String()}, nil
}
