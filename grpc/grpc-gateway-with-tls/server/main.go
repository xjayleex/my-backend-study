package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	pb "github.com/xjayleex/idl/protos/grpc-gateway-test"
	"github.com/xjayleex/my-backend-study/grpc/grpc-gateway-with-tls/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"os"
)

const (
	grpcAddr = ":10000"
	grpcServerCert = "/Users/ijaehyeon/keys/server.crt"
	grpcServerKey = "/Users/ijaehyeon/keys/server.pem"
)

type GrpcServer struct {
	server *grpc.Server
	userStore store.UserStore
	logger *logrus.Entry
}

var (
	redisAddress = flag.String("redisAddress","localhost", "Redis server host address")
	redisPort = flag.String("redisPort","6379","Redis tcp port")
	redisDB = flag.Int("redisDB",2,"Redis DB Number to store data(default:2, from 0 to 16 integer")
	interceptorLogger *logrus.Entry
)

type User struct {
	Mail				string	`json:"mail"`
	Username			string	`json:"name"`
}

func NewUser(username string, mail string) *User {
	return &User {
		Username: username,
		Mail: mail,
	}
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func NewGrpcServer(serverCrt string, serverKey string, userStore store.UserStore , opts ...grpc.ServerOption) (*GrpcServer , error) {
	if serverCrt == "" || serverKey == "" {
		return nil, errors.New("Server certificate path needed.")
	}

	cred, tlsErr := credentials.NewServerTLSFromFile(serverCrt, serverKey)

	if tlsErr != nil {
		return nil, tlsErr
	}

	opts = append(opts, grpc.Creds(cred))

	return &GrpcServer{
		server: grpc.NewServer(opts...),
		userStore: userStore,
		logger: logrus.WithFields(logrus.Fields{
			"Name": "gRPC-Server",
		}),
	}, nil
}
// Post Registration info.
// SendPost(context.Context, *TemplateRequest) (*TemplateResponse, error)
func (gs *GrpcServer) CheckEnrollment(ctx context.Context, req *pb.CheckEnrollmentRequest) (*pb.CommonResponseMsg, error) {
	v, err := gs.userStore.Find(req.Mail)
	if err != nil {
		se, _ := err.(*store.StoreError)
		if se.Code == store.ErrNoConnWithRedis {
			gs.logger.Error("No Connection with redis.")
			return nil, status.Error(codes.Internal, "Internal Error")
		} else if se.Code == store.ErrKeyNotExists {
			return nil, status.Error(codes.NotFound, "Mail Not Found")
		} else {
			// No matched StoreError here.
		}
	}

	asserted, ok := v.(*redis.StringCmd)
	if !ok {
		return nil, status.Error(codes.Internal, "Marshal/Unmarshal error")
	}

	unmarshaled, err := asserted.Bytes()
	if err != nil {
		return nil, status.Error(codes.Internal, "Marshal/Unmarshal error")
	}

	user := &User{}
	err = json.Unmarshal(unmarshaled, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "Marshal/Unmarshal error")
	}

	if user.Username != req.Name {
		return nil, status.Error(codes.NotFound, "No mathches with request")
	}

	result := &pb.CommonResponseMsg{
		Message: req.Mail + " is verified.",
	}
	return result, nil
}

func (gs *GrpcServer) Enroll(ctx context.Context, req *pb.EnrollmentRequest) (*pb.CommonResponseMsg, error) {
	err := gs.userStore.Save(req.Mail, NewUser(req.Name,req.Mail))
	if err != nil {
		se, _ := err.(*store.StoreError)
		if se.Code == store.ErrKeyExistsAlready {
			return nil, status.Error(codes.AlreadyExists, "The mail is already enrolled.")
		} else if se.Code == store.ErrNoConnWithRedis {
			gs.logger.Error("No Connection with redis.")
			return nil, status.Error(codes.Internal, "Internal Error")
		} else {
			return nil, status.Error(codes.Internal, "Internal Error")
		}
	}
	result := &pb.CommonResponseMsg{Message: "Enrolled successfully"}
	return result, nil
}

func (gs *GrpcServer) Serve(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		gs.logger.Fatalf("Faield to listen: %v", err)
		return errors.Errorf("Failed to listen: %v", err)
	}
	if gs.server == nil {
		return errors.New("gRPC server is nil")
	}
	if err = gs.server.Serve(lis); err != nil {
		return err
	}

	gs.logger.Info("Serving gRPC server ...")
	return nil
}

func serverInterceptor (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		interceptorLogger.Infof("Request from %s",p.Addr.String())
	}
	h, err := handler(ctx, req)

	if ok {
		interceptorLogger.Infof("Response to %s", p.Addr.String())
	}
	return h, err
}

func (gs *GrpcServer) Close() {
	gs.server.Stop()
}

func init(){
	flag.Parse()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	interceptorLogger = logrus.WithFields(logrus.Fields{
		"Name": "gRPC-Server-Interceptor",
	})
}

func main() {
	rs, err := store.NewRedisUserStore(
		&store.RedisClientOpts{
			Address: *redisAddress,
			Port:    *redisPort,
			DB:      *redisDB,
		})

	grpcServer, err := NewGrpcServer(grpcServerCert,grpcServerKey, rs , grpc.UnaryInterceptor(serverInterceptor))
	if err != nil {
		grpcServer.logger.Fatal("Failed with initializing gRPC Server : %v", err)
		os.Exit(1)
	}
	pb.RegisterEnrollmentServer(grpcServer.server, grpcServer)
	reflection.Register(grpcServer.server)
	defer grpcServer.Close()
	if err = grpcServer.Serve(grpcAddr); err != nil {
		grpcServer.logger.Fatalf("Failed to serve gRPC server : %v", err)
	}
}
