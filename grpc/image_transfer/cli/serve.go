package cli

import (
	clitool "gopkg.in/urfave/cli.v2"
	"github.com/xjayleex/my-backend-study/grpc_study/image_transfer/service"
)

var Serve = clitool.Command{
	Name: "serve",
	Usage: "To run a GRPC server.",
	Action: serveAction,
	Flags: []clitool.Flag{
		&clitool.IntFlag{
			Name: "port",
			Usage: "bind port.",
			Value: 8070,
		},
		/*&clitool.BoolFlag{
			Name: "http2",
			Usage: "If set true, use http2 instead GRPC.",
		},*/
		&clitool.StringFlag{
			Name: "certificate",
			Usage: "path to TLS certificate.",
		},
		&clitool.StringFlag{
			Name: "key",
			Usage: "path to key file.",
		},
		&clitool.StringFlag{
			Name: "storage",
			Value: "/home/jayhadoop/storage",
			Usage: "Storage path for save images on server",
		},
	},
}

func serveAction(c *clitool.Context) (err error) {
	var (
		port		= c.Int("port")
		key			= c.String("key")
		cert		= c.String("certificate")
		storage 	= c.String("storage")
		server		service.Server
	)
	grpcServer, err := service.NewGrpcServer(service.GrpcServerConfig{
		Port: port,
		Certificate: cert,
		Key: key,
		Storage: storage,
	})
	trap(err)
	server = &grpcServer
	err = server.Listen()
	trap(err)
	defer server.Close()

	return
}
