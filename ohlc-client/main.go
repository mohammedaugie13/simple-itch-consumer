package main

import (
	"context"
	"fmt"
	"log"
	"ohlc-client/models"
	s "ohlc-client/models/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer func() {
		if errConn := conn.Close(); errConn != nil {
			fmt.Println("Error when closing:", errConn)
		}
	}()

	c := s.NewOHLCServiceClient(conn)

	response, err := c.GetOHLC(context.Background(), &s.StockCode{StockCode: "TLKM"})
	if err != nil {
		return
	}

	ohlc := models.OHLC{
		PreviousPrice: response.PreviousPrice,
		OpenPrice:     response.OpenPrice,
		HighestPrice:  response.HighestPrice,
		LowestPrice:   response.LowestPrice,
		ClosePrice:    response.ClosePrice,
		AveragePrice:  response.AveragePrice,
		Volume:        response.Volume,
		Value:         response.Value,
	}
	ohlcJSON, _ := ohlc.ToJSON()
	log.Printf("Response from server: %s", string(ohlcJSON))
}
