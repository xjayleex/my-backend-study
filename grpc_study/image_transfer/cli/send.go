package cli

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/xjayleex/my-backend-study/grpc_study/image_transfer/service"
	clitool "gopkg.in/urfave/cli.v2"
	"os"
)

var Send = clitool.Command{
	Name: "send",
	Usage: "Sends a image file.",
	Action: sendAction,
	Flags: []clitool.Flag{
		&clitool.StringFlag{
			Name: "address",
			Value: "localhost:8070",
			Usage: "Address of the server to connect",
		},
		&clitool.IntFlag{
			Name: "chunk-size",
			Usage: "size of the chunk messages (grpc only)",
			Value: (1 << 12),
		},
		&clitool.StringFlag{
			Name: "file",
			Usage: "file path to upload.",
		},
		&clitool.StringFlag{
			Name: "root-certificate",
			Usage: "path of a certificate to add to the root CAs",
		},
		/*&clitool.BoolFlag{
			Name: "http2",
			Usage: "whether or not to use http2 = requires root-certificate",
		},*/
		&clitool.BoolFlag{
			Name: "compress",
			Usage: "Whether or not to enable payload compression",
		},
	},
}

func sendAction(c *clitool.Context) (err error) {
	var (
		chunkSize		= c.Int("chunk-size")
		address			= c.String("address")
		file			= c.String("file")
		rootCertificate = c.String("root-certificate")
		compress		= c.Bool("compress")
		client			service.Client
	)

	if address == "" {
		Trap(errors.New("Address required."))
	}
}

