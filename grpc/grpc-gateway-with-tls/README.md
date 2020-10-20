---
layout: post
comment: false
title: gRPC-gateway
date: 2020-10-12
category: "Golang"
tags: [gRPC-gateway]
update: 2020-10-12
author: Jaehyun Lee
---
> ### Contents
[**gRPC Gateway Introduction**](#grpc-gateway-Introduction)  
[**Installation**](#installation)  
[**User Database**](#user-database)  
[**.proto**](#.proto)  
[**Code Generation**](#code-generation)  
[**gRPC Server**](#grpc-server)  
[**gRPC Gateway**](#grpc-gateway)  
[**Client Request**](#client-request)  

#### gRPC Gateway Introduction
---
gRPC는 protobuf를 idl로써 사용해 통신에 필요한 데이터와 서비스를 정의하고, 정의된 idl에서 원하는 프로그래밍 언어의 클라이언트/서버 스텁 코드를 쉽게 생성 할 수 있도록 해준다. gRPC가 가져다주는 이점들이 있기는 하지만, 기존의 서비스에 쉽게 붙여질 수 없다면 문제가 될 것이다.  
JSON 기반의 REST 서비스들이 gRPC 서비스의 API를 호출하기 위해서는`gRPC-gateway`를 고려해볼 수 있다. 
gRPC-gateway는 protobuf에 정의된 서비스를 이용해, RESTful API 요청을 gRPC로 변환하는 Reverse Proxy이다. 이를 통해 gRPC 서버를 gRPC와 RESTful API 스타일로 동시에 제공 할 수 있다.

![Image](/assets/images/grpc-gateway.png){:style="width: 90%; margin: 0 auto; display: block;"}
[(*Figure from 'github::grpc-gatway'*)](https://github.com/grpc-ecosystem/grpc-gateway)

본문에서는 유저 등록 서비스를 gRPC를 통해 구현하고, gRPC-gateway를 이용해 간단한 RESTful 서비스를 제공하는 방법에 대해 다룬다.

#### Installation
---
아래 패키지들을 설치하고, 만약에 `$GOBIN`을 환경변수로 등록하지 않은 상태라면, `$GOPATH/bin`을 `$GOBIN` 환경변수로 설정해준다.

```bash
$ go get "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
$ go get "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
$ go get "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
$ go get "google.golang.org/protobuf/cmd/protoc-gen-go"

$ go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

#### User Database
---
유저는 메일 주소, 이름을 통해 시스템에 등록하고, 등록 여부를 확인 할 수 있도록 하고자 했다. 유저 데이터 구조는 아래와 같다.
```go
type User struct {
	Mail				string	`json:"mail"`
	Username			string	`json:"name"`
}
```
User 데이터베이스로는 Redis를 사용하기로 했으며, 유저의 Mail을 Key로, 객체를 Redis에 저장한다. 이 데이터 구조를 Redis에 저장하고, 값을 찾아오는 구현은 [go-redis 클라이언트 구현](https://xjayleex.github.io)에 따로 정리했다.
상기시켜볼만 한 부분을 살펴보자. 객체를 마샬링해서 Redis에 저장하려면 `encoding` 패키지의 BinaryMarshaler와 BinaryUnmarshaler를 구현해야한다. [store.go](https://https://github.com/xjayleex/my-backend-study/blob/master/grpc/grpc-gateway-with-tls/store/store.go)에 Redis의 값으로 들어갈 `RedisValue` 인터페이스를 선언하고, 위에서 언급한 두 인터페이스의 메소드인 MarshalBinary()와 UnmarshalBinary()를 RedisValue의 메소드로 정의했다. 여기에서는 User 구조체를 그대로 Redis에 저장할 것이고, 따라서 Redis를 사용할 gRPC 서버 코드에서 User 타입을 리시버로하는 MarshalBinary(), UnmarshalBinary 메소드를 구현하면 될 것이다.
```go
// store.go
type RedisValue interface {
	MarshalBinary() ([]byte, error)
	UnmarshalBinary(data[]byte) error
}

type UserStore interface {
	Save(key string, value RedisValue) error
	Find(key string) (interface{},error)
}

```

#### .proto
---
enrollment.proto를 정의한다. 기존 gRPC만 고려했을 때, 유저 등록 서비스는 다음과 같이 정의 할 수 있다.
```protobuf
service Enrollment {
  rpc CheckEnrollment (CheckEnrollmentRequest) returns (CommonResponseMsg); 
  rpc Enroll (EnrollmentRequest) returns (CommonResponseMsg); 
}

// The request message containing the user's name and email addr.
message CheckEnrollmentRequest {
  string name = 1;
  string mail = 2;
}

// The response message containing the Enrollment info.
message CommonResponseMsg {
  string message = 1;
}

message EnrollmentRequest {
  string name = 1;
  string mail = 2;
}
```

gRPC-gateway는 .proto 파일과 proto 코드에 작성된 [`google.api.http`](https://github.com/googleapis/googleapis/blob/master/google/api/http.proto#L46) 어노테이션 매핑에 따라서 스텁 코드를 생성한다. 위의 코드에 `google.api.http` 어노테이션을 추가해 http 요청을 받을 수 있는 url을 바인딩 해보자.

```protobuf
import "google/api/annotations.proto";
package enroll;

// The Enrollment service definition.
service Enrollment {
  // Check Enrollment info.
  rpc CheckEnrollment (CheckEnrollmentRequest) returns (CommonResponseMsg) {
    option (google.api.http) = {
      get: "/v1/users/{name}/{mail}"
      additional_bindings {
        get: "/v1/users/check/{name}/{mail}"
      }
    };
  }
  // Send Enrollment request which mapped with POST req.
  rpc Enroll (EnrollmentRequest) returns (CommonResponseMsg) {
    option (google.api.http) = {
      post: "/post"
      body: "*"
    };
  }
}

```

#### Code Generation
---
위에서 작성한 [`enrollment.proto`](https://github.com/xjayleex/idl/blob/master/protos/grpc-gateway-test/enrollment.proto) 파일을 컴파일해 다음의 코드들을 생성한다.
- gRPC 서버, 클라이언트 Stub 코드 (본문에서는 go언어를 타겟으로 코드를 생성했다, `enrollment.pb.go`).
- gRPC-gatway Stub 코드 (`enrollment.pb.gw`)
- API 문서화를 위한 swagger.json 파일 (`enrollment.swagger.json`)

```bash
$ protoc -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. enrollment.proto

$ protoc -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. enrollment.proto

$ protoc -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:. enrollment.proto
```

#### gRPC Server
---
gRPC Server 전체 코드는 이 곳 [`server/main.go`](https://github.com/xjayleex/my-backend-study/blob/master/grpc/grpc-gateway-with-tls/server/main.go)을 참고.

gRPC 서버 구조체를 정의한다. grpc.Server에 대한 포인터와, User 정보를 저장하고 조회할 UserStore 인터페이스, 그리고 로깅을 위한 Logger 패키지를 포함한다. 로깅 패키지로는 [`logrus`](https://github.com/sirupsen/logrus)를 사용했다.


```go
type GrpcServer struct {
	server *grpc.Server
	userStore store.UserStore
	logger *logrus.Entry
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

```

proto 파일에 정의해놓은 CheckEnrollment, Enroll 서비스를 구현한다. CheckEnrollment 서비스는 Request의 Mail와 Username로 내부 User DB를 조회한뒤, 등록 여부를 응답하는 서비스이다. 


```go
// CheckEnrollment Service
// which checks user enrollemnt.
func (gs *GrpcServer) CheckEnrollment(ctx context.Context, req *pb.CheckEnrollmentRequest) (*pb.CommonResponseMsg, error) {
	v, err := gs.userStore.Find(req.Mail)
	if err != nil {
		se, _ := err.(*store.StoreError)
		if se.Code == store.ErrNoConnWithRedis {
			gs.logger.Fatal("No Connection with redis.")
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
```
[`google.golang.org/grpc/codes`](https://godoc.org/google.golang.org/grpc/codes)에는 gRPC에서 사용하는 상태 코드가 표현되어 있다. 이 코드를 이용해 gRPC-gateway가 gRPC 서버에서 에러 코드를 리턴 받았을 때, [`grpc-gateway/runtime/error.go`](https://github.com/grpc-ecosystem/grpc-gateway/blob/master/runtime/errors.go)를 통해 http 상태 코드로 변환한다. 따라서 gRPC 서버에서 에러를 핸들링할 때, 꼭 gRPC 상태 코드로 핸들링하자.

Enroll 서비스 구현은 다음과 같다. Enroll 서비스는 Request에 포함된 Mail과 Username으로 User DB에 Mail을 키로해서 저장하는 서비스이다.

```go
func (gs *GrpcServer) Enroll(ctx context.Context, req *pb.EnrollmentRequest) (*pb.CommonResponseMsg, error) {
	err := gs.userStore.Save(req.Mail, NewUser(req.Name,req.Mail))
	if err != nil {
		se, _ := err.(*store.StoreError)
		if se.Code == store.ErrKeyExistsAlready {
			return nil, status.Error(codes.AlreadyExists, "The mail is already enrolled.")
		} else if se.Code == store.ErrNoConnWithRedis {
			gs.logger.Fatal("No Connection with redis.")
			return nil, status.Error(codes.Internal, "Internal Error")
		} else {
			return nil, status.Error(codes.Internal, "Internal Error")
		}
	}
	result := &pb.CommonResponseMsg{Message: "Enrolled successfully"}
	return result, nil
}
```

필수적이진 않지만, 클라이언트 요청과 서버 응답을 로깅하기 위해 서버측 Unary Interceptor를 구현했다. 클라이언트 요청에 실린 Context 객체로부터 클라이언트 정보를 추출할 수 있다.

```go
func serverInterceptor (ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	p, ok := peer.FromContext(ctx)
	if ok {
		interceptorLogger.Infof("Request from %s",p.Addr.String())
	}
	h, err := handler(ctx, req)
	// handle을 gRPC 서버로 넘기고, 받음

	if ok {
		interceptorLogger.Infof("Response to %s", p.Addr.String())
	}
	return h, err
}
```

유저 데이터 모델은 위에서 살펴본 바와 같고, 앞서 언급했듯이 MarshalBinary()와 UnmarshalBinary를 구현했다. 

```go
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
```

#### gRPC Gateway
---
gRPC-Gateway 전체 코드는 [`gateway/gateway.go`](https://github.com/xjayleex/my-backend-study/blob/master/grpc/grpc-gateway-with-tls/gateway/gateway.go)을 참고.

gRPC-Gateway를 사용하기 위한 기본적인 코드는 아주 간단한데, gRPC 서버를 엔드포인트로 하는 gRPC-Gateway의 멀티플렉서인 `runtime.ServeMux`를 http 멀티플렉서 핸들러로 붙여주는 것이 핵심이다. 
첫번째 과정인 gRPC 서버 엔드포인트를 `runtime.ServeMux`에 붙여주는 역할은, gateway 스텁 코드의 `Register[ServiceName]HandlerFromEndpoint` 메서드가 수행한다.
`runtime.ServerMux`는 `http.Handler`를 구현하고 있으므로, `http.ListenAndServe` 로 핸들러를 구동시켜주기만 하면 된다.

```go

func NewGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(opts...) // + runtime.WithMarshalerOption()
	grpcDialOpts := []grpc.DialOption{}

	if cred, err := credentials.NewClientTLSFromFile(grpcServerCert,""); err == nil {
		grpcDialOpts = append(grpcDialOpts, grpc.WithTransportCredentials(cred))
	} else {
		return nil, err
	}

	err := gw.RegisterEnrollmentHandlerFromEndpoint(ctx, mux, *getEndpoint, grpcDialOpts)
	if err != nil {
		return nil, errors.New("grpc Gateway : `GET` error")
	}

	err = gw.RegisterEnrollmentHandlerFromEndpoint(ctx, mux, *postEndpoint,grpcDialOpts)
	if err != nil {
		return nil, errors.New("grpc Gateway : `POST error")
	}

	return mux, nil
}

func Serve(address string, opts ...runtime.ServeMuxOption) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()

	mux.HandleFunc("/swagger/", swaggerHandler)
	//opts = append(opts, runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler))
	gwHandler, err := NewGateway(ctx, opts...)
	if err != nil {
		return err
	}

	mux.Handle("/", gwHandler)

	return http.ListenAndServeTLS(address, gwCert, gwKey, allowCORS(mux))
}
```

#### Client Request
---
gRPC 서버에서 http 응답을 받아내는 것이 목적이었으므로, gRPC 클라이언트는 구현할 필요가 없었다. 또한 http 클라이언트는 그냥 curl을 쓰자...

```bash
# 서버 실행
$ go run server/main.go
$ go run gateway/gateway.go
```

```bash
$ curl -X GET --cacert ~/keys/server.crt https://localhost:8080/v1/users/foo/foo@mail.co.kr -i
HTTP/2 404 
content-type: application/json
content-length: 62
date: Tue, 20 Oct 2020 08:13:44 GMT

{"error":"Mail Not Found","code":5,"message":"Mail Not Found"}
```
유저 정보를 아직 등록하지 않아서, 404 코드를 받았다. POST 메소드로 유저 정보를 등록해보자.

```bash
$ curl -i -X POST --cacert ~/keys/server.crt https://localhost:8080/post \
> -H "Content-Type: application/json" \
> -d '{"name":"foo", "mail":"foo@mail.co.kr"}'
HTTP/2 200 
content-type: application/json
grpc-metadata-content-type: application/grpc
content-length: 35
date: Tue, 20 Oct 2020 08:18:37 GMT

{"message":"Enrolled successfully"}
```

Redis 서버를 끄고 요청을 날려보면, 
```bash
curl -i -X POST --cacert ~/keys/server.crt https://localhost:8080/post \
-H "Content-Type: application/json" \
-d '{"name":"foo", "mail":"foo@mail.co.kr"}'
HTTP/2 500 
content-type: application/json
content-length: 63
date: Tue, 20 Oct 2020 08:29:38 GMT

{"error":"Internal Error","code":13,"message":"Internal Error"}
```
예상대로 Internal Error 상태 코드를 받을 수 있다.


