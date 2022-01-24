package service

import (
	"context"
	"fmt"
	"time"

	v1 "lb-server/api/helloworld/v1"
	"lb-server/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc"
)

// host := "localhost:9001"
var host = "lb-server.default:9001"
var connection *grpc.ClientConn

func init() {
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		fmt.Errorf("grpc.Dial err: %v", err)
		panic(err)
	}
	connection = conn
}

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc  *biz.GreeterUsecase
	log *log.Helper
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	return &GreeterService{uc: uc, log: log.NewHelper(logger)}
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
	client := v1.NewGreeterClient(connection)
	reply, err := client.SayHello(ctx, &v1.HelloRequest{Name: time.Now().Local().String()})
	if err != nil {
		s.log.WithContext(ctx).Errorf("Ping SayHello err: %v", err)
		return nil, err
	}
	return &v1.HelloReply{Message: reply.Message}, nil
}
