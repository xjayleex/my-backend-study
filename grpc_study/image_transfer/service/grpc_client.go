package service

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	messaging "github.com/xjayleex/idl/protos/imageproto"
	"google.golang.org/grpc/credentials"
	"io"
	"os"
	"time"
)

type GrpcClient struct {
	logger		zerolog.Logger
	conn		*grpc.ClientConn
	client		messaging.ImageTransferClient
	chunkSize	int
}

type GrpcClientConfig struct {
	Address			string
	ChunkSize		int
	RootCertificate string
	Compress 		bool
}

func NewGprcClient(cfg GrpcClientConfig) (c GrpcClient, err error) {
	var (
		grpcOpts = []grpc.DialOption{}
		grpcCreds credentials.TransportCredentials
	)

	if cfg.Address == "" {
		err = errors.Errorf("Address must be specified.")
		return
	}

	if cfg.Compress {
		grpcOpts = append(grpcOpts,
			grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
	}

	if cfg.RootCertificate == "" {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
	} else {
		grpcCreds, err = credentials.NewClientTLSFromFile(cfg.RootCertificate,"")
		if err != nil {
			err = errors.Wrapf(err,
				"failed to create grpc tls client with root-cert %s",
				cfg.RootCertificate)
			return
		}
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(grpcCreds))
	}

	switch {
	case cfg.ChunkSize == 0:
		err = errors.Errorf("ChunkSize must be specified.")
	case cfg.ChunkSize > (1 << 22):
		err = errors.Errorf("ChunkSize must be lower than 4MiB")
		return
	default:
		c.chunkSize = cfg.ChunkSize
	}

	c.logger = zerolog.New(os.Stdout).With().
		Str("from","client").Logger()
	c.conn, err = grpc.Dial(cfg.Address, grpcOpts...)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to start grpc connection with address %s",
			cfg.Address)
		return
	}
	c.client = messaging.NewImageTransferClient(c.conn)
	return
}

func (c *GrpcClient) TransferImageFile(ctx context.Context, path string) (stats Stats, err error) {
	var (
		writing = true
		buf		[]byte
		n		int
		file	*os.File
		status	*messaging.TransferStatus
	)

	file, err = os.Open(path)
	if err != nil {
		err = errors.Wrapf(err, "failed opening file %s", path)
		return
	}
	defer file.Close()

	stream, err := c.client.SendImage(ctx)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to create upload stream for file %s",
			path)
		return
	}
	defer stream.CloseSend()

	stats.StartTimeStamp = time.Now()
	buf = make([]byte, c.chunkSize)
	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}
			err = errors.Wrapf(err,
				"Error on copying file to byte buffer")
			return
		}

		err = stream.Send(&messaging.Chunk{
			Content: buf[:n],
		})

		if err != nil {
			err = errors.Wrapf(err, "failed on sending data via stream" )
			return
		}
	}

	stats.EndTimeStamp = time.Now()

	status, err = stream.CloseAndRecv()
	if err != nil {
		err = errors.Wrapf(err,
			"failed to receive status response")
		return
	}

	if status.StatusCode != messaging.TransStatCode_Ok {
		err = errors.Errorf(
			"Upload failed = msg : %s",
			status.Message)
		return
	}
	return
}

func (c *GrpcClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}