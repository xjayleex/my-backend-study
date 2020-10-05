package scatch

import (
	"fmt"
	"google.golang.org/grpc"
	"time"
)

type AuthClientInterceptor struct {
	authClient *AuthClient
	accessToken string
}

func NewAuthClientInterceptor (authClient *AuthClient, refreshDutarion time.Duration) (*AuthClientInterceptor, error) {
	interceptor := &AuthClientInterceptor{
		authClient: authClient,
	}
	err := interceptor.scheduleRefreshToken(refreshDutarion)
	if err != nil {
		return nil, err
	}
	return interceptor, nil
}

func (interceptor *AuthClientInterceptor) refreshToken() error {
	fmt.Println("It's `refreshToken()` method.")
	/*accessToken, err := interceptor.authClient.SignIn()
	if err != nil {
		return err
	}
	interceptor.accessToken = accessToken
	fmt.Println("token refreshed: %s", accessToken)*/
	return nil
}

func (interceptor *AuthClientInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {

	fmt.Println("It's `scheduleRefreshToken()` method.")
	/*err := interceptor.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration
		for {
			time.Sleep(wait)
			err := interceptor.refreshToken()
			if err != nil {
				wait = time.Second
			} else {
				wait = refreshDuration
			}
		}
	}()*/
	return nil
}

func (interceptor *AuthClientInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Printf("--> Unary interceptor : %s", method)
		return invoker(ctx,method, req, reply, cc, opts...)
	}
}



func main() {
	grpcClient, err := NewAuthClient()
}