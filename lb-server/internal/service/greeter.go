package service

import (
	"context"
	"time"

	v1 "lb-server/api/helloworld/v1"
	"lb-server/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// var host = "localhost:9001"

var host = "lb-server:9001"

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc  *biz.GreeterUsecase
	log *log.Helper

	greeterClient v1.GreeterClient
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	log := log.NewHelper(logger)
	conn, err := grpc.DialInsecure(context.Background(), grpc.WithEndpoint(host), grpc.WithTimeout(time.Second*5))
	if err != nil {
		log.Errorf("grpc.Dial err: %v", err)
		panic(err)
	}

	client := v1.NewGreeterClient(conn)
	return &GreeterService{uc: uc, log: log, greeterClient: client}
}

// SayHello implements helloworld.GreeterServer
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	s.log.WithContext(ctx).Infof("SayHello Received: %v", in.GetName())

	if in.GetName() == "error" {
		return nil, v1.ErrorUserNotFound("user not found: %s", in.GetName())
	}
	return &v1.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (s *GreeterService) Ping(ctx context.Context, in *v1.Empty) (*v1.HelloReply, error) {
	s.log.WithContext(ctx).Infof("Ping Received")
	reply, err := s.greeterClient.SayHello(ctx, &v1.HelloRequest{Name: time.Now().Local().String()})
	if err != nil {
		s.log.WithContext(ctx).Errorf("Ping SayHello err: %v", err)
		return nil, err
	}
	return &v1.HelloReply{Message: reply.Message}, nil
}

func (s *GreeterService) Healthz(ctx context.Context, in *v1.Empty) (*v1.HelloReply, error) {
	s.log.Debugf("Healthz Received")
	return &v1.HelloReply{Message: "ok"}, nil
}
