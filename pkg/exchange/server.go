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
	fn        func(string) bool
	SendErrCh chan error
	RecvErrCh chan error
	ReqCh     chan *pb.AlertReq
	ResCH     chan *pb.AlertRes
	Logger    *logging.Logger
}

const bufferSize = 10

func NewServer(callback func(string) bool) *Server {
	return &Server{
		fn:        callback,
		SendErrCh: make(chan error, bufferSize),
		RecvErrCh: make(chan error, bufferSize),
		ReqCh:     make(chan *pb.AlertReq, bufferSize),
		ResCH:     make(chan *pb.AlertRes, bufferSize),
		Logger:    util.Logger("Server"),
	}
}

func (s *Server) RequestPrice(stream pb.CryptoWatch_RequestPriceServer) error {
	for {
		go s.Recv(s.ReqCh, s.RecvErrCh, stream)
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
	if isSUpported := s.fn(exchange.Name); isSUpported {
		return &pb.Exchange{
			Name:      exchange.Name,
			Supported: isSUpported,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("Exchange %s is not supported", exchange.Name))
}

func (s *Server) Recv(reqCh chan *pb.AlertReq, errCh chan error, stream pb.CryptoWatch_RequestPriceServer) {
	req, err := stream.Recv()
	if err == io.EOF {
		errCh <- nil
	}

	if err != nil {
		s.Logger.WithError(err).Info("Pushed error into 'errCh'.")
		errCh <- err
		return
	}
	reqCh <- req
}

func (s *Server) Send(errCh chan error, stream pb.CryptoWatch_RequestPriceServer) {
	res := <-s.ResCH
	if err := stream.Send(res); err != nil {
		s.Logger.WithError(err).Info("Pushed error into 'errCh'.")
		errCh <- err
		return
	}
}

func (s *Server) StartServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		s.Logger.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCryptoWatchServer(grpcServer, s)

	s.Logger.Infof("Grpc server listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		s.Logger.Fatalf("failed serving server: %v", err)
	}
}

func (s *Server) handleConn(fn func(alertReq *pb.AlertReq)) {
	s.Logger.Info("Server handling conn.")
	for {
		select {
		case req := <-s.ReqCh:
			fn(req)
		}
	}
}
