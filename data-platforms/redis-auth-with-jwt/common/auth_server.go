package main

import (
	"fmt"
	"github.com/pkg/errors"
	pb "github.com/xjayleex/idl/protos/auth"
	"golang.org/x/net/context"
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

	res := &pb.SignInResponse{AccessToken: token}
	return res, nil
}