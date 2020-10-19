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
[(*Figure from 'github::grpc-gatway'*)](https://https://github.com/grpc-ecosystem/grpc-gateway)

본문에서는 유저 등록 서비스를 gRPC를 통해 구현하고, gRPC-gateway를 이용해 간단한 RESTful 서비스를 제공하는 방법에 대해 다룬다.

#### Installation
---
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
