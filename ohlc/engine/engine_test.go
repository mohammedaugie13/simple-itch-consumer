package engine

import (
	"context"
	"log"
	"net"
	"ohlc/eventprocessor"
	"ohlc/models"
	"ohlc/models/pb"
	"testing"

	"github.com/alphadose/haxmap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func server(ctx context.Context, hmap *haxmap.Map[string, *models.OHLC]) (pb.OHLCServiceClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pb.RegisterOHLCServiceServer(baseServer, &Server{Map: hmap})
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewOHLCServiceClient(conn)

	return client, closer
}

func TestServer(t *testing.T) {
	data := `{"type":"A","quantity":"0","price":"8000","stock_code":"BBCA"}
{"type":"P","quantity":"100","price":"8050","stock_code":"BBCA"}
{"type":"P","quantity":"500","price":"7950","stock_code":"BBCA"}
{"type":"A","quantity":"200","price":"8150","stock_code":"BBCA"}
{"type":"E","quantity":"300","price":"8100","stock_code":"BBCA"}
{"type":"A","quantity":"100","price":"8200","stock_code":"BBCA"}
`
	hmap := haxmap.New[string, *models.OHLC]()
	events := eventprocessor.ProcessEvent(data)
	eventprocessor.CalculateOHLC(events, hmap)
	ctx := context.Background()
	client, closer := server(ctx, hmap)
	defer closer()
	type expectation struct {
		out *pb.OHLC
		err error
	}
	tests := map[string]struct {
		in       *pb.StockCode
		expected expectation
	}{
		"MustSuccess": {
			in: &pb.StockCode{StockCode: "BBCA"},
			expected: expectation{
				out: &pb.OHLC{
					PreviousPrice: "8000",
					OpenPrice:     "8050",
					HighestPrice:  "8100",
					LowestPrice:   "7950",
					ClosePrice:    "8100",
					AveragePrice:  "8011",
					Volume:        "900",
					Value:         "7210000",
				},
				err: nil,
			},
		},
	}
	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetOHLC(ctx, tt.in)
			if err != nil {
				if tt.expected.err.Error() != err.Error() {
					t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
				}
			} else {
				if tt.expected.out.AveragePrice != out.AveragePrice ||
					tt.expected.out.LowestPrice != out.LowestPrice ||
					tt.expected.out.HighestPrice != out.HighestPrice {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}
		})
	}
}
