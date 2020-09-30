package service


import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	messaging "github.com/xjayleex/idl/protos/imageproto"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"net"
	"os"
	"strconv"
)

type GrpcServer struct {
	logger		zerolog.Logger
	server		*grpc.Server
	port		int
	cert		string
	key			string
	counter 	int
	storage		string
}

type GrpcServerConfig struct {
	Certificate string
	Key			string
	Port		int
	Storage  	string
}

func NewGrpcServer(cfg GrpcServerConfig) (s GrpcServer, err error) {
	s.logger = zerolog.New(os.Stdout).With().
		Str("from","server").Logger()
	if cfg.Port == 0 {
		err = errors.Errorf("Port must be specified.")
		return
	}
	// factory 패턴.
	// Grpc 서버 객체 생성시 내부 필드들 encapsulation.
	s.port = cfg.Port
	s.cert = cfg.Certificate
	s.key = cfg.Key
	s.storage = cfg.Storage
	return
}

func (s *GrpcServer) Listen() (err error) {
	var (
		listener net.Listener
		grpcOpts = []grpc.ServerOption{}
		grpcCreds credentials.TransportCredentials
	)

	listener, err = net.Listen("tcp", ":" + strconv.Itoa(s.port))
	if ErrExists(err) {
		err = errors.Wrapf(err, "failed to listen on port %d", s.port)
		return
	}

	if s.cert != "" && s.key != "" {
		grpcCreds, err = credentials.NewServerTLSFromFile(s.cert, s.key)

		if ErrExists(err) {
			err = errors.Wrapf(err,
				"failed to create tls service server using cert %s and key %s",
				s.cert, s.key)
			return
		}
		grpcOpts = append(grpcOpts, grpc.Creds(grpcCreds))
	}
	s.server = grpc.NewServer(grpcOpts...)
	messaging.RegisterImageTransferServer(s.server, s)

	err = s.server.Serve(listener)
	if ErrExists(err) {
		err = errors.Wrapf(err, "errored listening for service connections.")
		return
	}
	return
}

func (s *GrpcServer) SendImage(stream messaging.ImageTransfer_SendImageServer) (err error) {
	f, err := os.Create(s.storage + "/" + strconv.Itoa(s.counter) + ".jpeg")
	ErrExists(err)
	defer f.Close()
	w := bufio.NewWriter(f)
	var nn, n int

	for {
		b, err := stream.Recv()
		if ErrExists(err) {
			if err == io.EOF {
				goto END
			}
			err = errors.Wrapf(err,
				"failed unexpectedly while reading chunks from stream")
			return
		}

		nn, err = w.Write(b.GetContent())
		ErrExists(err)
		n += nn
	}
END:
	s.logger.Info().Msg("Data received.")
	w.Flush()
	fmt.Printf("Wrote file in byte %d\n",n)

	err = stream.SendAndClose(&messaging.TransferStatus{
		Message: "Image Transfer received with success.",
		StatusCode: messaging.TransStatCode_Ok,
	})

	if ErrExists(err) {
		err = errors.Wrapf(err, "failed to send status code.")
		return
	}
	return
}

func (s *GrpcServer) Close() {
	if s.server != nil {
		s.server.Stop()
	}
	return
}

func ErrExists (err error) bool{
	if err == nil {
		return false
	} else {
		return true
	}
}
