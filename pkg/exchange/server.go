package exchange

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/SkycoinProject/skycoin/src/util/logging"
	"google.golang.org/grpc"

	pb "github.com/Kifen/crypto-watch/pkg/proto"
	"github.com/Kifen/crypto-watch/pkg/util"
)

type Server struct {
	fn          func(string) bool
	SendErrCh   chan error
	RecvErrCh   chan error
	alertReqCh  chan *pb.AlertReq
	alertResCh  chan *pb.AlertRes
	symbolReqCH chan *pb.Symbol
	symbolResCH chan *pb.Symbol
	logger      *logging.Logger
}

const bufferSize = 10

func NewServer(callback func(string) bool) *Server {
	return &Server{
		fn:         callback,
		SendErrCh:  make(chan error, bufferSize),
		RecvErrCh:  make(chan error, bufferSize),
		alertReqCh: make(chan *pb.AlertReq, bufferSize),
		alertResCh: make(chan *pb.AlertRes, bufferSize),
		logger:     util.Logger("Server"),
	}
}

func (s *Server) RequestPrice(stream pb.CryptoWatch_RequestPriceServer) error {
	for {
		go s.Recv(s.alertReqCh, s.RecvErrCh, stream)
		go s.Send(s.SendErrCh, stream)

		select {
		case recvErr := <-s.RecvErrCh:
			return recvErr
		case sendErr := <-s.SendErrCh:
			return sendErr
		}
	}

	return nil
}

func (s *Server) IsExchangeSupported(ctx context.Context, exchange *pb.Exchange) (*pb.Exchange, error) {
	s.logger.Info("Sending response for supported exchange.")
	if isSupported := s.fn(exchange.Name); isSupported {
		return &pb.Exchange{
			Name:      exchange.Name,
			Supported: isSupported,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("Exchange %s is not supported", exchange.Name))
}

func (s *Server) IsSymbolValid(ctx context.Context, req *pb.Symbol) (*pb.Symbol, error) {
	s.symbolReqCH <- req
	res := <- s.symbolResCH

	return &pb.Symbol{
		Id: res.Id,
		ExchangeName: res.ExchangeName,
		Symbol: res.Symbol,
		Valid: res.Valid,
	}, nil
}

func (s *Server) Recv(reqCh chan *pb.AlertReq, errCh chan error, stream pb.CryptoWatch_RequestPriceServer) {
	req, err := stream.Recv()
	if err == io.EOF {
		errCh <- nil
	}

	if err != nil {
		s.logger.WithError(err).Info("Pushed error into 'errCh'.")
		errCh <- err
		return
	}
	reqCh <- req
}

func (s *Server) Send(errCh chan error, stream pb.CryptoWatch_RequestPriceServer) {
	res := <-s.alertResCh
	if err := stream.Send(res); err != nil {
		s.logger.WithError(err).Info("Pushed error into 'errCh'.")
		errCh <- err
		return
	}
}

func (s *Server) StartServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCryptoWatchServer(grpcServer, s)

	s.logger.Infof("Grpc server listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		s.logger.Fatalf("failed serving server: %v", err)
	}
}
