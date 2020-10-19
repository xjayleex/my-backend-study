package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/xjayleex/idl/protos/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"time"
)

type AuthClient struct {
	client auth.AuthServiceClient
	conn *grpc.ClientConn
	mail string
	password string
	accessToken string
}

type GrpcClientConfig struct {
	Address			string
	RootCertificate string
	UsingClientInc	bool
}



func NewAuthClient(mail string, password string, accessToken string, cfg GrpcClientConfig) (*AuthClient,error) {
	var (
		grpcOpts = []grpc.DialOption{}
		grpcCreds credentials.TransportCredentials
		err 	error
	)
	if cfg.Address == "" {
		return nil, errors.New("[Client Side Error] : No Address is set.")
	}
	if cfg.RootCertificate == "" {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	}
	if cfg.RootCertificate == "" {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	} else {
		grpcCreds, err = credentials.NewClientTLSFromFile(cfg.RootCertificate,"")
		if err != nil {
			err = errors.Wrapf(err,
				"failed to create grpc tls client with root-cert %s",
				cfg.RootCertificate)
			return nil, err
		}
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcCreds))
	}

	if cfg.UsingClientInc {
		grpcOpts = append(grpcOpts,grpc.WithUnaryInterceptor(authClientInterceptor))
	}


	ac := &AuthClient{
		mail: mail,
		password:  password,
		accessToken: accessToken,
	}

	ac.conn, err = grpc.Dial(cfg.Address, grpcOpts...)
	if err != nil {
		return nil, errors.New("failed to start grpc connection with address")
	}
	ac.client = auth.NewAuthServiceClient(ac.conn)
	return ac, nil
}

func authClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	fmt.Println(time.Now()," : Client Interceptor [After calling gRPC Server]" )
	err := invoker(ctx, method, req, reply, cc, opts...)
	fmt.Println(time.Now()," : Client Interceptor [Before get response from gRPC Server]" )
	return err
}

func (ac *AuthClient) SignIn(opts ...grpc.CallOption) (string, error) {
	ctx := context.Background()
	req := &auth.SignInRequest{
		Mail:        ac.mail,
		Password:    ac.password,
		AccessToken: ac.accessToken,
	}
	fmt.Println(time.Now(), " : Client Origin gRPC call")
	time.Sleep(1000 * time.Millisecond)
	res, err := ac.client.SignIn(ctx, req, opts...)
	fmt.Println(time.Now()," : Gotton Response" )
	if err != nil {
		return "", err
	}
	return res.GetAccessToken(), nil
}

func (ac *AuthClient) Close() {
	if ac.conn != nil {
		ac.conn.Close()
	}
}

func main() {
	 ac, err := NewAuthClient("jayleekau@gmail.com","somepassword",
	 	"",
	 	GrpcClientConfig{
			Address:         "127.0.0.1:9090",
			RootCertificate: "",
			UsingClientInc:  true,
		})
	 defer ac.Close()
	 if err != nil {
	 	fmt.Println(err)
	 	os.Exit(1)
	 }
	 result, err := ac.SignIn()
	 if err != nil {
	 	fmt.Println(err)
	 } else {
	 	fmt.Println(result)
	 }
}