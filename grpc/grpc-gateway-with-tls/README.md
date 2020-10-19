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
[**Introduction**](#introduction)  
[**Encoding**](#encoding)  
[**Decoding**](#decoding)  
[**Generic JSON with empty interface type**](#generic-json-with-empty-interface-type)  
[**Decoding arbitrary data**](#decoding-arbitrary-data)  
[**Reference Types**](#reference-types)  
[**Streaming Encoders and Decoders**](#streaming-encoders-and-decoders)  

#### gRPC Gateway
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

#### enrollment.proto
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

proto 파일에 정의해놓은 CheckEnrollment, Enroll 서비스를 구현한다. CheckEnrollment 서비스는 Request의 Mail와 Username로 내부 User DB를 조회한뒤, 등록 여부를 응답하는 서비스이다. Enroll 서비스는 Request에 포함된 Mail과 Username으로 User DB에 Mail을 키로해서 저장하는 서비스이다.
