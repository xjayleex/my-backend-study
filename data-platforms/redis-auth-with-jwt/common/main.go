package main

import (
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
	pb "github.com/xjayleex/idl/protos/auth"
)

func main() {
	fmt.Println("Auth Server is Starting...")
	var (
		listener net.Listener
	)
	listener, err := net.Listen("tcp", ":" + strconv.Itoa(9090))
	if err != nil {
		err = errors.Wrapf(err, "failed to listen on port %d", 9090)
		return
	}

	grpcOpts := []grpc.ServerOption{
	}
	rs := NewRedisUserStore(&RedisClientOpts{
		Address: "localhost",
		Port:    "6379",
		DB:	UserStoreDB,
	})
	authServer := NewAuthServer(rs, NewJWTManager("secret", 1 * time.Minute))
	if true {
		grpcOpts = append(grpcOpts, grpc.UnaryInterceptor(authServer.authServerInterceptor))
	}
	authGrpcServer := grpc.NewServer(grpcOpts...)
	pb.RegisterAuthServiceServer(authGrpcServer,authServer)
	err = authGrpcServer.Serve(listener)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Auth Server is On.")
	defer authGrpcServer.Stop()
	fmt.Println("Stopping Auth Server ...")
}