package main

import (
	"fmt"
	"github.com/pkg/errors"
	pb "github.com/xjayleex/idl/protos/auth"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)


type AuthServer struct {
	userStore UserStore
	jwtManager *JWTManager
}

func NewAuthServer (userStore UserStore, jwtManager *JWTManager ) *AuthServer {
	return &AuthServer{
		userStore: userStore,
		jwtManager: jwtManager,
	}
}



func (authServer *AuthServer) SignUp (ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error)  {
	if authServer == nil {
		return nil, ErrNoAuthServer
	}
	user, err := NewUser(req.Mail, req.Username, req.Password)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = authServer.userStore.SignUp(user)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &pb.SignUpResponse{}, nil
}

func (authServer *AuthServer) SignIn (ctx context.Context, req *pb.SignInRequest) (*pb.SignInResponse, error) {

	user, err := authServer.userStore.Find(req.GetMail())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("Nil user.")
	} else if !user.IsCorrectPassword(req.GetPassword()) {
		return nil, ErrIncorrectInfo
	}
	token, err := authServer.jwtManager.Generate(user)
	if err != nil {
		return nil, err
	}
	fmt.Println("Gerated token to",req.Mail," ",token)
	res := &pb.SignInResponse{AccessToken: token}
	return res, nil
}

func (as *AuthServer)authServerInterceptor (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Println("It's unary server auth interceptor.")
	r , ok := req.(*pb.SignInRequest)
	if ok && as != nil {
		_, err := as.jwtManager.Verify(r.AccessToken)
		if err == nil {
			return nil, errors.New("Token Not Expired , Yet")
		}
	}
	h, err := handler(ctx, req)
	return h, err
}
